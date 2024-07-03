package types

import (
	"time"
)

// datatype for both api and database
type Customer struct {
	//i think sql defaults to 32 bit int for keys, but sql returns 64
	Id          int64     `db:"id" json:"id"`
	NamePrefix  string    `db:"name_pref" json:"name_pref,omitempty"`
	NameFirst   string    `db:"name_first" json:"name_first,omitempty"`
	NameMiddle  string    `db:"name_middle" json:"name_middle,omitempty"`
	NameLast    string    `db:"name_last" json:"name_last,omitempty"`
	NameSuffix  string    `db:"name_suffix" json:"name_suffix,omitempty"`
	Email       string    `db:"email" json:"email,omitempty"`
	PhoneNumber string    `db:"phone_number" json:"phone_number,omitempty"`
	Deleted     bool      `db:"deleted" json:"-"` //hide in json
	CreatedBy   int64     `db:"created_by" json:"created_by,omitempty"`
	CreatedAt   time.Time `db:"created_at" json:"created_at,omitempty"`
	UpdatedBy   int64     `db:"updated_by" json:"updated_by,omitempty"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at,omitempty"`
}

// api response for customer or error
type CustResp struct {
	Customer *Customer `json:"customer,omitempty"`
	Error    string    `json:"error,omitempty"`
}

// api response for list customers
type CustList struct {
	Customers *[]Customer `json:"customers"`
	Next      int64       `json:"json "next_index"`
	Error     string      `json:"error,omitempty"`
}
