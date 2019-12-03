package api

import (
	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/testutil"
	"github.com/musicmash/musicmash/pkg/api/notifysettings"
	"github.com/stretchr/testify/assert"
)

func (t *testAPISuite) TestNotificationSettings_Create() {
	// action
	err := notifysettings.Create(t.client, testutil.UserObjque, &notifysettings.Settings{
		Service: "telegram",
		Data:    "chat-id-here",
	})

	// assert
	assert.NoError(t.T(), err)
	settings, err := notifysettings.List(t.client, testutil.UserObjque)
	assert.NoError(t.T(), err)
	assert.Len(t.T(), settings, 1)
	assert.Equal(t.T(), "telegram", settings[0].Service)
	assert.Equal(t.T(), "chat-id-here", settings[0].Data)
}

func (t *testAPISuite) TestNotificationSettings_Create_AlreadyExists() {
	// arrange
	assert.NoError(t.T(), db.DbMgr.EnsureNotificationSettingsExists(&db.NotificationSettings{
		UserName: testutil.UserObjque,
		Service:  "telegram",
		Data:     "chat-id-here",
	}))

	// action
	err := notifysettings.Create(t.client, testutil.UserObjque, &notifysettings.Settings{
		Service: "telegram",
		Data:    "chat-id-here",
	})

	// assert
	assert.Error(t.T(), err)
}

func (t *testAPISuite) TestNotificationSettings_Update() {
	// arrange
	assert.NoError(t.T(), db.DbMgr.EnsureNotificationSettingsExists(&db.NotificationSettings{
		UserName: testutil.UserObjque,
		Service:  "icq",
		Data:     "chat-id-here",
	}))
	assert.NoError(t.T(), db.DbMgr.EnsureNotificationSettingsExists(&db.NotificationSettings{
		UserName: testutil.UserObjque,
		Service:  "telegram",
		Data:     "chat-id-here",
	}))

	// action
	err := notifysettings.Update(t.client, testutil.UserObjque, &notifysettings.Settings{
		Service: "telegram",
		Data:    "new-chat-id-here",
	})

	// assert
	assert.NoError(t.T(), err)
	settings, err := notifysettings.List(t.client, testutil.UserObjque)
	assert.NoError(t.T(), err)
	assert.Len(t.T(), settings, 2)
	assert.Equal(t.T(), "icq", settings[0].Service)
	assert.Equal(t.T(), "chat-id-here", settings[0].Data)
	assert.Equal(t.T(), "telegram", settings[1].Service)
	assert.Equal(t.T(), "new-chat-id-here", settings[1].Data)
}
