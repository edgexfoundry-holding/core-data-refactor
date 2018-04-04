/*******************************************************************************
 * Copyright 2018 Dell Inc.
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
 * @author: Trevor Conn, Dell
 * @version: 0.5.0
 *******************************************************************************/

package routers

import (
	"net/http"

	"github.com/edgexfoundry/edgex-go/core/data/routers/internal"
	"github.com/gorilla/mux"
	)

type gorillaRouter struct {
	router *mux.Router
}

func (g *gorillaRouter) LoadRoutes() http.Handler {
	b := g.router.PathPrefix("/api/v1").Subrouter()

	// EVENTS
	// /api/v1/event
	b.HandleFunc("/event", internal.EventHandler).Methods(http.MethodGet, http.MethodPut, http.MethodPost)
	e := b.PathPrefix("/event").Subrouter()
	e.HandleFunc("/scrub", internal.ScrubHandler).Methods(http.MethodDelete)
	e.HandleFunc("/scruball", internal.ScrubAllHandler).Methods(http.MethodDelete)
	e.HandleFunc("/count", internal.EventCountHandler).Methods(http.MethodGet)
	e.HandleFunc("/count/{deviceId}", internal.EventCountByDeviceIdHandler).Methods(http.MethodGet)
	e.HandleFunc("/{id}", internal.GetEventByIdHandler).Methods(http.MethodGet)
	e.HandleFunc("/id/{id}", internal.EventIdHandler).Methods(http.MethodDelete, http.MethodPut)
	e.HandleFunc("/device/{deviceId}/{limit:[0-9]+}", internal.GetEventByDeviceHandler).Methods(http.MethodGet)
	e.HandleFunc("/device/{deviceId}", internal.DeleteByDeviceIdHandler).Methods(http.MethodDelete)
	e.HandleFunc("/removeold/age/{age:[0-9]+}", internal.EventByAgeHandler).Methods(http.MethodDelete)
	e.HandleFunc("/{start:[0-9]+}/{end:[0-9]+}/{limit:[0-9]+}", internal.EventByCreationTimeHandler).Methods(http.MethodGet)
	e.HandleFunc("/device/{deviceId}/valuedescriptor/{valueDescriptor}/{limit:[0-9]+}", internal.ReadingByDeviceFilteredValueDescriptor).Methods(http.MethodGet)

	// READINGS
	// /api/v1/reading
	b.HandleFunc("/reading", internal.ReadingHandler).Methods(http.MethodGet, http.MethodPut, http.MethodPost)
	rd := b.PathPrefix("/reading").Subrouter()
	rd.HandleFunc("/count", internal.ReadingCountHandler).Methods(http.MethodGet)
	rd.HandleFunc("/id/{id}", internal.DeleteReadingByIdHandler).Methods(http.MethodDelete)
	rd.HandleFunc("/{id}", internal.GetReadingByIdHandler).Methods(http.MethodGet)
	rd.HandleFunc("/device/{deviceId}/{limit:[0-9]+}", internal.ReadingByDeviceHandler).Methods(http.MethodGet)
	rd.HandleFunc("/name/{name}/{limit:[0-9]+}", internal.ReadingByValueDescriptorHandler).Methods(http.MethodGet)
	rd.HandleFunc("/uomlabel/{uomLabel}/{limit:[0-9]+}", internal.ReadingByUomLabelHandler).Methods(http.MethodGet)
	rd.HandleFunc("/label/{label}/{limit:[0-9]+}", internal.ReadingByLabelHandler).Methods(http.MethodGet)
	rd.HandleFunc("/type/{type}/{limit:[0-9]+}", internal.ReadingByTypeHandler).Methods(http.MethodGet)
	rd.HandleFunc("/{start:[0-9]+}/{end:[0-9]+}/{limit:[0-9]+}", internal.ReadingByCreationTimeHandler).Methods(http.MethodGet)
	rd.HandleFunc("/name/{name}/device/{device}/{limit:[0-9]+}", internal.ReadingByValueDescriptorAndDeviceHandler).Methods(http.MethodGet)

	// VALUE DESCRIPTORS
	// /api/v1/valuedescriptor
	b.HandleFunc("/valuedescriptor", internal.ValueDescriptorHandler).Methods(http.MethodGet, http.MethodPut, http.MethodPost)
	vd := b.PathPrefix("/valuedescriptor").Subrouter()
	vd.HandleFunc("/id/{id}", internal.DeleteValueDescriptorByIdHandler).Methods(http.MethodDelete)
	vd.HandleFunc("/name/{name}", internal.ValueDescriptorByNameHandler).Methods(http.MethodGet, http.MethodDelete)
	vd.HandleFunc("/{id}", internal.ValueDescriptorByIdHandler).Methods(http.MethodGet)
	vd.HandleFunc("/uomlabel/{uomLabel}", internal.ValueDescriptorByUomLabelHandler).Methods(http.MethodGet)
	vd.HandleFunc("/label/{label}", internal.ValueDescriptorByLabelHandler).Methods(http.MethodGet)
	vd.HandleFunc("/devicename/{device}", internal.ValueDescriptorByDeviceHandler).Methods(http.MethodGet)
	vd.HandleFunc("/deviceid/{id}", internal.ValueDescriptorByDeviceIdHandler).Methods(http.MethodGet)

	// Ping Resource
	// /api/v1/ping
	b.HandleFunc("/ping", internal.PingHandler)

	return g.router
}