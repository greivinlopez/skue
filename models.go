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

// Modeler represents any abstraction that can follow the CRUD operations.
// Create, Read, Update and Delete are the four basic operations
// of persistent storage.
type Modeler interface {
	Create() (err error)
	Read() (err error)
	Update() (err error)
	Delete() (err error)
}

var ErrNotFound = errors.New("Not Found")

// ----------------------------------------------------------------------------
// PERSISTANCE UTILS:  Handles modelers CRUD and interaction with HTTP

// Saves a model to the underlying storage.
// Internally it calls the Create method of the given modeler.
// The modeler is constructed from the JSON body of the given request.
// Writes to the http writer according to what happens with the modeler
// following the REST architectural style.
func Create(modeler Modeler, w http.ResponseWriter, r *http.Request) {
	err := FromJson(r, &modeler)

	if err != nil {
		ServiceResponse(w, http.StatusBadRequest, fmt.Sprintf("Failed reading JSON from request: %v", err))
	} else {
		err = modeler.Create()
		if err != nil {
			ServiceResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed saving the item: %v", err))
		} else {
			ToJson(w, http.StatusCreated, modeler)
		}
	}
}

// Reads the model from underlying storage.
// Internally it calls the Read method of the given modeler which assumes
// it knows it's id.
// Writes to the http writer according to what happens with the modeler
// following the REST architectural style.
func Read(modeler Modeler, w http.ResponseWriter) {
	err := modeler.Read()
	if err != nil {
		if err == ErrNotFound {
			NotFound(w)
		} else {
			ServiceResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed updating the item: %v", err))
		}
	} else {
		ToJson(w, http.StatusOK, modeler)
	}
}

// Updates the given model in the underlying storage
// Internally it calls the Update method of the given modeler.
// The modeler is constructed from the JSON body of the given request.
// Writes to the http writer according to what happens with the modeler
// following the REST architectural style.
func Update(modeler Modeler, w http.ResponseWriter, r *http.Request) {
	err := FromJson(r, &modeler)

	if err != nil {
		ServiceResponse(w, http.StatusBadRequest, fmt.Sprintf("Failed reading JSON from request: %v", err))
	} else {
		err = modeler.Update()
		if err != nil {
			if err == ErrNotFound {
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
// Internally it calls the Read method of the given modeler which assumes
// it knows it's id.
// If the modeler is created successfully then it calls the Delete method.
// Writes to the http writer according to what happens with the modeler
// following the REST architectural style.
func Delete(modeler Modeler, w http.ResponseWriter, r *http.Request) {
	err := modeler.Read()
	if err != nil {
		if err == ErrNotFound {
			NotFound(w)
		} else {
			ServiceResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed retrieving the item: %v", err))
		}
	} else {
		err = modeler.Delete()
		if err != nil {
			ServiceResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed deleting the item: %v", err))
		} else {
			ServiceResponse(w, http.StatusOK, "Successfully deleted")
		}
	}
}
