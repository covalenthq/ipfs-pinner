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

// PinataResponse response back to file pin request
type PinataResponse struct {
	// This is the IPFS multi-hash provided back for your content
	IpfsHash *string `json:"IpfsHash,omitempty"`
	// This is how large (in bytes) the content you just pinned is
	PinSize *int32 `json:"PinSize,omitempty"`
	// This is the timestamp for your content pinning (represented in ISO 8601 format)
	Timestamp *string `json:"Timestamp,omitempty"`
}

// NewPinataResponse instantiates a new PinataResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPinataResponse() *PinataResponse {
	this := PinataResponse{}
	return &this
}

// NewPinataResponseWithDefaults instantiates a new PinataResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPinataResponseWithDefaults() *PinataResponse {
	this := PinataResponse{}
	return &this
}

// GetIpfsHash returns the IpfsHash field value if set, zero value otherwise.
func (o *PinataResponse) GetIpfsHash() string {
	if o == nil || o.IpfsHash == nil {
		var ret string
		return ret
	}
	return *o.IpfsHash
}

// GetIpfsHashOk returns a tuple with the IpfsHash field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PinataResponse) GetIpfsHashOk() (*string, bool) {
	if o == nil || o.IpfsHash == nil {
		return nil, false
	}
	return o.IpfsHash, true
}

// HasIpfsHash returns a boolean if a field has been set.
func (o *PinataResponse) HasIpfsHash() bool {
	if o != nil && o.IpfsHash != nil {
		return true
	}

	return false
}

// SetIpfsHash gets a reference to the given string and assigns it to the IpfsHash field.
func (o *PinataResponse) SetIpfsHash(v string) {
	o.IpfsHash = &v
}

// GetPinSize returns the PinSize field value if set, zero value otherwise.
func (o *PinataResponse) GetPinSize() int32 {
	if o == nil || o.PinSize == nil {
		var ret int32
		return ret
	}
	return *o.PinSize
}

// GetPinSizeOk returns a tuple with the PinSize field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PinataResponse) GetPinSizeOk() (*int32, bool) {
	if o == nil || o.PinSize == nil {
		return nil, false
	}
	return o.PinSize, true
}

// HasPinSize returns a boolean if a field has been set.
func (o *PinataResponse) HasPinSize() bool {
	if o != nil && o.PinSize != nil {
		return true
	}

	return false
}

// SetPinSize gets a reference to the given int32 and assigns it to the PinSize field.
func (o *PinataResponse) SetPinSize(v int32) {
	o.PinSize = &v
}

// GetTimestamp returns the Timestamp field value if set, zero value otherwise.
func (o *PinataResponse) GetTimestamp() string {
	if o == nil || o.Timestamp == nil {
		var ret string
		return ret
	}
	return *o.Timestamp
}

// GetTimestampOk returns a tuple with the Timestamp field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PinataResponse) GetTimestampOk() (*string, bool) {
	if o == nil || o.Timestamp == nil {
		return nil, false
	}
	return o.Timestamp, true
}

// HasTimestamp returns a boolean if a field has been set.
func (o *PinataResponse) HasTimestamp() bool {
	if o != nil && o.Timestamp != nil {
		return true
	}

	return false
}

// SetTimestamp gets a reference to the given string and assigns it to the Timestamp field.
func (o *PinataResponse) SetTimestamp(v string) {
	o.Timestamp = &v
}

func (o PinataResponse) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.IpfsHash != nil {
		toSerialize["IpfsHash"] = o.IpfsHash
	}
	if o.PinSize != nil {
		toSerialize["PinSize"] = o.PinSize
	}
	if o.Timestamp != nil {
		toSerialize["Timestamp"] = o.Timestamp
	}
	return json.Marshal(toSerialize)
}

type NullablePinataResponse struct {
	value *PinataResponse
	isSet bool
}

func (v NullablePinataResponse) Get() *PinataResponse {
	return v.value
}

func (v *NullablePinataResponse) Set(val *PinataResponse) {
	v.value = val
	v.isSet = true
}

func (v NullablePinataResponse) IsSet() bool {
	return v.isSet
}

func (v *NullablePinataResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePinataResponse(val *PinataResponse) *NullablePinataResponse {
	return &NullablePinataResponse{value: val, isSet: true}
}

func (v NullablePinataResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePinataResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
