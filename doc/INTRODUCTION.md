# Instances to Play around with

## Definitions

* **IID aware** or **enabled** - We call a service **IID aware** (**IID enabled**) if it responds to an IID-request with an Instance-Identification, which includes at least its **Miid**
* **Ciid** - We call the complete Instance Identifier (IID) that is transmitted in a header field **Ciid** (Complete Instance Identifier).
* **Miid** - An IID that describes exactly only one instance we call **Miid** (Micro-Service IID).
* **Epoch time** - The **epoch time** is the number of seconds that have elapsed since the service has been started. The **epoch time** is mandatory in any Miid.
* **Caller** - The microservice that initiates a request
* **Callee** - The microservice that receives a request and provides a response
* **Call Chain** - If a service A (**caller**) calls a service B (**callee**) we call this a **call chain**.
* **(Complete) Call Graph** - If a call chain contains all services that have been called we call this the **call graph**, or **complete call graph**
* **known call graph** - If we would like to emphasize that the **call graph** contains only services that we could control we call this **known call graph**
* **directed call graph** - lists all calls in the order they occurred.
* **recording service** - We call a service that records the services it has contacted to build a call graph **recording service**
* **disclosing service** - We call a **recording service** a **disclosing service** if it sends back the information to the client on request
* **shallow** - A service that reports as **Ciid** only its own **Miid**
* **deep** - A service that reports as  **Ciid** in addition to its own **Miid** the **Ciid**s of contacted or called services.
* **contacted services** - If a service does not maintain the order of services used to respond to a request, it reports the **contacted services**.
* **called services** - If a service does maintain the complete sequence of services it has used, it reports the **called services**  

## Ciids

### Single Ciids (Miids)

Describe only the identity of one service. Syntactically undistinguishable from Miids.

    MsA/0.1%123s
    MsA/1.1/devl-0faaccdd%234s

### Call-Chain Ciids

Describes that _service a_ calls  _service b_, and that _service b_ calls _service c_, etc.

    MsA/0.1%123s(MsB/0.2%234s)
    MsA/0.1%123s(MsB/0.2%234s(MsC/0.1%345s(MsD/0.3%456s)))

### Ciids as Call Graphs

Describes that _service a_ calls  _service b_ and then _service c_

    MsA/0.1%123s(MsB/0.2%234s+MsC/0.1%234s)
    MsA/0.1%123s(MsB/0.2%234s(MsC/0.1%345s(MsD/0.3%456s)+MsE/1.1%789s))

## Miids

### Generic Service

    MsA/0.1/dev-ccaaffee%123s
    MsA/0.1%123s
    MsB/0.2%234s
    MsC/0.1%345s
    MsD/0.3%456s
    MsE/1.1%

### External service

When recording calls to external services, there is not much knowledge about the instance of a service. However we might just want to cope with the fact that a service has been called.
Several proposal have been discussed.

Currently under evaluation is:

* use a service name a descriptive name
* replace the version number by `x`
* encode the called url in way that does not interfere with the iid encoding (bacially avoid `/`, `+` and `%`). base64 encoding **can not** be used, as it can contain `/` and `+`. 
Instead base64url, urlencode or [base62x](https://ufqi.com/dev/base62x/) could be used. For the sake of simplicity we currently evaluate [base64url](https://de.wikipedia.org/wiki/Base64#Base64url). A simple implemenation in golang can be found [here](https://github.com/dvsekhvalnov/jose2go/blob/v1.5.0/base64url/base64url.go)

## Call Graphs

There are several options on how to create a call graph, or on how a call is to be interpreted.

A call graph contains _**all**_ calls that have been made to respond to the requesting call. More on the interpretation of _**all**_ later.

However, there are several options on how a service can behave:

* A service discloses
  * its identity only (**shallow**) _or_
  * in addition the identity of services it has contacted (**deep**)

* A service
  * enumerates the services it has contacted (**contacted**) _or_
  * lists the services it has contacted in the order they have been contacted (**called**)

* A service can include information for services that do not disclose their identity

### Non-iid aware services, internal vs. external

When implementing IID capabilities in a service 
Every service that is out of the control of the organization is considered to be an external service. **External services** are likely not to be able to report their Ciids.

### By expectation vs by confirmation

If a service records (enumerates or lists) contacted services it has to expect that a service can or will not disclose its identity. This is, in particular, relevant for external services. We call a service that records its contacted services to return them as part of its Ciid a **recording services**.

A recording service can create a record entry either

* by expectation _or_
* by confirmation

**By expectation** means that a call graph entry is based on information created by the _caller_ while be **by confirmation** means that the respective call graph entry is created by the information as provided by the _callee_.

A **confirming** service _should_ record a called to service if, and only if the *callee* has responded with a valid Ciid. An **expecting** service _should_ record every service it is contacting and _should_ replace the _expected_ information with the _confirmed_ information.

A service _should_ not combine or mix reporting **by expectation** and **by confirmation**.


### A comparison of base64, base64url and urlencode

base64: 	https://www.base64encode.org/
urlencode: 	https://www.urlencoder.org/

Referenz:

	plain: 		https://de.wikiquote.org/wiki/Kleobulos_von_Lindos
	base64: 	aHR0cHM6Ly9kZS53aWtpcXVvdGUub3JnL3dpa2kvS2xlb2J1bG9zX3Zvbl9MaW5kb3M=
	urlencode: 	https%3A%2F%2Fde.wikiquote.org%2Fwiki%2FKleobulos_von_Lindos


	plain:		https://docs.google.com/spreadsheets/d/1nyJWSFDRFDk7yjdFzmDLZLKvMhDEVWW8eG9mKCGkjbQ/edit#gid=0
	base64: 	aaHR0cHM6Ly9kb2NzLmdvb2dsZS5jb20vc3ByZWFkc2hlZXRzL2QvMW55SldTRkRSRkRrN3lqZEZ6bURMWkxLdk1oREVWV1c4ZUc5bUtDR2tqYlEvZWRpdCNnaWQ9MA==
	urlencode: 	https%3A%2F%2Fdocs.google.com%2Fspreadsheets%2Fd%2F1nyJWSFDRFDk7yjdFzmDLZLKvMhDEVWW8eG9mKCGkjbQ%2Fedit%23gid%3D0


	plain:		https://www.google.com/search?client=firefox-b-d&q=example+query+parameters+url
	base64: 	aHR0cHM6Ly93d3cuZ29vZ2xlLmNvbS9zZWFyY2g/Y2xpZW50PWZpcmVmb3gtYi1kJnE9ZXhhbXBsZStxdWVyeStwYXJhbWV0ZXJzK3VybA==
	urlencode: 	https%3A%2F%2Fwww.google.com%2Fsearch%3Fclient%3Dfirefox-b-d%26q%3Dexample%2Bquery%2Bparameters%2Burl
