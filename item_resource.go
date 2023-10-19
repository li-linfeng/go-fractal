package fractal

import (
	"fmt"
)

type ItemResource struct {
	data        interface{}
	transformed map[string]interface{}
	transformer TransformerInterface
	includes    Includes
}

func (i *ItemResource) SetData(data interface{}) {
	i.data = data
}

func (i *ItemResource) SetTransformer(transformer TransformerInterface) {
	i.transformer = transformer
}

func (i *ItemResource) SetIncludes(includes []string) {
	i.includes.setIncludes(includes)
}

func (i *ItemResource) SetNestedInclude(includes map[string]interface{}) {
	i.includes.setNestedInclude(includes)
}

func (i *ItemResource) checkDataValidate() error {
	if IsSlice(i.data) {
		return fmt.Errorf("input can not be  slice : %T", i.data)
	}

	if !IsMap(i.data) && !IsStruct(i.data) {
		return fmt.Errorf("input must be map or struct : %T", i.data)
	}
	return nil
}

func (i *ItemResource) ToJson() (string, error) {
	if i.transformed == nil {
		if err := i.TransformResource(); err != nil {
			return "", err
		}
	}

	return ToJson(i.transformed)
}

func (i *ItemResource) GetTransformedData() interface{} {
	return i.transformed
}

func (i *ItemResource) TransformResource() error {
	// 如果未设置data 则直接返回空map
	if i.data == nil {
		i.transformed = map[string]interface{}{}
		return nil
	}
	// 校验数据是否正确
	if err := i.checkDataValidate(); err != nil {
		return err
	}
	// 如果未设置transformer 则直接返回data
	if i.transformer == nil {
		return fmt.Errorf("can not find transformer")
	}

	i.transformer.setData(i.data)
	i.includes.ensureInclude()
	i.transformed = i.transformer.Transform()

	// 将include处理成层级状
	for name, value := range i.includes.nestedInclude {
		result, err := i.transformer.CallIncludeFunction(name)
		if err != nil {
			return err
		}
		nestedInclude := value.(map[string]interface{})

		if err = i.processResource(result, nestedInclude, name); err != nil {
			return err
		}
	}
	return nil
}

// 处理资源对象的通用逻辑
func (i *ItemResource) processResource(resource ResourceInterface, nestedInclude map[string]interface{}, name string) error {
	resource.SetNestedInclude(nestedInclude)
	err := resource.TransformResource()
	if err != nil {
		return err
	}
	i.transformed[name] = resource.GetTransformedData()
	return nil
}
