package jsonschema

func New() JsonSchema {
	return JsonSchema{
		Type:       "object",
		Properties: map[string]*JsonSchema{},
	}
}

type JsonSchema struct {
	Type        string                 `json:"type"`
	Description string                 `json:"description,omitempty"`
	Properties  map[string]*JsonSchema `json:"properties,omitempty"`
	Required    []string               `json:"required,omitempty"`
}

func (j *JsonSchema) AddString(name string, description string, required bool) {
	if j.Type != "object" {
		j.Type = "object"
	}

	if j.Properties == nil {
		j.Properties = map[string]*JsonSchema{}
	}

	j.Properties[name] = &JsonSchema{
		Type:        "string",
		Description: description,
	}
	if required {
		j.Required = append(j.Required, name)
	}
}

func (j *JsonSchema) AddNumber(name string, description string, required bool) {
	if j.Type != "object" {
		j.Type = "object"
	}

	if j.Properties == nil {
		j.Properties = map[string]*JsonSchema{}
	}

	j.Properties[name] = &JsonSchema{
		Type:        "number",
		Description: description,
	}
	if required {
		j.Required = append(j.Required, name)
	}
}

func (j *JsonSchema) AddObject(name string, description string, required bool) *JsonSchema {
	if j.Type != "object" {
		j.Type = "object"
	}

	if j.Properties == nil {
		j.Properties = map[string]*JsonSchema{}
	}

	schema := &JsonSchema{
		Type:        "object",
		Description: description,
		Properties:  map[string]*JsonSchema{},
	}

	if required {
		j.Required = append(j.Required, name)
	}

	j.Properties[name] = schema

	return schema
}

func (j JsonSchema) Generate() any {
	switch j.Type {
	case "string":
		return "foo"
	case "number":
		return float64(10)
	case "":
		fallthrough
	case "object":
		//recurse
		result := map[string]any{}
		for name, prop := range j.Properties {
			result[name] = prop.Generate()
		}
		return result
	default:
		return nil
	}
}
