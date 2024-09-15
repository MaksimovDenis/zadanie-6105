// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package queries

import (
	"database/sql/driver"
	"fmt"

	"github.com/kak-tus/nan"
)

type OrganizationType string

const (
	OrganizationTypeIE  OrganizationType = "IE"
	OrganizationTypeLLC OrganizationType = "LLC"
	OrganizationTypeJSC OrganizationType = "JSC"
)

func (e *OrganizationType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = OrganizationType(s)
	case string:
		*e = OrganizationType(s)
	default:
		return fmt.Errorf("unsupported scan type for OrganizationType: %T", src)
	}
	return nil
}

type NullOrganizationType struct {
	OrganizationType OrganizationType
	Valid            bool // Valid is true if OrganizationType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullOrganizationType) Scan(value interface{}) error {
	if value == nil {
		ns.OrganizationType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.OrganizationType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullOrganizationType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.OrganizationType), nil
}

type StatusType string

const (
	StatusTypeCreated   StatusType = "Created"
	StatusTypePublished StatusType = "Published"
	StatusTypeCanceled  StatusType = "Canceled"
	StatusTypeApproved  StatusType = "Approved"
	StatusTypeRejected  StatusType = "Rejected"
	StatusTypeClosed    StatusType = "Closed"
)

func (e *StatusType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = StatusType(s)
	case string:
		*e = StatusType(s)
	default:
		return fmt.Errorf("unsupported scan type for StatusType: %T", src)
	}
	return nil
}

type NullStatusType struct {
	StatusType StatusType
	Valid      bool // Valid is true if StatusType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullStatusType) Scan(value interface{}) error {
	if value == nil {
		ns.StatusType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.StatusType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullStatusType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.StatusType), nil
}

type Bid struct {
	ID             int32
	TenderID       nan.NullInt32
	OrganizationID nan.NullInt32
	UserID         nan.NullInt32
	Name           string
	Description    nan.NullString
	Status         NullStatusType
	AuthorType     nan.NullString
	Version        nan.NullInt32
	CreatedAt      nan.NullTime
	UpdatedAt      nan.NullTime
}

type BidFeedback struct {
	ID        int32
	BidID     nan.NullInt32
	UserID    nan.NullInt32
	Feedback  string
	CreatedAt nan.NullTime
}

type BidsDecision struct {
	ID           int32
	BidID        nan.NullInt32
	DecisionType nan.NullString
	UserID       nan.NullInt32
	CreatedAt    nan.NullTime
}

type BidsHistory struct {
	ID             int32
	TenderID       nan.NullInt32
	BidID          nan.NullInt32
	OrganizationID nan.NullInt32
	UserID         nan.NullInt32
	Name           string
	Description    nan.NullString
	Status         NullStatusType
	AuthorType     nan.NullString
	Version        nan.NullInt32
	CreatedAt      nan.NullTime
	UpdatedAt      nan.NullTime
}

type Employee struct {
	ID        int32
	Username  string
	FirstName nan.NullString
	LastName  nan.NullString
	CreatedAt nan.NullTime
	UpdatedAt nan.NullTime
}

type Organization struct {
	ID          int32
	Name        string
	Description nan.NullString
	Type        NullOrganizationType
	CreatedAt   nan.NullTime
	UpdatedAt   nan.NullTime
}

type OrganizationResponsible struct {
	ID             int32
	OrganizationID nan.NullInt32
	UserID         nan.NullInt32
}

type Tender struct {
	ID             int32
	OrganizationID nan.NullInt32
	CreatedBy      nan.NullInt32
	Name           string
	Description    nan.NullString
	Status         NullStatusType
	ServiceType    nan.NullString
	Version        nan.NullInt32
	CreatedAt      nan.NullTime
	UpdatedAt      nan.NullTime
}

type TenderHistory struct {
	ID          int32
	TenderID    nan.NullInt32
	Name        nan.NullString
	Description nan.NullString
	ServiceType nan.NullString
	Status      NullStatusType
	Version     nan.NullInt32
	CreatedAt   nan.NullTime
	UpdatedAt   nan.NullTime
}
