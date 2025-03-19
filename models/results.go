package models

import (
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/goccy/go-json"
	log "github.com/sirupsen/logrus"
	"rtcgw/clients"
	"rtcgw/config"
	"rtcgw/models/tracker"
)

type LabXpertResult struct {
	PatientID  string `json:"patient_id" binding:"required"`
	Lab        string `json:"lab,omitempty"`
	MTB        string `json:"mtb"`
	RR         string `json:"rr"`
	ResultDate string `json:"result_date"`
	FacilityID string `json:"facility_dhis2_id"`
}

func SearchTE(client *clients.Client, echisID, orgUnit, program string) (bool, []tracker.TrackedEntity) {
	params := make(map[string]string)
	// params["trackedEntity"] = echisID
	params["orgUnit"] = orgUnit
	params["program"] = program
	params["ouMode"] = "SELECTED"
	params["orgUnitMode"] = "SELECTED"
	params["filter"] = fmt.Sprintf("%s:EQ:%s", config.RTCGwConf.API.DHIS2SearchAttribute, echisID)
	//params["query"] = fmt.Sprintf("%s", echisID)
	resp, err := client.GetResource("/tracker/trackedEntities", params)
	if err != nil {
		log.Info("Error calling resource!!!")
		fmt.Printf("Error when calling GetResource: %v\n", err)
		return false, nil
	}

	v, _, _, err := jsonparser.Get(resp.Body(), "instances")
	if err != nil {
		log.Infof("Error getting instances: %v", err)
		return false, nil
	}
	var instances []tracker.TrackedEntity
	err = json.Unmarshal(v, &instances)
	if err != nil {
		log.Infof("Error unmarshalling instances: %v", err)
		return false, nil
	}
	return true, instances
}

// CheckDhis2Presence returns true if a TE is present for the given results
func (r *LabXpertResult) CheckDhis2Presence(c *clients.Client) bool {
	log.Info("Checking TE in DHIS2")
	//exists := utils.SearchTE(c,
	//	r.PatientID, r.FacilityID, config.RTCGwConf.API.DHIS2TrackerProgram)
	exists, _ := SearchTE(c,
		r.PatientID, r.FacilityID, config.RTCGwConf.API.DHIS2TrackerProgram)
	if !exists {
		log.Infof("TE not found in DHIS2 for patient: %s, facility: %s, program: %s",
			r.PatientID, r.FacilityID, config.RTCGwConf.API.DHIS2TrackerProgram)
	}

	return exists // Placeholder return value

}

func (r *LabXpertResult) InDhis2() (*SyncLog, bool) {
	syncLog, err := GetSyncLogByECHISID(r.PatientID)
	if err != nil {
		log.Infof("Error getting sync log for patient: %s: Error: %v", r.PatientID, err.Error())
		return nil, false
	}
	return syncLog, true
}

func (r *LabXpertResult) SaveResults(c *clients.Client) {
	//exists, instances := SearchTE(c,
	//	r.PatientID, r.FacilityID, config.RTCGwConf.API.DHIS2TrackerProgram)
	//if !exists {
	//	log.Infof("TE not found in DHIS2 for patient: %s, facility: %s, program: %s",
	//		r.PatientID, r.FacilityID, config.RTCGwConf.API.DHIS2TrackerProgram)
	//	return
	//}
	//if len(instances) > 0 {
	//	log.Infof("TE found in DHIS2 for patient: %s, facility: %s, program: %s",
	//		r.PatientID, r.FacilityID, config.RTCGwConf.API.DHIS2TrackerProgram)
	//	// Update tracked entity
	//	// TODO: Implement update tracked entity function
	//
	//}
	dhis2Log, present := r.InDhis2()
	if !present {
		// Alert that no match found in DHIS2 and return
		log.Infof("No match found in DHIS2 for eCHISID %s", r.PatientID)
		return
	}
	// create event update payload
	log.Infof("Updating event: %s", dhis2Log.EventID)

}

// GetResult returns the result as it is in DHIS2 and a Yes/No for diagnosis status
func (r *LabXpertResult) GetResult() (string, string) {
	switch r.MTB {
	case "DETECTED VERY LOW", "DETECTED LOW", "DETECTED MEDIUM", "DETECTED HIGH":
		switch r.RR {
		case "DETECTED":
			return "MTB detected, rifampicin resistance detected", "Yes"
		case "INDETERMINATE":
			return "MTB detected, rifampicin resistance indeterminate", "Yes"
		case "NOT DETECTED":
			return "MTB detected, rifampicin resistance not detected", "Yes"
		case "Invalid":
			return "Invalid", "No"
		case "Error":
			return "Error", "No"
		}
	case "ERROR":
		return "Error", "No"
	case "INVALID":
		return "Invalid", "No"
	case "NO RESULT":
		return "No result", "No"
	case "NOT DETECTED":
		return "MTB not detected", "No"
	default:
		return "No Result", "No"
	}
	return "No result", "No"
}
