package push_token

import (
	"cetm_booking/x/db/mongodb"
)

type PushToken struct {
	mongodb.BaseModel `bson:",inline"`
	UserId            string `bson:"user_id" json:"user_id" validate:"required"`
	Role              int    `bson:"role" json:"role"`
	IsRevoke          bool   `bson:"is_revoke" json:"is_revoke"`
	DeviceId          string `bson:"device_id" json:"device_id"`
	Platform          string `bson:"platform" json:"platform"`
	VersionApp        string `bson:"version_app" json:"version_app"`
	PushToken         string `bson:"push_token" json:"push_token"`
}

var PushTokenTable = mongodb.NewTable("push_token", "k", 80)
