# PinataResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**IpfsHash** | Pointer to **string** | This is the IPFS multi-hash provided back for your content | [optional] 
**PinSize** | Pointer to **int32** | This is how large (in bytes) the content you just pinned is | [optional] 
**Timestamp** | Pointer to **string** | This is the timestamp for your content pinning (represented in ISO 8601 format) | [optional] 

## Methods

### NewPinataResponse

`func NewPinataResponse() *PinataResponse`

NewPinataResponse instantiates a new PinataResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPinataResponseWithDefaults

`func NewPinataResponseWithDefaults() *PinataResponse`

NewPinataResponseWithDefaults instantiates a new PinataResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetIpfsHash

`func (o *PinataResponse) GetIpfsHash() string`

GetIpfsHash returns the IpfsHash field if non-nil, zero value otherwise.

### GetIpfsHashOk

`func (o *PinataResponse) GetIpfsHashOk() (*string, bool)`

GetIpfsHashOk returns a tuple with the IpfsHash field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpfsHash

`func (o *PinataResponse) SetIpfsHash(v string)`

SetIpfsHash sets IpfsHash field to given value.

### HasIpfsHash

`func (o *PinataResponse) HasIpfsHash() bool`

HasIpfsHash returns a boolean if a field has been set.

### GetPinSize

`func (o *PinataResponse) GetPinSize() int32`

GetPinSize returns the PinSize field if non-nil, zero value otherwise.

### GetPinSizeOk

`func (o *PinataResponse) GetPinSizeOk() (*int32, bool)`

GetPinSizeOk returns a tuple with the PinSize field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPinSize

`func (o *PinataResponse) SetPinSize(v int32)`

SetPinSize sets PinSize field to given value.

### HasPinSize

`func (o *PinataResponse) HasPinSize() bool`

HasPinSize returns a boolean if a field has been set.

### GetTimestamp

`func (o *PinataResponse) GetTimestamp() string`

GetTimestamp returns the Timestamp field if non-nil, zero value otherwise.

### GetTimestampOk

`func (o *PinataResponse) GetTimestampOk() (*string, bool)`

GetTimestampOk returns a tuple with the Timestamp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTimestamp

`func (o *PinataResponse) SetTimestamp(v string)`

SetTimestamp sets Timestamp field to given value.

### HasTimestamp

`func (o *PinataResponse) HasTimestamp() bool`

HasTimestamp returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


