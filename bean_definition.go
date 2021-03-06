package nuwa

import (
	"reflect"
)

// BeanDefinition is the definition of bean
type BeanDefinition interface {
	// Type return the bean type
	Type() reflect.Type

	// Name return the bean name
	Name() string

	// Scope return the scope of bean
	Scope() Scope

	// ConstructorArgumentValues return the argument for factory bean
	ConstructorArgumentValues() []reflect.Value

	// FactoryBeanName return the bean name for factory bean
	FactoryBeanName() string

	// InitMethodName return the init method for bean
	InitMethodName() string

	// DestroyMethodName return the destroy method for bean
	DestroyMethodName() string

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

// NewBeanDefinition return the BeanDefinition impl
func NewBeanDefinition() BeanDefinition {
	return &BeanDefinitionImpl{}
}

type BeanDefinitionImpl struct {
	Typ               reflect.Type
	name              string
	scope             Scope
	consArgs          []reflect.Value
	factoryBeanName   string
	initMethodName    string
	destroyMethodName string
	fieldDescriptors  []FieldDescriptor
}

func (b *BeanDefinitionImpl) Type() reflect.Type {
	return b.Typ
}

func (b *BeanDefinitionImpl) Name() string {
	return b.name
}

func (b *BeanDefinitionImpl) Scope() Scope {
	return b.scope
}

func (b *BeanDefinitionImpl) ConstructorArgumentValues() []reflect.Value {
	return b.consArgs
}

func (b *BeanDefinitionImpl) FactoryBeanName() string {
	return b.factoryBeanName
}

func (b *BeanDefinitionImpl) InitMethodName() string {
	return b.initMethodName
}

func (b *BeanDefinitionImpl) DestroyMethodName() string {
	return b.destroyMethodName
}

func (b *BeanDefinitionImpl) FieldDescriptors() []FieldDescriptor {
	return b.fieldDescriptors
}
