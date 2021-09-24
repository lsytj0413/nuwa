package nuwa

// Scope is the bean life scope
type Scope = string

const (
	// ScopeSingleton identifier for the standard singleton scope.
	ScopeSingleton Scope = "singleton"

	// ScopePrototype identifier for the standard prototype scope.
	ScopePrototype Scope = "prototype"
)
