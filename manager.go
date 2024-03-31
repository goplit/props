package props

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

type Properties struct {
	set setMap
	obj mapping
	ref interface{}
}

func New(reference interface{}) Properties {
	m := mapFieldData(reference)
	p := Properties{
		set: make(setMap),
		obj: m,
		ref: reference,
	}
	return p
}

func (p Properties) InitDefaults() error {
	for key, mapElem := range p.obj {
		// If already set by another pass, then mark as skip
		if _, alreadySet := p.set[mapElem.name]; alreadySet {
			mapElem.skip = true
			continue
		}
		// Apply default value
		err := valApply(p.obj, key, mapElem.def)
		if err != nil {
			return fmt.Errorf("init defaults fail, error: %w", err)
		}
	}
	return nil
}

func (p Properties) FromEnv() error {
	for key, mapElem := range p.obj {
		val, kind := getEnvOrDefault(mapElem.key, mapElem.def)
		// Check if env was not present, and it's a default value we try to reset
		if _, alreadySet := p.set[mapElem.name]; alreadySet && kind == prop_default {
			mapElem.skip = true
			continue
		}
		// Apply default value
		err := valApply(p.obj, key, val)
		if err != nil {
			return fmt.Errorf("from env fail, error: %w", err)
		}
	}
	return nil
}

func (p Properties) FromYamlFile(fileName string) error {
	// Pull yaml file information
	if len(fileName) == 0 {

		return fmt.Errorf("no file was provided")
	}
	yamlMap := make(map[string]interface{})
	file, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("cannot read yaml file %s, error %w", fileName, err)
	}
	err = yaml.Unmarshal(file, yamlMap)
	if err != nil {
		return fmt.Errorf("cannot unmarshal yaml file, error %w", err)
	}
	// Apply values to the mapping
	for key, mapElem := range p.obj {
		// Check if value from the properties presented in yaml
		if val, found := yamlMap[mapElem.key]; found {
			// Apply default value
			err := interfaceValApply(p.obj, key, val)
			if err != nil {
				return fmt.Errorf("from yaml fail, error: %w", err)
			}
		}
	}
	return nil
}

// FromArgs
// Will take presented os.Args and try to match and fill them to properties reference
func (p Properties) FromArgs() error {
	var err error

	args := getOsArgs()
	argMap := make(map[string]string)
	for _, arg := range args {
		kval := strings.Split(arg, "=")
		if len(kval) == 2 {
			argMap[strings.ToLower(kval[0])] = kval[1]
		} else if len(kval) == 1 {
			argMap[strings.ToLower(kval[0])] = ""
		}
	}

	for key, mapElem := range p.obj {
		if val, exists := argMap[strings.ToLower(mapElem.key)]; exists {
			err = valApply(p.obj, key, val)
			if err != nil {
				return fmt.Errorf("from args fail, error: %w", err)
			}
		}
	}

	return err
}

// FromCallback
// Should try to form properties from some callback,
// underlying function can retrieve data from any API
// @param fn func()(map[string]string, error, bool)
func (p Properties) FromCallback(fn func() (map[string]string, error)) error {
	if fn == nil {
		return fmt.Errorf("no func provided")
	}
	// Get map of parameters
	props, err := fn()
	if err != nil {
		return fmt.Errorf("from callback fail, error %w", err)
	}
	// Reassure those will be searchable
	for k, v := range props {
		props[strings.ToLower(k)] = v
	}
	// Match
	for key, mapElem := range p.obj {
		if val, exists := props[strings.ToLower(mapElem.key)]; exists {
			err = valApply(p.obj, key, val)
			if err != nil {
				return fmt.Errorf("from args fail, error: %w", err)
			}
		}
	}
	return nil
}

func (p Properties) Commit() {
	refApply(p.ref, p.obj, p.set)
}
