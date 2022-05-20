package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	timod "github.com/thingsdb/go-timod"
	"github.com/vmihailenco/msgpack"
)

type getHitsReq struct {
	Sid string `msgpack:"sid"`
}

type getHitsRes struct {
	Hits []interface{}
}

func handleGetHits(pkg *timod.Pkg) {
	var req getHitsReq
	err := msgpack.Unmarshal(pkg.Data, &req)
	if err != nil {
		timod.WriteEx(
			pkg.Pid,
			timod.ExBadData,
			"Failed to unpack get-hits request. Expecting a SID (string)")
		return
	}

	params := url.Values{}
	params.Set("sid", req.Sid)

	reqURL, err := url.Parse(cred.URI)
	if err != nil {
		timod.WriteEx(
			pkg.Pid,
			timod.ExBadData,
			fmt.Sprintf("Failed to parse URI (%s) (%s)", reqURL.String(), err))
		return
	}

	reqURL.Path = path.Join(reqURL.Path, "ticket/hit")
	reqURL.RawQuery = params.Encode()

	httpReq, err := http.NewRequest("GET", reqURL.String(), http.NoBody)
	if err != nil {
		timod.WriteEx(
			pkg.Pid,
			timod.ExBadData,
			fmt.Sprintf("Failed to create request (%s)", err))
		return
	}

	httpReq.SetBasicAuth(cred.Login, cred.Password)
	httpReq.Header.Set("Content-Type", "application/json")

	// Do the actual API request
	client := &http.Client{}
	res, err := client.Do(httpReq)
	if err != nil {
		timod.WriteEx(
			pkg.Pid,
			timod.ExOperation,
			fmt.Sprintf("Failed to perform the request (%s)", err))
		return
	}

	// Read the body
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		timod.WriteEx(
			pkg.Pid,
			timod.ExBadData,
			fmt.Sprintf("Failed to read bytes from response (%s)", err))
		return
	}

	// Check status code and return error if not 2XX
	if res.StatusCode/100 != 2 {
		timod.WriteEx(
			pkg.Pid,
			timod.ExBadData,
			fmt.Sprintf("%s (%d)", string(resBody), res.StatusCode))
		return
	}

	// Unpack the newTicket response
	var response getHitsRes
	err = json.Unmarshal(resBody, &response)
	if err != nil {
		timod.WriteEx(
			pkg.Pid,
			timod.ExBadData,
			fmt.Sprintf("Failed to unpack response (%s)", err))
		return
	}

	// Return the sid for the new ticket
	timod.WriteResponse(pkg.Pid, &response.Hits)
}
