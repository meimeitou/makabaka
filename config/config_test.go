package config

import (
	"fmt"
	"testing"
)

func TestConfig(t *testing.T) {
	c, err := ReadConfigFile("../examples/config.yaml")
	if err != nil {
		panic(err)
	}
	fmt.Println(c)
}
