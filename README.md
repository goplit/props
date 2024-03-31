# Props
Initialize configuration from defaults or environmental variables
using your own reference for the go struct and keep configuration
customizations to individual implementation

# Usage
## Install
```bash
go get github.com/goplit/props
```
## Initialize and pull options from Environment
```go
type Configuration struct {
	Host string         `key:"HOST" def:"localhost"`
	Port int            `key:"port" def:"4444"`
	IsUsingProxy bool   `key:"OPT_C" def:"false"`
}

c := &Configuration{}
props := New(c)
// Extract default values
p.InitDefaults()
// Extract environment vars
p.FromEnv()
// Commit to Configuration
p.Commit()
```
