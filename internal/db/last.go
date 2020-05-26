package db

import "time"

const (
	ActionFetch  = "fetch"
	ActionNotify = "notify"
)

type LastAction struct {
	ID     int       `db:"id"`
	Date   time.Time `db:"date"`
	Action string    `db:"action"`
}

type LastActionMgr interface {
	GetLastActionDate(action string) (*LastAction, error)
	SetLastActionDate(action string, time time.Time) error
}

func (mgr *AppDatabaseMgr) GetLastActionDate(action string) (*LastAction, error) {
	const query = "select * from last_actions where action = $1"

	last := LastAction{}
	err := mgr.newdb.Get(&last, query, action)
	if err != nil {
		return nil, err
	}

	return &last, nil
}

func (mgr *AppDatabaseMgr) CreateLastAction(action string, date time.Time) error {
	const query = "insert into last_actions (date, action) values ($1, $2)"

	_, err := mgr.newdb.Exec(query, date, action)

	return err
}

func (mgr *AppDatabaseMgr) SetLastActionDate(action string, date time.Time) error {
	// TODO (m.kalinin): replace with postgres upsert
	_, err := mgr.GetLastActionDate(action)
	if err != nil {
		return mgr.CreateLastAction(action, date)
	}

	const query = "update last_actions set date = $1 where action = $2"

	_, err = mgr.newdb.Exec(query, date, action)

	return err
}
