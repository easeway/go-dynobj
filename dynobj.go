package dynobj

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Query interface {
	Query(path []string, defVal interface{}) (interface{}, error)
}

type DynQuery struct {
	Object interface{}
}

func (q *DynQuery) Query(path []string, defVal interface{}) (interface{}, error) {
	depth := len(path)
	obj := q.Object
	for d := 0; d < depth; d++ {
		switch obj.(type) {
		case map[string]interface{}:
			if val, present := obj.(map[string]interface{})[path[d]]; present {
				obj = val
			} else {
				return defVal, notFound(path, d)
			}
		case []interface{}:
			index, err := strconv.ParseUint(path[d], 0, 0)
			if err != nil {
				return defVal, err
			}
			arr := obj.([]interface{})
			if index < uint64(len(arr)) {
				obj = arr[index]
			} else {
				return defVal, notFound(path, d)
			}
		default:
			return defVal, notFound(path, d)
		}
	}
	return obj, nil
}

func NewJsonQuery(reader io.Reader) (*DynQuery, error) {
	q := &DynQuery{}
	err := json.NewDecoder(reader).Decode(&q.Object)
	return q, err
}

func NewJsonStringQuery(json string) (*DynQuery, error) {
	return NewJsonQuery(bytes.NewBufferString(json))
}

func notFound(path []string, dep int) error {
	return errors.New("Not found: " + strings.Join(path[0:dep+1], "/"))
}

func ParsePath(path string) []string {
	tokens := strings.Split(path, "/")
	result := make([]string, 0, len(tokens))
	for _, v := range tokens {
		if len(v) > 0 {
			result = append(result, v)
		}
	}
	return result
}

type DynObj struct {
	Query Query
}

func NewJsonObj(reader io.Reader) (*DynObj, error) {
	if q, err := NewJsonQuery(reader); err == nil {
		return &DynObj{Query: q}, err
	} else {
		return nil, err
	}
}

func NewJsonStringObj(json string) (*DynObj, error) {
	return NewJsonObj(bytes.NewBufferString(json))
}

func (o *DynObj) AsAnyD(path string, defVal interface{}) interface{} {
	val, _ := o.Query.Query(ParsePath(path), defVal)
	return val
}

func (o *DynObj) AsStrD(path string, defVal string) string {
	if val, err := o.Query.Query(ParsePath(path), defVal); err != nil {
		return defVal
	} else {
		return fmt.Sprintf("%v", val)
	}
}

func (o *DynObj) AsIntD(path string, defVal int) int {
	if val, err := o.Query.Query(ParsePath(path), defVal); err != nil {
		return defVal
	} else {
		switch val.(type) {
		case int, int8, int16, int32, int64,
			uint, uint8, uint16, uint32, uint64:
			return val.(int)
		case float32, float64:
			return int(val.(float64))
		case string:
			if val64, err := strconv.ParseInt(val.(string), 0, 0); err != nil {
				return defVal
			} else {
				return int(val64)
			}
		default:
			return defVal
		}
	}
}

func (o *DynObj) AsBoolD(path string, defVal bool) bool {
	if val, err := o.Query.Query(ParsePath(path), defVal); err != nil {
		return defVal
	} else {
		switch val.(type) {
		case bool:
			return val.(bool)
		case int, int8, int16, int32, int64,
			uint, uint8, uint16, uint32, uint64:
			return val.(int64) != 0
		case float32, float64:
			return val.(float64) != 0
		case string:
			strVal := strings.ToLower(val.(string))
			if strVal == "true" || strVal == "yes" {
				return true
			} else if strVal == "false" || strVal == "no" {
				return false
			} else {
				return defVal
			}
		default:
			return defVal
		}
	}
}

func (o *DynObj) AsAny(path string) interface{} {
	return o.AsAnyD(path, nil)
}

func (o *DynObj) AsStr(path string) string {
	return o.AsStrD(path, "")
}

func (o *DynObj) AsInt(path string) int {
	return o.AsIntD(path, 0)
}

func (o *DynObj) AsBool(path string) bool {
	return o.AsBoolD(path, false)
}
