// package api handles api calls for data transformations
package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// API Action carries all the information required to
// make an API call
type ApiAction struct {
	// THe URI of the API
	API string
	// MRX ID of the API
	MrxID string
	// ResponseMIMEType of the output data
	ResponseMIMEType  string
	APISchemaLocation string
	APIParams         []Parameter
}

// The optional parameters for
// calling an API
type Parameter struct {
	Key         string `json:"key"`
	Optional    bool   `json:"optional"`
	Description string `json:"description"`
}

// Transform takes an array of Bytes and  makes a series of
// API calls to transform them into the target metadata.
func (a *ApiAction) Transform(in []byte, params url.Values) ([]byte, error) {

	// set up and extra parameters to send to the target URL
	apiParameters := ""
	if len(a.APIParams) > 0 {
		apiParameters = "?"
		for _, p := range a.APIParams {
			val := params.Get(p.Key)
			if val == "" {
				// don't include
				if !p.Optional {
					return nil, fmt.Errorf("missing required parameter: %v", p.Key)
				}
			} else {
				if apiParameters != "?" {
					apiParameters += "&"
				}
				apiParameters += p.Key + "=" + val
			}
		}
	}

	out, err := ApiExtractBytes(in, a.API+apiParameters, a.APISchemaLocation, a.ResponseMIMEType)

	// do any extra things here
	return out, err
}

// ActionType describes the API's action for logging
func (a *ApiAction) ActionType() string {
	return fmt.Sprintf("Making a POST Request to %v", a.API)
}

// DataID gives the metarex ID of the data being transformed
func (a *ApiAction) DataID() string {
	return a.MrxID
}

// the json format for any error messages
// returned by the api
type ErrMessage struct {
	Error string `json:"error"`
}

// API Extract takes bytes and makes the API call returning the bytes
// It validates the inputs and outputs against the OpenAPI specification of that API
func ApiExtractBytes(toTransform []byte, API, APISpec, dataFormat string) ([]byte, error) {

	// Load the OpenAPI Spec
	//loader := openapi3.NewLoader()
	//doc, err := loader.LoadFromFile(APISpec)
	//if err != nil {
	//	return nil, fmt.Errorf("error validating against API schema :%v", err.Error())
	//

	// convert to jsonbytes then make the call
	resp, err := http.Post(API, dataFormat, bytes.NewReader(toTransform))
	if err != nil {
		return nil, fmt.Errorf("error making POST request: %v", err.Error())
	}
	/*
		// API Schema Fluff
		// checking for errors afterwards
		// generate a copy of the request body that was made to check against
		ctx := context.Background()
		httpReq, err := http.NewRequest(http.MethodPost, API, bytes.NewReader(in))
		if err != nil {
			return nil, fmt.Errorf("error in POST request at data point %v : %v", i, err.Error())
		}
		// add the real header
		httpReq.Header.Add("Content-Type", dataFormat)

		router, err := gorillamux.NewRouter(doc)
		if err != nil {
			return nil, err
		}
		fmt.Println(httpReq)
		route, pathParams, err := router.FindRoute(httpReq)

		if err != nil {
			return nil, fmt.Errorf("error finding API route data point %v : %v", i, err.Error())
		}

		// Validate request
		requestValidationInput := &openapi3filter.RequestValidationInput{
			Request:    httpReq,
			PathParams: pathParams,
			Route:      route,
		}

		err = openapi3filter.ValidateRequest(ctx, requestValidationInput)
		switch dataFormat {
		// ignore responses that can't be decoded
		case "image/png", "image/jpeg", "application/xml",
			"application/octet-stream", "text/csv", "audio/wav":

		default:
			if err != nil {
				return nil, fmt.Errorf("error validating API request at data point %v : %v", i, err.Error())
			}
		}
	*/
	// get the result and return it as the required golang struct
	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading the data from the API response: %v", err.Error())
	}

	/*
		// Validate response using OpenAPI
		responseValidationInput := &openapi3filter.ResponseValidationInput{
			RequestValidationInput: requestValidationInput,
			Status:                 resp.StatusCode,
			Header:                 resp.Header,
		}
		responseValidationInput.SetBodyBytes(resBody)
		err = openapi3filter.ValidateResponse(ctx, responseValidationInput)
		ctypes := resp.Header[http.CanonicalHeaderKey("Content-Type")]
		ctype := ""
		if len(ctypes) > 0 {
			ctype = ctypes[0]
		}

		switch ctype {
		// ignore responses that can't be decoded
		case "image/png", "image/jpeg", "application/octet-stream",
			"text/plain", "application/xml":

		default:
			if err != nil {
				return nil, fmt.Errorf("error validating API response at data point %v : %v", i, err.Error())
			}
		}*/

	if resp.StatusCode != http.StatusOK {
		var e ErrMessage
		json.Unmarshal(resBody, &e)
		return nil, fmt.Errorf("%v", e.Error)
	}

	return resBody, nil
}
