package util

import (
	"net/url"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

func GetDomain(hostname string) (string, string) {

	pattern := "https://|http://"
	re := regexp.MustCompile(pattern)
	var parts []string

	if re.MatchString(hostname) {
		u, err := url.Parse(hostname)
		if err != nil {
			log.Fatal(err)
		}
		parts = strings.Split(u.Hostname(), ".")

	} else {

		parts = strings.Split(hostname, ".")
	}

	log.WithFields(
		log.Fields{
			"domain-parts": parts,
		}).Debug()

	domain := ""

	// for co.jp
	if parts[len(parts)-1] == "jp" && parts[len(parts)-2] == "co" {
		domain = parts[len(parts)-3] + "." + parts[len(parts)-2] + "." + parts[len(parts)-1]
	} else {
		domain = parts[len(parts)-2] + "." + parts[len(parts)-1]
	}

	host := strings.Join(parts, ".")
	return domain, host

}
