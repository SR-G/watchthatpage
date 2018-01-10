package core

import(
	"testing"
)

func TestTemplate(t *testing.T) {
	results := GlobalResults{}
	result1 := Result{}
	result2 := Result{}
	results.Results = append(results.Results, result1)
	results.Results = append(results.Results, result2)
	content, err := renderTemplatedBody("/go/bin/templates/mail.tmpl", &results)
	if err != nil {
		t.Log("Error : ", err)
	}
	t.Log("Generated template is : \n" + content)
}