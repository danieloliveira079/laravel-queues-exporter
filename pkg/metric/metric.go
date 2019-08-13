package metric

import (
	"fmt"
	"strings"
)

type Metric struct {
	Name  string
	Value int64
}

func (m *Metric) WithPrefix(prefix string) string {
	parsedName := strings.ReplaceAll(m.Name, ":", ".")
	return fmt.Sprintf("%s.%s", prefix, parsedName)
}

func (m *Metric) ValueToFloat64() float64 {
	return float64(m.Value)
}
