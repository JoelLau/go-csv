package gocsv

import (
	"regexp"
	"strings"
)

// splits a multi-line string into 2D array
func ReadAll(data []byte) ([][]string, error) {
	val := make([][]string, 0)

	lines := regexp.MustCompile("\r?\n").Split(string(data), -1)
	for _, line := range lines {
		row, err := ReadRow(string(line))
		if err != nil {
			return nil, err
		}

		val = append(val, row)
	}

	return val, nil
}

// TODO: get handle to singleton default row parser
// TODO: handle other delimeters
// TODO: handle escape characters and other crazy quote cases (see https://www.rfc-editor.org/rfc/rfc4180#section-2)3
const DELIMETER = ','

func ReadRow(row string) ([]string, error) {
	cells := make([]string, 0)

	var sb strings.Builder
	for _, r := range row {
		switch r {
		case DELIMETER:
			cells = append(cells, sb.String())
			sb.Reset()
		default:
			sb.WriteRune(r)
		}
	}
	// add whatever's left in the stringbuilder as the final element
	cells = append(cells, sb.String())

	return cells, nil
}
