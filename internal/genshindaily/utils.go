package genshindaily

import "net/http"

func getCookieByKey(cookies []*http.Cookie, name string) *http.Cookie {
	for _, v := range cookies {
		if v.Name == name {
			return v
		}
	}
	return nil
}

func parseCookies(s string) []*http.Cookie {
	request := http.Request{Header: http.Header{"Cookie": []string{s}}}
	return request.Cookies()
}
