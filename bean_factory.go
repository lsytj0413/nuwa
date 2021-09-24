package nuwa

import (
	"reflect"

	"github.com/lsytj0413/nuwa/property"
)

// BeanFactory providing the full capabilities of SPI.
type BeanFactory interface {
	// GetBean return an instance, which may be shared or independent, of the specified bean.
	GetBean(name string) (interface{}, error)

	AliasRegistry
	BeanDefinitionRegistry
	property.Properties
}

// NewBeanFactory return the BeanFactory impl
func NewBeanFactory() BeanFactory {
	return &beanFactoryImpl{
		AliasRegistry:          NewAliasRegistry(),
		BeanDefinitionRegistry: NewBeanDefinitionRegistry(),
		Properties:             property.NewProperties(),
	}
}

type beanFactoryImpl struct {
	AliasRegistry
	BeanDefinitionRegistry
	property.Properties
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
			if fd.Property != nil {
				fv := v.Elem().Field(fd.FieldIndex)

				// If the field is ptr, first set it to the ptr to zero value
				// Otherwise it will be cannot setable
				// NOTE: the property will take care of this
				// if fv.Kind() == reflect.Ptr {
				// 	fv.Set(reflect.New(v.Elem().Type().Field(fd.FieldIndex).Type.Elem()))
				// 	fv = fv.Elem()
				// }

				err = f.Retrive(fd.Property.Name, fv)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return v.Interface(), nil
}
