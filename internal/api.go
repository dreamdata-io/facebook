// A facebook graph api client in go.
// https://github.com/huandu/facebook/
//
// Copyright 2012, Huan Du
// Licensed under the MIT license
// https://github.com/huandu/facebook/blob/master/LICENSE

// Package facebook is a Go library fully supports Facebook Graph API with
// file upload and batch request. It can be used in Google App Engine.
//
// Library design is highly influenced by facebook official PHP/JS SDK.
// If you have experience with PHP/JS SDK, you may feel quite familiar with it.
//
// Go to project home page to see samples. Link: https://github.com/huandu/facebook
package internal

import (
	"net/http"
)

// Facebook graph api methods.
const (
	GET    Method = "GET"
	POST   Method = "POST"
	DELETE Method = "DELETE"
	PUT    Method = "PUT"
)

var (
	// Version is the default facebook api version.
	// It can be any valid version string (e.g. "v2.3") or empty.
	//
	// See https://developers.facebook.com/docs/apps/versions for details.
	Version string

	// Debug is the app level debug mode.
	// After setting DebugMode, all newly created session will use the mode
	// to communicate with graph API.
	//
	// See https://developers.facebook.com/docs/graph-api/using-graph-api/v2.3#debugging
	Debug DebugMode

	// RFC3339Timestamps will cause all requests to return RFC3339 formatted timestamps.
	// RFC3339 is the timestamp format that the built in UnmarshalJSON expects.
	RFC3339Timestamps bool
)

var (
	// default facebook session.
	defaultSession = &Session{}
)

// DebugMode is the debug mode of Graph API.
// See https://developers.facebook.com/docs/graph-api/using-graph-api/v2.3#graphapidebugmode
type DebugMode string

// Method is HTTP method for an API call.
// Can be GET, POST or DELETE.
type Method string

// Api makes a facebook graph api call with default session.
//
// Method can be GET, POST, DELETE or PUT.
//
// Params represents query strings in this call.
// Keys and values in params will be encoded into the URL automatically, so there is
// no need to encode keys or values in params manually. Params can be nil.
//
// If you want to get
//
//	https://graph.facebook.com/huandu?fields=name,username
//
// Api should be called as following
//
//	Api("/huandu", GET, Params{"fields": "name,username"})
//
// or in a simplified way
//
//	Get("/huandu", Params{"fields": "name,username"})
//
// Api is a wrapper of Session.Api(). It's designed for graph api that doesn't require
// app id, app secret and access token. It can be called in multiple goroutines.
//
// If app id, app secret or access token is required in graph api, caller should
// create a new facebook session through App instance instead.
func Api(path string, method Method, params Params) (Result, error) {
	return defaultSession.Api(path, method, params)
}

// Get is a short hand of Api(path, GET, params).
func Get(path string, params Params) (Result, error) {
	return Api(path, GET, params)
}

// Post is a short hand of Api(path, POST, params).
func Post(path string, params Params) (Result, error) {
	return Api(path, POST, params)
}

// Delete is a short hand of Api(path, DELETE, params).
func Delete(path string, params Params) (Result, error) {
	return Api(path, DELETE, params)
}

// Put is a short hand of Api(path, PUT, params).
func Put(path string, params Params) (Result, error) {
	return Api(path, PUT, params)
}

// BatchApi makes a batch facebook graph api call with default session.
//
// BatchApi supports most kinds of batch calls defines in facebook batch api document,
// except uploading binary data. Use Batch to do so.
//
// Note: API response is stored in "body" field of a Result.
//
//	results, _ := BatchApi(accessToken, Params{...}, Params{...})
//
//	// Use first batch api response.
//	var res1 *BatchResult
//	var err error
//	res1, err = results[0].Batch()
//
//	if err != nil {
//	    // this is not a valid batch api response.
//	}
//
//	// Use BatchResult#Result to get response body content as Result.
//	res := res1.Result
//
// Facebook document: https://developers.facebook.com/docs/graph-api/making-multiple-requests
func BatchApi(accessToken string, params ...Params) ([]Result, error) {
	return Batch(Params{"access_token": accessToken}, params...)
}

// Batch makes a batch facebook graph api call with default session.
// Batch is designed for more advanced usage including uploading binary files.
//
// An uploading files sample
//
//	// equivalent to following curl command (borrowed from facebook docs)
//	//     curl \
//	//         -F 'access_token=…' \
//	//         -F 'batch=[{"method":"POST","relative_url":"me/photos","body":"message=My cat photo","attached_files":"file1"},{"method":"POST","relative_url":"me/photos","body":"message=My dog photo","attached_files":"file2"}]' \
//	//         -F 'file1=@cat.gif' \
//	//         -F 'file2=@dog.jpg' \
//	//         https://graph.facebook.com
//	Batch(Params{
//	    "access_token": "the-access-token",
//	    "file1": File("cat.gif"),
//	    "file2": File("dog.jpg"),
//	}, Params{
//	    "method": "POST",
//	    "relative_url": "me/photos",
//	    "body": "message=My cat photo",
//	    "attached_files": "file1",
//	}, Params{
//	    "method": "POST",
//	    "relative_url": "me/photos",
//	    "body": "message=My dog photo",
//	    "attached_files": "file2",
//	})
//
// Facebook document: https://developers.facebook.com/docs/graph-api/making-multiple-requests
func Batch(batchParams Params, params ...Params) ([]Result, error) {
	return defaultSession.Batch(batchParams, params...)
}

// Request makes an arbitrary HTTP request with default session.
// It expects server responses a facebook Graph API response.
//
//	request, _ := http.NewRequest("https://graph.facebook.com/538744468", "GET", nil)
//	res, err := Request(request)
//	fmt.Println(res["gender"])  // get "male"
func Request(request *http.Request) (Result, error) {
	return defaultSession.Request(request)
}

// DefaultHttpClient returns the http client for default session.
func DefaultHttpClient() HttpClient {
	return defaultSession.HttpClient
}

// SetHttpClient updates the http client of default session.
func SetHttpClient(client HttpClient) {
	defaultSession.HttpClient = client
}
