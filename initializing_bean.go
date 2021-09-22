package nuwa

// InitializingBean is to be implemented by beans that need to react once all their properties
// have been set.
type InitializingBean interface {
	// AfterPropertiesSet is invoked after all properties have been set.
	AfterPropertiesSet() error
}
