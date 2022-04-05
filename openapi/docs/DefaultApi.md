# \DefaultApi

All URIs are relative to *https://pinning-service.example.com*

Method | HTTP request | Description
------------- | ------------- | -------------
[**PinningPinFileToIPFSPost**](DefaultApi.md#PinningPinFileToIPFSPost) | **Post** /pinning/pinFileToIPFS | Upload file to IPFS



## PinningPinFileToIPFSPost

> PinataResponse PinningPinFileToIPFSPost(ctx).PinataFilePinRequest(pinataFilePinRequest).Execute()

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
    pinataFilePinRequest := *openapiclient.NewPinataFilePinRequest("File_example") // PinataFilePinRequest | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.DefaultApi.PinningPinFileToIPFSPost(context.Background()).PinataFilePinRequest(pinataFilePinRequest).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.PinningPinFileToIPFSPost``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `PinningPinFileToIPFSPost`: PinataResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.PinningPinFileToIPFSPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiPinningPinFileToIPFSPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **pinataFilePinRequest** | [**PinataFilePinRequest**](PinataFilePinRequest.md) |  | 

### Return type

[**PinataResponse**](PinataResponse.md)

### Authorization

[accessToken](../README.md#accessToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

