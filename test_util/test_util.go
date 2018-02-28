package test_util

import "testing"

func Assert(t *testing.T, expected, got interface{}) {
	if expected != got {
		t.Errorf("Expected %v, got %v", expected, got)
	}
}
