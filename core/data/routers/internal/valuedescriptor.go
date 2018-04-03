/*******************************************************************************
 * Copyright 2017 Dell Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *
 * @microservice: core-data-go library
 * @author: Ryan Comer, Trevor Conn Dell
 * @version: 0.5.0
 *******************************************************************************/
package internal

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/edgexfoundry/edgex-go/core/domain/models"
	"github.com/gorilla/mux"
	"github.com/edgexfoundry/edgex-go/core/domain/errs"
	"github.com/edgexfoundry/edgex-go/core/aggregates/events"
)

const (
	maxExceededString string = "Error, exceeded the max limit as defined in config"
)


// GET, POST, and PUT for value descriptors
// api/v1/valuedescriptor
func ValueDescriptorHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	switch r.Method {
	case http.MethodGet:
		vList, err := events.GetAllValueDescriptors()
		if err != nil {
			getLogger().Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Check the limit
		if len(vList) > getConfiguration().ReadMaxLimit {
			http.Error(w, maxExceededString, http.StatusRequestEntityTooLarge)
			getLogger().Error(maxExceededString)
			return
		}

		encode(vList, w)
	case http.MethodPost:
		dec := json.NewDecoder(r.Body)
		v := models.ValueDescriptor{}
		err := dec.Decode(&v)
		// Problems decoding
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			getLogger().Error("Error decoding the value descriptor: " + err.Error())
			return
		}

		// Check the formatting
		match, err := events.ValidateFormatString(v)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			getLogger().Error("Error checking for format string for POSTed value descriptor")
			return
		}
		if !match {
			err := errors.New("error posting value descriptor. Format is not a valid printf format")
			http.Error(w, err.Error(), http.StatusConflict)
			getLogger().Error(err.Error())
			return
		}

		id, err := events.AddValueDescriptor(v)
		if err != nil {
			if err == errs.ErrNotUnique {
				http.Error(w, "Value Descriptor already exists", http.StatusConflict)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(id))
	case http.MethodPut:
		dec := json.NewDecoder(r.Body)
		from := models.ValueDescriptor{}
		err := dec.Decode(&from)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			getLogger().Error("Error decoding the value descriptor: " + err.Error())
			return
		}

		err = events.UpdateValueDescriptor(from)
		if err != nil {
			if err == errs.ErrNotUnique {
				http.Error(w, "Value descriptor name is not unique", http.StatusConflict)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			getLogger().Error(err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("true"))
	}
}

// Delete the value descriptor based on the ID
// DataValidationException (HTTP 409) - The value descriptor is still referenced by readings
// NotFoundException (404) - Can't find the value descriptor
// valuedescriptor/id/{id}
func DeleteValueDescriptorByIdHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	id := vars["id"]

	err := events.DeleteValueDescriptorById(id)
	if err != nil {
		getLogger().Error(err.Error())
		if err == errs.ErrNotFound { //One could make the case that this should not throw an error
			http.Error(w, "Value descriptor not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("true"))
}

// Value descriptors based on name
// api/v1/valuedescriptor/name/{name}
func ValueDescriptorByNameHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	name, err := url.QueryUnescape(vars["name"])

	// Problems unescaping
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		getLogger().Error("Error unescaping the value descriptor name: " + err.Error())
		return
	}

	switch r.Method {
	case http.MethodGet:
		v, err := events.GetValueDescriptorByName(name)
		if err != nil {
			if err == errs.ErrNotFound {
				http.Error(w, "Value Descriptor not found", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			getLogger().Error(err.Error())
			return
		}

		encode(v, w)
	case http.MethodDelete:
		// Check if the value descriptor exists
		vd, err := events.GetValueDescriptorByName(name)
		if err != nil {
			if err == errs.ErrNotFound { //One could make the argument that this should not throw an error
				http.Error(w, "Value Descriptor not found", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			getLogger().Error(err.Error())
			return
		}

		if err = events.DeleteValueDescriptorById(vd.Id.Hex()); err != nil {
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("true"))
	}
}

// Get a value descriptor based on the ID
// HTTP 404 not found if the ID isn't in the database
// api/v1/valuedescriptor/{id}
func ValueDescriptorByIdHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	id := vars["id"]

	switch r.Method {
	case http.MethodGet:
		v, err := events.GetValueDescriptorById(id)
		if err != nil {
			if err == errs.ErrNotFound {
				http.Error(w, "Value descriptor not found", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusServiceUnavailable)
			}
			getLogger().Error(err.Error())
			return
		}

		encode(v, w)
	}
}

// Get the value descriptor from the UOM label
// api/v1/valuedescriptor/uomlabel/{uomLabel}
func ValueDescriptorByUomLabelHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	uomLabel, err := url.QueryUnescape(vars["uomLabel"])

	// Problem unescaping
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		getLogger().Error("Error unescaping the UOM Label of the value descriptor: " + err.Error())
		return
	}

	switch r.Method {
	case http.MethodGet:
		v, err := events.GetValueDescriptorsByUomLabel(uomLabel)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			getLogger().Error(err.Error())
			return
		}

		encode(v, w)
	}
}

// Get value descriptors who have one of the labels
// api/v1/valuedescriptor/label/{label}
func ValueDescriptorByLabelHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	label, err := url.QueryUnescape(vars["label"])

	// Problem unescaping
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		getLogger().Error("Error unescaping label for the value descriptor: " + err.Error())
		return
	}

	switch r.Method {
	case http.MethodGet:
		v, err := events.GetValueDescriptorsByLabel(label)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			getLogger().Error(err.Error())
			return
		}

		encode(v, w)
	}
}

// Return the value descriptors that are asociated with a device
// The value descriptor is expected parameters on puts or expected values on get/put commands
// api/v1/valuedescriptor/devicename/{device}
func ValueDescriptorByDeviceHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)

	device, err := url.QueryUnescape(vars["device"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		getLogger().Error("Error unescaping the device: " + err.Error())
		return
	}

	vdList, err := events.GetValueDescriptorsByDeviceName(device)
	if err != nil {
		if err == errs.ErrNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		getLogger().Error(err.Error())
		return
	}

	encode(vdList, w)
}

// Return the value descriptors that are associated with the device specified by the device ID
// Associated value descripts are expected parameters of PUT commands and expected results of PUT/GET commands
// api/v1/valuedescriptor/deviceid/{id}
func ValueDescriptorByDeviceIdHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)

	deviceId, err := url.QueryUnescape(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		getLogger().Error("Error unescaping the device ID: " + err.Error())
		return
	}

	vdList, err := events.GetValueDescriptorsByDeviceId(deviceId)
	if err != nil {
		if err == errs.ErrNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		getLogger().Error(err.Error())
		return
	}

	encode(vdList, w)
}
