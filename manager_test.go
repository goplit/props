package props

import (
	"reflect"
	"testing"
)

type Configuration struct {
	OptionA string `key:"OPT_A" def:"unfilled"`
	OptionB int    `key:"OPT_B" def:"4444"`
	OptionC bool   `key:"OPT_C" def:"false"`
}

func assert(v1 interface{}, v2 interface{}, t *testing.T) {
	if !reflect.DeepEqual(v1, v2) {
		t.Errorf("%s val <%v> != <%v>", t.Name(), v1, v2)
	}
}

func TestProperties_FromEnv(t *testing.T) {
	// Add properties
	c := &Configuration{}
	props := New(c)
	// Replace with faker
	getEnvOrDefault = func(name string, def string) (string, propertyType) {
		switch name {
		case "OPT_A":
			return "filled", prop_value
		case "OPT_B":
			return "8888", prop_value
		case "OPT_C":
			return "true", prop_value
		}
		return def, prop_default
	}
	type fields struct {
		set setMap
		obj mapping
		ref interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"Test getting options from env",
			fields{
				set: props.set,
				obj: props.obj,
				ref: props.ref,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			p := Properties{
				set: tt.fields.set,
				obj: tt.fields.obj,
				ref: tt.fields.ref,
			}

			if err := p.FromEnv(); (err != nil) != tt.wantErr {
				t.Errorf("FromEnv() error = %v, wantErr %v", err, tt.wantErr)
			}
			p.Commit()

			assert(c.OptionA, "filled", t)
			assert(c.OptionB, 8888, t)
			assert(c.OptionC, true, t)
		})
	}
}
