package models

import (
	"database/sql"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/goccy/go-json"
	log "github.com/sirupsen/logrus"
	"rtcgw/clients"
	"rtcgw/config"
	"rtcgw/db"
	"rtcgw/utils"
	"time"
)

type SyncLog struct {
	ID                        int64        `db:"id" json:"id"`
	ECHISID                   string       `db:"echis_id" json:"echis_id"`
	EventID                   string       `db:"event_id" json:"event_id"`
	TrackedEntity             string       `db:"tracked_entity" json:"trackedEntityInstance"`
	EventDate                 sql.NullTime `db:"event_date" json:"event_date"`
	OrgUnit                   string       `db:"org_unit" json:"org_unit"`
	ECHISClientCreationErrors string       `db:"echis_clinet_creation_errors" json:"echisClientCreationErrors"`
	ResultsUpdated            bool         `db:"results_updated" json:"results_updated"`
	ResultsUpdateErrors       string       `db:"results_update_errors" json:"resultsUpdateErrors"`
	LabEvent                  string       `db:"lab_event" json:"lab_event"`
	LabEnrollment             string       `db:"lab_enrollment" json:"lab_enrollment"`
	Created                   time.Time    `db:"created" json:"created"`
	Updated                   time.Time    `db:"updated" json:"updated"`
}

func (s *SyncLog) Save() {
	dbConn := db.GetDB()
	_, err := dbConn.NamedExec(`INSERT INTO sync_log 
	(echis_id, event_id, event_date, tracked_entity, org_unit) 
		VALUES (:echis_id, :event_id, :event_date, :tracked_entity, :org_unit)`, s)
	if err != nil {
		log.WithError(err).Error("Failed to save sync log")
	}
}

func (s *SyncLog) SetResultUpdated() {
	s.ResultsUpdated = true
	dbConn := db.GetDB()
	_, err := dbConn.Exec(`UPDATE sync_log SET results_updated = $1 WHERE id = $2`, s.ResultsUpdated, s.ID)
	if err != nil {
		log.WithError(err).Error("Failed to update sync log")
	}
}

// SetResultsUpdateErrors ...
func (s *SyncLog) SetResultsUpdateErrors(errors string) {
	s.ResultsUpdateErrors = errors
	dbConn := db.GetDB()
	_, err := dbConn.Exec(`UPDATE sync_log SET results_update_errors = $1 WHERE id = $2`, s.ResultsUpdateErrors, s.ID)
	if err != nil {
		log.WithError(err).Error("Failed to update sync log")
	}
}

// SetECHISClientCreationErrors ...
func (s *SyncLog) SetECHISClientCreationErrors(errors string) {
	s.ECHISClientCreationErrors = errors
	dbConn := db.GetDB()
	_, err := dbConn.Exec(`UPDATE sync_log SET echis_client_creation_errors = $1 WHERE id = $2`, s.ECHISClientCreationErrors, s.ID)
	if err != nil {
		log.WithError(err).Error("Failed to update sync log")
	}
}

// SetLabEnrollment ...
func (s *SyncLog) SetLabEnrollment(enrollment string) {
	s.LabEnrollment = enrollment
	dbConn := db.GetDB()
	_, err := dbConn.Exec(`UPDATE sync_log SET lab_enrollment = $1 WHERE id = $2`, s.LabEnrollment, s.ID)
	if err != nil {
		log.WithError(err).Error("Failed to update sync log")
	}
}

// SetLabEvent ...
func (s *SyncLog) SetLabEvent(event string) {
	s.LabEvent = event
	dbConn := db.GetDB()
	_, err := dbConn.Exec(`UPDATE sync_log SET lab_event = $1 WHERE id = $2`, s.LabEvent, s.ID)
	if err != nil {
		log.WithError(err).Error("Failed to update sync log")
	}
}

// StringToNullTime converts string to sql.NullTime (Handles NULL values)
func StringToNullTime(s sql.NullString) sql.NullTime {
	if !s.Valid {
		return sql.NullTime{Valid: false}
	}

	// Define multiple time layouts to try
	layouts := []string{
		"2006-01-02 15:04:05.999999-07:00", // for timestamps with a timezone offset (e.g., +03:00)
		"2006-01-02 15:04:05.999999Z",      // for UTC timestamps with literal Z
		// Add more layouts here if needed
	}

	var parsedTime time.Time
	var err error
	for _, layout := range layouts {
		parsedTime, err = time.Parse(layout, s.String)
		if err == nil {
			return sql.NullTime{Time: parsedTime, Valid: true}
		}
	}

	// If no layout matches, log the error and return invalid NullTime
	log.Println("Error parsing time:", err)
	return sql.NullTime{Valid: false}
}

func DHIS2EventExists(event string) bool {
	resp, err := clients.Dhis2Client.GetResource(fmt.Sprintf("events/%s?fields=uid", event), nil)
	if err != nil || !resp.IsSuccess() {
		log.Infof("Error checking DHIS2 event existence: %v: %v", err, string(resp.Body()))
		return false
	}
	return resp.IsSuccess()
}
func (s *SyncLog) LabEventExists() bool {
	if s.LabEvent == "" {
		return false
	}
	return DHIS2EventExists(s.LabEvent)
}

func GetSyncLogByECHISID(echisID string) (*SyncLog, error) {
	logObj := SyncLog{}
	var eventDateStr, labEvent, labEnrollment sql.NullString

	err := db.GetDB().QueryRow(
		`SELECT id, echis_id, event_id, tracked_entity, event_date, results_updated, 
		lab_event, lab_enrollment, org_unit
		FROM sync_log WHERE echis_id = $1`, echisID).
		Scan(&logObj.ID, &logObj.ECHISID, &logObj.EventID,
			&logObj.TrackedEntity, &eventDateStr, &logObj.ResultsUpdated,
			&labEvent, &labEnrollment, &logObj.OrgUnit)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	logObj.EventDate = StringToNullTime(eventDateStr)
	logObj.LabEvent = labEvent.String
	logObj.LabEnrollment = labEnrollment.String
	return &logObj, nil
}

func SyncLogCreationLastXDays(numberOfDays int) []string {
	days := make([]string, numberOfDays)
	for i := range days {
		days[i] = fmt.Sprintf("0")
	}
	// select from sync_log table where echis_id and event_id and tracked_entity are not empty between now and x days back and return a count for each day
	interval := fmt.Sprintf("%d days", numberOfDays)
	rows, err := db.GetDB().Query(`SELECT count(distinct echis_id), to_char(created, 'YYYY-mm-dd') FROM sync_log 
		WHERE created > NOW() - $1::interval group by to_char(created, 'YYYY-mm-dd')`, interval)
	if err != nil {
		log.WithError(err).Error("Failed to get sync log creation last X days")
		return nil
	}
	defer rows.Close()
	dayDates := utils.LastXDays(numberOfDays)

	for rows.Next() {
		var day string
		var count int64
		err = rows.Scan(&count, &day)
		if err != nil {
			log.WithError(err).Error("Failed to scan sync log creation last X days")
			continue
		}
		idx := utils.IndexOf(dayDates, day)
		if idx != -1 {
			days[idx] = fmt.Sprintf("%d", count)
		} else {
			log.WithField("day", day).Error("Day not found in last X days")
		}
	}

	if err = rows.Err(); err != nil {
		log.WithError(err).Error("Error scanning sync log creation last X days")
		return nil
	}
	return days
}

func SyncLogUpdateLastXDays(numberOfDays int) []string {
	days := make([]string, numberOfDays)
	for i := range days {
		days[i] = fmt.Sprintf("0")
	}
	interval := fmt.Sprintf("%d days", numberOfDays)
	rows, err := db.GetDB().Query(`SELECT count(distinct echis_id), to_char(created, 'YYYY-mm-dd') FROM sync_log 
        WHERE results_updated = TRUE AND created > NOW() - $1::interval group by to_char(created, 'YYYY-mm-dd')`, interval)
	if err != nil {
		log.WithError(err).Error("Failed to get sync log update last X days")
		return nil
	}
	defer rows.Close()
	dayDates := utils.LastXDays(numberOfDays)

	for rows.Next() {
		var day string
		var count int64
		err = rows.Scan(&count, &day)
		if err != nil {
			log.WithError(err).Error("Failed to scan sync log creation last X days")
			continue
		}
		idx := utils.IndexOf(dayDates, day)
		if idx != -1 {
			days[idx] = fmt.Sprintf("%d", count)
		} else {
			log.WithField("day", day).Error("Day not found in last X days")
		}
	}
	if err = rows.Err(); err != nil {
		log.WithError(err).Error("Error scanning sync log creation last X days")
		return nil
	}
	return days
}

func (s *SyncLog) CheckLabProgramEnrollment() bool {
	params := map[string]string{
		"trackedEntityInstance": s.TrackedEntity,
		"program":               config.RTCGwConf.API.DHIS2LaboratoryProgram,
		"fields":                "enrollment",
		"ou":                    s.OrgUnit,
		"skipPaging":            "true",
	}
	log.Infof("Checking Enrollment for %v", params)

	resp, err := clients.Dhis2Client.GetResource("enrollments", params)
	if err != nil || !resp.IsSuccess() {
		log.Infof("Error checking lab program enrollment: %v: %v", err, string(resp.Body()))
		return false
	}
	v, _, _, err := jsonparser.Get(resp.Body(), "enrollments")
	if err != nil {
		log.Infof("Error getting enrollements: %v", err)
		return false
	}
	var enrollments []map[string]string
	err = json.Unmarshal(v, &enrollments)
	if err != nil {
		log.Infof("Error unmarshalling enrollments: Resp: %v -- %v, Error:", string(v), enrollments, err.Error())
		return false
	}
	if len(enrollments) > 0 {
		s.SetLabEnrollment(enrollments[0]["enrollment"])
		log.Infof("Found Lab Program enrollment for patient %v, TE: %v,  %v", s.ECHISID, s.TrackedEntity, s.LabEnrollment)
		return true
	}
	return false
}
