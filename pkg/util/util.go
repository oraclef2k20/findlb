package util

import (
	"log"
	"net/url"
	"regexp"
	"strings"
)

func GetDomain(hostname string) (string, string) {

	pattern := "https://|http://"
	re := regexp.MustCompile(pattern)
	parts := []string{}

	if re.MatchString(hostname) {
		u, err := url.Parse(hostname)
		if err != nil {
			log.Fatal(err)
		}
		parts = strings.Split(u.Hostname(), ".")

	} else {

		parts = strings.Split(hostname, ".")
	}

	domain := parts[len(parts)-2] + "." + parts[len(parts)-1]
	host := strings.Join(parts, ".")
	return domain, host

}
