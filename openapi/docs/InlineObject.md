# InlineObject

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**File** | ***os.File** | file you&#39;re attempting to upload to pinata | 
**PinataOptions** | Pointer to [**PinataOptions**](PinataOptions.md) |  | [optional] 
**PinataMetadata** | Pointer to [**PinataMetadata**](PinataMetadata.md) |  | [optional] 

## Methods

### NewInlineObject

`func NewInlineObject(file *os.File, ) *InlineObject`

NewInlineObject instantiates a new InlineObject object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewInlineObjectWithDefaults

`func NewInlineObjectWithDefaults() *InlineObject`

NewInlineObjectWithDefaults instantiates a new InlineObject object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetFile

`func (o *InlineObject) GetFile() *os.File`

GetFile returns the File field if non-nil, zero value otherwise.

### GetFileOk

`func (o *InlineObject) GetFileOk() (**os.File, bool)`

GetFileOk returns a tuple with the File field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFile

`func (o *InlineObject) SetFile(v *os.File)`

SetFile sets File field to given value.


### GetPinataOptions

`func (o *InlineObject) GetPinataOptions() PinataOptions`

GetPinataOptions returns the PinataOptions field if non-nil, zero value otherwise.

### GetPinataOptionsOk

`func (o *InlineObject) GetPinataOptionsOk() (*PinataOptions, bool)`

GetPinataOptionsOk returns a tuple with the PinataOptions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPinataOptions

`func (o *InlineObject) SetPinataOptions(v PinataOptions)`

SetPinataOptions sets PinataOptions field to given value.

### HasPinataOptions

`func (o *InlineObject) HasPinataOptions() bool`

HasPinataOptions returns a boolean if a field has been set.

### GetPinataMetadata

`func (o *InlineObject) GetPinataMetadata() PinataMetadata`

GetPinataMetadata returns the PinataMetadata field if non-nil, zero value otherwise.

### GetPinataMetadataOk

`func (o *InlineObject) GetPinataMetadataOk() (*PinataMetadata, bool)`

GetPinataMetadataOk returns a tuple with the PinataMetadata field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPinataMetadata

`func (o *InlineObject) SetPinataMetadata(v PinataMetadata)`

SetPinataMetadata sets PinataMetadata field to given value.

### HasPinataMetadata

`func (o *InlineObject) HasPinataMetadata() bool`

HasPinataMetadata returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


