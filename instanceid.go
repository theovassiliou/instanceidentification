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

// An IID-Request Option
type Option interface {
	Command() string
}

// A collection of IID-Request options
type Options map[string]Option
type IidRequest interface {

	// SetIidAuth set's the authorisation key value. Chainable
	// Value xyz is included literally as key=xyz
	// If empty string is passed no authorisation key will be send
	SetIidAuth(string) IidRequest

	// GetIidAuth returns the authorisation key value, or an empty string if not set.
	GetIidAuth() string

	// Options returns the options set
	Options() Options

	// HasOptions returns true, in case options have been indicated, false otherwise
	HasOptions() bool

	// SetOption sets an option for the Iid-Request header. Chainable
	SetOption(Option) IidRequest

	// String returns the canonical iid-request value string represenation
	String() string

	// GetHeader returns the canonical iid-request Header string represenation
	GetHeader() string
}
