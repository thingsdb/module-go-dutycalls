// ThingsDB module for communication with DutyCalls.
//
// For example:
//
//     // Create the module (@thingsdb scope)
//     new_module('dutycalls', 'github.com/thingsdb/module-go-dutycalls');
//
//     // Configure the module
//     set_module_conf('dutycalls', {
//         login: 'mylogin',
//         password: 'mysecret',
//     });
//
//     // Use the module
//     dutycalls.new_ticket("mychannel", {
//         title: "my title"
//     }).then(|sid| {
//         sid;  // the sid of the new ticket
//     });
//
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"sync"

	timod "github.com/thingsdb/go-timod"

	"github.com/vmihailenco/msgpack"
)

var mux sync.Mutex
var cred authDutyCalls

const DefaultURI = "https://dutycalls.me/api"

type authDutyCalls struct {
	Login    string `msgpack:"login"`
	Password string `msgpack:"password"`
	URI      string `msgpack:"uri"`
}

type dcRequest struct {
	Handler *string `msgpack:"handler"`
}

type newTicketReq struct {
	Channel *string      `msgpack:"channel"`
	Ticket  *interface{} `msgpack:"ticket"`
}

type newTicketRes struct {
	Sid     string `json:"sid"`
	Channel string `json:"channel"`
}

func handleConf(auth *authDutyCalls) {
	cred = *auth

	if cred.URI == "" {
		cred.URI = DefaultURI
	}

	timod.WriteConfOk()
}

func handleNewTicket(pkg *timod.Pkg) {
	var req newTicketReq
	err := msgpack.Unmarshal(pkg.Data, &req)
	if err != nil {
		timod.WriteEx(
			pkg.Pid,
			timod.ExBadData,
			"Error: Failed to unpack New Ticket request")
		return
	}

	params := url.Values{}
	jsonBody, err := json.Marshal(req.Ticket)
	if err != nil {
		timod.WriteEx(
			pkg.Pid,
			timod.ExBadData,
			fmt.Sprintf("Error: Failed to JSON marshal ticket (%s)", err))
		return
	}
	body := bytes.NewReader(jsonBody)

	params.Set("channel", *req.Channel)

	uri := filepath.Join(cred.URI, "tichet")

	reqURL, err := url.Parse(uri)

	if err != nil {
		timod.WriteEx(
			pkg.Pid,
			timod.ExBadData,
			fmt.Sprintf("Error: Failed to parse URI (%s) (%s)", uri, err))
		return
	}

	reqURL.RawQuery = params.Encode()

	httpReq, err := http.NewRequest("POST", reqURL.String(), body)
	if err != nil {
		timod.WriteEx(
			pkg.Pid,
			timod.ExBadData,
			fmt.Sprintf("Error: Failed to create request (%s)", err))
		return
	}

	httpReq.Header.Set("Content-Type", "application/json")

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

	keys := make([]newTicketRes, 0)
	err = json.Unmarshal(resBody, &keys)
	if err != nil {
		timod.WriteEx(
			pkg.Pid,
			timod.ExBadData,
			fmt.Sprintf("Error: Failed to unpack response (%s)", err))
		return
	}

	timod.WriteResponse(pkg.Pid, &keys[0].Sid)
}

func onModuleReq(pkg *timod.Pkg) {
	var req dcRequest
	err := msgpack.Unmarshal(pkg.Data, &req)
	if err != nil {
		timod.WriteEx(
			pkg.Pid,
			timod.ExBadData,
			"Error: Failed to unpack DutyCalls request")
		return
	}

	if *req.Handler == "new-ticket" {
		handleNewTicket(pkg)
	} else if req.Handler == nil {
		timod.WriteEx(
			pkg.Pid,
			timod.ExBadData,
			"Error: missing handler")
	} else {
		timod.WriteEx(
			pkg.Pid,
			timod.ExBadData,
			fmt.Sprintf("Error: unknown handler: %s", *req.Handler))
	}
}

func handler(buf *timod.Buffer, quit chan bool) {
	for {
		select {
		case pkg := <-buf.PkgCh:
			switch timod.Proto(pkg.Tp) {
			case timod.ProtoModuleConf:
				var auth authDutyCalls
				err := msgpack.Unmarshal(pkg.Data, &auth)
				if err == nil {
					handleConf(&auth)
				} else {
					log.Println("Error: Missing or invalid DutyCalls configuration")
					timod.WriteConfErr()
				}

			case timod.ProtoModuleReq:
				onModuleReq(pkg)

			default:
				log.Printf("Error: Unexpected package type: %d", pkg.Tp)
			}
		case err := <-buf.ErrCh:
			// In case of an error you probably want to quit the module.
			// ThingsDB will try to restart the module a few times if this
			// happens.
			log.Printf("Error: %s", err)
			quit <- true
		}
	}
}

func main() {
	// Starts the module
	timod.StartModule("dutycalls", handler)
}
