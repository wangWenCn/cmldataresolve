package cmldataresolve

import (
	"strconv"
	"strings"
)

func lookupPath(data map[string]any, path string) (any, bool) {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil, false
	}
	if value, ok := data[path]; ok {
		return value, true
	}
	parts := splitPath(path)
	if len(parts) == 0 {
		return nil, false
	}
	var current any = data
	for _, part := range parts {
		switch c := current.(type) {
		case map[string]any:
			next, ok := c[part]
			if !ok {
				return nil, false
			}
			current = next
		case []any:
			idx, err := strconv.Atoi(part)
			if err != nil || idx < 0 || idx >= len(c) {
				return nil, false
			}
			current = c[idx]
		default:
			return nil, false
		}
	}
	return current, true
}

func splitPath(path string) []string {
	var out []string
	var b strings.Builder
	for i := 0; i < len(path); i++ {
		ch := path[i]
		switch ch {
		case '.':
			if b.Len() > 0 {
				out = append(out, b.String())
				b.Reset()
			}
		case '[':
			if b.Len() > 0 {
				out = append(out, b.String())
				b.Reset()
			}
			j := i + 1
			for ; j < len(path) && path[j] != ']'; j++ {
			}
			if j <= len(path) {
				out = append(out, strings.TrimSpace(path[i+1:j]))
				i = j
			}
		default:
			b.WriteByte(ch)
		}
	}
	if b.Len() > 0 {
		out = append(out, b.String())
	}
	return out
}
