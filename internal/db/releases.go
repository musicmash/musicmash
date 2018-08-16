package db

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/objque/musicmash/internal/log"
)

type Release struct {
	ID         int64     `gorm:"primary_key" sql:"AUTO_INCREMENT" json:"-"`
	CreatedAt  time.Time `json:"-"`
	Date       time.Time `gorm:"not null" sql:"index" json:"date"`
	ArtistName string    `json:"artist_name"`
	StoreID    uint64    `sql:"index" json:"store_id"`
}

type ReleaseMgr interface {
	CreateRelease(release *Release) error
	FindRelease(artist string, storeID uint64) (*Release, error)
	IsReleaseExists(storeID uint64) bool
	GetAllReleases() ([]*Release, error)
	EnsureReleaseExists(release *Release) error
	GetReleasesForUserFilterByPeriod(userID string, since, till time.Time) ([]*Release, error)
	GetReleasesForUserSince(userID string, since time.Time) ([]*Release, error)
}

func (mgr *AppDatabaseMgr) FindRelease(artist string, storeID uint64) (*Release, error) {
	release := Release{}
	if err := mgr.db.Where("artist_name = ? and store_id = ?", artist, storeID).First(&release).Error; err != nil {
		return nil, err
	}
	return &release, nil
}

func (mgr *AppDatabaseMgr) IsReleaseExists(storeID uint64) bool {
	release := Release{}
	if err := mgr.db.Where("store_id = ?", storeID).First(&release).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return false
		}

		log.Error(err)
		return false
	}
	return true
}

func (mgr *AppDatabaseMgr) GetAllReleases() ([]*Release, error) {
	var releases = make([]*Release, 0)
	return releases, mgr.db.Find(&releases).Error
}

func (mgr *AppDatabaseMgr) CreateRelease(release *Release) error {
	return mgr.db.Create(release).Error
}

func (mgr *AppDatabaseMgr) EnsureReleaseExists(release *Release) error {
	_, err := mgr.FindRelease(release.ArtistName, release.StoreID)
	if err != nil {
		return mgr.CreateRelease(release)
	}
	return nil
}

func (mgr *AppDatabaseMgr) GetReleasesForUserFilterByPeriod(userID string, since, till time.Time) ([]*Release, error) {
	const sub = "select artist_name from subscriptions where user_id = ?"
	const sql = "select * from releases where artist_name in (" + sub + ") and date >= ? and date <= ? order by date desc"
	releases := []*Release{}
	if err := mgr.db.Raw(sql, userID, since, till).Scan(&releases).Error; err != nil {
		return nil, err
	}
	return releases, nil
}

func (mgr *AppDatabaseMgr) GetReleasesForUserSince(userID string, since time.Time) ([]*Release, error) {
	const sub = "select artist_name from subscriptions where user_id = ?"
	const sql = "select * from releases where artist_name in (" + sub + ") and date >= ? order by date desc"
	releases := []*Release{}
	if err := mgr.db.Raw(sql, userID, since).Scan(&releases).Error; err != nil {
		return nil, err
	}
	return releases, nil
}
