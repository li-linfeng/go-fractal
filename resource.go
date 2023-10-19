package fractal

type ResourceInterface interface {
	SetData(data interface{})
	SetIncludes(includes []string)
	SetNestedInclude(map[string]interface{})
	SetTransformer(transformer TransformerInterface)

	GetTransformedData() interface{}
	TransformResource() error
	ToJson() (string, error)
}

type Includes struct {
	includes      []string
	nestedInclude map[string]interface{}
}

func (i *Includes) setIncludes(includes []string) {
	i.includes = includes
	i.nestedInclude = BuildNestedMap(includes)
}

func (i *Includes) setNestedInclude(nestedInclude map[string]interface{}) {
	i.nestedInclude = nestedInclude
}

func (i *Includes) ensureInclude() {
	if i.nestedInclude == nil && len(i.includes) > 0 {
		i.nestedInclude = BuildNestedMap(i.includes)
	}
}
