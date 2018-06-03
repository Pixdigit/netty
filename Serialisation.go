package netty

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)


func Serialize(variable interface{}) (string, error) {
	var result string
	var err error

	switch t := variable.(type) {
	case bool:
		if t {
			result = "bool:true"
			} else {
				result = "bool:false"
			}
		case int32:
			result = "int:" + strconv.Itoa(int(t))
		case int:
			result = "int:" + strconv.Itoa(t)
		case float32:
			result = "float:" + strconv.FormatFloat(float64(t), 'G', -1, 64)
		case float64:
			result = "float:" + strconv.FormatFloat(t, 'G', -1, 64)
		case string:
			result = "str:" + t
		default:
			err = errors.New(fmt.Sprintf("Can not serialize %#v of type %T", t, t))
	}
	return result, err
}

func Deserialize(raw string) (interface{}, error) {
	var result interface{}
	raw = strings.TrimLeft(raw, " ")
	varType, err := typeOfSerialized(raw);	if err != nil {return nil, errors.Wrap(err, "could not deserialize \""+raw+"\"")}
	varValue := raw[strings.Index(raw, ":")+1:]

	switch varType {
	case "bool":
		if varValue == "true" {
			result = true
		} else if varValue == "false" {
			result = false
		} else {
			err = errors.New("could not deserialize " + varValue + " as boolean")
		}
	case "int":
		result, err = strconv.Atoi(varValue)
		errors.Wrap(err, "could not deserialize "+varValue+" as int")
	case "float":
		result, err = strconv.ParseFloat(varValue, 64)
		errors.Wrap(err, "could not deserialize "+varValue+" as float")
	case "str":
		result = varValue
	default:
		err = errors.New("unknown var type in deserialisation " + varType)
	}
	return result, err
}

func typeOfSerialized(s string) (string, error) {
	if !strings.Contains(s, ":") {
		return "", errors.New("value is untyped")
	}
	return s[:strings.Index(s, ":")], nil
}
