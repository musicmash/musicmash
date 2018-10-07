package db

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Release struct {
	ID         uint64
	CreatedAt  time.Time
	ArtistName string
	Title      string
	Poster     string
	Released   time.Time
	StoreName  string `gorm:"unique_index:idx_rel_store_name_store_id"`
	StoreID    string `gorm:"unique_index:idx_rel_store_name_store_id"`
}

type ReleaseMgr interface {
	EnsureReleaseExists(release *Release) error
	GetAllReleases() ([]*Release, error)
	GetReleasesForUserFilterByPeriod(userName string, since, till time.Time) ([]*Release, error)
	GetReleasesForUserSince(userName string, since time.Time) ([]*Release, error)
	FindNewReleases(date time.Time) ([]*Release, error)
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
	var releases = make([]*Release, 0)
	return releases, mgr.db.Find(&releases).Error
}

func (mgr *AppDatabaseMgr) GetReleasesForUserFilterByPeriod(userName string, since, till time.Time) ([]*Release, error) {
	// inner query: select artist_name from subscriptions where user_name = XXX
	// select * from releases where artist_name in (INNER) and and released >= ? and released <= ?
	releases := []*Release{}
	const query = "select artist_name from subscriptions where user_name = ?"
	innerQuery := mgr.db.Raw(query, userName).QueryExpr()
	where := mgr.db.Where("artist_name in (?) and released >= ? and released <= ?", innerQuery, since, till)
	if err := where.Find(&releases).Error; err != nil {
		return nil, err
	}
	return releases, nil
}

func (mgr *AppDatabaseMgr) GetReleasesForUserSince(userName string, since time.Time) ([]*Release, error) {
	// inner query: select artist_name from subscriptions where user_name = XXX
	// select * from releases where artist_name in (INNER) and and released >= ?
	releases := []*Release{}
	const query = "select artist_name from subscriptions where user_name = ?"
	innerQuery := mgr.db.Raw(query, userName).QueryExpr()
	where := mgr.db.Where("artist_name in (?) and released >= ?", innerQuery, since)
	if err := where.Find(&releases).Error; err != nil {
		return nil, err
	}
	return releases, nil
}

func (mgr *AppDatabaseMgr) FindNewReleases(date time.Time) ([]*Release, error) {
	releases := []*Release{}
	if err := mgr.db.Where("created_at >= ?", date).Find(&releases).Error; err != nil {
		return nil, err
	}
	return releases, nil
}
