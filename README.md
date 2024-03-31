# Props
Initialize configuration from defaults or environmental variables
using your own reference for the go struct and keep configuration
customizations to individual implementation

Single dependency on yaml reader to keep things clean in your
mod file

# Usage
## Install
```bash
go get github.com/goplit/props
```
## Initialize and pull options from Environment
```go
// Configuration struct referehce
// usage:
//     key<string> as what to expect from various sources
//     def<string> as to what to apply if source don't have a key-value pair
type Configuration struct {
	Host         string `key:"HOST"  def:"localhost"`
	Port         int    `key:"port"  def:"4444"`
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

## Pull options from run arguments
```go
c := &Configuration{}
props := New(c)
// Extract default values
p.InitDefaults()
// Extract from arguments
p.FromArgs()
// Commit to Configuration
p.Commit()
```

## Initialize and pull options from yaml file
```go
c := &Configuration{}
props := New(c)
// Extract default values
p.InitDefaults()
// Extract from yaml file
p.FromYamlFile()
// Commit to Configuration
p.Commit()
```

## Rollover mix&match
```go
c := &Configuration{}
props := New(c)
// Extract default values
p.InitDefaults()
// Extract from yaml
p.FromYamlFile()
// Roll with envs
p.FromEnv()
// Roll with args
p.FromArgs()
// Commit to Configuration
p.Commit()
```