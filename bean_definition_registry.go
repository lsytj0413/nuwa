package nuwa

import (
	"fmt"
	"sync"
)

// BeanDefinitionRegistry is the interface that hold bean definitions.
type BeanDefinitionRegistry interface {
	// RegisterBeanDefinition register a new bean definition with this registry.
	// It will return err if BeanDefinition is invalid or beanName is already exists.
	RegisterBeanDefinition(beanName string, beanDefinition BeanDefinition) error

	// RemoveBeanDefinition remove the bean definition for the given name.
	// It will return err if there is no such bean definition.
	RemoveBeanDefinition(beanName string) error

	// GetBeanDefinition return the bean definition for the given bean name.
	// It will return err if there is no such bean definition.
	GetBeanDefinition(beanName string) (BeanDefinition, error)
}

// NewBeanDefinitionRegistry return the BeanDefinitionRegistry impl
func NewBeanDefinitionRegistry() BeanDefinitionRegistry {
	return &beanDefinitionRegistryImpl{
		beanDefinitionMap: make(map[string]BeanDefinition),
	}
}

type beanDefinitionRegistryImpl struct {
	beanDefinitionMap map[string]BeanDefinition
	lock              sync.RWMutex
}

func (r *beanDefinitionRegistryImpl) RegisterBeanDefinition(
	beanName string,
	beanDefinition BeanDefinition,
) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	_, ok := r.beanDefinitionMap[beanName]
	if ok {
		return fmt.Errorf("Cannot register bean '%v': It is already registered", beanName)
	}

	r.beanDefinitionMap[beanName] = beanDefinition
	return nil
}

func (r *beanDefinitionRegistryImpl) RemoveBeanDefinition(beanName string) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	_, ok := r.beanDefinitionMap[beanName]
	if !ok {
		return fmt.Errorf("No bean '%v' registered", beanName)
	}

	delete(r.beanDefinitionMap, beanName)
	return nil
}

func (r *beanDefinitionRegistryImpl) GetBeanDefinition(beanName string) (BeanDefinition, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	beanDefinition, ok := r.beanDefinitionMap[beanName]
	if !ok {
		return nil, fmt.Errorf("No bean '%v' registered", beanName)
	}

	return beanDefinition, nil
}
