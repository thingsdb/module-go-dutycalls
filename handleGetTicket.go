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

type getTicketReq struct {
	Sid *string `msgpack:"sid"`
}

type getTicketRes struct {
	Tickets []interface{}
}

func handleGetTicket(pkg *timod.Pkg) {
	var req getTicketReq
	err := msgpack.Unmarshal(pkg.Data, &req)
	if err != nil {
		timod.WriteEx(
			pkg.Pid,
			timod.ExBadData,
			"Error: Failed to unpack Get Ticket request")
		return
	}

	params := url.Values{}
	params.Set("sid", *req.Sid)

	reqURL, err := url.Parse(cred.URI)
	reqURL.Path = path.Join(reqURL.Path, "ticket")

	if err != nil {
		timod.WriteEx(
			pkg.Pid,
			timod.ExBadData,
			fmt.Sprintf("Error: Failed to parse URI (%s) (%s)", reqURL.String(), err))
		return
	}

	reqURL.RawQuery = params.Encode()

	httpReq, err := http.NewRequest("GET", reqURL.String(), http.NoBody)
	if err != nil {
		timod.WriteEx(
			pkg.Pid,
			timod.ExBadData,
			fmt.Sprintf("Error: Failed to create request (%s)", err))
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
			fmt.Sprintf("Error: Failed to perform the request (%s)", err))
		return
	}

	// Read the body
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		timod.WriteEx(
			pkg.Pid,
			timod.ExBadData,
			fmt.Sprintf("Error: Failed to read bytes from response (%s)", err))
		return
	}

	// Check status code and return error if not 200
	if res.StatusCode != 200 {
		timod.WriteEx(
			pkg.Pid,
			timod.ExBadData,
			fmt.Sprintf("Error: %s (%d)", string(resBody), res.StatusCode))
		return
	}

	// Unpack the newTicket response
	var response getTicketRes
	err = json.Unmarshal(resBody, &response)
	if err != nil {
		timod.WriteEx(
			pkg.Pid,
			timod.ExBadData,
			fmt.Sprintf("Error: Failed to unpack response (%s)", err))
		return
	}

	// Return the sid for the new ticket
	timod.WriteResponse(pkg.Pid, &response.Tickets)
}
