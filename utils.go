package gonertia

import (
	crypto "crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

// set is a helper that creates a map whose keys are slice values.
// Values of set are empty structs.
//
// Example:
// []string{"foo", "bar"} -> map[string]{"foo": struct{}{}, "bar": struct{}{}}
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

// md5 creates a md5 hash based on the bytes of the string.
func md5(str string) string {
	hash := crypto.Sum([]byte(str))
	return hex.EncodeToString(hash[:])
}

// md5File creates a md5 hash based on the bytes of the file.
func md5File(path string) (string, error) {
	h := crypto.New()

	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

// firstOr returns first element of slice, or fallback if slice is empty.
func firstOr[T any](items []T, fallback T) T {
	if len(items) > 0 {
		return items[0]
	}
	return fallback
}
