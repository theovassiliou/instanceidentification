package instanceid

import "time"

// Std header-name for HTTP-Requests
const XINSTANCEID = "X-Instance-Id"

// Defines the capabilities of a CIID
type Ciid interface {

	// Returns the miid part of the ciids
	Miid() Miid

	// Returns the call stack
	Ciids() Stack

	// Returns the canoncical instance id representation
	String() string

	// Sets the epoch to now, with time being startTime of service. Chainable
	SetEpoch(time.Time) Ciid

	// Sets the call stack. Chainable
	SetCiids(Stack) Ciid
}

// Stack represents a list of services that have been called by the Ciid
type Stack []Ciid

// Defines the capabilities of a MIID
type Miid interface {
	// Returns the service name
	Sn() string

	// Retunrs the version number
	Vn() string

	// Returns the application specific part of the miid
	Va() string

	// Returns the epoch of the miid
	T() int

	// Returns the canoncical instance id representation
	String() string

	// Sets the epoch of the miid in s. Chainable
	SetT(int) Miid

	// Sets the epoch to now, with time being startTime of service. Chainable
	SetEpoch(time.Time) Miid
}
