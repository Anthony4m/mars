package evaluator

import "fmt"

type Binding struct {
	Value     Value
	IsMutable bool
}

// Environment stores variables in the current scope
type Environment struct {
	store map[string]Binding
	outer *Environment
}

// NewEnvironment creates a new environment
func NewEnvironment() *Environment {
	return &Environment{
		store: make(map[string]Binding),
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
func (e *Environment) Get(name string) (Binding, bool) {
	binding, ok := e.store[name]
	//check current scope
	if ok {
		return binding, true
	}
	//check outer scope
	if e.outer != nil {
		binding, ok = e.outer.Get(name)
		if ok {
			return binding, true
		}
	}

	return Binding{}, false
}

// Set stores a value in the environment
func (e *Environment) Set(name string, val Value, isMutable bool) Binding {
	bind := Binding{val, isMutable}
	e.store[name] = bind
	return bind
}

// Update updates an existing variable (for mutable variables)
func (e *Environment) Update(name string, val Value) error {
	// Check current scope
	if binding, ok := e.store[name]; ok {
		if binding.IsMutable == false {
			return fmt.Errorf("cannot assign to immutable variable '%s'", name)
		}
		bind := e.Set(name, val, binding.IsMutable)
		e.store[name] = bind
		return nil
	}

	// Check outer scopes
	if e.outer != nil {
		return e.outer.Update(name, val)
	}

	return fmt.Errorf("undefined variable '%s'", name)
}
