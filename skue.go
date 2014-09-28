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
	"io"
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

// Producer is intended to be an encoder that writes a value to http writers
// for a particular MIME type.
type Producer interface {
	MimeType() string
	Out(w http.ResponseWriter, statusCode int, value interface{})
}

// Consumer is intended to be a decoder of HTTP requests that uses a particular
// MIME type to decode the intended object into a value.
type Consumer interface {
	MimeType() string
	In(r *http.Request, value interface{}) error
}

func Produce(producer Producer, w http.ResponseWriter, r *http.Request, status int, value interface{}) {
	acceptEncoding := r.Header.Get(HEADER_AcceptEncoding)
	// According to HTTP/1.1 protocol section 14.1 about Accept header field
	// "If an Accept header field is present, and if the server cannot send
	// a response which is acceptable according to the combined Accept field
	// value, then the server SHOULD send a 406 (not acceptable) response."
	if acceptEncoding != "" && !strings.Contains(acceptEncoding, "*/*") {
		if !strings.Contains(acceptEncoding, producer.MimeType()) {
			w.WriteHeader(http.StatusNotAcceptable)
			io.WriteString(w, "Not Acceptable")
		}
	}
	producer.Out(w, status, value)
}

func Consume(consumer Consumer, w http.ResponseWriter, r *http.Request, value interface{}) error {
	contentType := r.Header.Get(HEADER_ContentType)
	// According to HTTP/1.1 protocol section 14.17 about Content-Type header
	if !strings.Contains(contentType, consumer.MimeType()) {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		io.WriteString(w, "Unsupported Media Type")
	}
	return consumer.In(r, value)
}

// ServiceResponse is a convenience function to create an http response
// encoded as JSON with a simple message
func ServiceResponse(producer Producer, w http.ResponseWriter, r *http.Request, httpStatus int, message string) {
	simpleMessage := &SimpleMessage{httpStatus, message}
	Produce(producer, w, r, httpStatus, simpleMessage)
}

// ----------------------------------------------------------------------------
// HANDLERS

// NotAllowed handler will response with a "405 Method Not Allowed" response
// It is a convenience handler to route all not allowed services
func NotAllowed(producer Producer, w http.ResponseWriter, r *http.Request) {
	ServiceResponse(producer, w, r, http.StatusMethodNotAllowed, "Method Not Allowed")
}

// NotFound handler will respond with a "404 Not Found" response
func NotFound(producer Producer, w http.ResponseWriter, r *http.Request) {
	ServiceResponse(producer, w, r, http.StatusNotFound, "Item not found")
}
