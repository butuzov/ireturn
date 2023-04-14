package example

type Numeric interface {
	~int
}

func Min[T Numeric](x T, y T) (T, T) {
	if x > y {
		return y, y
	}
	return x, x
}

func MixedReturnParameters[T, K Numeric](x T, y K) (T, K, T, K, K, T) {
	return x, y, x, y, y, x
}

func Max[foobar Numeric](x foobar, y foobar) foobar {
	if x < y {
		return y
	}
	return x
}

func Foo[GENERIC any](foobar GENERIC) (variable GENERIC) {
	variable = foobar
	return
}

// SumIntsOrFloats sums the values of map m. It supports both int64 and float64
// as types for map values.
func SumIntsOrFloats[K comparable, V int64 | float64](m map[K]V) V {
	var s V
	for _, v := range m {
		s += v
	}
	return s
}
