package tracker

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
	Imported        int             `json:"imported"`
	Updated         int             `json:"updated"`
	Deleted         int             `json:"deleted"`
	Ignored         int             `json:"ignored"`
	ImportOptions   ImportOptions   `json:"importOptions"`
	ImportSummaries []ImportSummary `json:"importSummaries"`
	Total           int             `json:"total"`
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

// GetEventReference retrieves the first event reference under enrollments → importSummaries → events.
func (r *RootResponse) GetTrackedEntityAndEventReferences() (string, string, bool) {
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
							return trackedEntityInstance, eventSummary.Reference, true
						}
					}
				}
			}
		}
	}
	return "", "", false
}
