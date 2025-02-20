package goreq

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/html/charset"
	"io"
	"net/http"
	urllib "net/url"
	"strings"
)

func QueryStringToMap(s string) map[string]string {
	query := make(map[string]string)
	pairs := strings.Split(s, "&")
	for _, pair := range pairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) == 2 {
			query[urllib.QueryEscape(kv[0])] = urllib.QueryEscape(kv[1])
		} else {
			query[kv[0]] = ""
		}
	}
	return query
}

//func QueryMapToString(query map[string]string) string {
//	if query == nil {
//		return ""
//	}
//	var sb strings.Builder
//	for key, value := range query {
//		if sb.Len() > 0 {
//			sb.WriteString("&")
//		}
//		sb.WriteString(urllib.QueryEscape(key))
//		sb.WriteString("=")
//		sb.WriteString(urllib.QueryEscape(value))
//	}
//	return sb.String()
//}

func QueryMapToString(m map[string]interface{}) string {
	var queryParams []string

	for key, value := range m {
		encodedKey := urllib.QueryEscape(key)
		var encodedValue string

		switch v := value.(type) {
		case string:
			encodedValue = urllib.QueryEscape(v)
		case int, int32, int64, float64:
			encodedValue = fmt.Sprintf("%v", v)
		default:
			encodedValue = fmt.Sprintf("%v", v)
		}

		queryParams = append(queryParams, fmt.Sprintf("%s=%s", encodedKey, encodedValue))
	}

	return strings.Join(queryParams, "&")
}

func QueriesEscape(s string) string {
	pairs := strings.Split(s, "&")
	encodedPairs := make([]string, len(pairs))
	for i, pair := range pairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) == 2 {
			encodedPairs[i] = fmt.Sprintf("%s=%s", urllib.QueryEscape(kv[0]), urllib.QueryEscape(kv[1]))
		} else {
			encodedPairs[i] = pair
		}
	}
	return strings.Join(encodedPairs, "&")
}

func GetResSuffix(res *http.Response) string {
	ct := res.Header.Get("Content-Type")
	ss := strings.Split(ct, ";")
	if ss == nil || len(ss) < 1 {
		return ""
	}
	return MimeToExtensions[ss[0]]
}

func NewReaderFromRes(res *http.Response) (io.Reader, error) {
	ct := res.Header.Get("Content-Type")
	return charset.NewReader(res.Body, ct)
}

func ResToJson(res *http.Response) (map[string]interface{}, error) {
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return result, nil
}
