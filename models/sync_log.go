package models

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
	"rtcgw/db"
	"time"
)

type SyncLog struct {
	ID            int64        `db:"id" json:"id"`
	ECHISID       string       `db:"echis_id" json:"echis_id"`
	EventID       string       `db:"event_id" json:"event_id"`
	TrackedEntity string       `db:"tracked_entity" json:"trackedEntityInstance"`
	EventDate     sql.NullTime `db:"event_date" json:"event_date"`
	Created       time.Time    `db:"created" json:"created"`
	Updated       time.Time    `db:"updated" json:"updated"`
}

func (s *SyncLog) Save() {
	dbConn := db.GetDB()
	_, err := dbConn.NamedExec(`INSERT INTO sync_log 
	(echis_id, event_id, event_date, tracked_entity) 
		VALUES (:echis_id, :event_id, :event_date, :tracked_entity)`, s)
	if err != nil {
		log.WithError(err).Error("Failed to save sync log")
	}
}

// StringToNullTime converts string to sql.NullTime (Handles NULL values)
func StringToNullTime(s sql.NullString) sql.NullTime {
	if !s.Valid {
		return sql.NullTime{Valid: false}
	}

	// Use custom format to match the actual timestamp structure
	const customTimeFormat = "2006-01-02 15:04:05.999999Z"
	parsedTime, err := time.Parse(customTimeFormat, s.String)
	if err != nil {
		log.Println("Error parsing time:", err)
		return sql.NullTime{Valid: false} // Return invalid NullTime if parsing fails
	}

	return sql.NullTime{Time: parsedTime, Valid: true}
}

func GetSyncLogByECHISID(echisID string) (*SyncLog, error) {
	logObj := SyncLog{}
	var eventDateStr sql.NullString

	err := db.GetDB().QueryRow(
		`SELECT id, echis_id, event_id, tracked_entity, event_date FROM sync_log WHERE echis_id = $1`, echisID).
		Scan(&logObj.ID, &logObj.ECHISID, &logObj.EventID, &logObj.TrackedEntity, &eventDateStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	logObj.EventDate = StringToNullTime(eventDateStr)

	return &logObj, nil
}
