package structenv

import (
	"os"
	"reflect"
	"testing"
	"time"
)

func TestReadmeExample(t *testing.T) {
	type ServiceConfig struct {
		BindAddr       string        `env:"BIND_ADDR" default:"0.0.0.0:8080"`
		RequestTimeout time.Duration `env:"TIMEOUT" default:"3s"`
		LogDebug       bool          `env:"LOG_DEBUG"`
	}

	os.Setenv("TIMEOUT", "5s")
	os.Setenv("LOG_DEBUG", "yes")

	var v ServiceConfig
	if err := LoadFromEnv(&v); err != nil {
		t.Fatal(err)
	}

	want := ServiceConfig{
		BindAddr:       "0.0.0.0:8080",
		RequestTimeout: 5 * time.Second,
		LogDebug:       true,
	}
	if !reflect.DeepEqual(want, v) {
		t.Fatalf("unexpected result: want=%#v got=%#v", want, v)
	}
}

type LoadFromEnvFixture struct {
	FieldStr     string        `env:"key_str" default:"abc"`
	FieldInt     int           `env:"key_int" default:"12"`
	FieldFloat64 float64       `env:"key_float64" default:"11.2"`
	FieldDur     time.Duration `env:"key_dur" default:"3s"`
	FieldBool    bool          `env:"key_bool"`
}

func TestLoadFromEnv_Defaults(t *testing.T) {
	var got LoadFromEnvFixture
	if err := LoadFromEnv(&got); err != nil {
		t.Fatalf("failed loading values from env: %v", err)
	}

	want := LoadFromEnvFixture{
		FieldStr:     "abc",
		FieldInt:     12,
		FieldFloat64: 11.2,
		FieldDur:     3 * time.Second,
		FieldBool:    false, // no default set
	}

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("unexpected result: want=%+v got=%+v", want, got)
	}
}

func TestLoadFromEnv_Explicit(t *testing.T) {
	var got LoadFromEnvFixture

	os.Setenv("key_str", "xyz")
	os.Setenv("key_int", "22")
	os.Setenv("key_float64", "44.4")
	os.Setenv("key_dur", "4ms")
	os.Setenv("key_bool", "yes")

	if err := LoadFromEnv(&got); err != nil {
		t.Fatalf("failed loading values from env: %v", err)
	}

	want := LoadFromEnvFixture{
		FieldStr:     "xyz",
		FieldInt:     22,
		FieldFloat64: 44.4,
		FieldDur:     4 * time.Millisecond,
		FieldBool:    true,
	}

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("unexpected result: want=%+v got=%+v", want, got)
	}
}
