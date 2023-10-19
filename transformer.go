package fractal

import "errors"

type TransformerInterface interface {
	Transform() map[string]interface{}
	CallIncludeFunction(funName string) (ResourceInterface, error)
	setData(data interface{})
	GetData() interface{}
	SetAvailableIncludeFunctions()
}

type Transformer struct {
	data                      interface{}
	AvailableIncludeFunctions map[string]func() ResourceInterface
}

func (t *Transformer) SetAvailableIncludeFunctions() {
	t.AvailableIncludeFunctions = nil
}

func (t *Transformer) setData(data interface{}) {
	t.data = data
}

func (t *Transformer) GetData() interface{} {
	return t.data
}

func (t *Transformer) CallIncludeFunction(funName string) (ResourceInterface, error) {
	if fun, ok := t.AvailableIncludeFunctions[funName]; ok {
		return fun(), nil
	}
	return nil, errors.New("include function not found")
}
