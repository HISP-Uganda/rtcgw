package utils

import (
	"fmt"
	"github.com/goccy/go-json"
	log "github.com/sirupsen/logrus"
	"rtcgw/clients"
	"rtcgw/config"
)

// SearchTE searches for the existence of a TrackedEntity in DHIS2 that matches a given tracked entity attribute value, orgUnit and Program
// also uses our client *
func SearchTE(client *clients.Client, echisID, orgUnit, program string) bool {
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
		return false
	}

	var responseMap interface{}

	err = json.Unmarshal(resp.Body(), &responseMap)
	if err != nil {
		log.Info("Error unmarshalling response!!!")
		fmt.Printf("Error unmarshalling response: %v\n", err)
		return false
	}
	ret, err := PrintResponse(responseMap, true)
	if err != nil {
		log.Info("Error printing response!!!")
		fmt.Printf("Error printing response: %v\n", err)
		return false
	} else {
		log.Infof("Response: %s\n", ret)
	}
	return false
}

func PrintResponse(responseMap any, pretty bool) (string, error) {
	if pretty {
		prettyJSON, err := json.MarshalIndent(responseMap, "", "  ")
		if err != nil {
			return "", err
		}
		return string(prettyJSON), nil
	} else {
		retJson, err := json.Marshal(responseMap)
		if err != nil {
			return "", err
		}
		return string(retJson), nil
	}
}
