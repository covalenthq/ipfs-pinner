# \PinsApi

All URIs are relative to *https://pinning-service.example.com*

Method | HTTP request | Description
------------- | ------------- | -------------
[**PinsGet**](PinsApi.md#PinsGet) | **Get** /pins | List pin objects
[**PinsPost**](PinsApi.md#PinsPost) | **Post** /pins | Add pin object
[**PinsRequestidDelete**](PinsApi.md#PinsRequestidDelete) | **Delete** /pins/{requestid} | Remove pin object
[**PinsRequestidGet**](PinsApi.md#PinsRequestidGet) | **Get** /pins/{requestid} | Get pin object
[**PinsRequestidPost**](PinsApi.md#PinsRequestidPost) | **Post** /pins/{requestid} | Replace pin object



## PinsGet

> PinResults PinsGet(ctx).Cid(cid).Name(name).Match(match).Status(status).Before(before).After(after).Limit(limit).Meta(meta).Execute()

List pin objects



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    "time"
    openapiclient "./openapi"
)

func main() {
    cid := []string{"Inner_example"} // []string | Return pin objects responsible for pinning the specified CID(s); be aware that using longer hash functions introduces further constraints on the number of CIDs that will fit under the limit of 2000 characters per URL  in browser contexts (optional)
    name := "PreciousData.pdf" // string | Return pin objects with specified name (by default a case-sensitive, exact match) (optional)
    match := openapiclient.TextMatchingStrategy("exact") // TextMatchingStrategy | Customize the text matching strategy applied when the name filter is present; exact (the default) is a case-sensitive exact match, partial matches anywhere in the name, iexact and ipartial are case-insensitive versions of the exact and partial strategies (optional) (default to "exact")
    status := []openapiclient.Status{openapiclient.Status("queued")} // []Status | Return pin objects for pins with the specified status (optional)
    before := time.Now() // time.Time | Return results created (queued) before provided timestamp (optional)
    after := time.Now() // time.Time | Return results created (queued) after provided timestamp (optional)
    limit := int32(56) // int32 | Max records to return (optional) (default to 10)
    meta := map[string]string{"key": "Inner_example"} // map[string]string | Return pin objects that match specified metadata keys passed as a string representation of a JSON object; when implementing a client library, make sure the parameter is URL-encoded to ensure safe transport (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.PinsApi.PinsGet(context.Background()).Cid(cid).Name(name).Match(match).Status(status).Before(before).After(after).Limit(limit).Meta(meta).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `PinsApi.PinsGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `PinsGet`: PinResults
    fmt.Fprintf(os.Stdout, "Response from `PinsApi.PinsGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiPinsGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **cid** | **[]string** | Return pin objects responsible for pinning the specified CID(s); be aware that using longer hash functions introduces further constraints on the number of CIDs that will fit under the limit of 2000 characters per URL  in browser contexts | 
 **name** | **string** | Return pin objects with specified name (by default a case-sensitive, exact match) | 
 **match** | [**TextMatchingStrategy**](TextMatchingStrategy.md) | Customize the text matching strategy applied when the name filter is present; exact (the default) is a case-sensitive exact match, partial matches anywhere in the name, iexact and ipartial are case-insensitive versions of the exact and partial strategies | [default to &quot;exact&quot;]
 **status** | [**[]Status**](Status.md) | Return pin objects for pins with the specified status | 
 **before** | **time.Time** | Return results created (queued) before provided timestamp | 
 **after** | **time.Time** | Return results created (queued) after provided timestamp | 
 **limit** | **int32** | Max records to return | [default to 10]
 **meta** | **map[string]string** | Return pin objects that match specified metadata keys passed as a string representation of a JSON object; when implementing a client library, make sure the parameter is URL-encoded to ensure safe transport | 

### Return type

[**PinResults**](PinResults.md)

### Authorization

[accessToken](../README.md#accessToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PinsPost

> PinStatus PinsPost(ctx).Pin(pin).Execute()

Add pin object



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
    pin := *openapiclient.NewPin("QmCIDToBePinned") // Pin | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.PinsApi.PinsPost(context.Background()).Pin(pin).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `PinsApi.PinsPost``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `PinsPost`: PinStatus
    fmt.Fprintf(os.Stdout, "Response from `PinsApi.PinsPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiPinsPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **pin** | [**Pin**](Pin.md) |  | 

### Return type

[**PinStatus**](PinStatus.md)

### Authorization

[accessToken](../README.md#accessToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PinsRequestidDelete

> PinsRequestidDelete(ctx, requestid).Execute()

Remove pin object



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
    requestid := "requestid_example" // string | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.PinsApi.PinsRequestidDelete(context.Background(), requestid).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `PinsApi.PinsRequestidDelete``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**requestid** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiPinsRequestidDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

[accessToken](../README.md#accessToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PinsRequestidGet

> PinStatus PinsRequestidGet(ctx, requestid).Execute()

Get pin object



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
    requestid := "requestid_example" // string | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.PinsApi.PinsRequestidGet(context.Background(), requestid).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `PinsApi.PinsRequestidGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `PinsRequestidGet`: PinStatus
    fmt.Fprintf(os.Stdout, "Response from `PinsApi.PinsRequestidGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**requestid** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiPinsRequestidGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**PinStatus**](PinStatus.md)

### Authorization

[accessToken](../README.md#accessToken)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PinsRequestidPost

> PinStatus PinsRequestidPost(ctx, requestid).Pin(pin).Execute()

Replace pin object



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
    requestid := "requestid_example" // string | 
    pin := *openapiclient.NewPin("QmCIDToBePinned") // Pin | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.PinsApi.PinsRequestidPost(context.Background(), requestid).Pin(pin).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `PinsApi.PinsRequestidPost``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `PinsRequestidPost`: PinStatus
    fmt.Fprintf(os.Stdout, "Response from `PinsApi.PinsRequestidPost`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**requestid** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiPinsRequestidPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **pin** | [**Pin**](Pin.md) |  | 

### Return type

[**PinStatus**](PinStatus.md)

### Authorization

[accessToken](../README.md#accessToken)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

