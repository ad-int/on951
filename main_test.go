package main

import (
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	code, err := run(m)
	if err != nil {
		panic(err)
	}
	os.Exit(code)
}
func run(m *testing.M) (int, error) {

	defer func() {
		log.Println("testing is done!")
	}()

	return m.Run(), nil // OK
}
