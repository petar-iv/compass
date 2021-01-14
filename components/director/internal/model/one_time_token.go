package model

import "time"

type OneTimeToken struct {
	Token        string
	ConnectorURL string
	CreatedAt    time.Time
	Used         bool
	UsedAt       time.Time
}
