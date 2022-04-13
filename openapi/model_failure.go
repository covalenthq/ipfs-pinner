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

// Failure Response for a failed request
type Failure struct {
	Error FailureError `json:"error"`
}

// NewFailure instantiates a new Failure object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewFailure(error_ FailureError) *Failure {
	this := Failure{}
	this.Error = error_
	return &this
}

// NewFailureWithDefaults instantiates a new Failure object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewFailureWithDefaults() *Failure {
	this := Failure{}
	return &this
}

// GetError returns the Error field value
func (o *Failure) GetError() FailureError {
	if o == nil {
		var ret FailureError
		return ret
	}

	return o.Error
}

// GetErrorOk returns a tuple with the Error field value
// and a boolean to check if the value has been set.
func (o *Failure) GetErrorOk() (*FailureError, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Error, true
}

// SetError sets field value
func (o *Failure) SetError(v FailureError) {
	o.Error = v
}

func (o Failure) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["error"] = o.Error
	}
	return json.Marshal(toSerialize)
}

type NullableFailure struct {
	value *Failure
	isSet bool
}

func (v NullableFailure) Get() *Failure {
	return v.value
}

func (v *NullableFailure) Set(val *Failure) {
	v.value = val
	v.isSet = true
}

func (v NullableFailure) IsSet() bool {
	return v.isSet
}

func (v *NullableFailure) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableFailure(val *Failure) *NullableFailure {
	return &NullableFailure{value: val, isSet: true}
}

func (v NullableFailure) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableFailure) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
