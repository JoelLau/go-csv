package gocsv

import (
	"strings"
)

func ReadAll(b []byte) ([][]string, error) {
	rows := [][]string{}

	index := 0
	runes := []rune(string(b))
	isWithinQuotes := false

	row := []string{}
	var sb strings.Builder

	for ; index < len(runes); index++ {
		curr := runes[index]

		if !isWithinQuotes && sb.Len() == 0 && curr == '"' {
			isWithinQuotes = true
			continue
		}

		if isWithinQuotes && curr == '"' {
			// escape consecutive double quotes
			if (index+1 < len(runes)) && (runes[index+1] == '"') {
				sb.WriteRune(curr)
				index++
				continue
			}

			isWithinQuotes = false
			continue
		}

		if !isWithinQuotes && curr == '\n' {
			row = append(row, sb.String())
			sb.Reset()

			rows = append(rows, row)
			row = []string{}

			continue
		}

		if !isWithinQuotes && curr == ',' {
			row = append(row, sb.String())
			sb.Reset()

			isWithinQuotes = false
			continue
		}

		sb.WriteRune(curr)
	}

	if sb.String() != "" && len(row) != 0 {
		row = append(row, sb.String())
		rows = append(rows, row)
	}

	return rows, nil
}
