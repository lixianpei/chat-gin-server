package test

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("TestMain...")
	os.Exit(m.Run())
}

func TestA(t *testing.T) {
	fmt.Println("TestA...")
}
