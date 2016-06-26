package libwebawk

import (
	"testing"
)

func SimpleMatch(t *testing.T) {
	_, _, err := ParseWebawkProgram("/a/ {print $0}")

	if err != nil {
		t.Error("Expected non-error1")
	}
}

func TwoLevelMatch(t *testing.T) {
	_, _, err := ParseWebawkProgram("/a/ {print $0}")

	if err != nil {
		t.Error("Expected non-error2")
	}
}
