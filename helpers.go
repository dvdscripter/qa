package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func jsonFromRequest(dst interface{}, r *http.Request) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &dst); err != nil {
		return err
	}
	return nil
}

func idFromRequest(param string, r *http.Request) (int, error) {
	params := mux.Vars(r)
	rawID, exist := params[param]
	if !exist {
		return -1, fmt.Errorf("Missing %s parameter", param)
	}

	id, err := strconv.Atoi(rawID)
	if err != nil {
		return -1, err
	}
	return id, nil
}
