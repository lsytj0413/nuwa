package props

// Properties hold the configurable properties name & value.
type Properties interface {
	// Get return the value for key, the key is case-sensitive.
	Get(key string) (string, error)

	// Retrive will return the value for key, and set it to i.
	// NOTE: the i must be setable
	Retrive(key string, i interface{}) error

	// Set the value for key, it will overwrite the old value if key is already exists.
	// The val will been transform to string to store with the container, so when
	//   user try to get the val of key, the string returned.
	Set(key string, val interface{}) error
}
