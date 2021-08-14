package helpers

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
)

/**
NullableFieldsToStruct fills all non-nil fields from source to target
*/
func NullableFieldsToStruct(source, target interface{}) (isFound bool, err error) {
	if reflect.TypeOf(target).Kind() != reflect.Ptr {
		return false, fmt.Errorf("target is not a pointer")
	}
	if reflect.TypeOf(source).Kind() != reflect.Struct {
		return false, fmt.Errorf("source is not a struct")
	}
	if reflect.TypeOf(target).Elem().Kind() != reflect.Struct {
		return false, fmt.Errorf("target is not a pointer to struct")
	}

	targetFields := reflect.ValueOf(target)
	sourceFields := reflect.TypeOf(source)
	sourceValues := reflect.ValueOf(source)
	sourceFieldsNum := sourceFields.NumField()
	isFound = false
	for i := 0; i < sourceFieldsNum; i++ {
		field := sourceFields.Field(i)
		value := sourceValues.Field(i)
		if value.Kind() == reflect.Ptr && !value.IsNil() { //applicable only for non-nil pointers
			f := targetFields.Elem().FieldByName(field.Name)
			if !f.IsValid() {
				continue
			}

			isFound = true
			if !f.CanSet() {
				//TODO continue may be (for private fields)?
				return isFound, fmt.Errorf("cannot set given field `%s`(%s) with old value `%s`",
					field.Name, f.Kind(), f.String())
			}

			if f.Kind() == reflect.Ptr {
				continue
			}

			if f.Kind() != field.Type.Elem().Kind() {
				return isFound, fmt.Errorf("bad type of field `%s`: expect %s, but have %s",
					field.Name, field.Type.Elem().Kind().String(), f.Kind().String())
			}

			//simple Set isn't used because of custom types (aliases)
			switch f.Kind() {
			case reflect.String:
				f.SetString(value.Elem().String())
			case reflect.Bool:
				f.SetBool(value.Elem().Bool())
			case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int8, reflect.Int16:
				f.SetInt(value.Elem().Int())
			case reflect.Uint, reflect.Uint32, reflect.Uint64, reflect.Uint8, reflect.Uint16:
				f.SetUint(value.Elem().Uint())
			case reflect.Float64, reflect.Float32:
				f.SetFloat(value.Elem().Float())
			case reflect.Struct, reflect.Map:
				f.Set(value.Elem())
			case reflect.Array, reflect.Slice:
				f.Set(value.Elem())
			default:
				return isFound, fmt.Errorf("unknown type of field `%s`: %s",
					field.Name, field.Type.Elem().Kind().String())
			}
		}
	}

	return isFound, nil
}

func SetTagsSqlTypeShard(shardId int64, target interface{}) (interface{}, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()
	if reflect.TypeOf(target).Kind() != reflect.Ptr {
		return nil, fmt.Errorf("target is not a pointer")
	}

	if reflect.TypeOf(target).Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("target is not a pointer to struct")
	}

	targetFields := reflect.TypeOf(target).Elem()
	targetFieldsNum := targetFields.NumField()
	var fs []reflect.StructField
	for i := 0; i < targetFieldsNum; i++ {
		field := targetFields.Field(i)
		tag := field.Tag
		sqlTag := tag.Get("sql")
		if len(sqlTag) > 0 {
			var re = regexp.MustCompile(`type:\s*(\?SHARD)`)
			s := re.ReplaceAllString(sqlTag, "type: shard"+strconv.FormatInt(shardId, 10))
			field.Tag = reflect.StructTag(s)
		}
		fs = append(fs, field)
	}

	r := reflect.ValueOf(target).Elem().Convert(reflect.StructOf(fs))
	return r.Interface(), nil
}
