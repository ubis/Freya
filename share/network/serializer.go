package network

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
)

// struct tags
const (
	tagLength = "ll"
)

// serialize with specified endian
var endian binary.ByteOrder = binary.LittleEndian

// Serialize converts a Go value to a byte slice.
// It uses reflection for complex types with special serialization rules.
func Serialize(v interface{}) ([]byte, error) {
	if needsReflection(v) {
		return serializeWithReflection(v)
	}
	return serializeWithoutReflection(v)
}

// needsReflection checks if a value requires reflection-based serialization.
// It returns true for structs with fields tagged for special handling.
func needsReflection(v interface{}) bool {
	valueType := reflect.TypeOf(v)
	if valueType.Kind() != reflect.Struct {
		return false
	}

	for i := 0; i < valueType.NumField(); i++ {
		if valueType.Field(i).Tag.Get(tagLength) != "" {
			return true
		}
	}
	return false
}

// serializeWithReflection handles serialization for complex structs.
// It respects tagLength tags for custom serialization logic.
// serializeWithReflection handles complex struct serialization.
func serializeWithReflection(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Struct {
		return nil, fmt.Errorf("only structs supported; got %s", value.Kind())
	}

	typeLengthFields := make(map[string]int)

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldInfo := value.Type().Field(i)
		tag := fieldInfo.Tag.Get(tagLength)

		if tag == "" {
			err := binary.Write(&buf, endian, field.Interface())
			if err != nil {
				return nil, fmt.Errorf("failed to serialize field %s: %w",
					fieldInfo.Name, err)
			}

			continue
		}

		lenValue, err := getLengthFieldValue(value, tag, typeLengthFields)
		if err != nil {
			return nil, err
		}

		if lenValue > field.Len() {
			return nil, fmt.Errorf("length greater than slice length")
		}

		for j := 0; j < lenValue; j++ {
			err := binary.Write(&buf, endian, field.Index(j).Interface())
			if err != nil {
				return nil, fmt.Errorf("failed to serialize slice: %w", err)
			}
		}
	}

	return buf.Bytes(), nil
}

// getLengthFieldValue retrieves the length field's value for serialization.
func getLengthFieldValue(value reflect.Value, tag string,
	typeLengthFields map[string]int) (int, error) {

	if lenValue, ok := typeLengthFields[tag]; ok {
		return lenValue, nil
	}

	lenFieldValue := value.FieldByName(tag)
	if !lenFieldValue.IsValid() || lenFieldValue.Kind() != reflect.Uint8 {
		return 0, fmt.Errorf("invalid length field: %s", tag)
	}

	lenValue := int(lenFieldValue.Uint())
	typeLengthFields[tag] = lenValue
	return lenValue, nil
}

// serializeWithoutReflection handles serialization for simple types.
func serializeWithoutReflection(v interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, endian, v); err != nil {
		return nil, fmt.Errorf("serialization failed: %w", err)
	}

	return buf.Bytes(), nil
}

// Deserialize converts a byte slice back into a Go value.
// It supports both simple types and complex structs with special rules.
func Deserialize(r *Reader, v interface{}) error {
	buf := bytes.NewReader(r.buffer)
	value := reflect.ValueOf(v)

	if value.Kind() != reflect.Ptr || value.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("only pointers to structs supported; got %T", v)
	}

	value = value.Elem()
	valueType := value.Type()

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldInfo := valueType.Field(i)

		if fieldInfo.Tag.Get(tagLength) == "" {
			err := binary.Read(buf, endian, field.Addr().Interface())
			if err != nil {
				return fmt.Errorf("failed to deserialize field %s: %w",
					fieldInfo.Name, err)
			}

			continue
		}

		length, err := readLengthField(buf, value, fieldInfo)
		if err != nil {
			return err
		}

		slice, err := deserializeSlice(buf, field, length)
		if err != nil {
			return err
		}

		field.Set(slice)
	}

	return nil
}

// readLengthField reads the length of a variable-length field.
func readLengthField(buf *bytes.Reader, value reflect.Value,
	fieldInfo reflect.StructField) (int, error) {

	lenField := value.FieldByName(fieldInfo.Tag.Get(tagLength))

	if !lenField.IsValid() || lenField.Kind() != reflect.Uint8 {
		return 0, fmt.Errorf("invalid length field: %s", fieldInfo.Tag.Get("ml"))
	}

	return int(lenField.Uint()), nil
}

// deserializeSlice deserializes a slice of bytes.
func deserializeSlice(buf *bytes.Reader, field reflect.Value, length int) (
	reflect.Value, error) {

	slice := reflect.MakeSlice(field.Type(), length, length)

	for j := 0; j < length; j++ {
		err := binary.Read(buf, endian, slice.Index(j).Addr().Interface())
		if err != nil {
			return reflect.Value{},
				fmt.Errorf("failed to deserialize slice: %w", err)
		}
	}

	return slice, nil
}
