package processor

import (
	"bytes"
	_ "embed"
	"io"
	"regexp"
	"strings"
	"text/template"
)

func UrlToWebsite(url string) string {
	protocol := regexp.MustCompile(`https?://`)
	url = protocol.ReplaceAllLiteralString(url, "")

	path := regexp.MustCompile(`/.*`)
	url = path.ReplaceAllLiteralString(url, "")
	return strings.ReplaceAll(url, ".", "-")
}

func GenerateManifests(website string, body string) (*ManifestReaders, error) {
	dep, err := generateManifest(deploymentTemplate, website, body)
	if err != nil {
		return nil, err
	}
	ser, err := generateManifest(serviceTemplate, website, body)
	if err != nil {
		return nil, err
	}
	ing, err := generateManifest(ingressTemplate, website, body)
	if err != nil {
		return nil, err
	}

	return &ManifestReaders{
		deploymentReader: dep,
		serviceReader:    ser,
		ingressReader:    ing,
	}, nil
}

//go:embed templates/deployment.yml.tmpl
var deploymentTemplate string

//go:embed templates/service.yml.tmpl
var serviceTemplate string

//go:embed templates/ingress.yml.tmpl
var ingressTemplate string

func generateManifest(template string, website string, body string) (*bytes.Buffer, error) {
	m := manifestApplier{templateApplier: templateApplier{template: template}}
	buf := new(bytes.Buffer)
	err := m.writeManifest(buf, website, body)
	return buf, err
}

func (m *manifestApplier) writeManifest(w io.Writer, website string, body string) error {
	body = strings.ReplaceAll(body, "'", `\'`)
	body = strings.ReplaceAll(body, "\n", `\n`)
	params := manifestParams{Website: website, Body: body}
	return m.templateApplier.Write(w, params)
}

type manifestParams struct {
	Website string
	Body    string
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
