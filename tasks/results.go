package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	log "github.com/sirupsen/logrus"
	"rtcgw/clients"
	"rtcgw/config"
	"rtcgw/models"
	"rtcgw/models/tracker"
	"time"
)

const (
	TypeSendResults = "results:send"
)

func NewResultsTask(request models.LabXpertResult) (*asynq.Task, error) {
	payload, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeSendResults, payload), nil
}

func HandleResultsTask(cxt context.Context, task *asynq.Task) error {
	var result models.LabXpertResult
	if err := json.Unmarshal(task.Payload(), &result); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	// Send the results to DHIS2
	log.Printf("Sending result to DHIS2: %v", result.PatientID)
	if patientLog, ok := result.InDhis2(); ok {
		tbResult, diagnosed := result.GetResult()
		log.Infof("Patient found: %v with result: %s and event: %s", patientLog.ECHISID, tbResult, patientLog.EventID)
		resultsDate, err := time.Parse("2006-01-02 15:04:05", result.ResultDate)
		if err != nil {
			fmt.Println("Error parsing date:", err)
			return err
		}
		dv := []tracker.DataValue{
			{
				DataElement: config.RTCGwConf.API.DHIS2Mapping["data_elements"]["results"],
				Value:       tbResult,
			},
			{
				DataElement: config.RTCGwConf.API.DHIS2Mapping["data_elements"]["results_date"],
				Value:       resultsDate.Format("2006-01-02"),
			},
			{
				DataElement: config.RTCGwConf.API.DHIS2Mapping["data_elements"]["diagnosed"],
				Value:       diagnosed,
			},
		}
		// Iterate through dv and create EventUpdatePayload and send to DHIS2
		for _, v := range dv {
			ep := tracker.EventUpdatePayload{
				Event:         patientLog.EventID,
				Program:       config.RTCGwConf.API.DHIS2TrackerProgram,
				OrgUnit:       result.FacilityID,
				Status:        "ACTIVE",
				ProgramStage:  config.RTCGwConf.API.DHIS2TrackerProgramStage,
				DataValues:    []tracker.DataValue{v},
				TrackedEntity: patientLog.TrackedEntity,
			}
			// utils
			jsonData, err := json.MarshalIndent(ep, "", "  ")
			if err != nil {
				log.Infof("Error marshaling JSON: %v", err)
				return err
			}
			log.Infof("JSON EventUpdate Payload: %s", jsonData)

			resp, err := clients.Dhis2Client.PutResource(
				fmt.Sprintf("events/%s/%s", patientLog.EventID, v.DataElement), ep)
			if err != nil || !resp.IsSuccess() {
				log.Infof("Error sending result to DHIS2: %v: %v", err, string(resp.Body()))
				continue
			}
		}

	}

	log.Printf("Done sending result to DHIS2: %v", result.PatientID)
	return nil
}
