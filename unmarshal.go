package gocsv

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const StructTagCSV = "csv"

// TODO: consolidate type checks
func Unmarshal(data []byte, v any) error {
	// reflect value of `v` (expecting `&S[]`)
	vVal := reflect.ValueOf(v)
	if vVal.Kind() != reflect.Pointer {
		return fmt.Errorf("`v` must be pointer")
	}

	// expecting `S[]`
	vRef := vVal.Elem()
	if vRef.Kind() != reflect.Slice {
		return fmt.Errorf("`v` must be pointer to slice")
	}
	vRefType := vRef.Type()

	// expecting `S`
	vRefElem := vRefType.Elem()
	if vRefElem.Kind() != reflect.Struct {
		return fmt.Errorf("`v` must be pointer to slice of struct types")
	}

	m, err := newHeaderToIndexMap(v)
	if err != nil {
		return err
	}

	lines, err := ReadAll(data)
	if err != nil {
		return err
	}

	// assume first line contains headers
	var headers []string
	if len(lines) >= 1 {
		headers = lines[0]
	}

	// key: csv column index, value: struct field index
	mm := map[int]int{}
	for idx, h := range headers {
		if i, ok := m[h]; ok {
			mm[idx] = i
		}
	}

	var rows [][]string
	if len(lines) >= 1 {
		rows = lines[1:]
	}

	newSlice := reflect.MakeSlice(vRefType, 0, len(rows))

	for _, row := range rows {
		if len(row) <= 0 {
			continue
		}

		newElem := reflect.New(vRefElem).Elem()

		// iterate through map so we skip csv columns that don't need parsing
		for columnIndex, fieldIndex := range mm {
			field := newElem.Field(fieldIndex)
			rawVal := row[columnIndex]

			if err := setFieldValue(field, rawVal); err != nil {
				return err
			}
		}

		newSlice = reflect.Append(newSlice, newElem)
	}

	vRef.Set(newSlice)
	return nil
}

type CSVUnmarshaller interface {
	UnmarshalCSV(data []byte) error
}

func setFieldValue(field reflect.Value, strVal string) error {
	csvUnmarshallerType := reflect.TypeOf((*CSVUnmarshaller)(nil)).Elem()

	fieldType := field.Type()

	if reflect.PointerTo(fieldType).Implements(csvUnmarshallerType) {
		if !field.CanAddr() {
			return fmt.Errorf("field %s is not addressable and cannot be unmarshalled with a pointer receiver", fieldType.Name())
		}

		fieldPtr := field.Addr()

		if unmarshaller, ok := fieldPtr.Interface().(CSVUnmarshaller); ok {
			return unmarshaller.UnmarshalCSV([]byte(strVal))
		}
	}

	switch field.Kind() {
	case reflect.Bool:
		s := strings.ToLower(strings.TrimSpace(strVal))
		switch s {
		case "true", "1":
			field.SetBool(true)
		case "false", "0":
			field.SetBool(false)
		default:
			if s == "" {
				field.SetBool(reflect.Zero(field.Type()).Bool())
				break
			}
			return fmt.Errorf("error attempting to set value '%+v' to bool", strVal)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, err := strconv.ParseUint(strVal, 10, 64)
		if err != nil {
			if strings.TrimSpace(strVal) == "" {
				field.SetUint(reflect.Zero(field.Type()).Uint())
				break
			}
			return fmt.Errorf("error parsing csv cell value to uint: %+v", err)
		}

		field.SetUint(i)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(strVal, 10, 64)
		if err != nil {
			if strings.TrimSpace(strVal) == "" {
				field.SetInt(reflect.Zero(field.Type()).Int())
				break
			}
			return fmt.Errorf("error parsing csv cell value to integer: %+v", err)
		}

		field.SetInt(i)
	case reflect.Float32:
		i, err := strconv.ParseFloat(strVal, 32)
		if err != nil {
			if strings.TrimSpace(strVal) == "" {
				field.SetFloat(reflect.Zero(field.Type()).Float())
				break
			}
			return fmt.Errorf("error parsing csv cell value to float32: %+v", err)
		}

		field.SetFloat(i)
	case reflect.Float64:
		i, err := strconv.ParseFloat(strVal, 64)
		if err != nil {
			if strings.TrimSpace(strVal) == "" {
				field.SetFloat(reflect.Zero(field.Type()).Float())
				break
			}
			return fmt.Errorf("error parsing csv cell value to float64: %+v", err)
		}

		field.SetFloat(i)
	case reflect.String:
		field.SetString(strVal)
	default:
		field.Set(reflect.Zero(field.Type()))
		return fmt.Errorf("could not find parser for type: %+v", field.Kind())
	}

	return nil
}

// func setFieldValue(field reflect.Value, strVal string) error {

// 	csvUnmarshallerType := reflect.TypeOf((*CSVUnmarshaller)(nil)).Elem()
// 	fieldType := reflect.TypeOf(field)
// 	if reflect.PointerTo(fieldType).Implements(csvUnmarshallerType) {
// 		// Check if the field's value is addressable. This is necessary
// 		// to get a pointer to it.
// 		if !field.CanAddr() {
// 			return fmt.Errorf("field %s is not addressable and cannot be unmarshalled with a pointer receiver", fieldType.Name())
// 		}

// 		// Get the address of the field.
// 		fieldPtr := field.Addr()

// 		// Perform the type assertion on the pointer.
// 		if unmarshaller, ok := fieldPtr.Interface().(CSVUnmarshaller); ok {
// 			if err := unmarshaller.UnmarshalCSV([]byte(strVal)); err != nil {
// 				return err
// 			}
// 		}
// 	}

// 	return nil
// }

// build a map with (key: csv header), (val: struct field index)
func newHeaderToIndexMap(v any) (map[string]int, error) {
	vType := reflect.TypeOf(v)
	if vType.Kind() != reflect.Pointer {
		return nil, fmt.Errorf("`v` must be pointer")
	}

	vRef := vType.Elem()
	if vRef.Kind() != reflect.Slice {
		return nil, fmt.Errorf("`v` must be pointer to slice")
	}

	vRefElem := vRef.Elem()
	if vRefElem.Kind() != reflect.Struct {
		return nil, fmt.Errorf("`v` must be pointer to slice of struct types")
	}

	m := map[string]int{}
	for i := range vRefElem.NumField() {
		field := vRefElem.Field(i)
		m[field.Tag.Get(StructTagCSV)] = i
	}

	return m, nil
}
