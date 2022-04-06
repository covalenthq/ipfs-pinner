/*
IPFS Pinning Service API

some notes

API version: 1.0.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"encoding/json"
)

// FailureError struct for FailureError
type FailureError struct {
	// Mandatory string identifying the type of error
	Reason string `json:"reason"`
	// Optional, longer description of the error; may include UUID of transaction for support, links to documentation etc
	Details *string `json:"details,omitempty"`
}

// NewFailureError instantiates a new FailureError object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewFailureError(reason string) *FailureError {
	this := FailureError{}
	this.Reason = reason
	return &this
}

// NewFailureErrorWithDefaults instantiates a new FailureError object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewFailureErrorWithDefaults() *FailureError {
	this := FailureError{}
	return &this
}

// GetReason returns the Reason field value
func (o *FailureError) GetReason() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Reason
}

// GetReasonOk returns a tuple with the Reason field value
// and a boolean to check if the value has been set.
func (o *FailureError) GetReasonOk() (*string, bool) {
	if o == nil  {
		return nil, false
	}
	return &o.Reason, true
}

// SetReason sets field value
func (o *FailureError) SetReason(v string) {
	o.Reason = v
}

// GetDetails returns the Details field value if set, zero value otherwise.
func (o *FailureError) GetDetails() string {
	if o == nil || o.Details == nil {
		var ret string
		return ret
	}
	return *o.Details
}

// GetDetailsOk returns a tuple with the Details field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *FailureError) GetDetailsOk() (*string, bool) {
	if o == nil || o.Details == nil {
		return nil, false
	}
	return o.Details, true
}

// HasDetails returns a boolean if a field has been set.
func (o *FailureError) HasDetails() bool {
	if o != nil && o.Details != nil {
		return true
	}

	return false
}

// SetDetails gets a reference to the given string and assigns it to the Details field.
func (o *FailureError) SetDetails(v string) {
	o.Details = &v
}

func (o FailureError) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["reason"] = o.Reason
	}
	if o.Details != nil {
		toSerialize["details"] = o.Details
	}
	return json.Marshal(toSerialize)
}

type NullableFailureError struct {
	value *FailureError
	isSet bool
}

func (v NullableFailureError) Get() *FailureError {
	return v.value
}

func (v *NullableFailureError) Set(val *FailureError) {
	v.value = val
	v.isSet = true
}

func (v NullableFailureError) IsSet() bool {
	return v.isSet
}

func (v *NullableFailureError) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableFailureError(val *FailureError) *NullableFailureError {
	return &NullableFailureError{value: val, isSet: true}
}

func (v NullableFailureError) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableFailureError) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

