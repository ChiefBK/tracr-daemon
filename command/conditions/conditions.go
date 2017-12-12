package conditions

var ConditionFunctions = make(map[string]func() bool)

func TrueFunction() bool {
	return true
}
