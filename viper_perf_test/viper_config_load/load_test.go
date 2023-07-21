package viper_config_load

import (
	"testing"

	"github.com/spf13/viper"
)

func checkValuesViaViper(t interface{}) {

	type logFatal interface {
		Fatal(args ...any)
	}
	b, ok := t.(logFatal)
	if !ok {
		panic("input param doesn't support logFatal")
	}

	if viper.GetString("some_string") != "a big string here" {
		b.Fatal("some_string key didn't match expected value")
	}
	// if viper.GetInt("some_int") != 1000000 {
	// 	b.Fatal("some_int key didn't match expected value")
	// }
	// if viper.GetFloat64("some_float") != 1.101010101010 {
	// 	b.Fatal("some_float key didn't match expected value")
	// }
}

type Config struct {
	SomeString string  `mapstructure:"some_string"`
	SomeInt    int     `mapstructure:"some_int"`
	SomeFloat  float64 `mapstructure:"some_float"`
	SomeStruct struct {
		FirstField  string `mapstructure:"first_field"`
		SecondField int    `mapstructure:"second_field"`
		ThirdField  bool   `mapstructure:"third_field"`
	} `mapstructure:"some_struct"`
}

var appConfig *Config

func init() {
	conf := &Config{}
	err := viper.Unmarshal(conf)
	if err != nil {
		panic("err in viper.Unmarshal:" + err.Error())
	}

	appConfig = conf
}

func checkValuesViaStruct(t interface{}) {

	type logFatal interface {
		Fatal(args ...any)
	}
	b, ok := t.(logFatal)
	if !ok {
		panic("input param doesn't support logFatal")
	}

	if appConfig.SomeString != "a big string here" {
		b.Fatal("some_string key didn't match expected value")
	}
	// if appConfig.SomeInt != 1000000 {
	// 	b.Fatal("some_int key didn't match expected value")
	// }
	// if appConfig.SomeFloat != 1.101010101010 {
	// 	b.Fatal("some_float key didn't match expected value")
	// }
}

// func TestCheclValues(t *testing.T) {
// 	checkValuesViaViper(t)
// }

// func TestCheckValuesViaStruct(t *testing.T) {
// 	checkValuesViaStruct(t)
// 	t.Log(appConfig)
// }

func BenchmarkViperGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		checkValuesViaViper(b)
	}
}
func BenchmarkViperGetViaStruct(b *testing.B) {
	for i := 0; i < b.N; i++ {
		checkValuesViaStruct(b)
	}
}
