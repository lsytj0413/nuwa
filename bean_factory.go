package nuwa

import (
	"reflect"
)

// BeanFactory providing the full capabilities of SPI.
type BeanFactory interface {
	// GetBean return an instance, which may be shared or independent, of the specified bean.
	GetBean(name string) (interface{}, error)

	AliasRegistry
	BeanDefinitionRegistry
	Properties
}

// NewBeanFactory return the BeanFactory impl
func NewBeanFactory() BeanFactory {
	return &beanFactoryImpl{
		AliasRegistry:          NewAliasRegistry(),
		BeanDefinitionRegistry: NewBeanDefinitionRegistry(),
		Properties:             NewProperties(),
	}
}

type beanFactoryImpl struct {
	AliasRegistry
	BeanDefinitionRegistry
	Properties
}

func (f *beanFactoryImpl) GetBean(name string) (interface{}, error) {
	beanDefinition, err := f.GetBeanDefinition(name)
	if err != nil {
		return nil, err
	}

	v, err := NewValue(beanDefinition.Type())
	if err != nil {
		return nil, err
	}

	typ := beanDefinition.Type()
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	switch typ.Kind() {
	case reflect.Struct:
		for _, fd := range beanDefinition.FieldDescriptors() {
			setter, err := NewFieldValueSetter(v, fd.FieldIndex)
			if err != nil {
				return nil, err
			}

			if fd.Property != nil {
				setter.Set(f.Get(fd.Property.Name))
				// v.Elem().Field(fd.FieldIndex).Set(reflect.ValueOf(f.Get(fd.Property.Name)))
				continue
			}
		}
	}

	return v.Interface(), nil
}
