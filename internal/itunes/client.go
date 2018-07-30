package itunes

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/objque/musicmash/internal/config"
	"github.com/objque/musicmash/internal/log"
	"github.com/pkg/errors"
)

const (
	htmlTagTime      = `<time data-test-we-datetime datetime="`
	htmlTagReleaseID = `class="featured-album targeted-link"`
)

var exp = regexp.MustCompile(`.*\/(\d+)`)

func decode(buffer []byte) (*LastRelease, error) {
	parts := strings.Split(string(buffer), htmlTagTime)
	if len(parts) != 2 {
		return nil, errors.New("after split by a time-html tag we have not 2 parts")
	}

	// Jul 18, 2018" aria-label="July 18 ...
	released := strings.Split(parts[1], `"`)[0]
	t, err := parseTime(released)
	if err != nil {
		return nil, errors.Wrapf(err, "can't parse time '%s'", released)
	}

	parts = strings.Split(strings.Split(parts[0], htmlTagReleaseID)[0], `<a href="`)
	releaseURL := parts[len(parts)-1]
	releaseID := exp.FindStringSubmatch(releaseURL)
	if len(releaseID) != 2 {
		return nil, fmt.Errorf("found too many substrings by regex in '%s'", releaseURL)
	}

	id, err := strconv.ParseUint(releaseID[1], 10, 64)
	if err != nil {
		return nil, errors.Wrapf(err, "can't parse uint from '%s', fullURL: '%s'", releaseID[1], releaseURL)
	}
	return &LastRelease{
		ID:   id,
		Date: *t,
	}, nil
}

func GetArtistInfo(id uint64) (*LastRelease, error) {
	url := fmt.Sprintf("%s/us/artist/%d", config.Config.StoreURL, id)
	log.Debugf("Requesting '%s'...", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrapf(err, "can't receive page '%s'", url)
	}

	buffer, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, errors.Wrapf(err, "can't read response '%s'", url)
	}

	info, err := decode(buffer)
	if err != nil {
		return nil, errors.Wrapf(err, "can't decode '%s'", url)
	}
	log.Debugf("Last release on '%s'", info.Date)
	return info, nil
}
