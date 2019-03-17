package pgot

import (
	"text/template"
)

var (
	// This should include all publicly accessible custom got functions
	funcMap = template.FuncMap{
		"lnp":  lnp,
	}
)

// The lnp (link new page) function converts the supplied url into a link which
// opens a new page.  If label is blank (aka "") then the url is displayed,
// otherwise the label text is used when displaying the link.
func lnp(label, url string) string {
	if label == "" {
		return "<a href=\"" + url + "\" target=\"_blank\">" + url + "</a>"
	} else {
		return "<a href=\"" + url + "\" target=\"_blank\">" + label + "</a>"
	}
}
