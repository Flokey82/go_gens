package aitree

// ForceFailure returns failure on completion.
type ForceFailure struct {
	Name  string
	Child Node
}

// NewForceFailure returns a new ForceFailure decorator.
func NewForceFailure(name string, child Node) *ForceFailure {
	return &ForceFailure{
		Name:  name,
		Child: child,
	}
}

// Tick will execute the child and return failure on completion.
func (node *ForceFailure) Tick() State {
	if node.Child.Tick() == StateRunning {
		return StateRunning
	}
	return StateFailure
}

// ForceSuccess returns success on completion.
type ForceSuccess struct {
	Name  string
	Child Node
}

// NewForceSuccess returns a new ForceSuccess decorator.
func NewForceSuccess(name string, child Node) *ForceSuccess {
	return &ForceSuccess{
		Name:  name,
		Child: child,
	}
}

// Tick will execute the child and return success on completion.
func (node *ForceSuccess) Tick() State {
	if node.Child.Tick() == StateRunning {
		return StateRunning
	}
	return StateSuccess
}

// Inverter implements an inverter decorator (success on failure and vice versa).
type Inverter struct {
	Name  string
	Child Node
}

// NewInverter returns a new Inverter decorator.
func NewInverter(name string, child Node) *Inverter {
	return &Inverter{
		Name:  name,
		Child: child,
	}
}

// Tick will execute the child and return success on failure and vice versa.
func (node *Inverter) Tick() State {
	result := node.Child.Tick()
	if result == StateRunning {
		return StateRunning
	}
	if result == StateSuccess {
		return StateFailure
	}
	return StateSuccess
}

// Retry implements a retry logic decorator.
type Retry struct {
	Name     string
	Child    Node
	Current  int
	Attempts int
}

// NewRetry returns a new retry decorator.
func NewRetry(name string, child Node, n int) *Retry {
	return &Retry{
		Name:     name,
		Child:    child,
		Attempts: n,
	}
}

// Tick will attempt up to n times to execute the child until it succeeds or will return failure if the retries are exhausted.
func (node *Retry) Tick() State {
	if node.Current < node.Attempts {
		result := node.Child.Tick()
		if result == StateRunning {
			return StateRunning
		}
		if result == StateSuccess {
			node.Current = 0
			return StateSuccess
		}
		node.Current++
		if node.Current < node.Attempts {
			return StateRunning
		}
		node.Current = 0
		return StateFailure
	}
	return StateFailure
}

// Repeater implements a repeater logic decorator.
type Repeater struct {
	Name     string
	Child    Node
	Current  int
	Attempts int
}

// NewRepeater returns a new repeater decorator.
func NewRepeater(name string, child Node, n int) *Repeater {
	return &Repeater{
		Name:     name,
		Child:    child,
		Attempts: n,
	}
}

// Tick will execute the child up to n times and will return success if all attempts succeed.
func (node *Repeater) Tick() State {
	if node.Current < node.Attempts {
		result := node.Child.Tick()
		if result == StateRunning {
			return StateRunning
		}
		if result == StateFailure {
			node.Current = 0
			return StateFailure
		}
		node.Current++
		if node.Current < node.Attempts {
			return StateRunning
		}
		node.Current = 0
		return StateSuccess
	}
	return StateSuccess
}
