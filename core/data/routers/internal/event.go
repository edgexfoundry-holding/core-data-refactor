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
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/edgexfoundry/edgex-go/core/aggregates/events"
	"github.com/edgexfoundry/edgex-go/core/data/config"
	"github.com/edgexfoundry/edgex-go/core/data/log"
	"github.com/edgexfoundry/edgex-go/core/domain/errs"
	"github.com/edgexfoundry/edgex-go/core/domain/models"
	"github.com/edgexfoundry/edgex-go/support/logging-client"
	"github.com/gorilla/mux"
)

func getLogger() logger.LoggingClient {
	return log.Logger
}

func getConfiguration() *config.ConfigurationStruct {
	return config.Configuration
}

// Undocumented feature to remove all readings and events from the database
// This should primarily be used for debugging purposes
func ScrubAllHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	switch r.Method {
	case http.MethodDelete:
		getLogger().Info("Deleting all events from database")

		err := events.Purge()
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		encode(true, w)
	}
}

/*
Handler for the event API
Status code 404 - event not found
Status code 413 - number of events exceeds limit
Status code 503 - unanticipated issues
api/v1/event
*/
func EventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Body != nil {
		defer r.Body.Close()
	}

	switch r.Method {
	// Get all events
	case http.MethodGet:
		events, err := events.GetAllEvents()
		if err != nil {
			getLogger().Error(err.Error())
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		// Check max limit
		if len(events) > getConfiguration().ReadMaxLimit {
			getLogger().Error(maxExceededString)
			http.Error(w, maxExceededString, http.StatusRequestEntityTooLarge)
			return
		}

		encode(events, w)
		break
	// Post a new event
	case http.MethodPost:
		var e models.Event
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&e)

		// Problem Decoding Event
		if err != nil {
			getLogger().Error("Error decoding event: " + err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		getLogger().Info("Posting Event: " + e.String())

		id, err := events.AddNewEvent(e)
		if err != nil {
			getLogger().Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
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

		break
	// Do not update the readings
	case http.MethodPut:
		var from models.Event
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&from)

		// Problem decoding event
		if err != nil {
			getLogger().Error("Error decoding the event: " + err.Error())
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		err = events.UpdateEvent(from)

		if err != nil {
			if err == errs.ErrNotFound {
				http.Error(w, fmt.Errorf("event id %s: %v", from.ID.Hex(), err.Error()).Error(), http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusServiceUnavailable)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("true"))
	}
}

//GET
//Return the event specified by the event ID
///api/v1/event/{id}
//id - ID of the event to return
func GetEventByIdHandler(w http.ResponseWriter, r *http.Request) {
	if r.Body != nil {
		defer r.Body.Close()
	}

	switch r.Method {
	case http.MethodGet:
		// URL parameters
		vars := mux.Vars(r)
		id := vars["id"]

		// Get the event
		e, err := events.GetEventById(id)
		if err != nil {
			if err == errs.ErrNotFound {
				http.Error(w, fmt.Errorf("event id %s: %v", id, err.Error()).Error(), http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusServiceUnavailable)
			}
			return
		}

		// Return the result
		encode(e, w)
	}
}

/*
Return number of events in Core Data
/api/v1/event/count
*/
func EventCountHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	switch r.Method {
	case http.MethodGet:
		count, err := events.CountEvents()
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		// Return result
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(strconv.Itoa(count)))
		if err != nil {
			getLogger().Error(err.Error(), "")
		}
	}
}

/*
Return number of events for a given device in Core Data
deviceID - ID of the device to get count for
/api/v1/event/count/{deviceId}
*/
func EventCountByDeviceIdHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	id, err := url.QueryUnescape(vars["deviceId"])
	if err != nil {
		getLogger().Error("Problem unescaping URL: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case http.MethodGet:

		count, err := events.CountByDevice(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Return result
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(strconv.Itoa(count)))
		break
	}
}

/*
DELETE, PUT
Handle events specified by an ID
/api/v1/event/id/{id}
404 - ID not found
*/
func EventIdHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	id := vars["id"]

	switch r.Method {
	// Set the 'pushed' timestamp for the event to the current time - event is going to another (not fuse) service
	case http.MethodPut:
		err := events.Touch(id)
		if err != nil {
			if err == errs.ErrNotFound {
				http.Error(w, fmt.Errorf("event id %s: %v", id, err.Error()).Error(), http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("true"))
		break
	// Delete the event and all of it's readings
	case http.MethodDelete:
		// Check if the event exists
		err := events.DeleteEventById(id)
		if err != nil {
			if err == errs.ErrNotFound { //One could make the argument that this shouldn't throw an error
				http.Error(w, fmt.Errorf("event id %s: %v", id, err.Error()).Error(), http.StatusNotFound)
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

// Get event by device id
// Returns the events for the given device sorted by creation date and limited by 'limit'
// {deviceId} - the device that the events are for
// {limit} - the limit of events
// api/v1/event/device/{deviceId}/{limit}
func GetEventByDeviceHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	limit := vars["limit"]
	deviceId, err := url.QueryUnescape(vars["deviceId"])

	// Problems unescaping URL
	if err != nil {
		getLogger().Error("Error unescaping URL: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// Convert limit to int
		limitNum, err := strconv.Atoi(limit)
		if err != nil {
			getLogger().Error("Error converting to integer: " + err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		eventList, err := events.GetEventsByDevice(deviceId, limitNum)
		if err != nil {
			getLogger().Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		encode(eventList, w)
	}
}

// Delete all of the events associated with a device
// api/v1/event/device/{deviceId}
// 404 - device ID not found in metadata
// 503 - service unavailable
func DeleteByDeviceIdHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	deviceId, err := url.QueryUnescape(vars["deviceId"])
	// Problems unescaping URL
	if err != nil {
		getLogger().Error("Error unescaping the URL: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case http.MethodDelete:
		count, err := events.DeleteByDevice(deviceId)
		if err != nil {
			getLogger().Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(strconv.Itoa(count)))
	}
}

// Get events by creation time
// {start} - start time, {end} - end time, {limit} - max number of results
// Sort the events by creation date
// 413 - number of results exceeds limit
// 503 - service unavailable
// api/v1/event/{start}/{end}/{limit}
func EventByCreationTimeHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	start, err := strconv.ParseInt(vars["start"], 10, 64)
	// Problems converting start time
	if err != nil {
		getLogger().Error("Problem converting start time: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	end, err := strconv.ParseInt(vars["end"], 10, 64)
	// Problems converting end time
	if err != nil {
		getLogger().Error("Problem converting end time: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	limit, err := strconv.Atoi(vars["limit"])
	// Problems converting limit
	if err != nil {
		getLogger().Error("Problem converting limit: " + strconv.Itoa(limit))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case http.MethodGet:
		e, err := events.GetEventsByCreateTime(start, end, limit)
		if err != nil {
			getLogger().Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		encode(e, w)
	}
}

// Get the readings for a device and filter them based on the value descriptor
// Only those readings whos name is the value descriptor should get through
// /event/device/{deviceId}/valuedescriptor/{valueDescriptor}/{limit}
// 413 - number exceeds limit
func ReadingByDeviceFilteredValueDescriptor(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	limit := vars["limit"]

	valueDescriptor, err := url.QueryUnescape(vars["valueDescriptor"])
	// Problems unescaping URL
	if err != nil {
		getLogger().Error("Problem unescaping value descriptor: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	deviceId, err := url.QueryUnescape(vars["deviceId"])
	// Problems unescaping URL
	if err != nil {
		getLogger().Error("Problem unescaping device ID: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	limitNum, err := strconv.Atoi(limit)
	// Problem converting the limit
	if err != nil {
		getLogger().Error("Problem converting limit to integer: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	switch r.Method {
	case http.MethodGet:
		readings, err := events.GetReadingsByDeviceAndValueDescriptor(deviceId, valueDescriptor, limitNum)
		if err != nil {
			getLogger().Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		encode(readings, w)
	}
}

// Remove all the old events and associated readings (by age)
// event/removeold/age/{age}
func EventByAgeHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	age, err := strconv.ParseInt(vars["age"], 10, 64)

	// Problem converting age
	if err != nil {
		getLogger().Error("Error converting the age to an integer")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case http.MethodDelete:
		getLogger().Info("Deleting events by age: " + vars["age"])

		count, err := events.DeleteByAge(age)
		if err != nil {
			getLogger().Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Return the count
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(strconv.Itoa(count)))
	}
}

// Scrub all the events that have been pushed
// Also remove the readings associated with the events
// api/v1/event/scrub
func ScrubHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	switch r.Method {
	case http.MethodDelete:
		getLogger().Info("Scrubbing events.  Deleting all events that have been pushed")

		count, err := events.PurgeIfPublished()
		if err != nil {
			getLogger().Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(strconv.Itoa(count)))
	}
}
