package middleware

import (
	lib "go_boilerplate/lib"
	"net"
	"net/http"
	"net/url"
	"strings"
)

var log = lib.GetLogger()

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getRealAddr(r)
		url := parseURL(r.URL)
		headers := getHeaders(r)

		log.Debugf("[Request]: %v \n"+
			"[Method]: %v \n"+
			"[IP]: %v \n"+
			"[Agent]: %v \n"+
			"[Proto]: %v \n"+
			"%v", url, r.Method, ip, r.UserAgent(), r.Proto, headers)
		next.ServeHTTP(w, r)
	})
}

func parseURL(u *url.URL) string {
	url := u.Scheme + "://"
	if u.Opaque != "" {
		url += u.Opaque
	}
	if u.Host != "" {
		if host, port, err := net.SplitHostPort(u.Host); err == nil {
			url += host + ":" + port
		} else {
			url += host
		}
	}
	if u.Path != "" {
		url += u.Path
	}
	if u.Fragment != "" {
		url += u.Fragment
	}
	if u.RawQuery != "" {
		url += u.RawQuery
	}
	return url
}

func getRealAddr(r *http.Request) string {

	remoteIP := ""
	// the default is the originating ip. but we try to find better options because this is almost
	// never the right IP
	if parts := strings.Split(r.RemoteAddr, ":"); len(parts) == 2 {
		remoteIP = parts[0]
	}
	// If we have a forwarded-for header, take the address from there
	if xff := strings.Trim(r.Header.Get("X-Forwarded-For"), ","); len(xff) > 0 {
		addrs := strings.Split(xff, ",")
		lastFwd := addrs[len(addrs)-1]
		if ip := net.ParseIP(lastFwd); ip != nil {
			remoteIP = ip.String()
		}
	} else if xri := r.Header.Get("X-Real-Ip"); len(xri) > 0 { // parse X-Real-Ip header

		if ip := net.ParseIP(xri); ip != nil {
			remoteIP = ip.String()
		}
	} else { // if doesn't match any case just pass RemoteAddr as it is
		remoteIP = r.RemoteAddr
	}

	return remoteIP

}

func getHeaders(r *http.Request) string {
	hd := ""
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			hd += "[" + name + "]: " + h + " \n"
		}
	}
	return hd
}
