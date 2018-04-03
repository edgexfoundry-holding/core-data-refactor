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
 * @author: Ryan Comer, Dell
 * @version: 0.5.0
 *******************************************************************************/
package internal

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/edgexfoundry/edgex-go/core/aggregates/events"
	"github.com/edgexfoundry/edgex-go/core/domain/models"
	"github.com/gorilla/mux"
	"github.com/edgexfoundry/edgex-go/core/domain/errs"
)


// Reading handler
// GET, PUT, and POST readings
func ReadingHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	switch r.Method {
	case http.MethodGet:
		r, err := events.GetAllReadings()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Check max limit
		if len(r) > getConfiguration().ReadMaxLimit {
			http.Error(w, maxExceededString, http.StatusRequestEntityTooLarge)
			getLogger().Error(maxExceededString)
			return
		}

		encode(r, w)
	case http.MethodPost:
		reading := models.Reading{}
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&reading)

		// Problem decoding
		if err != nil {
			getLogger().Error("Error decoding the reading: " + err.Error())
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		id, err := events.AddNewReading(reading)
		if err != nil {
			getLogger().Error(err.Error())
			if err == errs.ErrNotFound {
				http.Error(w, "Value descriptor not found for reading", http.StatusConflict)
			} else {
				http.Error(w, err.Error(), http.StatusServiceUnavailable)
			}
			return
		}

		//This is kind of whack, held over from previous logic for now
		if id != "unsaved" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(id))
		} else {
			err = encode(id, w)
			if err != nil {
				getLogger().Error(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	case http.MethodPut:
		from := models.Reading{}
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&from)

		// Problem decoding
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			getLogger().Error("Error decoding the reading: " + err.Error())
			return
		}

		err = events.UpdateReading(from)
		if err != nil {
			if err == errs.ErrNotFound {
				http.Error(w, "Reading not found", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("true"))
	}
}

// Get a reading by id
// HTTP 404 not found if the reading can't be found by the ID
// api/v1/reading/{id}
func GetReadingByIdHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	id := vars["id"]

	switch r.Method {
	case http.MethodGet:
		reading, err := events.GetReadingById(id)
		if err != nil {
			if err == errs.ErrNotFound {
				http.Error(w, "Reading not found", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusServiceUnavailable)
			}
			return
		}

		encode(reading, w)
	}
}

// Return a count for the number of readings in core data
// api/v1/reading/count
func ReadingCountHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	switch r.Method {
	case http.MethodGet:
		count, err := events.CountReadings()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(strconv.Itoa(count)))
		if err != nil {
			getLogger().Error(err.Error(), "")
		}
	}
}

// Delete a reading by its id
// api/v1/reading/id/{id}
func DeleteReadingByIdHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	id := vars["id"]

	switch r.Method {
	case http.MethodDelete:
		// Check if the reading exists
		err := events.DeleteReadingById(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("true"))
	}
}

// Get all the readings for the device - sort by creation date
// 404 - device ID or name doesn't match
// 413 - max count exceeded
// api/v1/reading/device/{deviceId}/{limit}
func ReadingByDeviceHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	limit, err := strconv.Atoi(vars["limit"])
	// Problems converting limit to int
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		getLogger().Error("Error converting the limit to an integer: " + err.Error())
		return
	}
	deviceId, err := url.QueryUnescape(vars["deviceId"])
	// Problems unescaping URL
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		getLogger().Error("Error unescaping the device ID: " + err.Error())
		return
	}

	switch r.Method {
	case http.MethodGet:
		readings, err := events.GetReadingsByDevice(deviceId, limit)
		if err != nil {
			if err == errs.ErrNotFound {
				http.Error(w, err.Error(), http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		encode(readings, w)
	}
}

// Return a list of all readings associated with a value descriptor, limited by limit
// HTTP 413 (limit exceeded) if the limit is greater than max limit
// api/v1/reading/name/{name}/{limit}
func ReadingByValueDescriptorHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	name, err := url.QueryUnescape(vars["name"])
	// Problems with unescaping URL
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		getLogger().Error("Error unescaping value descriptor name: " + err.Error())
		return
	}
	limit, err := strconv.Atoi(vars["limit"])
	// Problems converting limit to int
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		getLogger().Error("Error converting the limit to an integer: " + err.Error())
		return
	}

	// Check for value descriptor
	readings, err := events.GetReadingsByValueDescriptor(name, limit)
	if err != nil {
		if err == errs.ErrNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	encode(readings, w)
}

// Return a list of readings based on the UOM label for the value decriptor
// api/v1/reading/uomlabel/{uomLabel}/{limit}
func ReadingByUomLabelHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)

	uomLabel, err := url.QueryUnescape(vars["uomLabel"])
	// Problems unescaping URL
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		getLogger().Error("Error unescaping the UOM Label: " + err.Error())
		return
	}

	limit, err := strconv.Atoi(vars["limit"])
	// Problems converting limit to int
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		getLogger().Error("Error converting the limit to an integer: " + err.Error())
		return
	}

	readings, err := events.GetReadingsByUomLabel(uomLabel, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encode(readings, w)
}

// Get readings by the value descriptor (specified by the label)
// 413 - limit exceeded
// api/v1/reading/label/{label}/{limit}
func ReadingByLabelHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	label, err := url.QueryUnescape(vars["label"])
	// Problem unescaping
	if err != nil {
		getLogger().Error("Error unescaping the label of the value descriptor: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	limit, err := strconv.Atoi(vars["limit"])
	// Problems converting to int
	if err != nil {
		getLogger().Error("Error converting the limit to an integer: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	readings, err := events.GetReadingsByValueDescriptorLabel(label, limit)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encode(readings, w)
}

// Return a list of readings who's value descriptor has the type
// 413 - number exceeds the current limit
// /reading/type/{type}/{limit}
func ReadingByTypeHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)

	t, err := url.QueryUnescape(vars["type"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		getLogger().Error("Error escaping the type: " + err.Error())
		return
	}

	l, err := strconv.Atoi(vars["limit"])
	// Problem converting to int
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		getLogger().Error("Error converting the limit to an integer: " + err.Error())
		return
	}

	readings, err := events.GetReadingsByType(t, l)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encode(readings, w)
}

// Return a list of readings between the start and end (creation time)
// /reading/{start}/{end}/{limit}
func ReadingByCreationTimeHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	s, err := strconv.ParseInt((vars["start"]), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		getLogger().Error("Error converting the start time to an integer: " + err.Error())
		return
	}
	e, err := strconv.ParseInt((vars["end"]), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		getLogger().Error("Error converting the end time to an integer: " + err.Error())
		return
	}
	l, err := strconv.Atoi(vars["limit"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		getLogger().Error("Error converting the limit to an integer: " + err.Error())
		return
	}

	switch r.Method {
	case http.MethodGet:

		readings, err := events.GetReadingsByCreateTime(s, e, l)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		encode(readings, w)
	}
}

// Return a list of redings associated with the device and value descriptor
// Limit exceeded exception 413 if the limit exceeds the max limit
// api/v1/reading/name/{name}/device/{device}/{limit}
func ReadingByValueDescriptorAndDeviceHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)

	// Get the variables from the URL
	name, err := url.QueryUnescape(vars["name"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		getLogger().Error("Error unescaping the value descriptor name: " + err.Error())
		return
	}

	device, err := url.QueryUnescape(vars["device"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		getLogger().Error("Error unescaping the device: " + err.Error())
		return
	}

	limit, err := strconv.Atoi(vars["limit"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		getLogger().Error("Error converting limit to an integer: " + err.Error())
		return
	}

	switch r.Method {
	case http.MethodGet:
		readings, err := events.GetReadingsByDeviceAndValueDescriptor(device, name, limit)
		if err != nil {
			getLogger().Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		encode(readings, w)
	}
}
