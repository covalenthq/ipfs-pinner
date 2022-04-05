# PinataOptions

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CidVersion** | Pointer to **string** | CID version IPFS will use when creating a hash for your content | [optional] 
**WrapWithDirectory** | Pointer to **bool** | Wrap your content inside of a directory when adding to IPFS. | [optional] 

## Methods

### NewPinataOptions

`func NewPinataOptions() *PinataOptions`

NewPinataOptions instantiates a new PinataOptions object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPinataOptionsWithDefaults

`func NewPinataOptionsWithDefaults() *PinataOptions`

NewPinataOptionsWithDefaults instantiates a new PinataOptions object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCidVersion

`func (o *PinataOptions) GetCidVersion() string`

GetCidVersion returns the CidVersion field if non-nil, zero value otherwise.

### GetCidVersionOk

`func (o *PinataOptions) GetCidVersionOk() (*string, bool)`

GetCidVersionOk returns a tuple with the CidVersion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCidVersion

`func (o *PinataOptions) SetCidVersion(v string)`

SetCidVersion sets CidVersion field to given value.

### HasCidVersion

`func (o *PinataOptions) HasCidVersion() bool`

HasCidVersion returns a boolean if a field has been set.

### GetWrapWithDirectory

`func (o *PinataOptions) GetWrapWithDirectory() bool`

GetWrapWithDirectory returns the WrapWithDirectory field if non-nil, zero value otherwise.

### GetWrapWithDirectoryOk

`func (o *PinataOptions) GetWrapWithDirectoryOk() (*bool, bool)`

GetWrapWithDirectoryOk returns a tuple with the WrapWithDirectory field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWrapWithDirectory

`func (o *PinataOptions) SetWrapWithDirectory(v bool)`

SetWrapWithDirectory sets WrapWithDirectory field to given value.

### HasWrapWithDirectory

`func (o *PinataOptions) HasWrapWithDirectory() bool`

HasWrapWithDirectory returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


