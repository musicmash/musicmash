package api

import (
	"time"

	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/testutils/vars"
	"github.com/musicmash/musicmash/internal/utils/ptr"
	"github.com/musicmash/musicmash/pkg/api/releases"
	"github.com/stretchr/testify/assert"
)

func (t *testAPISuite) fillRelease(release *db.Release) {
	assert.NoError(t.T(), db.Mgr.EnsureReleaseExists(release))
}

func (t *testAPISuite) prepareReleasesTestCase() {
	r := time.Now().UTC()
	monthAgo := r.AddDate(0, -1, 0)
	yearAgo := r.AddDate(-1, 0, 0)
	assert.NoError(t.T(), db.Mgr.EnsureArtistExists(&db.Artist{Name: vars.ArtistAlgorithm}))
	assert.NoError(t.T(), db.Mgr.EnsureArtistExists(&db.Artist{Name: vars.ArtistSkrillex}))
	assert.NoError(t.T(), db.Mgr.SubscribeUser(vars.UserObjque, []int64{1})) // ArtistAlgorithm
	assert.NoError(t.T(), db.Mgr.EnsureStoreExists(vars.StoreApple))
	assert.NoError(t.T(), db.Mgr.EnsureStoreExists(vars.StoreSpotify))
	assert.NoError(t.T(), db.Mgr.EnsureStoreExists(vars.StoreDeezer))
	// first release
	t.fillRelease(&db.Release{ArtistID: 1, Title: vars.ReleaseAlgorithmFloatingIP, Poster: vars.PosterSimple, Released: r, SpotifyID: "3000", Type: vars.ReleaseTypeAlbum, Explicit: true, TracksCount: 10, DurationMs: 25})
	// second release
	t.fillRelease(&db.Release{ArtistID: 1, Title: vars.ReleaseArchitectsHollyHell, Poster: vars.PosterMiddle, Released: monthAgo, SpotifyID: "1100", Type: vars.ReleaseTypeSong, Explicit: false, TracksCount: 10, DurationMs: 25})
	// third release from another artist
	t.fillRelease(&db.Release{ArtistID: 2, Title: vars.ReleaseRitaOraLouder, Poster: vars.PosterGiant, Released: yearAgo, SpotifyID: "1110", Type: vars.ReleaseTypeVideo, Explicit: true, TracksCount: 10, DurationMs: 25})
}

func (t *testAPISuite) TestReleases_Get_All() {
	// arrange
	t.prepareReleasesTestCase()

	// action
	rels, err := releases.List(t.client, nil)

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), rels, 3)
	// releases are sort by release date desc
	assert.Equal(t.T(), int64(1), rels[0].ArtistID)
	assert.Equal(t.T(), vars.ReleaseAlgorithmFloatingIP, rels[0].Title)
	assert.Equal(t.T(), vars.PosterSimple, rels[0].Poster)
	assert.Equal(t.T(), "3000", rels[0].SpotifyID)
	assert.Equal(t.T(), vars.ReleaseTypeAlbum, rels[0].Type)
	assert.Equal(t.T(), int32(10), rels[0].TracksCount)
	assert.Equal(t.T(), int64(25), rels[0].DurationMs)
	assert.True(t.T(), rels[0].Explicit)

	assert.Equal(t.T(), int64(1), rels[1].ArtistID)
	assert.Equal(t.T(), vars.ReleaseArchitectsHollyHell, rels[1].Title)
	assert.Equal(t.T(), vars.PosterMiddle, rels[1].Poster)
	assert.Equal(t.T(), "1100", rels[1].SpotifyID)
	assert.Equal(t.T(), vars.ReleaseTypeSong, rels[1].Type)
	assert.Equal(t.T(), int32(10), rels[1].TracksCount)
	assert.Equal(t.T(), int64(25), rels[1].DurationMs)
	assert.False(t.T(), rels[1].Explicit)

	assert.Equal(t.T(), int64(2), rels[2].ArtistID)
	assert.Equal(t.T(), vars.ReleaseRitaOraLouder, rels[2].Title)
	assert.Equal(t.T(), vars.PosterGiant, rels[2].Poster)
	assert.Equal(t.T(), "1110", rels[2].SpotifyID)
	assert.Equal(t.T(), vars.ReleaseTypeVideo, rels[2].Type)
	assert.Equal(t.T(), int32(10), rels[2].TracksCount)
	assert.Equal(t.T(), int64(25), rels[2].DurationMs)
	assert.True(t.T(), rels[2].Explicit)
}

func (t *testAPISuite) TestReleases_Get_All_ChangeSortType() {
	// arrange
	t.prepareReleasesTestCase()

	// action

	rels, err := releases.List(t.client, &releases.GetOptions{
		SortType: releases.SortTypeASC,
	})

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), rels, 3)
	// releases are sort by release date ASC!
	assert.Equal(t.T(), int64(2), rels[0].ArtistID)
	assert.Equal(t.T(), vars.ReleaseRitaOraLouder, rels[0].Title)
	assert.Equal(t.T(), vars.PosterGiant, rels[0].Poster)
	assert.Equal(t.T(), "1110", rels[0].SpotifyID)
	assert.Equal(t.T(), vars.ReleaseTypeVideo, rels[0].Type)
	assert.Equal(t.T(), int32(10), rels[0].TracksCount)
	assert.Equal(t.T(), int64(25), rels[0].DurationMs)
	assert.True(t.T(), rels[0].Explicit)

	assert.Equal(t.T(), int64(1), rels[1].ArtistID)
	assert.Equal(t.T(), vars.ReleaseArchitectsHollyHell, rels[1].Title)
	assert.Equal(t.T(), vars.PosterMiddle, rels[1].Poster)
	assert.Equal(t.T(), "1100", rels[1].SpotifyID)
	assert.Equal(t.T(), vars.ReleaseTypeSong, rels[1].Type)
	assert.Equal(t.T(), int32(10), rels[1].TracksCount)
	assert.Equal(t.T(), int64(25), rels[1].DurationMs)
	assert.False(t.T(), rels[1].Explicit)

	assert.Equal(t.T(), int64(1), rels[2].ArtistID)
	assert.Equal(t.T(), vars.ReleaseAlgorithmFloatingIP, rels[2].Title)
	assert.Equal(t.T(), vars.PosterSimple, rels[2].Poster)
	assert.Equal(t.T(), "3000", rels[2].SpotifyID)
	assert.Equal(t.T(), vars.ReleaseTypeAlbum, rels[2].Type)
	assert.Equal(t.T(), int32(10), rels[2].TracksCount)
	assert.Equal(t.T(), int64(25), rels[2].DurationMs)
	assert.True(t.T(), rels[2].Explicit)
}

func (t *testAPISuite) TestReleases_Get_FilterByLimit() {
	// arrange
	t.prepareReleasesTestCase()

	// action
	rels, err := releases.List(t.client, &releases.GetOptions{
		Limit: ptr.Uint(1),
	})

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), rels, 1)
	assert.Equal(t.T(), int64(1), rels[0].ArtistID)
	assert.Equal(t.T(), vars.ReleaseAlgorithmFloatingIP, rels[0].Title)
	assert.Equal(t.T(), vars.PosterSimple, rels[0].Poster)
	assert.Equal(t.T(), "3000", rels[0].SpotifyID)
	assert.Equal(t.T(), vars.ReleaseTypeAlbum, rels[0].Type)
	assert.Equal(t.T(), int32(10), rels[0].TracksCount)
	assert.Equal(t.T(), int64(25), rels[0].DurationMs)
	assert.True(t.T(), rels[0].Explicit)
}

func (t *testAPISuite) TestReleases_Get_FilterByArtistID_ReleaseType() {
	// arrange
	t.prepareReleasesTestCase()

	// action
	rels, err := releases.List(t.client, &releases.GetOptions{
		ReleaseType: vars.ReleaseTypeSong,
	})

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), rels, 1)
	assert.Equal(t.T(), int64(1), rels[0].ArtistID)
	assert.Equal(t.T(), vars.ReleaseArchitectsHollyHell, rels[0].Title)
	assert.Equal(t.T(), vars.PosterMiddle, rels[0].Poster)
	assert.Equal(t.T(), "1100", rels[0].SpotifyID)
	assert.Equal(t.T(), vars.ReleaseTypeSong, rels[0].Type)
	assert.Equal(t.T(), int32(10), rels[0].TracksCount)
	assert.Equal(t.T(), int64(25), rels[0].DurationMs)
	assert.False(t.T(), rels[0].Explicit)
}

func (t *testAPISuite) TestReleases_Get_FilterBySince() {
	// arrange
	t.prepareReleasesTestCase()

	// action
	since := time.Now().UTC().Truncate(time.Hour)
	rels, err := releases.List(t.client, &releases.GetOptions{
		Since: &since,
	})

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), rels, 1)
	assert.Equal(t.T(), int64(1), rels[0].ArtistID)
	assert.Equal(t.T(), vars.ReleaseAlgorithmFloatingIP, rels[0].Title)
	assert.Equal(t.T(), vars.PosterSimple, rels[0].Poster)
	assert.Equal(t.T(), "3000", rels[0].SpotifyID)
	assert.Equal(t.T(), vars.ReleaseTypeAlbum, rels[0].Type)
	assert.Equal(t.T(), int32(10), rels[0].TracksCount)
	assert.Equal(t.T(), int64(25), rels[0].DurationMs)
	assert.True(t.T(), rels[0].Explicit)
}

func (t *testAPISuite) TestReleases_Get_FilterByTill() {
	// arrange
	t.prepareReleasesTestCase()

	// action
	till := time.Now().UTC().Truncate(time.Hour).AddDate(0, -1, 0)
	rels, err := releases.List(t.client, &releases.GetOptions{
		Till: &till,
	})

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), rels, 1)
	assert.Equal(t.T(), int64(2), rels[0].ArtistID)
	assert.Equal(t.T(), vars.ReleaseRitaOraLouder, rels[0].Title)
	assert.Equal(t.T(), vars.PosterGiant, rels[0].Poster)
	assert.Equal(t.T(), "1110", rels[0].SpotifyID)
	assert.Equal(t.T(), vars.ReleaseTypeVideo, rels[0].Type)
	assert.Equal(t.T(), int32(10), rels[0].TracksCount)
	assert.Equal(t.T(), int64(25), rels[0].DurationMs)
	assert.True(t.T(), rels[0].Explicit)
}

func (t *testAPISuite) TestReleases_Get_FilterBy_SinceAndTill() {
	// arrange
	t.prepareReleasesTestCase()

	// action
	since := time.Now().UTC().Truncate(time.Hour).AddDate(-2, 0, 0)
	till := time.Now().UTC().Truncate(time.Hour).AddDate(0, -2, 0)
	rels, err := releases.List(t.client, &releases.GetOptions{
		Since: &since,
		Till:  &till,
	})

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), rels, 1)
	assert.Equal(t.T(), int64(2), rels[0].ArtistID)
	assert.Equal(t.T(), vars.ReleaseRitaOraLouder, rels[0].Title)
	assert.Equal(t.T(), vars.PosterGiant, rels[0].Poster)
	assert.Equal(t.T(), "1110", rels[0].SpotifyID)
	assert.Equal(t.T(), vars.ReleaseTypeVideo, rels[0].Type)
	assert.Equal(t.T(), int32(10), rels[0].TracksCount)
	assert.Equal(t.T(), int64(25), rels[0].DurationMs)
	assert.True(t.T(), rels[0].Explicit)
}

func (t *testAPISuite) TestReleases_Get_FilterBy_Explicit() {
	// arrange
	t.prepareReleasesTestCase()

	// action
	explicit := false
	rels, err := releases.List(t.client, &releases.GetOptions{
		Explicit: &explicit,
	})

	// assert
	assert.NoError(t.T(), err)
	assert.Len(t.T(), rels, 1)
	// releases are sort by release date desc
	assert.Equal(t.T(), int64(1), rels[0].ArtistID)
	assert.Equal(t.T(), vars.ReleaseArchitectsHollyHell, rels[0].Title)
	assert.Equal(t.T(), vars.PosterMiddle, rels[0].Poster)
	assert.Equal(t.T(), "1100", rels[0].SpotifyID)
	assert.Equal(t.T(), vars.ReleaseTypeSong, rels[0].Type)
	assert.Equal(t.T(), int32(10), rels[0].TracksCount)
	assert.Equal(t.T(), int64(25), rels[0].DurationMs)
	assert.False(t.T(), rels[0].Explicit)
}
