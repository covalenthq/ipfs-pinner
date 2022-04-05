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

// PinataMetadata struct for PinataMetadata
type PinataMetadata struct {
	// null
	Name *string `json:"name,omitempty"`
	Keyvalues *map[string]string `json:"keyvalues,omitempty"`
}

// NewPinataMetadata instantiates a new PinataMetadata object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPinataMetadata() *PinataMetadata {
	this := PinataMetadata{}
	return &this
}

// NewPinataMetadataWithDefaults instantiates a new PinataMetadata object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPinataMetadataWithDefaults() *PinataMetadata {
	this := PinataMetadata{}
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *PinataMetadata) GetName() string {
	if o == nil || o.Name == nil {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PinataMetadata) GetNameOk() (*string, bool) {
	if o == nil || o.Name == nil {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *PinataMetadata) HasName() bool {
	if o != nil && o.Name != nil {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *PinataMetadata) SetName(v string) {
	o.Name = &v
}

// GetKeyvalues returns the Keyvalues field value if set, zero value otherwise.
func (o *PinataMetadata) GetKeyvalues() map[string]string {
	if o == nil || o.Keyvalues == nil {
		var ret map[string]string
		return ret
	}
	return *o.Keyvalues
}

// GetKeyvaluesOk returns a tuple with the Keyvalues field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PinataMetadata) GetKeyvaluesOk() (*map[string]string, bool) {
	if o == nil || o.Keyvalues == nil {
		return nil, false
	}
	return o.Keyvalues, true
}

// HasKeyvalues returns a boolean if a field has been set.
func (o *PinataMetadata) HasKeyvalues() bool {
	if o != nil && o.Keyvalues != nil {
		return true
	}

	return false
}

// SetKeyvalues gets a reference to the given map[string]string and assigns it to the Keyvalues field.
func (o *PinataMetadata) SetKeyvalues(v map[string]string) {
	o.Keyvalues = &v
}

func (o PinataMetadata) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Name != nil {
		toSerialize["name"] = o.Name
	}
	if o.Keyvalues != nil {
		toSerialize["keyvalues"] = o.Keyvalues
	}
	return json.Marshal(toSerialize)
}

type NullablePinataMetadata struct {
	value *PinataMetadata
	isSet bool
}

func (v NullablePinataMetadata) Get() *PinataMetadata {
	return v.value
}

func (v *NullablePinataMetadata) Set(val *PinataMetadata) {
	v.value = val
	v.isSet = true
}

func (v NullablePinataMetadata) IsSet() bool {
	return v.isSet
}

func (v *NullablePinataMetadata) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePinataMetadata(val *PinataMetadata) *NullablePinataMetadata {
	return &NullablePinataMetadata{value: val, isSet: true}
}

func (v NullablePinataMetadata) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePinataMetadata) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


