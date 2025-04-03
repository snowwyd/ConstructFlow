package utils

import (
	"path/filepath"
	"regexp"
	"strings"
)

func ParseBaseName(originalKey string) (string, string) {
	ext := filepath.Ext(originalKey)
	base := strings.TrimSuffix(originalKey, ext)

	// Удаляем все суффиксы версий
	re := regexp.MustCompile(`_v\d+$`)
	base = re.ReplaceAllString(base, "")

	return base, ext
}
