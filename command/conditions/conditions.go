package conditions

func TrueFunction() func() bool {
	return func() bool {
		return true
	}
}