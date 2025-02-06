package tracker

import (
	"time"
)

// NestedPayload represents the top-level structure of a nested payload for DHIS2 tracker data import.
type NestedPayload struct {
	TrackedEntities []NestedTrackedEntity `json:"trackedEntities,omitempty"`
}

// NestedTrackedEntity represents a tracked entity with nested enrollments.
type NestedTrackedEntity struct {
	Enrollments       []NestedEnrollment `json:"enrollments,omitempty"`
	OrgUnit           string             `json:"orgUnit"`
	TrackedEntityType string             `json:"trackedEntityType"`
}

// NestedEnrollment represents an enrollment with nested attributes and events.
type NestedEnrollment struct {
	Attributes        []NestedAttribute `json:"attributes,omitempty"`
	EnrolledAt        time.Time         `json:"enrolledAt"`
	Events            []NestedEvent     `json:"events,omitempty"`
	OccurredAt        time.Time         `json:"occurredAt"`
	OrgUnit           string            `json:"orgUnit"`
	Program           string            `json:"program"`
	Status            string            `json:"status"`
	TrackedEntityType string            `json:"trackedEntityType"`
}

// NestedAttribute represents an attribute within an enrollment.
type NestedAttribute struct {
	Attribute   string `json:"attribute"`
	DisplayName string `json:"displayName,omitempty"`
	Value       string `json:"value"`
}

// NestedEvent represents an event with nested data values and notes.
type NestedEvent struct {
	AttributeCategoryOptions string      `json:"attributeCategoryOptions,omitempty"`
	AttributeOptionCombo     string      `json:"attributeOptionCombo,omitempty"`
	DataValues               []DataValue `json:"dataValues,omitempty"`
	EnrollmentStatus         string      `json:"enrollmentStatus,omitempty"`
	Notes                    []Note      `json:"notes,omitempty"`
	OccurredAt               time.Time   `json:"occurredAt"`
	OrgUnit                  string      `json:"orgUnit"`
	Program                  string      `json:"program"`
	ProgramStage             string      `json:"programStage"`
	ScheduledAt              time.Time   `json:"scheduledAt,omitempty"`
	Status                   string      `json:"status"`
}
