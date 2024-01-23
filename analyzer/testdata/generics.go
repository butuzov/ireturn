package example

type typeConstraints interface {
	~int | ~float64 | ~float32
}

func Min[T typeConstraints](x T, y T) (T, T) {
	if x > y {
		return y, y
	}
	return x, x
}

func MixedReturnParameters[T, K typeConstraints](x T, y K) (T, K, T, K, K, T) {
	return x, y, x, y, y, x
}

func Max[foobar typeConstraints](x foobar, y foobar) foobar {
	if x < y {
		return y
	}
	return x
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

func FuncWithGenericAny_NamedReturn[T_ANY any](foobar T_ANY) (variable T_ANY) {
	variable = foobar
	return
}

func FuncWithGenericAny[T_ANY any](foo T_ANY) T_ANY {
	return foo
}

func is_test() {
	if FuncWithGenericAny[int64](65) == 65 {
		print("yeap")
	}
}

// ISSUE: https://github.com/butuzov/ireturn/issues/37
type Map1[K comparable, V_COMPARABLE comparable] map[K]V_COMPARABLE

func (m Map1[K, V_COMPARABLE]) Get(key K) V_COMPARABLE { return m[key] }

type Map2[K comparable, V_ANY any] map[K]V_ANY

func (m Map2[K, V_ANY]) Get(key K) V_ANY { return m[key] }

// Empty Interface return
func FunctionAny() any {
	return nil
}

func FunctionInterface() interface{} {
	return nil
}
