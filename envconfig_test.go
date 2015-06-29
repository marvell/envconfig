package envconfig

import (
	"reflect"
	"testing"
)

type TestConfig struct {
	StrVar  string   `env:"STR_VAR" default:"StrVar" usage:"StrVar usage"`
	StrVars []string `env:"STR_VARS" default:"StrVar1,StrVar2" usage:"StrVars usage"`
	BoolVar bool     `env:"BOOL_VAR" default:"true" usage:"BoolVar usage"`
}

func TestMain(t *testing.T) {
	cfg := &TestConfig{}
	Parse(cfg)

	want1 := "StrVar"
	if cfg.StrVar != want1 {
		t.Errorf("StrVar == %v; want %v", cfg.StrVar, want1)
	}

	want2 := []string{"StrVar1", "StrVar2"}
	if !reflect.DeepEqual(cfg.StrVars, want2) {
		t.Errorf("StrVars == %v; want %v", cfg.StrVars, want2)
	}

	want3 := true
	if cfg.BoolVar != want3 {
		t.Errorf("BoolVar == %v; want %v", cfg.BoolVar, want3)
	}
}
