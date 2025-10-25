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
	var sb strings.Builder // holds contents of the field we're currently building

	for ; index < len(runes); index++ {
		curr := runes[index]

		if !isWithinQuotes && sb.Len() == 0 && curr == '"' {
			isWithinQuotes = true
			continue
		}

		if isWithinQuotes && curr == '"' {
			// if field is quoted, replace 2 double quotes with 1 double quote (i.e. escape them)
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

			// skip the first space after commas - some programs export with ", " as a delimeter
			if (index+1 < len(runes)) && (runes[index+1] == ' ') {
				index++
			}
			continue
		}

		sb.WriteRune(curr)
	}

	// put whatever's left in "buffers" (sb, row) to the output if its not an empty row - we ignore those
	row = append(row, sb.String())
	if strings.TrimSpace(strings.Join(row, "")) != "" {
		row = append(row, sb.String())
		rows = append(rows, row)
	}

	return rows, nil
}
