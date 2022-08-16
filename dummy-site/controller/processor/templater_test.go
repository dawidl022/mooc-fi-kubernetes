package processor

import (
	"bytes"
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUrlToWebsite_StripsProtocol(t *testing.T) {
	expected := []string{"example", "website"}
	inputs := []string{"http://example", "https://website"}
	assertUrlToWebsite(t, inputs, expected)
}

func TestUrlToWebsite_ReplacesDotsInDomainWithDashes(t *testing.T) {
	inputs := []string{"example.com", "www.website.co.uk"}
	expected := []string{"example-com", "www-website-co-uk"}
	assertUrlToWebsite(t, inputs, expected)
}

// TODO consider keeping path, replacing '/' with '-'
// and removing special chars and query parameters
func TestUrlToWebsite_StripsPath(t *testing.T) {
	inputs := []string{"example/hello", "website/some/path/to/article.html"}
	expected := []string{"example", "website"}
	assertUrlToWebsite(t, inputs, expected)
}

func TestUrlToWebsite_ReturnValidK8sIdentifier(t *testing.T) {
	inputs := []string{
		"http://example.com/hello",
		"https://www.website.co.uk/some/path/to/article.html",
	}
	expected := []string{"example-com", "www-website-co-uk"}
	assertUrlToWebsite(t, inputs, expected)
}

func assertUrlToWebsite(t *testing.T, inputs []string, expected []string) {
	for i, input := range inputs {
		assert.Equal(t, expected[i], UrlToWebsite(input))
	}
}

func TestWriteTemplate_WritesUsingGivenTemplate(t *testing.T) {
	message := "Hello, world!"
	templateApplier := templateApplier{template: message}

	buf := new(bytes.Buffer)
	err := templateApplier.Write(buf, nil)

	assert.NoError(t, err)
	assert.Equal(t, message, buf.String())
}

func TestWriteTemplate_WritesPlaceholders(t *testing.T) {
	params := struct {
		Message string
		Sender  string
	}{
		Message: "Hello, World!",
		Sender:  "Dawid",
	}
	template := "{{ .Message }}, from {{ .Sender }}"
	expected := "Hello, World!, from Dawid"

	templateApplier := templateApplier{template: template}
	buf := new(bytes.Buffer)

	err := templateApplier.Write(buf, params)
	assert.NoError(t, err)
	assert.Equal(t, expected, buf.String())
}

//go:embed testdata/manifest.yml
var expectedManifest string

func TestGenerateManifests_EscapesQuotesAndNewlines(t *testing.T) {
	applier := newManifestApplier()
	message := "Hello, world!\nThis is a cool test 'n' all."
	buf := new(bytes.Buffer)

	applier.writeManifest(buf, "website", message)
	assert.Equal(t, expectedManifest, buf.String())
}

//go:embed templates/manifest.yml.tmpl
var manifestTemplate string

func newManifestApplier() manifestApplier {
	return manifestApplier{templateApplier: templateApplier{template: manifestTemplate}}
}
