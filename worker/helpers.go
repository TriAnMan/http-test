package worker

import "regexp"


var exp = regexp.MustCompile("^http(s)?://")

func PrepareUrl(rawUrl string) string {
	if hasScheme := exp.MatchString(rawUrl); !hasScheme {
		rawUrl = "http://" + rawUrl
	}
	return rawUrl
}


