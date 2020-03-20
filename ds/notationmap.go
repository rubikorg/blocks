package ds

import (
	"errors"
	"fmt"
)

// NotationMap is a type of map structure that can get you the value of a
// embedded key inside a map
type NotationMap struct {
	m        map[string]interface{}
	isFlat   bool
	editable bool
}

// NewNotationMap ...
func NewNotationMap(canEdit bool) NotationMap {
	return NotationMap{}
}

// Assign reassigns the holding map `m` inside the struct
func (nm NotationMap) Assign(m map[string]interface{}) error {
	if nm.editable {
		nm.m = m
		return nil
	}
	return errors.New("NotationMapAssignError: Cannot assign a non-editable notation map")
}

// Flatten functions creates a flat map of accessor with dot notations
func (nm NotationMap) Flatten() NotationMap {
	final := make(map[string]interface{})

	for k, v := range nm.m {
		if interm, ok := v.(map[string]interface{}); ok {
			final[k] = interm
			childs := traverseObjects(interm, k)
			final = override(final, childs)
		} else {
			final[k] = v
		}

	}

	return NotationMap{
		m:        final,
		editable: nm.editable,
		isFlat:   true,
	}
}

// Map returns the holding map instance for population
func (nm NotationMap) Map() map[string]interface{} {
	return nm.m
}

// Length returns the length of the holding map
func (nm NotationMap) Length() int {
	return len(nm.m)
}

func traverseObjects(target map[string]interface{}, parent string) map[string]interface{} {
	newSource := make(map[string]interface{})

	for k, v := range target {
		baseKey := fmt.Sprintf("%s.%s", parent, k)
		if interm, ok := v.(map[string]interface{}); ok {
			newSource[baseKey] = interm
			childs := traverseObjects(interm, baseKey)
			newSource = override(newSource, childs)
		} else {
			newSource[baseKey] = v
		}

	}

	return newSource
}

func override(host, source map[string]interface{}) map[string]interface{} {
	for k, v := range source {
		host[k] = v
	}
	return host
}

// Get values of key using dot notations from NotationMap
func (nm NotationMap) Get(accessor string) interface{} {
	return nm.m[accessor]
}

// Set value of a accessor using dot notations from NotationMap
func (nm NotationMap) Set(accessor string, value interface{}) error {
	if nm.editable {
		nm.m[accessor] = value
		return nil
	}
	return errors.New("NotationMapSetError: Cannot edit a non-editable NotationMap")
}

func noSuchKeyErr(key, acc string) error {
	return errors.New("no such key: " + key + " for notation: " + acc)
}
