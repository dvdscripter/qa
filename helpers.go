package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func jsonFromRequest(dst interface{}, r *http.Request) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if !json.Valid(body) {
		return errors.Errorf("Invalid body content")
	}

	if err := json.Unmarshal(body, &dst); err != nil {
		return errors.Wrap(err, "cannot unmarshal json body")
	}
	return nil
}

func idFromRequest(param string, r *http.Request) (int, error) {
	params := mux.Vars(r)
	rawID, exist := params[param]
	if !exist {
		return -1, errors.Errorf("Missing %s parameter", param)
	}

	id, err := strconv.Atoi(rawID)
	if err != nil {
		return -1, errors.Wrapf(err, "cannot convert %s to int", id)
	}
	return id, nil
}
