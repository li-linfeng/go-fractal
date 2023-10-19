package fractal

type CollectionResource struct {
	sliceData        []interface{}
	sliceTransformed []map[string]interface{}
	transformer      TransformerInterface
	includes         Includes
	err              error
}

func (c *CollectionResource) SetData(data interface{}) {
	sliceData, err := ConvertToInterfaceSlice(data)
	c.err = err
	c.sliceData = sliceData
}

func (c *CollectionResource) SetTransformer(transformer TransformerInterface) {
	c.transformer = transformer
}

func (c *CollectionResource) SetIncludes(includes []string) {
	c.includes.setIncludes(includes)
}

func (c *CollectionResource) SetNestedInclude(includes map[string]interface{}) {
	c.includes.setNestedInclude(includes)
}

func (c *CollectionResource) ToJson() (string, error) {
	if c.sliceTransformed == nil {
		if err := c.TransformResource(); err != nil {
			return "", err
		}
	}
	return ToJson(c.sliceTransformed)
}

func (c *CollectionResource) TransformResource() error {
	// 判断是否存在异常
	if c.err != nil {
		return c.err
	}

	c.includes.ensureInclude()
	c.sliceTransformed = []map[string]interface{}{}
	for _, value := range c.sliceData {
		i := &ItemResource{}
		i.SetData(value)
		// 设置 transformer
		i.SetTransformer(c.transformer)
		i.includes.setNestedInclude(c.includes.nestedInclude)
		if err := i.TransformResource(); err != nil {
			return err
		}
		c.sliceTransformed = append(c.sliceTransformed, i.transformed)
	}

	return nil
}

func (c *CollectionResource) GetTransformedData() interface{} {
	return c.sliceTransformed
}
