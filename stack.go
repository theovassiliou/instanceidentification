package instanceid

// Stack represents a list of services that have been called by the Ciid
type Stack []Ciid

// Push a new value onto the stack
func (s *Stack) Push(str Ciid) {
	*s = append(*s, str) // Simply append the new value to the end of the stack
}

// IsEmpty: check if stack is empty
func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

// Remove and return top element of stack. Return false if stack is empty.
func (s *Stack) Pop() (Ciid, bool) {
	if s.IsEmpty() {
		return Ciid{}, false
	} else {
		index := len(*s) - 1   // Get the index of the top most element.
		element := (*s)[index] // Index into the slice and obtain the element.
		*s = (*s)[:index]      // Remove it from the stack by slicing it off.
		return element, true
	}
}

func (myself *Ciid) SetStack(callStack Stack) *Ciid {
	myself.Ciids = callStack
	return myself
}

func (myself *Ciid) ClearStack() *Ciid {
	myself.Ciids = nil
	return myself
}
