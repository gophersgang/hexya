// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package po

import (
	"bytes"
	"strings"
)

func decodePoString(text string) string {
	lines := strings.Split(text, "\n")
	for i := 0; i < len(lines); i++ {
		left := strings.Index(lines[i], `"`)
		right := strings.LastIndex(lines[i], `"`)
		if left < 0 || right < 0 || left == right {
			lines[i] = ""
			continue
		}
		line := lines[i][left+1 : right]
		data := make([]byte, 0, len(line))
		for i := 0; i < len(line); i++ {
			if line[i] != '\\' {
				data = append(data, line[i])
				continue
			}
			if i+1 >= len(line) {
				break
			}
			switch line[i+1] {
			case 'n': // \\n -> \n
				data = append(data, '\n')
				i++
			case 't': // \\t -> \n
				data = append(data, '\t')
				i++
			case '\\': // \\\ -> ?
				data = append(data, '\\')
				i++
			}
		}
		lines[i] = string(data)
	}
	return strings.Join(lines, "")
}

func encodePoString(text string) string {
	if text == "" {
		return "\"\"\n"
	}
	var buf bytes.Buffer
	lines := strings.Split(text, "\n")
	for i := 0; i < len(lines); i++ {
		if lines[i] == "" {
			if i != len(lines)-1 {
				buf.WriteString(`"\n"` + "\n")
			}
			continue
		}
		if i == 0 && len(lines) > 1 {
			buf.WriteString("\"\"\n")
		}
		buf.WriteRune('"')
		for _, r := range lines[i] {
			switch r {
			case '\\':
				buf.WriteString(`\\`)
			case '"':
				buf.WriteString(`\"`)
			case '\n':
				buf.WriteString(`\n`)
			case '\t':
				buf.WriteString(`\t`)
			default:
				buf.WriteRune(r)
			}
		}
		if i != len(lines)-1 {
			buf.WriteString(`\n`)
		}
		buf.WriteString("\"\n")
	}
	return buf.String()
}

func encodeCommentPoString(text string) string {
	var buf bytes.Buffer
	lines := strings.Split(text, "\n")
	if len(lines) > 1 {
		buf.WriteString(`""` + "\n")
	}
	for i := 0; i < len(lines); i++ {
		if len(lines) > 0 {
			buf.WriteString("#| ")
		}
		buf.WriteRune('"')
		for _, r := range lines[i] {
			switch r {
			case '\\':
				buf.WriteString(`\\`)
			case '"':
				buf.WriteString(`\"`)
			case '\n':
				buf.WriteString(`\n`)
			case '\t':
				buf.WriteString(`\t`)
			default:
				buf.WriteRune(r)
			}
		}
		if i < len(lines)-1 {
			buf.WriteString(`\n"` + "\n")
		} else {
			buf.WriteString(`"`)
		}
	}
	return buf.String()
}
