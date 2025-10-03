package gocsv

import (
	"fmt"
	"reflect"
	"strconv"
)

const StructTagCSV = "csv"

func Unmarshal(data []byte, v any) error {
	// reflect type of `v` (expecting `&S[]`)
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

	m, err := StructTagFieldIndexMap(v)
	if err != nil {
		return err
	}

	lines, err := ReadAll(data)
	if err != nil {
		return err
	}

	// assume first line csontains headers
	headers := lines[0]

	// key: csv column index, value: struct field index
	mm := map[int]int{}
	for idx, h := range headers {
		mm[idx] = m[h]
	}

	rows := lines[1:]

	newSlice := reflect.MakeSlice(vRefType, 0, len(rows))

	for _, row := range rows {

		newElem := reflect.New(vRefElem).Elem()

		// iterate through map so we skip csv columns that don't need parsing
		for columnIndex, fieldIndex := range mm {
			field := newElem.Field(fieldIndex)
			rawVal := row[columnIndex]

			switch field.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				i, err := strconv.ParseInt(rawVal, 10, 64)
				if err != nil {
					return fmt.Errorf("error parsing csv cell value to integer: %+v", err)
				}

				field.SetInt(i)
			case reflect.Float32:
				i, err := strconv.ParseFloat(rawVal, 32)
				if err != nil {
					return fmt.Errorf("error parsing csv cell value to float32: %+v", err)
				}

				field.SetFloat(i)
			case reflect.Float64:
				i, err := strconv.ParseFloat(rawVal, 64)
				if err != nil {
					return fmt.Errorf("error parsing csv cell value to float64: %+v", err)
				}

				field.SetFloat(i)
			case reflect.String:
				field.SetString(rawVal)
			default:
			}
		}

		newSlice = reflect.Append(newSlice, newElem)
	}

	vRef.Set(newSlice)
	return nil
}

// build a map with (key: csv header), (val: struct field index)
func StructTagFieldIndexMap(v any) (map[string]int, error) {
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
