package processor

import (
	_ "embed"
	"io"
	"regexp"
	"strings"
	"text/template"
)

func urlToWebsite(url string) string {
	protocol := regexp.MustCompile(`https?://`)
	url = protocol.ReplaceAllLiteralString(url, "")

	path := regexp.MustCompile(`/.*`)
	url = path.ReplaceAllLiteralString(url, "")
	return strings.ReplaceAll(url, ".", "-")
}

type templateApplier struct {
	template string
}

func (t *templateApplier) Write(w io.Writer, params interface{}) error {
	tmpl, err := template.New("tmpl").Parse(t.template)
	if err != nil {
		return err
	}
	return tmpl.Execute(w, params)
}

type manifestApplier struct {
	templateApplier
}

//go:embed templates/manifest.yml.tmpl
var manifestTemplate string

func newManifestApplier() manifestApplier {
	return manifestApplier{templateApplier: templateApplier{template: manifestTemplate}}
}

type manifestParams struct {
	Website string
	Body    string
}

func (m *manifestApplier) GenerateManifests(w io.Writer, website string, body string) {
	body = strings.ReplaceAll(body, "'", `\'`)
	body = strings.ReplaceAll(body, "\n", `\n`)
	params := manifestParams{Website: website, Body: body}
	m.templateApplier.Write(w, params)
}
