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
	log.Printf("Sending result to DHIS2 for PatientID: %v", result.PatientID)
	resultUpdatde := false
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
			//jsonData, err := json.MarshalIndent(ep, "", "  ")
			//if err != nil {
			//	log.Infof("Error marshaling JSON: %v", err)
			//	return err
			//}
			//log.Infof("JSON EventUpdate Payload: %s", jsonData)

			resp, err := clients.Dhis2Client.PutResource(
				fmt.Sprintf("events/%s/%s", patientLog.EventID, v.DataElement), ep)
			if err != nil || !resp.IsSuccess() {
				log.Infof("Error sending result to DHIS2: %v: %v", err, string(resp.Body()))
				var data tracker.RootResponse
				err = json.Unmarshal(resp.Body(), &data)
				if err != nil {
					log.Infof("Error unmarshalling response: %v", err)
					continue
				}
				conflictMsg := tracker.ConflictsToError(data.Response.Conflicts).Error()
				if conflictMsg != "" {
					patientLog.ResultsUpdateErrors = conflictMsg
					patientLog.SetResultsUpdateErrors(conflictMsg)
				} else {
					patientLog.ResultsUpdateErrors = conflictMsg
					patientLog.SetResultsUpdateErrors("")
				}
				resultUpdatde = false
				continue
			}
			resultUpdatde = true
		}
		if resultUpdatde {
			patientLog.SetResultUpdated()
		}
		// Create Enrollment into Lab Program
		if diagnosed == "Yes" {
			// If no Enrollment exists in Lab Program
			if !patientLog.CheckLabProgramEnrollment() {
				enrollment := tracker.EnrollmentPayload{
					Program:               config.RTCGwConf.API.DHIS2LaboratoryProgram,
					Status:                "ACTIVE",
					OrgUnit:               result.FacilityID,
					EnrollmentDate:        resultsDate.Format("2006-01-02"),
					IncidentDate:          resultsDate.Format("2006-01-02"),
					TrackedEntityInstance: patientLog.TrackedEntity,
				}
				enrollmentID, err := enrollment.Create()
				if err != nil {
					log.Infof("Error creating Enrollment: %v", err.Error())
					return err
				}
				patientLog.SetLabEnrollment(enrollmentID)
			} else {
				// If Enrollment exists in Lab Program
				// Update event data values
				dataValues := []tracker.DataValue{
					{
						DataElement: config.RTCGwConf.API.DHIS2Mapping["data_elements"]["lab_results"],
						Value:       tbResult,
					},
					{
						DataElement: config.RTCGwConf.API.DHIS2Mapping["data_elements"]["lab_results_date"],
						Value:       resultsDate.Format("2006-01-02"),
					},
					{
						DataElement: config.RTCGwConf.API.DHIS2Mapping["data_elements"]["lab_diagnosis"],
						Value:       "true",
					},
					{
						DataElement: config.RTCGwConf.API.DHIS2Mapping["data_elements"]["lab_sample_referred_from_community"],
						Value:       "true",
					},
				}
				if patientLog.LabEventExists() {
					log.Infof("Lab Program Event exists for Patient: %v, TE: %v, in Enrollment: %v, Event: %v, Datavalues: %v",
						patientLog.ECHISID, patientLog.TrackedEntity, patientLog.LabEnrollment, patientLog.LabEvent, dataValues)
					for _, v := range dataValues {
						ep := tracker.EventUpdatePayload{
							Event:         patientLog.EventID,
							Program:       config.RTCGwConf.API.DHIS2LaboratoryProgram,
							OrgUnit:       result.FacilityID,
							Status:        "ACTIVE",
							ProgramStage:  config.RTCGwConf.API.DHIS2LaboratoryProgramStage,
							DataValues:    []tracker.DataValue{v},
							TrackedEntity: patientLog.TrackedEntity,
						}
						dataValuePutURL := fmt.Sprintf("events/%s/%s/", patientLog.LabEvent, v.DataElement)
						resp, err := clients.Dhis2Client.PutResource(dataValuePutURL, ep)
						if err != nil || !resp.IsSuccess() {
							log.Infof("Error sending result to DHIS2: %v: %v", err, string(resp.Body()))
							continue
						}
					}

				} else {
					log.Infof("Lab program Event missing for Patient: %v, TE: %v, in Enrollment: %v, Event: %v, Datavalues: %v",
						patientLog.ECHISID, patientLog.TrackedEntity, patientLog.LabEnrollment, patientLog.LabEvent, dataValues)
					// Create event and send results
					event := tracker.EventCreationPayload{
						Program:               config.RTCGwConf.API.DHIS2LaboratoryProgram,
						Status:                "ACTIVE",
						OrgUnit:               result.FacilityID,
						ProgramStage:          config.RTCGwConf.API.DHIS2LaboratoryProgramStage,
						Enrollment:            patientLog.LabEnrollment,
						EventDate:             resultsDate.Format("2006-01-02"),
						TrackedEntityInstance: patientLog.TrackedEntity,
						DataValues:            dataValues,
					}
					eventID, err := event.Create()
					if err != nil {
						log.Infof("Error creating Event: %v", err)
						return err
					}
					patientLog.SetLabEvent(eventID)
				}
			}

		}

	}

	log.Printf("Done sending result to DHIS2 for patient: %v", result.PatientID)
	return nil
}
