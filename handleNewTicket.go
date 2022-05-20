package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	timod "github.com/thingsdb/go-timod"
	"github.com/vmihailenco/msgpack"
)

type newTicketReq struct {
	Channel string      `msgpack:"channel"`
	Ticket  interface{} `msgpack:"ticket"`
}

type newTicketRes struct {
	Tickets []struct {
		Sid     string
		Channel string
	}
}

func handleNewTicket(pkg *timod.Pkg) {
	var req newTicketReq
	err := msgpack.Unmarshal(pkg.Data, &req)
	if err != nil {
		timod.WriteEx(
			pkg.Pid,
			timod.ExBadData,
			"Failed to unpack new-ticket request. Expecting a CHANNEL (string) and TICKET (thing)")
		return
	}

	jsonBody, err := json.Marshal(req.Ticket)
	if err != nil {
		timod.WriteEx(
			pkg.Pid,
			timod.ExBadData,
			fmt.Sprintf("Failed to JSON marshal ticket (%s)", err))
		return
	}
	body := bytes.NewReader(jsonBody)

	params := url.Values{}
	params.Set("channel", req.Channel)

	reqURL, err := url.Parse(cred.URI)
	if err != nil {
		timod.WriteEx(
			pkg.Pid,
			timod.ExBadData,
			fmt.Sprintf("Failed to parse URI (%s) (%s)", reqURL.String(), err))
		return
	}

	reqURL.Path = path.Join(reqURL.Path, "ticket")
	reqURL.RawQuery = params.Encode()

	httpReq, err := http.NewRequest("POST", reqURL.String(), body)
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
	var response newTicketRes
	err = json.Unmarshal(resBody, &response)
	if err != nil {
		timod.WriteEx(
			pkg.Pid,
			timod.ExBadData,
			fmt.Sprintf("Failed to unpack response (%s)", err))
		return
	}

	// Return the sid for the new ticket
	timod.WriteResponse(pkg.Pid, &response.Tickets[0].Sid)
}
