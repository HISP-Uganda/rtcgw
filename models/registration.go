package models

import (
	"database/sql"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	log "github.com/sirupsen/logrus"
	"rtcgw/clients"
	"rtcgw/config"
	"rtcgw/models/tracker"
	"rtcgw/utils"
	"time"
)

type ECHISRequest struct {
	ECHISID             string `use_as:"attr" json:"echis_patient_id" binding:"required"`
	NIN                 string `use_as:"attr" json:"national_identification_number" binding:"omitempty,ugandaNIN"`
	Name                string `use_as:"attr" json:"patient_name" binding:"required"`
	Sex                 string `use_as:"attr" json:"patient_gender" binding:"omitempty,maleFemale"`
	FacilityID          string `use_as:"" json:"facility_id"`
	FacilityDHIS2ID     string `use_as:"" json:"facility_dhis2_id" binding:"required,dhis2UID"`
	PatientPhone        string `use_as:"attr" json:"patient_phone"`
	PatientCategory     string `use_as:"attr" json:"patient_category"`
	PatientAgeInYears   string `use_as:"attr" json:"patient_age_in_years"`
	PatientAgeInMonths  string `use_as:"attr" json:"patient_age_in_months,omitempty"`
	PatientAgeInDays    string `use_as:"attr" json:"patient_age_in_days,omitempty"`
	ClientCategory      string `use_as:"attr" json:"client_category,omitempty"`
	Cough               string `use_as:"de" json:"cough,omitempty" binding:"omitempty,yesNo"`
	Fever               string `use_as:"de" json:"fever,omitempty" binding:"omitempty,yesNo"`
	WeightLoss          string `use_as:"de" json:"weight_loss,omitempty" binding:"omitempty,yesNo"`
	ExcessiveNightSweat string `use_as:"de" json:"excessive_night_sweat,omitempty" binding:"omitempty,yesNo"`
	IsOnTBTreatment     string `use_as:"de" json:"is_on_tb_treatment,omitempty" binding:"omitempty,yesNo"`
	PoorWeightGain      string `use_as:"de" json:"poor_weight_gain,omitempty"`
}

// FormatValidationError translates validation errors
func FormatValidationError(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			// Customize messages for required fields
			switch e.Field() {
			case "echis_parent_id":
				errors["echis_parent_id"] = "echis_patient_id is required and cannot be empty."
			case "facility_dhis2_id":
				errors["facility_dhis2_id"] = "facility_dhis2_id must be a valid DHIS2 UID."
			case "patient_name":
				errors[e.Field()] = "patient_name is required and must be provided."
			case "national_identification_number":
				errors["national_identification_number"] = "invalid national_identification_number provided."
			case "patient_gender":
				errors["patient_gender"] = "patient_gender should be Male or Female"
			case "fever", "cough", "weight_loss", "excessive_night_sweat":
				errors[e.Field()] = fmt.Sprintf("Should be exactly 'Yes' or 'No'. Notice the case")
			default:

				errors[e.Field()] = fmt.Sprintf("Validation failed on '%s' condition", e.Tag())
			}
		}
	}

	return errors
}

func (r ECHISRequest) SaveClient(client *clients.Client) {
	attr := utils.GetFieldsByTag(r, "attr")
	des := utils.GetFieldsByTag(r, "de")

	attributesConf, exists := config.RTCGwConf.API.DHIS2Mapping["attributes"]
	if !exists {
		log.Infof("DHIS2Mapping not found for attributes in config")
		return
	}
	var attributes []tracker.NestedAttribute
	for k, v := range attributesConf {
		if v == "" {
			continue
		}
		val, e := attr[k]
		if e {
			attributes = append(attributes, tracker.NestedAttribute{
				Attribute: v,
				Value:     val,
			})

		}
	}
	log.Infof("attributes: %v", attributes)

	dataElementsConf, exists := config.RTCGwConf.API.DHIS2Mapping["data_elements"]
	var dataValues []tracker.DataValue
	if !exists {
		log.Infof("DHIS2Mapping not found for data_elements in config")
		return
	}
	for k, v := range dataElementsConf {
		if v == "" {
			continue
		}
		val, e := des[k]
		if e {
			dataValues = append(dataValues, tracker.DataValue{
				DataElement: v,
				Value:       val,
			})
		}
	}
	log.Infof("dataElements: %v The des: %v", dataValues, des)
	events := []tracker.NestedEvent{{
		DataValues:   dataValues,
		OrgUnit:      r.FacilityDHIS2ID,
		Program:      config.RTCGwConf.API.DHIS2TrackerProgram,
		ProgramStage: config.RTCGwConf.API.DHIS2TrackerProgramStage,
		OccurredAt:   utils.GetCurrentDate(),
		Status:       "ACTIVE",
	}}
	enrollment := tracker.NestedEnrollment{
		EnrolledAt: utils.GetCurrentDate(),
		Status:     "ACTIVE",
		Events:     events,
		OccurredAt: utils.GetCurrentDate(),
		OrgUnit:    r.FacilityDHIS2ID,
		Program:    config.RTCGwConf.API.DHIS2TrackerProgram,
		// TrackedEntityType: config.RTCGwConf.API.DHIS2TrackedEntityType,
	}

	nestedPayload := tracker.NestedPayload{
		TrackedEntities: []tracker.NestedTrackedEntity{{
			Attributes:        attributes,
			Enrollments:       []tracker.NestedEnrollment{enrollment},
			OrgUnit:           r.FacilityDHIS2ID,
			TrackedEntityType: config.RTCGwConf.API.DHIS2TrackedEntityType,
		}},
	}
	// turn nestedPayload to json and print it
	jsonData, err := json.MarshalIndent(nestedPayload, "", "  ")
	if err != nil {
		log.Infof("Error marshaling JSON: %v", err)
		return
	}
	log.Infof("JSON NestedPayload: %s", jsonData)
	resp, err := client.PostResource("trackedEntityInstances", nil, nestedPayload)
	if err != nil {
		log.Infof("Error sending request: %v", err)
		return
	}
	// if resp status code is 200
	if resp.IsSuccess() {
		log.Infof("Patient saved successfully in DHIS2")
		log.Infof("Response: %s", resp.Body())
		var data tracker.RootResponse
		err = json.Unmarshal(resp.Body(), &data)
		if err != nil {
			log.Infof("Error unmarshalling response: %v", err)
			return
		}
		trackedEntity, eventID, found, conflicts := data.GetTrackedEntityAndEventReferences()
		if found {
			currentTime := sql.NullTime{Time: time.Now(), Valid: true}
			conflictMsg := ""
			if conflicts != nil {
				conflictMsg = conflicts.Error()
			}
			synclog := SyncLog{
				ECHISID:                   r.ECHISID,
				EventID:                   eventID,
				EventDate:                 currentTime,
				TrackedEntity:             trackedEntity,
				ECHISClientCreationErrors: conflictMsg,
				OrgUnit:                   r.FacilityDHIS2ID,
			}
			synclog.Save()
			if conflictMsg != "" {
				log.Infof("Conflicts during DHIS2 sync: %v", conflicts.Error())
				synclog.SetECHISClientCreationErrors(conflicts.Error())
				return
			} else {
				synclog.SetECHISClientCreationErrors("")
			}
			log.Infof("Event ID: %s", eventID)
		} else {
			log.Infof("Error retrieving Event Reference from response: %v", data.Response.ImportSummaries)
			return
		}
		// print indented data
		jsonData, err = json.MarshalIndent(data, "", "  ")
		if err != nil {
			log.Infof("Error marshaling JSON: %v", err)
			return
		}
		log.Infof("JSON Response: %s", jsonData)

		return
	} else {
		log.Infof("Error saving patient in DHIS2: %s", resp.Body())
		return
	}

}

func (r ECHISRequest) UpdateClient(client *clients.Client, syncLog *SyncLog) {
	attr := utils.GetFieldsByTag(r, "attr")
	des := utils.GetFieldsByTag(r, "de")

	attributesConf, exists := config.RTCGwConf.API.DHIS2Mapping["attributes"]
	if !exists {
		log.Infof("DHIS2Mapping not found for attributes in config")
		return
	}
	var attributes []tracker.NestedAttribute
	for k, v := range attributesConf {
		if v == "" {
			continue
		}
		val, e := attr[k]
		if e {
			attributes = append(attributes, tracker.NestedAttribute{
				Attribute: v,
				Value:     val,
			})

		}
	}
	// log.Infof("attributes: %v: >> %v", attributes, attr)
	teUpdatePayload := tracker.TrackedEntityUpdatePayload{
		TrackedEntityInstance: syncLog.TrackedEntity,
		TrackedEntityType:     config.RTCGwConf.API.DHIS2TrackedEntityType,
		Attributes:            attributes,
		OrgUnit:               r.FacilityDHIS2ID,
	}
	putURL := fmt.Sprintf("trackedEntityInstances/%s?program=%s", syncLog.TrackedEntity, config.RTCGwConf.API.DHIS2TrackerProgram)
	resp, err := client.PutResource(putURL, teUpdatePayload)
	if err != nil || !resp.IsSuccess() {

		log.Infof("Error updating trackedEntity attributes in DHIS2: %v: %v", err, string(resp.Body()))
	}

	dataElementsConf, exists := config.RTCGwConf.API.DHIS2Mapping["data_elements"]
	var dataValues []tracker.DataValue
	if !exists {
		log.Infof("DHIS2Mapping not found for data_elements in config")
		return
	}
	for k, v := range dataElementsConf {
		if v == "" {
			continue
		}
		val, e := des[k]
		if e {
			dataValues = append(dataValues, tracker.DataValue{
				DataElement: v,
				Value:       val,
			})
		}
	}
	log.Infof("dataElements: %v The des: %v", dataValues, des)
	for _, v := range dataValues {
		ep := tracker.EventUpdatePayload{
			Event:         syncLog.EventID,
			Program:       config.RTCGwConf.API.DHIS2TrackerProgram,
			OrgUnit:       r.FacilityDHIS2ID,
			Status:        "ACTIVE",
			ProgramStage:  config.RTCGwConf.API.DHIS2TrackerProgramStage,
			DataValues:    []tracker.DataValue{v},
			TrackedEntity: syncLog.TrackedEntity,
		}
		dataValuePutURL := fmt.Sprintf("events/%s/%s/", syncLog.EventID, v.DataElement)
		resp, err := client.PutResource(dataValuePutURL, ep)
		if err != nil || !resp.IsSuccess() {
			log.Infof("Error sending result to DHIS2: %v: %v", err, string(resp.Body()))
			continue
		}
	}

}
