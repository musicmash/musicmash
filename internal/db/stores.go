package db

import (
	"github.com/jinzhu/gorm"
	"github.com/musicmash/musicmash/internal/log"
)

type Store struct {
	Name string `gorm:"primary_key" json:"name"`
}

type StoreMgr interface {
	IsStoreExists(name string) bool
	EnsureStoreExists(name string) error
	GetAllStores() ([]*Store, error)
}

func (mgr *AppDatabaseMgr) IsStoreExists(name string) bool {
	store := Store{}
	if err := mgr.db.Where("name = ?", name).First(&store).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return false
		}

		log.Error(err)
		return false
	}
	return true
}

func (mgr *AppDatabaseMgr) EnsureStoreExists(name string) error {
	if !mgr.IsStoreExists(name) {
		return mgr.db.Create(&Store{Name: name}).Error
	}
	return nil
}

func (mgr *AppDatabaseMgr) GetAllStores() ([]*Store, error) {
	stores := []*Store{}
	err := mgr.db.Find(&stores).Error
	if err != nil {
		return nil, err
	}
	return stores, nil
}
