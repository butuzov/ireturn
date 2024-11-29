package example

/*
	Actual disallow directive.
*/
//nolint:ireturn
func dissAllowDirective1() interface{} { return 1 }

/*
	Some other linter in disallow directive.
*/
//nolint:return1
func dissAllowDirective2() interface{} { return 1 }

/*
	Coma separate linters in nolint directive. golangci-lint compatible mode
*/
//nolint:ireturn,nocode
func dissAllowDirective3() interface{} { return 1 }

/*
	Example of the good example if interlapsing name and actual ireturn with dot at the end.
*/
//nolint:ireturnlong,ireturn.
func dissAllowDirective4() interface{} { return 1 }

/*
	Example of the good example if interlapsing name and actual ireturn with dot at the end.
*/
//nolint:ireturnlong,ireturn
func dissAllowDirective5() interface{} { return 1 }

/*
	Not works!
*/
//nolint:ireturnlong,itertutireturn1.
func dissAllowDirective6() interface{} { return 1 }
