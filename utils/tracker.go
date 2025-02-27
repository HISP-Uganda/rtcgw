package utils

import (
	"fmt"
	"github.com/goccy/go-json"
	log "github.com/sirupsen/logrus"
	"reflect"
	"rtcgw/clients"
	"rtcgw/config"
	"strings"
	"time"
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

func GetFieldsByTag(s interface{}, targetUseAs string) map[string]string {
	val := reflect.ValueOf(s)
	typ := reflect.TypeOf(s)

	// Ensure it's a struct
	if val.Kind() != reflect.Struct {
		fmt.Println("Provided value is not a struct")
		return nil
	}

	result := make(map[string]string)

	// Iterate through struct fields
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		useAsTag := field.Tag.Get("use_as") // Extract use_as tag
		jsonTag := field.Tag.Get("json")    // Extract json tag

		// Ensure jsonTag exists; otherwise, use field name as fallback
		jsonTag = strings.Split(jsonTag, ",")[0]
		if jsonTag == "" {
			jsonTag = field.Name
		}

		// If the use_as tag matches the desired value, add to result
		if useAsTag == targetUseAs {
			result[jsonTag] = val.Field(i).String()
		}
	}

	return result
}

func GetCurrentDate() time.Time {
	return time.Now().UTC()
}
