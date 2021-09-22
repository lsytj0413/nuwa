package nuwa

import (
	"fmt"
	"sync"
)

// AliasRegistry is the interface for managing aliases.
type AliasRegistry interface {
	// RegisterAlias register an alias for an name.
	// It will return err if the alias is already in use.
	RegisterAlias(name string, alias string) error

	// RemoveAlias remove the specified alias from this registry.
	// It will return err if no such alias was found.
	RemoveAlias(alias string) error

	// IsAlias determine whether the given name is defines as an alias.
	IsAlias(alias string) bool

	// GetAliases return the aliases for the given name, if it was defined.
	GetAliases(name string) []string

	// CanonicalName determine the raw name, resolving aliases to canonical names.
	CanonicalName(alias string) string
}

// NewAliasRegistry return the AliasRegistry impl
func NewAliasRegistry() AliasRegistry {
	return &aliasRegistryImpl{
		aliasMap: make(map[string]string),
	}
}

type aliasRegistryImpl struct {
	aliasMap map[string]string
	lock     sync.RWMutex
}

func (r *aliasRegistryImpl) RegisterAlias(name string, alias string) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if name == alias {
		delete(r.aliasMap, alias)
		return nil
	}

	registeredName, ok := r.aliasMap[alias]
	if ok {
		if registeredName == name {
			return nil
		}

		return fmt.Errorf("Cannot define alias '%v' for name '%v': It is already registered for name '%v'", alias, name, registeredName)
	}

	// Check for alias circle
	r.aliasMap[alias] = name
	return nil
}

func (r *aliasRegistryImpl) RemoveAlias(alias string) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	_, ok := r.aliasMap[alias]
	if !ok {
		return fmt.Errorf("No alias '%v' registered", alias)
	}

	delete(r.aliasMap, alias)
	return nil
}

func (r *aliasRegistryImpl) IsAlias(alias string) bool {
	r.lock.RLock()
	defer r.lock.RUnlock()

	_, ok := r.aliasMap[alias]
	return ok
}

func (r *aliasRegistryImpl) GetAliases(name string) []string {
	r.lock.RLock()
	defer r.lock.RUnlock()

	ret := []string{}
	var retrieveAliases func(string)
	retrieveAliases = func(name string) {
		for alias, registeredName := range r.aliasMap {
			if registeredName == name {
				ret = append(ret, alias)
				retrieveAliases(alias)
			}
		}
	}
	retrieveAliases(name)
	return ret
}

func (r *aliasRegistryImpl) CanonicalName(alias string) string {
	r.lock.RLock()
	defer r.lock.RUnlock()

	canonicalName := alias
	for name, ok := r.aliasMap[alias]; ok; {
		canonicalName = name
	}
	return canonicalName
}
