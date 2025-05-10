# go-structenv

This library can be used to populate struct fields from environment variables.
For example, given a struct like so:

```
type ServiceConfig struct {
    BindAddr       string        `env:"BIND_ADDR" default:"0.0.0.0:8080"`
    RequestTimeout time.Duration `env:"TIMEOUT" default:"3s"`
    LogDebug       bool          `env:"LOG_DEBUG"`
}
```

And assuming you have set some of the relevant environment variables:

```
TIMEOUT=5s
LOG_DEBUG=yes
```

You can easily populate those fields:

```
var v ServiceConfig
if err := envstruct.LoadFromEnv(&v); err != nil {
    panic(err)
}
```

In summary, this library provides several useful things
* Support for many stdlib types: bool, string, int, float64, time.Duration, etc
* Handling of "truthy" boolean values, like "yes" vs "no" or "on" vs "off
* Optional ability to set default values
