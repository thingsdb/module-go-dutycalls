# DutyCalls ThingsDB Module (Go)

DutyCalls module written using the [Go language](https://golang.org).


## Installation

Install the module by running the following command in the `@thingsdb` scope:

```javascript
new_module("dutycalls", "github.com/thingsdb/module-go-dutycalls");
```

Optionally, you can choose a specific version by adding a `@` followed with the release tag. For example: `@v0.1.0`.

## Configuration

The DutyCalls module requires configuration with the following properties:

Property | Type            | Description
-------- | --------------- | -----------
login    | str (required)  | Login to authenticate with.
password | str (required)  | Password / secret for the user.


Example configuration:

```javascript
set_module_conf("dutycalls", {
    login: "abcdefgh...",
    password: "hgfedcba...",
});
```

## Exposed functions

Name                            | Description
------------------------------- | -----------
[new_ticket](#new-ticket)       | Create a new ticket.
[get_ticket](#get-ticket)       | Get a ticket.
[get_tickets](#get-tickets)     | Get multiple tickets with a single request.
[close_ticket](#close-ticket)   | Close a ticket.
[close_tickets](#close-tickets) | Close a list of tickets.
[unack_ticket](#unack-ticket)   | Un-acknowledge a ticket.
[unack_tickets](#unack-tickets) | Un-acknowledge a list of tickets.
[new-hit](#new-hit)             | Create a new hit.
[get-hits](#get-hits)           | Get ticket hists.

### new ticket

Syntax: `new_ticket(channel, ticket)`

#### Arguments

- `channel`: _(str)_ Destination Channel to crete the ticket in.
- `ticket`: _(thing)_ Ticket, at least a title and body are required.

#### Example:

```javascript
ticket = {
    title: "Example ticket",
    body: "This is an example ticket."
};

dutycalls.new_ticket("mychannel", ticket).then(|sid| {
    sid;  // the SID (string) of the created ticket
}).else(|err| {
    err;  // some error has occurred
})
```

### get ticket

Syntax: `get_ticket(sid)`

#### Arguments

- `sid`: _(str)_ SID of the ticket.

#### Example:

```javascript
sid = "Ai782xf...";  // Some SID

dutycalls.get_ticket(sid).then(|ticket| {
    ticket;  // the ticket (thing) with the given SID
});
```

### get tickets

Syntax: `get_tickets([sid, ...])`

#### Arguments

- `[sid, ...]`: _(list of str)_ List with SIDs.

#### Example:

```javascript
sids = ["Ai782xf...", ["Aj35dwe..."];  // Some SIDs

dutycalls.get_tickets(sids).then(|tickets| {
    tickets;  // list with tickets for the given SIDs
});
```

### close ticket

Syntax: `close_ticket(sid)`

#### Arguments

- `sid`: _(str)_ SID of the ticket.

#### Example:

```javascript
sid = "Ai782xf...";  // Some SID

// Returns nil in case of success
dutycalls.close_ticket(sid).else(|err| {
    err;  // some error has occurred
});
```

### close tickets

Syntax: `close_tickets([sid, ...])`

#### Arguments

- `[sid, ...]`: _(list of str)_ List of SIDs.

#### Example:

```javascript
sids = ["Ai782xf...", "Aj35dwe..."];  // Some SIDs

// Returns nil in case of success
dutycalls.close_tickets(sids).else(|err| {
    err;  // some error has occurred
});
```

### unack ticket

Syntax: `unack_ticket(sid)`

#### Arguments

- `sid`: _(str)_ SID of the ticket.

#### Example:

```javascript
sid = "Ai782xf...";  // Some SID

// Returns nil in case of success
dutycalls.unack_ticket(sid).else(|err| {
    err;  // some error has occurred
});
```

### unack tickets

Syntax: `unack_tickets([sid, ...])`

#### Arguments

- `[sid, ...]`: _(list of str)_ List of SIDs.

#### Example:

```javascript
sids = ["Ai782xf...", "Aj35dwe..."];  // Some SIDs

// Returns nil in case of success
dutycalls.unack_tickets(sids).else(|err| {
    err;  // some error has occurred
});
```

### new hit

Syntax: `new_hit(sid, ticket)`

#### Arguments

- `sid`: _(str)_ Destination Ticket (SID) to crete the hit for.
- `hit`: _(thing)_ Hit, at least a summary is required.

#### Example:

```javascript
hit = {
    summary: "Example hist"
};

// Returns nil in case of success
dutycalls.new_hit("Ai782xf...", hit).else(|err| {
    err;  // some error has occurred
});
```


### get hits

Syntax: `get_hits(sid)`

#### Arguments

- `sid`: _(str)_ SID of the ticket.

#### Example:

```javascript
sid = "Ai782xf...";

dutycalls.get_hits(sid).then(|hits| {
    hits;  // some error has occurred
});
```
