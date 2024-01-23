package example

func errorReturn() error {
	return nil
}

type errorAlias = error

func errorAliasReturn() error {
	return nil
}

type errorType error

func errorTypeReturn() error {
	return nil
}

type errorInterface struct{}

func (e errorInterface) Error() string {
	return "i am an error"
}

func newErrorInterface() error {
	return errorInterface{}
}

var _ error = (*errorInterface)(nil)

// never suppose to be linted
func newErrorInterfaceConcrete() *errorInterface {
	return &errorInterface{}
}
