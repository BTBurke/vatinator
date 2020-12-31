package vatinator

import (
	"bytes"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const expect = `Hi Test,

There was a problem processing your forms.  This email has also been sent to Bryan so he can fix it.  Sorry about that!

Run log:
This is a test
And another line
`

func TestErrorTemplate(t *testing.T) {
	temp, err := template.New("email").Parse(errorEmail)
	require.NoError(t, err)
	var b bytes.Buffer
	assert.NoError(t, temp.Execute(&b, EmailData{
		FormData: FormData{FirstName: "Test"},
		RunLog:   "This is a test\nAnd another line",
	}))
	assert.Equal(t, expect, b.String())

}
