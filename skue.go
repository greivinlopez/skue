// The MIT License (MIT)
//
// Copyright (c) 2013 Greivin López
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
package skue

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// ----------------------------------------------------------------------------
// LIST OF CONSTANTS

const (
	MIME_XML  = "application/xml"
	MIME_JSON = "application/json"

	HEADER_Allow                         = "Allow"
	HEADER_Accept                        = "Accept"
	HEADER_Origin                        = "Origin"
	HEADER_ContentType                   = "Content-Type"
	HEADER_LastModified                  = "Last-Modified"
	HEADER_AcceptEncoding                = "Accept-Encoding"
	HEADER_ContentEncoding               = "Content-Encoding"
	HEADER_AccessControlExposeHeaders    = "Access-Control-Expose-Headers"
	HEADER_AccessControlRequestMethod    = "Access-Control-Request-Method"
	HEADER_AccessControlRequestHeaders   = "Access-Control-Request-Headers"
	HEADER_AccessControlAllowMethods     = "Access-Control-Allow-Methods"
	HEADER_AccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	HEADER_AccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	HEADER_AccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	HEADER_XRateLimitLimit               = "X-Rate-Limit-Limit"
	HEADER_XRateLimitRemaining           = "X-Rate-Limit-Remaining"
)

// SimpleMessage represents an HTTP simple response with
// an HTTP status code and a response message
type SimpleMessage struct {
	Status  int
	Message string
}

// ToJson is a convenience method for writing a value in json encoding
func ToJson(w http.ResponseWriter, statusCode int, value interface{}) {
	output, err := json.MarshalIndent(value, " ", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())
	} else {
		w.Header().Set(HEADER_ContentType, MIME_JSON)
		w.WriteHeader(statusCode)
		w.Write(output)
	}
}

// FromJson checks the Accept header and reads the content into the entityPointer
func FromJson(r *http.Request, entityPointer interface{}) error {
	contentType := r.Header.Get(HEADER_ContentType)
	buffer, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if strings.Contains(contentType, MIME_JSON) {
		err = json.Unmarshal(buffer, entityPointer)
	} else {
		err = errors.New("Unable to unmarshal content of type:" + contentType)
	}
	return err
}

// ServiceResponse is a convenience function to create an http response
// encoded as JSON with a simple message
func ServiceResponse(w http.ResponseWriter, httpStatus int, message string) {
	simpleMessage := &SimpleMessage{httpStatus, message}
	ToJson(w, httpStatus, simpleMessage)
}

// ----------------------------------------------------------------------------
// HANDLERS

// NotAllowed handler will response with a "405 Method Not Allowed" response
// It is a convenience handler to route all not allowed services
func NotAllowed(w http.ResponseWriter) {
	ServiceResponse(w, http.StatusMethodNotAllowed, "Method Not Allowed")
}

// NotFound handler will respond with a "404 Not Found" response
func NotFound(w http.ResponseWriter) {
	ServiceResponse(w, http.StatusNotFound, "Item not found")
}
