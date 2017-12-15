package msstore

import "testing"

func TestDummy(t *testing.T) {
	_, err := NewClient("test")
	if err != nil {
		t.Fatalf("%q\n", err)
	}
}
