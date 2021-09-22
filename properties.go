package nuwa

// Properties hold the configurable properties name & value
type Properties interface {
	// Get return the value for key, the key is case-sensitive.
	Get(key string) interface{}

	// Set the value for key, it will overwrite the old value if key is already exists.
	Set(key string, val interface{})
}

// NewProperties return the Properties impl
func NewProperties() Properties {
	return propertiesImpl(make(map[string]interface{}))
}

type propertiesImpl map[string]interface{}

func (p propertiesImpl) Get(key string) interface{} {
	val, ok := p[key]
	if ok {
		return val
	}

	return nil
}

func (p propertiesImpl) Set(key string, val interface{}) {
	p[key] = val
}
