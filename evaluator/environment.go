package evaluator

// Environment stores variables in the current scope
type Environment struct {
	store map[string]Value
	outer *Environment
}

// NewEnvironment creates a new environment
func NewEnvironment() *Environment {
	return &Environment{
		store: make(map[string]Value),
		outer: nil,
	}
}

// NewEnclosedEnvironment creates a new environment with an outer scope
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

// Get retrieves a value from the environment
func (e *Environment) Get(name string) (Value, bool) {
	value, ok := e.store[name]
	if !ok && e.outer != nil {
		return e.outer.Get(name)
	}
	return value, ok
}

// Set stores a value in the environment
func (e *Environment) Set(name string, val Value) Value {
	e.store[name] = val
	return val
}

// Update updates an existing variable (for mutable variables)
func (e *Environment) Update(name string, val Value) (Value, bool) {
	// Check current scope
	if _, ok := e.store[name]; ok {
		e.store[name] = val
		return val, true
	}

	// Check outer scopes
	if e.outer != nil {
		return e.outer.Update(name, val)
	}

	return nil, false
}
