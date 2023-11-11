package service

import (
	"fmt"
	"testing"
)

func TestCode(t *testing.T) {
	t.Log(fmt.Sprintf("%06d", 1))
}
