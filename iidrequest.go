package instanceid

import (
	"sort"
	"strings"
	"unicode"
)

type IRequest struct {
	key     string
	options Options
}

type IOption struct {
	commandName string
}

func (o IOption) Command() string {
	return o.commandName
}

// NewIRequestFromString creates a new IRequest from a value header passed as string
func NewIRequestFromString(v string) *IRequest {
	r := &IRequest{
		key: "empty",
	}
	return r.parseIidRequest(v)
}

// SetIidAuth set's the authorisation key value. Chainable
// Value xyz is included literally as key=xyz
// If empty string is passed no authorisation key will be send
func (r *IRequest) SetIidAuth(id string) IidRequest {
	if id == "" {
		r.key = "empty"
		return r
	}
	r.key = id
	return r
}

// GetIidAuth returns the authorisation key value, or an empty string if not set.
func (r IRequest) GetIidAuth() string {
	return r.key
}

// Options returns the options set
func (r IRequest) Options() Options {
	return r.options
}

// HasOptions returns true, in case options have been indicated, false otherwise
func (r IRequest) HasOptions() bool {
	return r.options != nil && len(r.options) > 0
}

// HasKey returns true, if a key element was present. False if empty, or no keys at all
func (r IRequest) HasKey() bool {
	return r.key != "empty" && r.key != ""
}

// SetOption sets an option for the Iid-Request header. Chainable
func (r *IRequest) SetOption(o Option) IidRequest {
	if o != nil {
		r.options[o.Command()] = o
	}
	return r
}

// String returns the canonical iid-request value string represenation
func (r IRequest) String() string {
	sB := strings.Builder{}
	if r.key != "empty" && r.key != "" {
		sB.WriteString("key=")
		sB.WriteString(r.key)
		sB.WriteString(" ")
	} else {
		sB.WriteString("empty")
		sB.WriteString(" ")
	}

	if r.options != nil && (len(r.options) > 0) {
		sB.WriteString("options=")

		keys := make([]string, 0, len(r.options))
		for k := range r.options {
			keys = append(keys, k)
		}

		// while technically not required, eases testing
		sort.Strings(keys)

		for _, k := range keys {
			sB.WriteString(k)
		}
	}
	return strings.Trim(sB.String(), " \t")
}

// GetHeader returns the canonical iid-request Header string represenation
func (r IRequest) GetHeader() string {

	return XINSTANCEID + ": " + r.String()
}

// parseIidRequest fills r based on given Iid header value
func (r *IRequest) parseIidRequest(id string) *IRequest {
	r.key = "empty"
	r.options = Options{}

	lastQuote := rune(0)
	f := func(c rune) bool {
		switch {
		case c == lastQuote:
			lastQuote = rune(0)
			return false
		case lastQuote != rune(0):
			return false
		case unicode.In(c, unicode.Quotation_Mark):
			lastQuote = c
			return false
		default:
			return unicode.IsSpace(c)

		}
	}

	// splitting string by space but considering quoted section
	items := strings.FieldsFunc(id, f)

	// create and fill the map
	m := make(map[string]string)
	for _, item := range items {
		x := strings.Split(item, "=")
		if len(x) == 2 {
			m[x[0]] = x[1]
		} else {
			m[x[0]] = ""
		}
	}

	// determine key
	for k, v := range m {
		if k == "empty" {
			r.SetIidAuth("empty")
		} else if k == "key" {
			r.SetIidAuth(v)
		} else if k == "options" {
			r.options = parseOption(v)
		}
	}

	return r
}

// parseIidRequest fills r based on given Iid header value
func parseOption(id string) (o Options) {
	if len(id) > 0 {
		o = make(Options)
	}
	for _, char := range id {
		o[string(char)] = IOption{commandName: string(char)}
	}
	return o
}
