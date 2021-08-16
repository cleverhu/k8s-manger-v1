package core

import (
	"fmt"
	"sort"
)

func GetLabels(m map[string]string) string {
	keys := make([]string, 0)
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	labels := ""
	for _, key := range keys {
		if labels != "" {
			labels += ","
		}
		labels += fmt.Sprintf("%s=%s", key, m[key])
	}

	return labels
}
