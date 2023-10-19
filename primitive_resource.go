package fractal

import (
	"fmt"
)

type PrimitiveResource struct {
	data                 interface{}
	primitiveTransformed interface{}
}

func (p *PrimitiveResource) SetData(data interface{}) {
	p.data = data
}

func (p *PrimitiveResource) SetIncludes(includes []string) {
}

func (p *PrimitiveResource) SetTransformer(transformer TransformerInterface) {
}

func (p *PrimitiveResource) GetTransformedData() interface{} {
	return p.primitiveTransformed
}

func (p *PrimitiveResource) SetNestedInclude(includes map[string]interface{}) {
}

func (p *PrimitiveResource) ToJson() (string, error) {
	if !isBasicOrBasicSlice(p.data) {
		return "", fmt.Errorf("not a valid basic type or basic slice")
	}

	return ToJson(p.data)
}

func (p *PrimitiveResource) TransformResource() error {
	if !isBasicOrBasicSlice(p.data) {
		return fmt.Errorf("not a valid basic type or basic slice")
	}
	p.primitiveTransformed = p.data
	return nil
}
