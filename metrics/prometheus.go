package metrics

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// WritePrometheus writes samples in Prometheus text format.
func WritePrometheus(w io.Writer, samples []Sample) {
	byName := make(map[string][]Sample)
	for _, s := range samples {
		baseName := s.Name
		for _, suffix := range []string{"_bucket", "_sum", "_count"} {
			if strings.HasSuffix(baseName, suffix) {
				baseName = strings.TrimSuffix(baseName, suffix)
				break
			}
		}
		byName[baseName] = append(byName[baseName], s)
	}

	names := make([]string, 0, len(byName))
	for name := range byName {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		for _, sample := range byName[name] {
			writePrometheusSample(w, sample)
		}
	}
}

func writePrometheusSample(w io.Writer, s Sample) {
	var sb strings.Builder
	sb.WriteString(s.Name)

	if s.Labels.Len() > 0 {
		sb.WriteByte('{')
		for i, key := range s.Labels.keys {
			if i > 0 {
				sb.WriteByte(',')
			}
			sb.WriteString(key)
			sb.WriteString(`="`)
			sb.WriteString(escapeLabel(s.Labels.values[i]))
			sb.WriteByte('"')
		}
		sb.WriteByte('}')
	}

	sb.WriteByte(' ')
	sb.WriteString(fmt.Sprintf("%g", s.Value))
	sb.WriteByte('\n')

	w.Write([]byte(sb.String()))
}

func escapeLabel(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	s = strings.ReplaceAll(s, "\n", `\n`)
	return s
}
