package db

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Release struct {
	ID        uint64    `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	ArtistID  int64     `db:"artist_id"`
	Title     string    `db:"title"`
	Poster    string    `db:"poster"`
	Released  time.Time `db:"released"`
	StoreName string    `db:"store_name"`
	StoreID   string    `db:"store_id"`
	Type      string    `db:"type"`
	Explicit  bool      `db:"explicit"`
}

type ReleaseMgr interface {
	EnsureReleaseExists(release *Release) error
	GetAllReleases() ([]*Release, error)
	FindReleases(condition map[string]interface{}) ([]*Release, error)
	FindNewReleases(date time.Time) ([]*Release, error)
	InsertBatchNewReleases(releases []*Release) error
}

func (r *Release) IsComing() bool {
	// if release day tomorrow or later, than that means coming release is here
	return r.Released.After(time.Now().UTC().Truncate(24 * time.Hour))
}

func (mgr *AppDatabaseMgr) EnsureReleaseExists(release *Release) error {
	res := Release{}
	err := mgr.db.Where("store_id = ? and store_name = ?", release.StoreID, release.StoreName).First(&res).Error
	if gorm.IsRecordNotFoundError(err) {
		return mgr.db.Create(release).Error
	}
	return err
}

func (mgr *AppDatabaseMgr) GetAllReleases() ([]*Release, error) {
	var releases = []*Release{}
	return releases, mgr.db.Find(&releases).Error
}

func (mgr *AppDatabaseMgr) FindNewReleases(date time.Time) ([]*Release, error) {
	releases := []*Release{}
	if err := mgr.db.Where("created_at >= ?", date).Find(&releases).Error; err != nil {
		return nil, err
	}
	return releases, nil
}

func (mgr *AppDatabaseMgr) FindReleases(condition map[string]interface{}) ([]*Release, error) {
	releases := []*Release{}
	err := mgr.db.Where(condition).Find(&releases).Error
	if err != nil {
		return nil, err
	}
	return releases, nil
}
