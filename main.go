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
	"fmt"
	"log"
	"sync"

	timod "github.com/thingsdb/go-timod"

	"github.com/vmihailenco/msgpack"
)

var mux sync.Mutex
var cred authDutyCalls

// DefaultURI is the default api URI to use
// The UIR can be changed using:
//
//     set_module_conf('dutycalls', {
//         login: 'mylogin',
//         password: 'mysecret',
//         uri: 'https://playground.dutycalls.me/api'
//     });
const DefaultURI = "https://dutycalls.me/api"

type authDutyCalls struct {
	Login    string `msgpack:"login"`
	Password string `msgpack:"password"`
	URI      string `msgpack:"uri"`
}

type dcRequest struct {
	Handler *string `msgpack:"handler"`
}

func handleConf(auth *authDutyCalls) {
	cred = *auth

	if cred.URI == "" {
		cred.URI = DefaultURI
	}

	timod.WriteConfOk()
}

func onModuleReq(pkg *timod.Pkg) {
	var req dcRequest
	err := msgpack.Unmarshal(pkg.Data, &req)
	if err != nil {
		timod.WriteEx(
			pkg.Pid,
			timod.ExBadData,
			"Failed to unpack DutyCalls request")
		return
	}

	if *req.Handler == "new-ticket" {
		handleNewTicket(pkg)
	} else if *req.Handler == "get-ticket" {
		handleGetTicket(pkg)
	} else if *req.Handler == "get-tickets" {
		handleGetTickets(pkg)
	} else if *req.Handler == "close-ticket" {
		handleStatusTicket(pkg, "closed")
	} else if *req.Handler == "close-tickets" {
		handleStatusTickets(pkg, "closed")
	} else if *req.Handler == "unack-ticket" {
		handleStatusTicket(pkg, "unacknowledged")
	} else if *req.Handler == "unack-tickets" {
		handleStatusTickets(pkg, "unacknowledged")
	} else if *req.Handler == "new-hit" {
		handleNewHit(pkg)
	} else if *req.Handler == "get-hits" {
		handleGetHits(pkg)
	} else if req.Handler == nil {
		timod.WriteEx(
			pkg.Pid,
			timod.ExBadData,
			"Missing handler")
	} else {
		timod.WriteEx(
			pkg.Pid,
			timod.ExBadData,
			fmt.Sprintf("Unknown handler: %s", *req.Handler))
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
					log.Println("Missing or invalid DutyCalls configuration")
					timod.WriteConfErr()
				}

			case timod.ProtoModuleReq:
				onModuleReq(pkg)

			default:
				log.Printf("Unexpected package type: %d", pkg.Tp)
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
