package nuwa

import (
	"reflect"

	"github.com/lsytj0413/nuwa/property"
	"github.com/lsytj0413/nuwa/utils"
	"github.com/lsytj0413/nuwa/xerrors"
)

// BeanFactory providing the full capabilities of SPI.
type BeanFactory interface {
	// GetBean return an instance, which may be shared or independent, of the specified bean.
	GetBean(name string) (interface{}, error)
	RetriveBean(name string, bean interface{}) error
	RetriveBeans(beans interface{}) error

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

	err = f.RetriveBean(name, v)
	if err != nil {
		return nil, err
	}
	return v.Interface(), nil
}

func (f *beanFactoryImpl) RetriveBean(name string, bean interface{}) error {
	v, err := utils.IndirectToSetableValue(bean)
	if err != nil {
		return err
	}

	beanDefinition, err := f.GetBeanDefinition(name)
	if err != nil {
		return err
	}

	typ := beanDefinition.Type()
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	switch typ.Kind() {
	case reflect.Struct:
		for _, fd := range beanDefinition.FieldDescriptors() {
			if fd.Property != nil {
				fv := v.Field(fd.FieldIndex)

				// If the field is ptr, first set it to the ptr to zero value
				// Otherwise it will be cannot setable
				// NOTE: the property will take care of this
				// if fv.Kind() == reflect.Ptr {
				// 	fv.Set(reflect.New(v.Elem().Type().Field(fd.FieldIndex).Type.Elem()))
				// 	fv = fv.Elem()
				// }

				err = f.Retrive(fd.Property.Name, fv)
				if err != nil {
					return err
				}

				continue
			}

			if fd.Bean != nil {
				fv := v.Field(fd.FieldIndex)

				obj, err := f.GetBean(fd.Bean.Name)
				if err != nil {
					return err
				}
				fv.Set(reflect.ValueOf(obj))
			}
		}
	}

	return nil
}

func (f *beanFactoryImpl) RetriveBeans(bean interface{}) error {
	v := reflect.ValueOf(bean)
	if v.Kind() != reflect.Ptr {
		return xerrors.Errorf("cannot retrive beans to '%T', it must be pointer", bean)
	}
	v = v.Elem()
	if v.Kind() != reflect.Slice {
		return xerrors.Errorf("cannot retrive beans to '%T', is must be *slice", bean)
	}

	typ := v.Type().Elem()
	beanDefinitions := f.GetAllBeanDefinition()
	beanNames := []string{}

	for k, v := range beanDefinitions {
		if typ.Kind() == reflect.Interface {
			if v.Type().Implements(typ) {
				beanNames = append(beanNames, k)
			}
		} else {
			if v.Type() == typ {
				beanNames = append(beanNames, k)
			}
		}
	}

	ret := reflect.MakeSlice(v.Type(), 0, 0)
	for _, name := range beanNames {
		b, err := f.GetBean(name)
		if err != nil {
			return err
		}

		ret = reflect.Append(ret, reflect.ValueOf(b))
	}
	v.Set(ret)
	return nil
}
