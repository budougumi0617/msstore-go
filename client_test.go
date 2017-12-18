// Copyright 2017 budougumi0617 All Rights Reserved.

package msstore

import "testing"

func TestDummy(t *testing.T) {
	_, err := NewClient("test", "test", "test")
	if err != nil {
		t.Fatalf("%q\n", err)
	}
}
