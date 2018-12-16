package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"strings"
)

func ParseOpeResult(body io.Reader) string {
	z := html.NewTokenizer(body)
	content := []string{}

	// 苦肉の HTML 解析...
	for z.Token().Data != "html" {
		tt := z.Next()
		if tt == html.StartTagToken {
			t := z.Token()
			// div タグに含まれる文字列を取得したい
			if t.Data == "div" {
				inner := z.Next()
				if inner == html.TextToken {
					text := (string)(z.Text())
					t := strings.TrimSpace(text)
					content = append(content, t)
				}
			}
		}
	}
	// 以下の文字列を取得出来ることを期待
	// External Command Interface
	// Your command request was successfully submitted to Nagios for processing.
	// fmt.Println(strings.Contains(content[1], "successfully"))
	resultMessage := content[1]
	if strings.Contains(resultMessage, "successfully") {
		return "ok"
	} else {
		return "ng"
	}
}

func ParseCheckHostStatus(targetHost string, body io.Reader) {
	z := html.NewTokenizer(body)
	keys := []string{}
	values := []string{}

	// 苦肉の HTML 解析...
	for z.Token().Data != "html" {
		tt := z.Next()
		if tt == html.StartTagToken {
			t := z.Token()
			var key string
			var value string
			if t.Data == "td" {
				inner1 := z.Next()
				if inner1 == html.TextToken {
					text := ((string)(z.Text()))
					if strings.TrimSpace(text) != "" {
						key = strings.TrimSpace(text)
						if strings.HasSuffix(key, ":") || strings.HasSuffix(key, "?") {
							keys = append(keys, strings.Trim(key, ":?"))
						} else {
							values = append(values, strings.Trim(key, ":?"))
						}
					}
				}

				if inner1 == html.StartTagToken {
					if z.Token().Data == "div" {
						inner2 := z.Next()
						if inner2 == html.TextToken {
							text := (string)(z.Text())
							if strings.TrimSpace(text) != "" {
								value = strings.TrimSpace(text)
								values = append(values, strings.Trim(value, ":?"))
							}
						}
					}
				}
			}
		}
	}

	m := map[string]string{}
	for i := 0; i < len(keys); i++ {
		m[keys[i]] = values[i]
	}
	// fmt.Println(m)
	data := map[string]interface{}{"Host": targetHost, "Status": m}
	data_json, _ := json.Marshal(data)
	fmt.Println(string(data_json))
}
