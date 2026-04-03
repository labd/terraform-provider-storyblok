package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// NormalizeTransport wraps an http.RoundTripper and normalizes JSON responses
// to handle cases where the Storyblok API returns numbers where strings are
// expected (e.g. allowed_paths).
type NormalizeTransport struct {
	inner http.RoundTripper
}

// NewNormalizeTransport creates a new NormalizeTransport that wraps the given
// transport.
func NewNormalizeTransport(inner http.RoundTripper) *NormalizeTransport {
	if inner == nil {
		inner = http.DefaultTransport
	}
	return &NormalizeTransport{inner: inner}
}

func (t *NormalizeTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := t.inner.RoundTrip(req)
	if err != nil || resp == nil {
		return resp, err
	}

	// Only process JSON responses.
	if !strings.Contains(resp.Header.Get("Content-Type"), "application/json") {
		return resp, nil
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		resp.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		return resp, nil
	}

	normalized := normalizeStringArrayFields(bodyBytes)
	resp.Body = io.NopCloser(bytes.NewReader(normalized))
	resp.ContentLength = int64(len(normalized))

	return resp, nil
}

// stringArrayPaths lists JSON paths whose arrays should always contain strings.
// The Storyblok API occasionally returns integer values in these arrays instead
// of strings.
var stringArrayPaths = [][]string{
	{"space_role", "allowed_paths"},
	{"space_role", "resolved_allowed_paths"},
}

// normalizeStringArrayFields parses body as JSON and converts any numeric
// values to strings in the fields listed in stringArrayPaths. The original
// bytes are returned unchanged when the body is not valid JSON or no
// conversion is needed.
func normalizeStringArrayFields(body []byte) []byte {
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return body
	}

	changed := false
	for _, path := range stringArrayPaths {
		if normalizeNumericArrayField(data, path) {
			changed = true
		}
	}

	if !changed {
		return body
	}

	result, err := json.Marshal(data)
	if err != nil {
		return body
	}
	return result
}

// normalizeNumericArrayField navigates the data map along path and converts
// any numeric values in the target array to their string representations.
// Returns true when at least one value was converted.
func normalizeNumericArrayField(data map[string]interface{}, path []string) bool {
	if len(path) == 0 {
		return false
	}

	val, ok := data[path[0]]
	if !ok {
		return false
	}

	if len(path) == 1 {
		arr, ok := val.([]interface{})
		if !ok {
			return false
		}
		changed := false
		for i, v := range arr {
			if f, ok := v.(float64); ok {
				arr[i] = strconv.FormatInt(int64(f), 10)
				changed = true
			}
		}
		return changed
	}

	nested, ok := val.(map[string]interface{})
	if !ok {
		return false
	}
	return normalizeNumericArrayField(nested, path[1:])
}
