package nuwa

import (
	"reflect"
	"strings"

	"github.com/lsytj0413/nuwa/xerrors"
)

// BeanDefinition is the definition of bean
type BeanDefinition interface {
	// Type return the bean type
	Type() reflect.Type

	// Name return the bean name
	Name() string

	// FieldDescriptors return the struct field descriptors
	FieldDescriptors() []FieldDescriptor
}

// FieldDescriptor is the descriptor for struct field
type FieldDescriptor struct {
	FieldIndex int
	Name       string
	Typ        reflect.Type
	Unexported bool

	// Property is the field property descriptor.
	// The field should be marked as value=${name}
	Property *PropertyFieldDescriptor

	// Bean is the field bean descriptor.
	// The field should be marked as autowire=${name}
	Bean *BeanFieldDescriptor
}

// PropertyFieldDescriptor is the descriptor for property autowired value.
type PropertyFieldDescriptor struct {
	Name string
}

// BeanFieldDescriptor is the descriptor for bean autowired value.
type BeanFieldDescriptor struct {
	Name string
}

type BeanDefinitionOption interface {
	Apply(*beanDefinitionOption)
}

type beanDefinitionOption struct {
	name string
}

type WithName string

func (w WithName) Apply(opt *beanDefinitionOption) {
	opt.name = string(w)
}

type fieldTag struct {
	Value    *fieldValueTag
	Autowire *fieldAutowireTag
}

type fieldValueTag struct {
	Name string
}

type fieldAutowireTag struct {
	Name string
}

// MustNewBeanDefinition return the BeanDefinition impl, and it panics when some error happend.
func MustNewBeanDefinition(typ reflect.Type, opts ...BeanDefinitionOption) BeanDefinition {
	b, err := NewBeanDefinition(typ, opts...)
	if err != nil {
		panic(err)
	}

	return b
}

// NewBeanDefinition return the BeanDefinition impl
func NewBeanDefinition(typ reflect.Type, opts ...BeanDefinitionOption) (BeanDefinition, error) {
	opt := &beanDefinitionOption{}
	for _, o := range opts {
		o.Apply(opt)
	}

	if typ.Kind() != reflect.Ptr || typ.Elem().Kind() != reflect.Struct {
		return nil, xerrors.Errorf("Cannot build bean definition from '%T', it must be *struct", typ.String())
	}

	beanName := opt.name
	if beanName == "" {
		beanName = typ.Elem().Name()
	}

	fieldDescriptors := []FieldDescriptor{}
	for idx := 0; idx < typ.Elem().NumField(); idx++ {
		field := typ.Elem().Field(idx)
		// format: `nuwa:"value:${};autowire:xxx"`
		tag, ok := field.Tag.Lookup("nuwa")
		if !ok {
			continue
		}

		fd := &FieldDescriptor{
			FieldIndex: idx,
			Name:       field.Name,
			Typ:        field.Type,
			Unexported: false,
		}
		if field.PkgPath != "" {
			fd.Unexported = true
		}

		tags := strings.Split(tag, ";")
		for _, tag := range tags {
			switch {
			case strings.HasPrefix(tag, "value:"):
				v := tag[6:]
				if !(strings.HasPrefix(v, "${") && strings.HasSuffix(v, "}")) {
					return nil, xerrors.Errorf("could not parse value '%v', it must format as '${x}'", tag)
				}

				v = v[2 : len(v)-1]
				if len(v) == 0 {
					return nil, xerrors.Errorf("could not parse value '%v', the content is empty", tag)
				}
				fd.Property = &PropertyFieldDescriptor{
					Name: v,
				}
			case strings.HasPrefix(tag, "autowire:"):
				v := tag[9:]
				if len(v) == 0 {
					return nil, xerrors.Errorf("could not parse autowire '%v', the content is empty", tag)
				}
				fd.Bean = &BeanFieldDescriptor{
					Name: v,
				}
			default:
				return nil, xerrors.Errorf("could not parse tag '%v', it must start with 'value:' or 'autowire:'", tag)
			}
		}

		if fd.Property != nil && fd.Bean != nil {
			return nil, xerrors.Errorf("could not both set value & bean on field '%v'", fd.Name)
		}
		if fd.Property == nil && fd.Bean == nil {
			return nil, xerrors.Errorf("could not both unset value & bean on field '%v'", fd.Name)
		}
		fieldDescriptors = append(fieldDescriptors, *fd)
	}

	b := &beanDefinitionHolder{
		Typ:              typ,
		name:             beanName,
		fieldDescriptors: fieldDescriptors,
	}
	return b, nil
}

type beanDefinitionHolder struct {
	Typ  reflect.Type
	name string

	fieldDescriptors []FieldDescriptor
}

func (b *beanDefinitionHolder) Type() reflect.Type {
	return b.Typ
}

func (b *beanDefinitionHolder) Name() string {
	return b.name
}

func (b *beanDefinitionHolder) FieldDescriptors() []FieldDescriptor {
	return b.fieldDescriptors
}
