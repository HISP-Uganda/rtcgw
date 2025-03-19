package tracker

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
)

// RootResponse represents the top-level response structure.
type RootResponse struct {
	HttpStatus     string   `json:"httpStatus"`
	HttpStatusCode int      `json:"httpStatusCode"`
	Status         string   `json:"status"`
	Message        string   `json:"message"`
	Response       Response `json:"response"`
}

// Response represents the main response data.
type Response struct {
	ResponseType    string          `json:"responseType"`
	Status          string          `json:"status"`
	Imported        int             `json:"imported,omitempty"`
	Updated         int             `json:"updated,omitempty"`
	Deleted         int             `json:"deleted,omitempty"`
	Ignored         int             `json:"ignored,omitempty"`
	ImportCount     ImportCount     `json:"importCount,omitempty"`
	Conflicts       []any           `json:"conflicts,omitempty"`
	ImportOptions   ImportOptions   `json:"importOptions"`
	ImportSummaries []ImportSummary `json:"importSummaries"`
	Total           int             `json:"total,omitempty"`
}

// ImportOptions represents various import configuration options.
type ImportOptions struct {
	IdSchemes                   map[string]interface{} `json:"idSchemes"`
	DryRun                      bool                   `json:"dryRun"`
	Async                       bool                   `json:"async"`
	ImportStrategy              string                 `json:"importStrategy"`
	MergeMode                   string                 `json:"mergeMode"`
	ReportMode                  string                 `json:"reportMode"`
	SkipExistingCheck           bool                   `json:"skipExistingCheck"`
	Sharing                     bool                   `json:"sharing"`
	SkipNotifications           bool                   `json:"skipNotifications"`
	SkipAudit                   bool                   `json:"skipAudit"`
	DatasetAllowsPeriods        bool                   `json:"datasetAllowsPeriods"`
	StrictPeriods               bool                   `json:"strictPeriods"`
	StrictDataElements          bool                   `json:"strictDataElements"`
	StrictCategoryOptionCombos  bool                   `json:"strictCategoryOptionCombos"`
	StrictAttributeOptionCombos bool                   `json:"strictAttributeOptionCombos"`
	StrictOrganisationUnits     bool                   `json:"strictOrganisationUnits"`
	StrictDataSetApproval       bool                   `json:"strictDataSetApproval"`
	StrictDataSetLocking        bool                   `json:"strictDataSetLocking"`
	StrictDataSetInputPeriods   bool                   `json:"strictDataSetInputPeriods"`
	RequireCategoryOptionCombo  bool                   `json:"requireCategoryOptionCombo"`
	RequireAttributeOptionCombo bool                   `json:"requireAttributeOptionCombo"`
	SkipPatternValidation       bool                   `json:"skipPatternValidation"`
	IgnoreEmptyCollection       bool                   `json:"ignoreEmptyCollection"`
	Force                       bool                   `json:"force"`
	FirstRowIsHeader            bool                   `json:"firstRowIsHeader"`
	SkipLastUpdated             bool                   `json:"skipLastUpdated"`
	MergeDataValues             bool                   `json:"mergeDataValues"`
	SkipCache                   bool                   `json:"skipCache"`
}

// ImportSummary represents a summary of an import operation.
type ImportSummary struct {
	ResponseType    string        `json:"responseType"`
	Status          string        `json:"status"`
	ImportOptions   ImportOptions `json:"importOptions"`
	ImportCount     ImportCount   `json:"importCount"`
	Conflicts       []any         `json:"conflicts"`
	RejectedIndexes []any         `json:"rejectedIndexes"`
	Reference       string        `json:"reference"`
	Href            string        `json:"href,omitempty"`
	Enrollments     *Response     `json:"enrollments,omitempty"`
	Events          *Response     `json:"events,omitempty"`
}

// ImportCount represents the count of imported, updated, ignored, and deleted items.
type ImportCount struct {
	Imported int `json:"imported"`
	Updated  int `json:"updated"`
	Ignored  int `json:"ignored"`
	Deleted  int `json:"deleted"`
}

func ConflictsToError(conflicts []any) error {
	if len(conflicts) == 0 {
		return nil
	}

	var conflictMessages []string
	for _, conflict := range conflicts {
		conflictMessages = append(conflictMessages, fmt.Sprintf("%v", conflict))
	}

	// Combine all conflicts into a single error message.
	return fmt.Errorf("conflicts: %s", strings.Join(conflictMessages, "; "))
}

// GetTrackedEntityAndEventReferences retrieves the first event reference under enrollments → importSummaries → events.
func (r *RootResponse) GetTrackedEntityAndEventReferences() (string, string, bool, error) {
	trackedEntityInstance := ""
	// Traverse top-level import summaries
	for _, summary := range r.Response.ImportSummaries {
		if summary.Reference != "" {
			trackedEntityInstance = summary.Reference
		}
		if summary.Enrollments != nil {
			// Traverse enrollments importSummaries
			for _, enrollmentSummary := range summary.Enrollments.ImportSummaries {
				if enrollmentSummary.Events != nil {
					// Traverse event importSummaries
					for _, eventSummary := range enrollmentSummary.Events.ImportSummaries {
						if eventSummary.Reference != "" {
							// status is not ERROR
							if eventSummary.Status != "ERROR" {
								return trackedEntityInstance, eventSummary.Reference, true, nil
							}
							return trackedEntityInstance, eventSummary.Reference, true, ConflictsToError(eventSummary.Conflicts)
						}
					}
				}
			}
		}
	}
	return "", "", false, errors.New("not found")
}

func (r *RootResponse) GetEventIDReferenceAfterCreatingEvent() (string, error) {
	for _, summary := range r.Response.ImportSummaries {
		if summary.Status == "SUCCESS" {
			return summary.Reference, nil
		}
	}
	log.Printf("%v", r.Response.ImportSummaries)
	return "", errors.New("not found")
}

func (r *RootResponse) GetEnrolmentIDReferenceAfterCreatingEnrolment() (string, error) {
	for _, summary := range r.Response.ImportSummaries {
		if summary.Status == "SUCCESS" {
			return summary.Reference, nil
		} else {

		}
	}
	return "", errors.New("not found")
}
