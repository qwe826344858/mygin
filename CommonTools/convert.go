package CommonTools

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Mapping map[string]string
type ProtoMapping map[string][]string

var notDefineErr = errors.New("未定义转化类型")

func MapToTable(tType reflect.Type, tValue reflect.Value, m map[string]string) error {
	for i := 0; i < tType.NumField(); i++ {
		tag := tType.Field(i).Tag.Get("json")
		if tag == "" {
			continue
		}

		name := tType.Field(i).Name
		target := tValue.FieldByName(name)

		err := reflectSetValByInterface(target, m[tag])
		if err != nil {
			return err
		}
	}

	return nil
}

func reflectSetValByInterface(target reflect.Value, in string) error {
	switch target.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, err := strconv.Atoi(in)
		if err != nil {
			return err
		}
		target.SetInt(int64(v))

	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, err := strconv.ParseUint(in, 10, 64)
		if err != nil {
			return err
		}
		target.SetUint(uint64(v))

	case reflect.String:
		target.SetString(string(in[:]))
	default:
		return notDefineErr
	}

	return nil
}

func InitProtoToTable(refType reflect.Type, mapping Mapping) {
	for i := 0; i < refType.NumField(); i++ {
		tag := refType.Field(i).Tag.Get("proto")
		if tag == "" {
			continue
		}

		name := refType.Field(i).Name
		tags := strings.Split(tag, ",")
		for _, t := range tags {
			mapping[t] = name
		}
	}
}

func InitTableToProto(refType reflect.Type, mapping ProtoMapping) {
	for i := 0; i < refType.NumField(); i++ {
		protoTag := refType.Field(i).Tag.Get("proto")
		if protoTag == "" {
			continue
		}

		name := refType.Field(i).Name
		tags := strings.Split(protoTag, ",")
		for _, tag := range tags {
			mapping[name] = append(mapping[name], tag)
		}
	}
}

func MapProtoToTable(mapping Mapping, target, source reflect.Value, stype reflect.Type, prefix string) {
	for i := 0; i < stype.NumField(); i++ {
		sname := stype.Field(i).Name
		if sname == "" {
			continue
		}

		//logger.Debugf("#### sname:%s", sname)
		name, ok := mapping[prefix+sname]
		if !ok {
			continue
		}

		//logger.Debugf("#### name:%s", name)

		t := target.Elem().FieldByName(name)

		v := source.Field(i).Interface()
		if v != reflect.Zero(source.Field(i).Type()).Interface() {
			reflectSetVal(t, source.Field(i))
			continue
		}

		bval := source.FieldByName(sname + "B")
		if !bval.IsValid() || bval.Kind() != reflect.Bool {
			continue
		} else if bval.Bool() {
			reflectSetVal(t, source.Field(i))
		}
	}
}

func MapTableToProto(target, source reflect.Value, mapping ProtoMapping, prefix, replace string) {
	for name, tag := range mapping {
		field := getReflectFieldByPrefix(tag, prefix, replace)
		if len(field) <= 0 {
			continue
		}

		v := reflectIterAndInit(target, field)
		reflectSetVal(v, source.FieldByName(name))
	}
}

// 通过名字前缀匹配获取字段名称
func getReflectFieldByPrefix(tag []string, prefix, replace string) string {
	for _, s := range tag {
		if !strings.HasPrefix(s, prefix) {
			continue
		}
		if len(replace) > 0 {
			return strings.Replace(s, replace, "", 1)
		} else {
			return s[len(prefix):]
		}
	}
	return ""
}

// reflectIterVal 遍历结构, 如果是 nil, 就初始化结构
func reflectIterAndInit(target reflect.Value, field string) reflect.Value {
	var val reflect.Value

	str := strings.Split(field, ".")
	for _, s := range str {
		val = target.FieldByName(s)
		switch val.Kind() {
		case reflect.Ptr, reflect.Array, reflect.Slice, reflect.Map:
			if val.IsNil() { // 空指针就初始化数据
				tmp := val.Interface()
				vtype := reflect.TypeOf(tmp).Elem()
				val.Set(reflect.New(vtype))
			}
		}
	}

	return val
}

func reflectSetVal(target, source reflect.Value) error {
	//logger.Debugf("#### target:%v, source:%v", target, source)
	switch target.Kind() {
	case reflect.Bool:
		if source.Kind() == reflect.Bool {
			target.SetBool(source.Bool())
			return nil
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch source.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			target.SetInt(source.Int())
			return nil

		case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			target.SetInt(int64(source.Uint()))
			return nil

		case reflect.String:
			if val, err := strconv.Atoi(source.String()); err == nil {
				target.SetInt(int64(val))
				return nil
			}
		}

	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch source.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			target.SetUint(uint64(source.Int()))
			return nil

		case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			target.SetUint(source.Uint())
			return nil

		case reflect.String:
			if val, err := strconv.Atoi(source.String()); err == nil {
				target.SetUint(uint64(val))
				return nil
			}
		}

	case reflect.String:
		if source.Kind() == reflect.String {
			target.SetString(source.String())
			return nil
		}
	}
	return nil
}

// 根据字段名取结构体字符串字段值
func GetStructStringField(input interface{}, key string) (value string, err error) {
	v, err := GetStructField(input, key)
	if err != nil {
		return
	}
	value, ok := v.(string)
	if !ok {
		return value, errors.New("can't convert key'v to string")
	}
	return
}

// 根据字段名取结构体整型字段值
func GetStructIntField(input interface{}, key string) (value int64, err error) {
	v, err := GetStructField(input, key)
	if err != nil {
		return
	}
	value, ok := v.(int64)
	if !ok {
		return value, errors.New("can't convert key'v to int64")
	}
	return
}

// 获取结构体特定字段值,必须传入struct，不能是指针
func GetStructField(inputStruct interface{}, filedName string) (interface{}, error) {
	rv := reflect.ValueOf(inputStruct)
	rt := reflect.TypeOf(inputStruct)
	var value interface{}
	if rt.Kind() != reflect.Struct {
		return value, errors.New("input must be struct")
	}
	filedExist := false
	for i := 0; i < rt.NumField(); i++ {
		curField := rv.Field(i)
		if rt.Field(i).Name == filedName {
			switch curField.Kind() {
			case reflect.Bool,
				reflect.Int,
				reflect.Int8,
				reflect.Int16,
				reflect.Int32,
				reflect.Int64,
				reflect.Uint,
				reflect.Uint8,
				reflect.Uint16,
				reflect.Uint32,
				reflect.Uint64,
				reflect.Uintptr,
				reflect.Float32,
				reflect.Float64,
				reflect.Complex64,
				reflect.Complex128,
				// reflect.Array,
				// reflect.Chan,
				// reflect.Func,
				reflect.Interface,
				reflect.Map,
				// reflect.Ptr,
				reflect.Slice,
				reflect.String:
				filedExist = true
				value = curField.Interface()
			default:
				return value, errors.New("filedName:" + filedName + " filedType not support")
			}
		}
	}
	if !filedExist {
		return value, fmt.Errorf("filedName:%s not found in %s's field", filedName, rt)
	}
	return value, nil
}
