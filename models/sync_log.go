package models

import (
	log "github.com/sirupsen/logrus"
	"rtcgw/db"
	"time"
)

type SyncLog struct {
	ID      int64     `db:"id" json:"id"`
	ECHISID string    `db:"echis_id" json:"echis_id"`
	EVENTID string    `db:"event_id" json:"event_id""`
	Created time.Time `db:"created" json:"created"`
	Updated time.Time `db:"updated" json:"updated"`
}

func (s *SyncLog) Save() {
	dbConn := db.GetDB()
	_, err := dbConn.NamedExec(`INSERT INTO sync_log (echis_id, event_id) VALUES (:echis_id, :event_id)`, s)
	if err != nil {
		log.WithError(err).Error("Failed to save sync log")
	}
}

func GetSyncLogByECHISID(echisID string) (*SyncLog, error) {
	logObj := SyncLog{}
	err := db.GetDB().Get(&logObj, "SELECT * FROM sync_log WHERE echis_id = $1", echisID)
	if err != nil {
		return nil, err
	}
	return &logObj, nil
}
