package gonertia

import (
	crypto "crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

func set[T comparable](data []T) map[T]struct{} {
	if len(data) == 0 {
		return nil
	}

	set := make(map[T]struct{}, len(data))
	for _, v := range data {
		set[v] = struct{}{}
	}

	return set
}

func md5(str string) string {
	hash := crypto.Sum([]byte(str))
	return hex.EncodeToString(hash[:])
}

func md5File(path string) (string, error) {
	h := crypto.New()

	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func firstOr[T any](items []T, fallback T) T {
	if len(items) > 0 {
		return items[0]
	}
	return fallback
}
