<!-- Code generated by gomarkdoc. DO NOT EDIT -->

# utils

```go
import "github.com/atropos112/atrogolib/utils"
```

## Index

- [func ArrContains\[T constraints.Ordered\]\(arr \[\]T, obj T\) bool](<#ArrContains>)
- [func ArrContainsArr\[T constraints.Ordered\]\(arr \[\]T, subArr \[\]T\) bool](<#ArrContainsArr>)
- [func GetCred\(value string\) \(string, error\)](<#GetCred>)
- [func GetCredUnsafe\(value string\) string](<#GetCredUnsafe>)
- [func MakeAPIRequest\(client \*http.Client, kind, apiBaseURL, endpoint, token string, request, response interface\{\}\) error](<#MakeAPIRequest>)
- [func MakeDeleteRequest\(client \*http.Client, apiBaseURL, endpoint, token string, response any\) error](<#MakeDeleteRequest>)
- [func MakeGetRequest\(client \*http.Client, apiBaseURL, endpoint, token string, response any\) error](<#MakeGetRequest>)
- [func MakePostRequest\(client \*http.Client, apiBaseURL, endpoint, token string, request, response any\) error](<#MakePostRequest>)
- [func MakePutRequest\(client \*http.Client, apiBaseURL, endpoint, token string, request, response any\) error](<#MakePutRequest>)
- [func RunAPIServer\(port int\) error](<#RunAPIServer>)
- [type APIError](<#APIError>)
  - [func \(e \*APIError\) Error\(\) string](<#APIError.Error>)
- [type AuthenticatedAPIClient](<#AuthenticatedAPIClient>)
  - [func NewAPIClient\(baseURL, token string\) AuthenticatedAPIClient](<#NewAPIClient>)
  - [func \(c \*AuthenticatedAPIClient\) Delete\(endpoint string, response interface\{\}\) error](<#AuthenticatedAPIClient.Delete>)
  - [func \(c \*AuthenticatedAPIClient\) Get\(endpoint string, response interface\{\}\) error](<#AuthenticatedAPIClient.Get>)
  - [func \(c \*AuthenticatedAPIClient\) Post\(endpoint string, request, response interface\{\}\) error](<#AuthenticatedAPIClient.Post>)
  - [func \(c \*AuthenticatedAPIClient\) Put\(endpoint string, request, response interface\{\}\) error](<#AuthenticatedAPIClient.Put>)
- [type DeveloperError](<#DeveloperError>)
  - [func \(e \*DeveloperError\) Error\(\) string](<#DeveloperError.Error>)
- [type GPTDoesntListenError](<#GPTDoesntListenError>)
  - [func \(e \*GPTDoesntListenError\) Error\(\) string](<#GPTDoesntListenError.Error>)
- [type NoCredFoundError](<#NoCredFoundError>)
  - [func \(e \*NoCredFoundError\) Error\(\) string](<#NoCredFoundError.Error>)


<a name="ArrContains"></a>
## func [ArrContains](<https://github.com/atropos112/atrogolib/blob/main/utils/arr.go#L16>)

```go
func ArrContains[T constraints.Ordered](arr []T, obj T) bool
```

ArrContains checks if an array contains obj

<a name="ArrContainsArr"></a>
## func [ArrContainsArr](<https://github.com/atropos112/atrogolib/blob/main/utils/arr.go#L6>)

```go
func ArrContainsArr[T constraints.Ordered](arr []T, subArr []T) bool
```

ArrContainsArr checks if an array contains all elements of another array

<a name="GetCred"></a>
## func [GetCred](<https://github.com/atropos112/atrogolib/blob/main/utils/creds.go#L20>)

```go
func GetCred(value string) (string, error)
```

GetCred is a function that gets a credential from the environment variables. If the credential is not found, it will return an error.

<a name="GetCredUnsafe"></a>
## func [GetCredUnsafe](<https://github.com/atropos112/atrogolib/blob/main/utils/creds.go#L10>)

```go
func GetCredUnsafe(value string) string
```

GetCredUnsafe is a function that gets a credential from the environment variables. If the credential is not found, it will log a fatal error.

<a name="MakeAPIRequest"></a>
## func [MakeAPIRequest](<https://github.com/atropos112/atrogolib/blob/main/utils/api.go#L27>)

```go
func MakeAPIRequest(client *http.Client, kind, apiBaseURL, endpoint, token string, request, response interface{}) error
```

MakeAPIRequest is a generic function to make an API request. It supports GET, POST, PUT, and DELETE requests.

<a name="MakeDeleteRequest"></a>
## func [MakeDeleteRequest](<https://github.com/atropos112/atrogolib/blob/main/utils/api.go#L99>)

```go
func MakeDeleteRequest(client *http.Client, apiBaseURL, endpoint, token string, response any) error
```

MakeDeleteRequest is a helper function to make a DELETE request to the specified endpoint. If token is not "" it will be added to the request as a Bearer token.

<a name="MakeGetRequest"></a>
## func [MakeGetRequest](<https://github.com/atropos112/atrogolib/blob/main/utils/api.go#L104>)

```go
func MakeGetRequest(client *http.Client, apiBaseURL, endpoint, token string, response any) error
```

MakeGetRequest is a helper function to make a GET request to the specified endpoint. If token is not "" it will be added to the request as a Bearer token.

<a name="MakePostRequest"></a>
## func [MakePostRequest](<https://github.com/atropos112/atrogolib/blob/main/utils/api.go#L109>)

```go
func MakePostRequest(client *http.Client, apiBaseURL, endpoint, token string, request, response any) error
```

MakePostRequest is a helper function to make a POST request to the specified endpoint. If token is not "" it will be added to the request as a Bearer token.

<a name="MakePutRequest"></a>
## func [MakePutRequest](<https://github.com/atropos112/atrogolib/blob/main/utils/api.go#L114>)

```go
func MakePutRequest(client *http.Client, apiBaseURL, endpoint, token string, request, response any) error
```

MakePutRequest is a helper function to make a PUT request to the specified endpoint. If token is not "" it will be added to the request as a Bearer token.

<a name="RunAPIServer"></a>
## func [RunAPIServer](<https://github.com/atropos112/atrogolib/blob/main/utils/api.go#L119>)

```go
func RunAPIServer(port int) error
```

RunAPIServer attaches logging middleware to the default http server and starts it on the specified port.

<a name="APIError"></a>
## type [APIError](<https://github.com/atropos112/atrogolib/blob/main/utils/api.go#L17-L20>)

APIError is an error type that is returned when an API request fails.

```go
type APIError struct {
    StatusCode int
    Message    string
}
```

<a name="APIError.Error"></a>
### func \(\*APIError\) [Error](<https://github.com/atropos112/atrogolib/blob/main/utils/api.go#L22>)

```go
func (e *APIError) Error() string
```



<a name="AuthenticatedAPIClient"></a>
## type [AuthenticatedAPIClient](<https://github.com/atropos112/atrogolib/blob/main/utils/types.go#L9-L13>)

AuthenticatedAPIClient is a struct that contains the base URL of the API and the token to use for requests.

```go
type AuthenticatedAPIClient struct {
    BaseURL string
    Token   string
    Client  *http.Client
}
```

<a name="NewAPIClient"></a>
### func [NewAPIClient](<https://github.com/atropos112/atrogolib/blob/main/utils/types.go#L16>)

```go
func NewAPIClient(baseURL, token string) AuthenticatedAPIClient
```

NewAPIClient creates a new AuthenticatedAPIClient with the specified base URL and token.

<a name="AuthenticatedAPIClient.Delete"></a>
### func \(\*AuthenticatedAPIClient\) [Delete](<https://github.com/atropos112/atrogolib/blob/main/utils/types.go#L25>)

```go
func (c *AuthenticatedAPIClient) Delete(endpoint string, response interface{}) error
```

Delete is a helper function to make a DELETE request to the specified endpoint. If token is not "" it will be added to the request as a Bearer token.

<a name="AuthenticatedAPIClient.Get"></a>
### func \(\*AuthenticatedAPIClient\) [Get](<https://github.com/atropos112/atrogolib/blob/main/utils/types.go#L30>)

```go
func (c *AuthenticatedAPIClient) Get(endpoint string, response interface{}) error
```

Get is a helper function to make a GET request to the specified endpoint. If token is not "" it will be added to the request as a Bearer token.

<a name="AuthenticatedAPIClient.Post"></a>
### func \(\*AuthenticatedAPIClient\) [Post](<https://github.com/atropos112/atrogolib/blob/main/utils/types.go#L35>)

```go
func (c *AuthenticatedAPIClient) Post(endpoint string, request, response interface{}) error
```

Post is a helper function to make a POST request to the specified endpoint. If token is not "" it will be added to the request as a Bearer token.

<a name="AuthenticatedAPIClient.Put"></a>
### func \(\*AuthenticatedAPIClient\) [Put](<https://github.com/atropos112/atrogolib/blob/main/utils/types.go#L40>)

```go
func (c *AuthenticatedAPIClient) Put(endpoint string, request, response interface{}) error
```

Put is a helper function to make a PUT request to the specified endpoint. If token is not "" it will be added to the request as a Bearer token.

<a name="DeveloperError"></a>
## type [DeveloperError](<https://github.com/atropos112/atrogolib/blob/main/utils/types.go#L54-L56>)

DeveloperError represents an error that is caused by a developer mistake

```go
type DeveloperError struct {
    Message string
}
```

<a name="DeveloperError.Error"></a>
### func \(\*DeveloperError\) [Error](<https://github.com/atropos112/atrogolib/blob/main/utils/types.go#L58>)

```go
func (e *DeveloperError) Error() string
```



<a name="GPTDoesntListenError"></a>
## type [GPTDoesntListenError](<https://github.com/atropos112/atrogolib/blob/main/utils/types.go#L63-L66>)

GPTDoesntListenError represents an error when GPT doesn't listen

```go
type GPTDoesntListenError struct {
    UserMessage string
    SysMessage  string
}
```

<a name="GPTDoesntListenError.Error"></a>
### func \(\*GPTDoesntListenError\) [Error](<https://github.com/atropos112/atrogolib/blob/main/utils/types.go#L68>)

```go
func (e *GPTDoesntListenError) Error() string
```



<a name="NoCredFoundError"></a>
## type [NoCredFoundError](<https://github.com/atropos112/atrogolib/blob/main/utils/types.go#L45-L47>)

NoCredFoundError represents an error when no credentials are found

```go
type NoCredFoundError struct {
    CredentialName string
}
```

<a name="NoCredFoundError.Error"></a>
### func \(\*NoCredFoundError\) [Error](<https://github.com/atropos112/atrogolib/blob/main/utils/types.go#L49>)

```go
func (e *NoCredFoundError) Error() string
```



Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)
