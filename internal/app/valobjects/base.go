package valobjects

import (
	"fmt"
)

type Value interface {
	fmt.Stringer
	Equals(value Value) bool
}
