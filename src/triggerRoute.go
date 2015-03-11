package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	//"log"
)

func getTriggerQueryParam(query url.Values) (*trigger, int64, string) {
	var (
		lat        string = ""
		lon        string = ""
		unit       string = "km"
		radius     int64  = 5
		appId      string = ""
		identifier string = ""
		err        error
	)

	if query.Get("lat") != "" {
		_, err = strconv.ParseFloat(query.Get("lat"), 64)
		if err == nil {
			lat = query.Get("lat")
		}
	}

	if query.Get("lon") != "" {
		_, err = strconv.ParseFloat(query.Get("lon"), 64)
		if err == nil {
			lon = query.Get("lon")
		}
	}

	if query.Get("radius") != "" {
		radius, err = strconv.ParseInt(query.Get("radius"), 10, 64)
		if err != nil {
			// Default query to 50 meters
			radius = 50
		}
	}

	unit = query.Get("unit")
	appId = query.Get("appId")
	identifier = query.Get("identifier")

	return &trigger{
		AppId:      appId,
		Identifier: identifier,
		Coords:     []string{lat, lon},
	}, radius, unit
}

func triggerGET(w http.ResponseWriter, r *http.Request, t *Trigger) error {
	triggerInfo, radius, unit := getTriggerQueryParam(r.URL.Query())

	coords := (*triggerInfo).Coords

	// lat, lon
	if coords[0] == "" || coords[1] == "" {
		return &InvalidRequest{
			"Required parameter 'lat' and 'lon' must be provided.",
		}
	}

	triggers, err := t.findNearBy(triggerInfo, radius, unit)

	if err != nil {
		return err
	}

	jsonRes, err := json.Marshal(triggers)

	if err != nil {
		return err
	}

	w.Write(jsonRes)
	return nil
}

func triggerDELETE(w http.ResponseWriter, r *http.Request, t *Trigger) error {
	qParams := r.URL.Query()
	appId := qParams["appId"][0]
	identifier := qParams["identifier"][0]

	if appId == "" || identifier == "" {
		return &InvalidRequest{
			"Required parameter 'appId' and 'identifier' must be provided.",
		}
	}

	if err := t.remove(appId, identifier); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte(""))
	return nil
}

func triggerPOST(w http.ResponseWriter, r *http.Request, t *Trigger) error {
	data := &[]trigger{}

	var (
		body []byte
		err  error
	)

	if body, err = ioutil.ReadAll(r.Body); err != nil {
		return err
	}

	if err = json.Unmarshal([]byte(body), data); err != nil {
		return &InvalidRequest{
			"Bad JSON - " + err.Error(),
		}
	}

	if err = t.add(data); err != nil {
		return err
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(body)
	return nil
}

func triggerPUT(w http.ResponseWriter, r *http.Request, t *Trigger) error {
	triggerInfo, _, _ := getTriggerQueryParam(r.URL.Query())
	data := &trigger{}

	var (
		body []byte
		err  error
	)

	if body, err = ioutil.ReadAll(r.Body); err != nil {
		return err
	}

	if err = json.Unmarshal([]byte(body), data); err != nil {
		return &InvalidRequest{
			"Bad JSON - " + err.Error(),
		}
	}

	err = t.updateIdentifier(triggerInfo.AppId, triggerInfo.Identifier, data.Identifier)

	if err != nil {
		return err
	}

	return nil
}

func makeTriggerHandlers(db *sql.DB) handler {
	triggerHandlers := map[string]handlerWithDB_and_Error{
		"DELETE": triggerDELETE,
		"GET":    triggerGET,
		"POST":   triggerPOST,
		"PUT":    triggerPUT,
	}

	return requestHandlerWithDB(db, triggerHandlers)
}
