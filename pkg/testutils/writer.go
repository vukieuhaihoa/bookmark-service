package testutils

import (
	"bytes"
	"io"
	"mime/multipart"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// CreateMultipartRequest creates a multipart/form-data request body with a single file part containing the provided content.
// It returns the multipart writer, the request body buffer, and the path to the temporary file created for the test.
// The caller is responsible for deleting the temporary file after use.
func CreateMultipartRequest(t *testing.T, content string) (*multipart.Writer, *bytes.Buffer) {

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "test.csv")
	assert.NoError(t, err)

	_, err = io.Copy(part, strings.NewReader(content))
	assert.NoError(t, err)

	err = writer.Close()
	assert.NoError(t, err)

	return writer, body
}
