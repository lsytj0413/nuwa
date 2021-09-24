package nuwa

// BeanPostProcessor hook that allows for custom modification of new bean instances.
// For example, checking for marker interfaces or wrapping beans with proxies.
type BeanPostProcessor interface {
	// PostProcessBeforeInitialization will been apply to the given new bean instance before any bean initialization callbacks.
	// It will called before AfterPropertiesSet & custom init-method, and the bean will already be populated with property values.
	PostProcessBeforeInitialization(obj interface{}, beanName string) (v interface{}, err error)

	// PostProcessAfterInitialization will been apply to the given new bean instance after any bean initialization callbacks.
	PostProcessAfterInitialization(obj interface{}, beanName string) (v interface{}, err error)
}
