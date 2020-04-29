package task

import (
	"testing"
)

func TestVMUniqueIdProperties(t *testing.T) {
	id1 := generateVMID()
	id2 := generateVMID()

	if id1 == id2 {
		t.Fatalf("Two auto-generated VM ID should not be equivalent")
	}
	if len(id1) != 11 || len(id2) != 11 {
		t.Fatalf("VM ID was expected to be 10")
	}
}
