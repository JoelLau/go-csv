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
	Unmarshal(data []byte) error
}

func setFieldValue(field reflect.Value, strVal string) error {
	// Custom unmarshaler check (including pointer receivers)
	var unmarshalerValue reflect.Value
	if field.Type().Implements(reflect.TypeOf((*CSVUnmarshaller)(nil)).Elem()) {
		unmarshalerValue = field
	} else if field.CanAddr() && reflect.PointerTo(field.Type()).Implements(reflect.TypeOf((*CSVUnmarshaller)(nil)).Elem()) {
		unmarshalerValue = field.Addr()
	}

	if unmarshalerValue.IsValid() {
		if u, ok := unmarshalerValue.Interface().(CSVUnmarshaller); ok {
			if err := u.Unmarshal([]byte(strVal)); err != nil {
				return fmt.Errorf("error parsing csv cell value to custom unmarshaller: %w", err)
			}

			return nil // Important: Stop processing after a custom unmarshaler handles the field.
		}
	}

	switch field.Kind() {
	case reflect.Bool:
		switch strVal {
		case "True", "true", "1":
			field.SetBool(true)
		case "False", "false", "0":
			field.SetBool(false)
		default:
			if strings.TrimSpace(strVal) == "" {
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
