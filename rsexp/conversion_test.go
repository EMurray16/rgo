package rsexp

import (
	"errors"
	"testing"
)

func TestNewGoSEXP(t *testing.T) {
	f := 3.14
	_, err := NewGoSEXP(f)
	if !errors.Is(err, NotASEXP) {
		t.Errorf("Got err %v instead of the expected NotASEXP\n", err)
	}
}
