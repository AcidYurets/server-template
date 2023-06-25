package utils

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

func NewURL(host, port string) (*url.URL, error) {

	// Отрежем / - если он есть
	host = strings.TrimRight(host, "/")

	// проверил есть ли в пути схема
	if !(strings.HasPrefix(host, "http://") ||
		strings.HasPrefix(host, "https://")) {
		host = "http://" + host
	}

	newUrl, err := url.Parse(host)
	if err != nil {
		return nil, fmt.Errorf("не получилось создать url: %s", err)
	}

	if len(port) > 0 {
		host, _, _ := net.SplitHostPort(newUrl.Host)
		if host == "" {
			host = newUrl.Host
		}
		newUrl.Host = net.JoinHostPort(host, port)
	}

	return newUrl, nil
}
