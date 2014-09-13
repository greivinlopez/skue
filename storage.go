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
	"errors"
	"fmt"
	"net/http"
)

// MemoryCacher represents an abstraction of any memory caching system used
// to speed up data driven systems by caching data in RAM instead of HD.
// Memory caching system examples:
// - Memcached: http://www.memcached.org
// - Redis: http://redis.io/
// More info here:
// Caching: http://en.wikipedia.org/wiki/Cache_(computing)
type MemoryCacher interface {
	Set(key interface{}, value interface{}) error
	Get(key interface{}, value interface{}) error
	Delete(key interface{}) error
}

// DatabasePersistor represents any abstraction that can follow the CRUD operations.
// Create, Read, Update and Delete are the four basic operations
// of persistent storage.
type DatabasePersistor interface {
	Create() (err error)
	Read(cache MemoryCacher) (err error)
	Update(cache MemoryCacher) (err error)
	Delete(cache MemoryCacher) (err error)
}

var ErrNotFound = errors.New("not found")

// ----------------------------------------------------------------------------
// PERSISTANCE UTILS:  Handles models CRUD and interaction with HTTP

// Saves a model to the underlying storage.
// Internally it calls the Create method of the given model.
// The model is constructed from the JSON body of the given request.
// Writes to the http writer according to what happens with the model
// following the REST architectural style.
func Create(model DatabasePersistor, w http.ResponseWriter, r *http.Request) {
	err := FromJson(r, &model)

	if err != nil {
		ServiceResponse(w, http.StatusBadRequest, fmt.Sprintf("Failed reading JSON from request: %v", err))
	} else {
		err = model.Create()
		if err != nil {
			ServiceResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed saving the item: %v", err))
		} else {
			ToJson(w, http.StatusCreated, model)
		}
	}
}

// Reads the model from underlying storage.
// Internally it calls the Read method of the given model which assumes
// it knows it's id.
// Writes to the http writer according to what happens with the model
// following the REST architectural style.
func Read(model DatabasePersistor, cache MemoryCacher, w http.ResponseWriter) {
	err := model.Read(cache)
	if err != nil {
		if err.Error() == "not found" {
			NotFound(w)
		} else {
			ServiceResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed updating the item: %v", err))
		}
	} else {
		ToJson(w, http.StatusOK, model)
	}
}

// Updates the given model in the underlying storage
// Internally it calls the Update method of the given model.
// The model is constructed from the JSON body of the given request.
// Writes to the http writer according to what happens with the model
// following the REST architectural style.
func Update(model DatabasePersistor, cache MemoryCacher, w http.ResponseWriter, r *http.Request) {
	err := FromJson(r, &model)

	if err != nil {
		ServiceResponse(w, http.StatusBadRequest, fmt.Sprintf("Failed reading JSON from request: %v", err))
	} else {
		err = model.Update(cache)
		if err != nil {
			if err.Error() == "not found" {
				NotFound(w)
			} else {
				ServiceResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed updating the item: %v", err))
			}
		} else {
			ServiceResponse(w, http.StatusOK, "Successfully updated")
		}
	}
}

// Deletes the model in the underlying storage.
// Internally it calls the Read method of the given model which assumes
// it knows it's id.
// If the model is created successfully then it calls the Delete method.
// Writes to the http writer according to what happens with the model
// following the REST architectural style.
func Delete(model DatabasePersistor, cache MemoryCacher, w http.ResponseWriter, r *http.Request) {
	err := model.Read(cache)
	if err != nil {
		if err.Error() == "not found" {
			NotFound(w)
		} else {
			ServiceResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed retrieving the item: %v", err))
		}
	} else {
		err = model.Delete(cache)
		if err != nil {
			ServiceResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed deleting the item: %v", err))
		} else {
			ServiceResponse(w, http.StatusOK, "Successfully deleted")
		}
	}
}
