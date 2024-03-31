package props

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
)

// Functions provided through wrapper for faker test cases
var getEnvOrDefault = envOrDefault
var getOsArgs = osArgs
var getOpenFile = openFile

type propertyType uint8

const (
	prop_default propertyType = 0
	prop_value   propertyType = 1
)

type setMap map[string]bool
type mapping map[string]*mapData
type mapData struct {
	// golang name of the value from the interface{}
	name string
	// golang type of the value derived from the interface{}
	typ3 reflect.Type
	// key which represents value in the provider data
	key string
	// default value available in the tag
	def string
	// intermediate value extracted from the initialization provider
	val interface{}
	// do not apply value to the ref object if object was set by any high level
	// configs than env vars setup, meaning [args,file,api] > envs
	skip bool
}

// Will map reference object fields to the interim structure
func mapFieldData(ref interface{}) mapping {
	search := make(mapping)
	v := reflect.ValueOf(ref)
	elem := v.Type().Elem()
	for i := 0; i < elem.NumField(); i++ {
		search[elem.Field(i).Name] = &mapData{
			name: elem.Field(i).Name,
			typ3: elem.Field(i).Type,
			key:  elem.Field(i).Tag.Get("key"),
			def:  elem.Field(i).Tag.Get("def"),
		}
	}
	return search
}

// Apply string type value to mapping interim structure
func valApply(m mapping, key string, val string) error {
	if mData, exists := m[key]; exists {
		var err error
		if len(val) > 0 {
			switch mData.typ3.Kind() {
			case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int8, reflect.Int16:
				mData.val, err = strconv.Atoi(val)
				if err != nil {
					return fmt.Errorf("cannot convert value %v to type of %v for key %v", val, mData.typ3.Kind(), mData.key)
				}
				break
			case reflect.String:
				mData.val = val
			case reflect.Bool:
				mData.val, err = strconv.ParseBool(val)
				if err != nil {
					return fmt.Errorf("cannot convert value %v to type of %v for key %v", val, mData.typ3.Kind(), mData.key)
				}
				break
			default:
				return fmt.Errorf("unsupported type of %v", mData.typ3.Kind())
			}
		}
	}
	return nil
}

// Apply mapping back to the referenced structure
func refApply(ref interface{}, m mapping, s setMap) {
	// Get object behind pointer
	ptr := reflect.ValueOf(ref)
	obj := ptr.Elem()
	for _, mappedData := range m {

		// Continue if property was meant to skip
		if mappedData.skip {
			continue
		}

		// reflect Value
		val := reflect.ValueOf(mappedData.val)

		// reflect Check if writable
		if val.IsValid() && obj.FieldByName(mappedData.name).CanSet() {
			obj.FieldByName(mappedData.name).Set(val)
			s[mappedData.name] = true
		}
	}
}

// Apply type of interface (multi-type) value to the mapping
func interfaceValApply(m mapping, key string, val interface{}) error {
	var err error
	if mData, exists := m[key]; exists {
		switch tVal := val.(type) {
		// String value is common value, similar procedure defined above
		case string:
			switch mData.typ3.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				mData.val, err = strconv.Atoi(tVal)
				if err != nil {
					return fmt.Errorf("cannot convert value %v to type of %v for key %v", val, mData.typ3.Kind(), mData.key)
				}
				break
			case reflect.String:
				mData.val = tVal
			case reflect.Bool:
				mData.val, err = strconv.ParseBool(tVal)
				if err != nil {
					return fmt.Errorf("cannot convert value %v to type of %v for key %v", val, mData.typ3.Kind(), mData.key)
				}
				break
			default:
				return fmt.Errorf("unsupported type of %v", mData.typ3.Kind())
			}
		// Integer cases will have special case for bool type
		case int:
			switch mData.typ3.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				mData.val = tVal
				if err != nil {
					return fmt.Errorf("cannot convert value %v to type of %v for key %v", val, mData.typ3.Kind(), mData.key)
				}
				break
			case reflect.String:
				mData.val = strconv.Itoa(tVal)
			case reflect.Bool:
				if tVal == 0 {
					mData.val = false
				} else if tVal == 1 {
					mData.val = true
				} else {
					return fmt.Errorf("cannot convert value %v to type of %v for key %v", val, mData.typ3.Kind(), mData.key)
				}
				break
			default:
				return fmt.Errorf("unsupported type of %v", mData.typ3.Kind())
			}
		case bool:
			switch mData.typ3.Kind() {
			case reflect.Bool:
				mData.val = tVal
				break
			default:
				return fmt.Errorf("unsupported bool into %v", mData.typ3.Kind())
			}
		}
	}
	return nil
}

func envOrDefault(name string, def string) (string, propertyType) {
	e, exists := os.LookupEnv(name)
	if !exists || len(e) == 0 {
		return def, prop_default
	}
	return e, prop_value
}

func osArgs() []string {
	return os.Args
}

func openFile(fileName string) ([]byte, error) {
	file, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("cannot read file %s, error %w", fileName, err)
	}
	return file, nil
}
