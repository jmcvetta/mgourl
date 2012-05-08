// This is Free Software.  See LICENSE.txt for details.

// Package mgourl provides functionality from the 'mgo' MongoDB driver's 
// unexported parseUrl() function.
package mgourl

import (
	"strings"
	"errors"
)


type AuthInfo struct {
	Db, User, Pass string
}

func isOptSep(c rune) bool {
	return c == ';' || c == '&'
}

func ParseURL(url string) (servers []string, auth AuthInfo, options map[string]string, err error) {
	if strings.HasPrefix(url, "mongodb://") {
		url = url[10:]
	}
	options = make(map[string]string)
	if c := strings.Index(url, "?"); c != -1 {
		for _, pair := range strings.FieldsFunc(url[c+1:], isOptSep) {
			l := strings.SplitN(pair, "=", 2)
			if len(l) != 2 || l[0] == "" || l[1] == "" {
				err = errors.New("Connection option must be key=value: " + pair)
				return
			}
			options[l[0]] = l[1]
		}
		url = url[:c]
	}
	if c := strings.Index(url, "@"); c != -1 {
		pair := strings.SplitN(url[:c], ":", 2)
		if len(pair) != 2 || pair[0] == "" {
			err = errors.New("Credentials must be provided as user:pass@host")
			return
		}
		auth.User = pair[0]
		auth.Pass = pair[1]
		url = url[c+1:]
		auth.Db = "admin"
	}
	if c := strings.Index(url, "/"); c != -1 {
		if c != len(url)-1 {
			auth.Db = url[c+1:]
		}
		url = url[:c]
	}
	if auth.User == "" {
		if auth.Db != "" {
			err = errors.New("Database name only makes sense with credentials")
			return
		}
	} else if auth.Db == "" {
		auth.Db = "admin"
	}
	servers = strings.Split(url, ",")
	// XXX This is untested. The test suite doesn't use the standard port.
	for i, server := range servers {
		p := strings.LastIndexAny(server, "]:")
		if p == -1 || server[p] != ':' {
			servers[i] = server + ":27017"
		}
	}
	return
}

