package jsonschema

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSchema(t *testing.T) {
	schema := New()
	schema.AddString("name", "The users name", true)
	schema.AddNumber("age", "The users age", false)
	prefs := schema.AddObject("preferences", "The users preferences", false)
	prefs.AddString("orientation", "landscape or portrait", true)

	json, err := json.MarshalIndent(schema, "", " ")
	assert.Nil(t, err)

	assert.Equal(t, string(json), expectedSchema)
}

func TestGenerate(t *testing.T) {
	schema := New()
	schema.AddString("name", "The users name", true)
	schema.AddNumber("age", "The users age", false)
	prefs := schema.AddObject("preferences", "The users preferences", false)
	prefs.AddString("orientation", "landscape or portrait", true)

	json, err := json.MarshalIndent(schema.Generate(), "", " ")
	assert.Nil(t, err)

	assert.Equal(t, string(json), expectedGeneratedObject)
}

const expectedGeneratedObject = `{
 "age": 10,
 "name": "foo",
 "preferences": {
  "orientation": "foo"
 }
}`

const expectedSchema = `{
 "type": "object",
 "properties": {
  "age": {
   "type": "number",
   "description": "The users age"
  },
  "name": {
   "type": "string",
   "description": "The users name"
  },
  "preferences": {
   "type": "object",
   "description": "The users preferences",
   "properties": {
    "orientation": {
     "type": "string",
     "description": "landscape or portrait"
    }
   },
   "required": [
    "orientation"
   ]
  }
 },
 "required": [
  "name"
 ]
}`
