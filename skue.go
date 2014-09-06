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
	"fmt"
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

type Modeler interface {
	CollectionName() string
	Create() (err error)
	Update() (err error)
	Delete() (err error)
}

type SimpleMessage struct {
	Status  int
	Message string
}

var ErrNotFound = errors.New("Not Found")

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
// encoded as JSON.
func ServiceResponse(w http.ResponseWriter, httpStatus int, message string) {
	simpleMessage := &SimpleMessage{httpStatus, message}
	ToJson(w, httpStatus, simpleMessage)
}

// ----------------------------------------------------------------------------
// PERSISTANCE UTILS

// Saves a model to the underlying storage
func SaveModel(modeler Modeler, w http.ResponseWriter, r *http.Request) {
	err := FromJson(r, &modeler)

	if err != nil {
		ServiceResponse(w, http.StatusBadRequest, fmt.Sprintf("Failed creating the owner: %v", err))
	} else {
		err = modeler.Create()
		if err != nil {
			ServiceResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed saving the object: %v", err))
		} else {
			ToJson(w, http.StatusCreated, modeler)
		}
	}
}

// Updates the given model in the underlying storage
func UpdateModel(modeler Modeler, w http.ResponseWriter, r *http.Request) {
	err := FromJson(r, &modeler)

	if err != nil {
		ServiceResponse(w, http.StatusBadRequest, fmt.Sprintf("Failed reading JSON from request: %v", err))
	} else {
		err = modeler.Update()
		if err != nil {
			if err == ErrNotFound {
				ServiceResponse(w, http.StatusNotFound, "Item not found")
			} else {
				ServiceResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed updating the object: %v", err))
			}
		} else {
			ServiceResponse(w, http.StatusOK, "Successfully updated")
		}
	}
}

// ----------------------------------------------------------------------------
// HANDLERS

// NotAllowed handler will response with a "405 Method Not Allowed" response
// It is a convenience handler to route all not allowed services
func NotAllowed(w http.ResponseWriter) {
	ServiceResponse(w, http.StatusMethodNotAllowed, "Method Not Allowed")
}
