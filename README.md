# DutyCalls ThingsDB Module (Go)

DutyCalls module written using the [Go language](https://golang.org).


## Installation

Install the module by running the following command in the `@thingsdb` scope:

```javascript
new_module('dutycalls', 'github.com/thingsdb/module-go-dutycalls');
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
set_module_conf('dutycalls', {
    login: 'iris',
    password: 'siri',
});
```

## Exposed functions

Name                        | Description
--------------------------- | -----------
[new_ticket](#new-ticket)   | Create a new ticket.

### new ticket

Syntax: `new_ticket(channel, ticket)`

#### Arguments

- `channel`: Destination Channel to crete the ticket in.
- `ticket`: Ticket

#### Example:

```javascript
siridb.new_ticket("mychannel").then(|res| {
    res;  // just return the response.
});
```
