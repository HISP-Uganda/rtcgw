package tracker

import (
	"fmt"
	"github.com/goccy/go-json"
	log "github.com/sirupsen/logrus"
	"rtcgw/clients"
	"time"
)

// User represents a user in DHIS2.
type User struct {
	UID       string `json:"uid"`
	Username  string `json:"username"`
	FirstName string `json:"firstName,omitempty"`
	Surname   string `json:"surname,omitempty"`
}

// TrackedEntity represents a tracked entity instance.
type TrackedEntity struct {
	TrackedEntity      string      `json:"trackedEntity"`
	TrackedEntityType  string      `json:"trackedEntityType,omitempty"`
	CreatedAt          time.Time   `json:"createdAt,omitempty"`
	CreatedAtClient    time.Time   `json:"createdAtClient,omitempty"`
	UpdatedAt          time.Time   `json:"updatedAt,omitempty"`
	UpdatedAtClient    time.Time   `json:"updatedAtClient,omitempty"`
	OrgUnit            string      `json:"orgUnit,omitempty"`
	Inactive           bool        `json:"inactive,omitempty"`
	Deleted            bool        `json:"deleted,omitempty"`
	PotentialDuplicate bool        `json:"potentialDuplicate,omitempty"`
	Geometry           string      `json:"geometry,omitempty"`
	StoredBy           string      `json:"storedBy,omitempty"`
	CreatedBy          string      `json:"createdBy,omitempty"`
	UpdatedBy          string      `json:"updatedBy,omitempty"`
	Attributes         []Attribute `json:"attributes,omitempty"`
}

// Attribute represents a tracked entity attribute in DHIS2.
type Attribute struct {
	Attribute   string    `json:"attribute"`
	Code        string    `json:"code,omitempty"`
	DisplayName string    `json:"displayName,omitempty"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty"`
	StoredBy    string    `json:"storedBy,omitempty"`
	ValueType   string    `json:"valueType,omitempty"`
	Value       string    `json:"value"`
}

// Enrollment represents an enrollment of a tracked entity instance.
type Enrollment struct {
	Enrollment      string         `json:"enrollment"`
	Program         string         `json:"program"`
	TrackedEntity   string         `json:"trackedEntity"`
	Status          string         `json:"status"`
	OrgUnit         string         `json:"orgUnit"`
	CreatedAt       time.Time      `json:"createdAt,omitempty"`
	CreatedAtClient time.Time      `json:"createdAtClient,omitempty"`
	UpdatedAt       time.Time      `json:"updatedAt,omitempty"`
	UpdatedAtClient time.Time      `json:"updatedAtClient,omitempty"`
	EnrolledAt      time.Time      `json:"enrolledAt"`
	OccurredAt      time.Time      `json:"occurredAt,omitempty"`
	CompletedAt     time.Time      `json:"completedAt,omitempty"`
	CompletedBy     string         `json:"completedBy,omitempty"`
	FollowUp        bool           `json:"followUp,omitempty"`
	Deleted         bool           `json:"deleted,omitempty"`
	Geometry        string         `json:"geometry,omitempty"`
	StoredBy        string         `json:"storedBy,omitempty"`
	CreatedBy       string         `json:"createdBy,omitempty"`
	UpdatedBy       string         `json:"updatedBy,omitempty"`
	Attributes      []Attribute    `json:"attributes,omitempty"`
	Events          []Event        `json:"events,omitempty"`
	Relationships   []Relationship `json:"relationships,omitempty"`
	Notes           []Note         `json:"notes,omitempty"`
}

// Event represents an event in DHIS2.
type Event struct {
	Event                    string         `json:"event"`
	ProgramStage             string         `json:"programStage"`
	Enrollment               string         `json:"enrollment"`
	Program                  string         `json:"program"`
	TrackedEntity            string         `json:"trackedEntityInstance,omitempty"`
	Status                   string         `json:"status"`
	EnrollmentStatus         string         `json:"enrollmentStatus,omitempty"`
	OrgUnit                  string         `json:"orgUnit"`
	CreatedAt                time.Time      `json:"created"`
	CreatedAtClient          time.Time      `json:"createdAtClient,omitempty"`
	UpdatedAt                time.Time      `json:"lastUpdated"`
	UpdatedAtClient          time.Time      `json:"updatedAtClient,omitempty"`
	ScheduledAt              time.Time      `json:"dueDate,omitempty"`
	OccurredAt               time.Time      `json:"eventDate"`
	CompletedAt              time.Time      `json:"completedDate,omitempty"`
	CompletedBy              string         `json:"completedBy,omitempty"`
	FollowUp                 bool           `json:"followUp,omitempty"`
	Deleted                  bool           `json:"deleted,omitempty"`
	Geometry                 string         `json:"geometry,omitempty"`
	StoredBy                 string         `json:"storedBy,omitempty"`
	CreatedBy                string         `json:"createdBy,omitempty"`
	UpdatedBy                string         `json:"updatedBy,omitempty"`
	AttributeOptionCombo     string         `json:"attributeOptionCombo,omitempty"`
	AttributeCategoryOptions string         `json:"attributeCategoryOptions,omitempty"`
	AssignedUser             string         `json:"assignedUser,omitempty"`
	DataValues               []DataValue    `json:"dataValues,omitempty"`
	Relationships            []Relationship `json:"relationships,omitempty"`
	Notes                    []Note         `json:"notes,omitempty"`
}

// DataValue represents a data value in an event or entity in DHIS2.
type DataValue struct {
	DataElement       string     `json:"dataElement"`
	Value             string     `json:"value"`
	ProvidedElsewhere bool       `json:"providedElsewhere,omitempty"`
	CreatedAt         *time.Time `json:"createdAt,omitempty"`
	UpdatedAt         *time.Time `json:"updatedAt,omitempty"`
	StoredBy          string     `json:"storedBy,omitempty"`
	CreatedBy         string     `json:"createdBy,omitempty"`
	UpdatedBy         string     `json:"updatedBy,omitempty"`
}

// Relationship represents a relationship between entities in DHIS2.
type Relationship struct {
	ID               string           `json:"id"`
	RelationshipType string           `json:"relationshipType"`
	RelationshipName string           `json:"relationshipName,omitempty"`
	CreatedAt        time.Time        `json:"createdAt"`
	UpdatedAt        time.Time        `json:"updatedAt"`
	CreatedAtClient  time.Time        `json:"createdAtClient,omitempty"`
	Bidirectional    bool             `json:"bidirectional,omitempty"`
	From             RelationshipItem `json:"from"`
	To               RelationshipItem `json:"to"`
}

// RelationshipItem represents a 'from' or 'to' item in a relationship.
type RelationshipItem struct {
	TrackedEntity string `json:"trackedEntity,omitempty"`
	Enrollment    string `json:"enrollment,omitempty"`
	Event         string `json:"event,omitempty"`
}

// Note represents a note attached to an enrollment, event, or other entity in DHIS2.
type Note struct {
	ID        string    `json:"id"`
	Value     string    `json:"value"`
	StoredAt  time.Time `json:"storedAt"`
	StoredBy  string    `json:"storedBy,omitempty"`
	CreatedBy string    `json:"createdBy,omitempty"`
}

// FlatPayload represents the top-level structure of a flat payload for DHIS2 tracker data import.
type FlatPayload struct {
	TrackedEntities []TrackedEntity `json:"trackedEntities,omitempty"`
	Enrollments     []Enrollment    `json:"enrollments,omitempty"`
	Events          []Event         `json:"events,omitempty"`
	Relationships   []Relationship  `json:"relationships,omitempty"`
}

type EventUpdatePayload struct {
	Event string `json:"event"`
	//EventDate     string      `json:"eventDate"`
	Program       string      `json:"program"`
	OrgUnit       string      `json:"orgUnit"`
	ProgramStage  string      `json:"programState"`
	Status        string      `json:"status,omitempty"`
	TrackedEntity string      `json:"trackedEntityInstance,omitempty"`
	DataValues    []DataValue `json:"dataValues"`
}

type TrackedEntityUpdatePayload struct {
	TrackedEntityInstance string            `json:"trackedEntityInstance"`
	OrgUnit               string            `json:"orgUnit"`
	TrackedEntityType     string            `json:"trackedEntityType"`
	Attributes            []NestedAttribute `json:"attributes"`
	Relationships         []Relationship    `json:"relationships,omitempty"`
}

type EnrollmentPayload struct {
	Program               string `json:"program"`
	Status                string `json:"status"`
	OrgUnit               string `json:"orgUnit"`
	EnrollmentDate        string `json:"enrollmentDate"`
	IncidentDate          string `json:"incidentDate"`
	TrackedEntityInstance string `json:"trackedEntityInstance"`
}

type EventCreationPayload struct {
	TrackedEntityInstance string      `json:"trackedEntityInstance"`
	Program               string      `json:"program"`
	ProgramStage          string      `json:"programStage"`
	Enrollment            string      `json:"enrollment"`
	OrgUnit               string      `json:"orgUnit"`
	Status                string      `json:"status"`
	EventDate             string      `json:"eventDate"`
	DataValues            []DataValue `json:"dataValues"`
}

func (e *EnrollmentPayload) Create() (string, error) {
	resp, err := clients.Dhis2Client.PostResource("enrollments", nil, e)
	if err != nil {
		return "", err
	}
	if resp.IsSuccess() {
		var response RootResponse
		err := json.Unmarshal(resp.Body(), &response)
		if err != nil {
			log.Infof("Error unmarshalling enrollment response: %v", string(resp.Body()))
			return "", err
		}
		enID, err := response.GetEnrolmentIDReferenceAfterCreatingEnrolment()
		if err != nil {
			log.Infof("Error getting EnrolmentID %v", string(resp.Body()))
			return "", err
		}
		return enID, nil
	} else {
		return "", fmt.Errorf("failed to create enrollment: %v", string(resp.Body()))
	}
	return "", fmt.Errorf("failed to create enrollment")
}

func (e *EventCreationPayload) Create() (string, error) {
	eventData := map[string][]EventCreationPayload{
		"events": {*e},
	}
	resp, err := clients.Dhis2Client.PostResource("events", nil, eventData)
	if err != nil {
		return "", err
	}
	if resp.IsSuccess() {
		// Get Event ID
		var response RootResponse
		err := json.Unmarshal(resp.Body(), &response)
		if err != nil {
			return "", err
		}
		evID, err := response.GetEventIDReferenceAfterCreatingEvent()
		if err != nil {
			return "", err
		}
		return evID, nil
	} else {
		log.Debugf("failed to create event: %v", string(resp.Body()))
		// return "", fmt.Errorf("failed to create event: %v", string(resp.Body()))
	}
	return "", fmt.Errorf("failed to create event")
}
