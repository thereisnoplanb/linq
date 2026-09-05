package linq

import "golang.org/x/exp/constraints"

// Real is a constraint that permits any integer or floating-point type.
type Real interface {
	constraints.Integer | constraints.Float
}

// Number is a constraint that permits any real or complex numeric type.
type Number interface {
	Real | constraints.Complex
}
