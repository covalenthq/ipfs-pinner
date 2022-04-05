# \FilepinApi

All URIs are relative to *https://pinning-service.example.com*

Method | HTTP request | Description
------------- | ------------- | -------------
[**PinataFileUpload**](FilepinApi.md#PinataFileUpload) | **Post** /pinning/pinFileToIPFS | Upload file to IPFS



## PinataFileUpload

> PinataResponse PinataFileUpload(ctx).File(file).PinataOptions(pinataOptions).PinataMetadata(pinataMetadata).Execute()

Upload file to IPFS



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    file := os.NewFile(1234, "some_file") // *os.File | file you're attempting to upload to pinata
    pinataOptions := *openapiclient.NewPinataOptions() // PinataOptions |  (optional)
    pinataMetadata := *openapiclient.NewPinataMetadata() // PinataMetadata |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.FilepinApi.PinataFileUpload(context.Background()).File(file).PinataOptions(pinataOptions).PinataMetadata(pinataMetadata).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `FilepinApi.PinataFileUpload``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `PinataFileUpload`: PinataResponse
    fmt.Fprintf(os.Stdout, "Response from `FilepinApi.PinataFileUpload`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiPinataFileUploadRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **file** | ***os.File** | file you&#39;re attempting to upload to pinata | 
 **pinataOptions** | [**PinataOptions**](PinataOptions.md) |  | 
 **pinataMetadata** | [**PinataMetadata**](PinataMetadata.md) |  | 

### Return type

[**PinataResponse**](PinataResponse.md)

### Authorization

[accessToken](../README.md#accessToken)

### HTTP request headers

- **Content-Type**: multipart/form-data
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

