package testutils

import (
	"bytes"
	"io"
	"mime/multipart"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// CreateMultipartRequest creates a multipart/form-data request body with a single file part containing the provided content.
// It returns the multipart writer, the request body buffer, and the path to the temporary file created for the test.
// The caller is responsible for deleting the temporary file after use.
func CreateMultipartRequest(t *testing.T, content string) (*multipart.Writer, *bytes.Buffer, string) {
	// Create test file
	tmpFilePath := createTempFile(t, content)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.csv")
	assert.NoError(t, err)

	file, err := os.Open(tmpFilePath)
	assert.NoError(t, err)
	defer file.Close()

	_, err = io.Copy(part, file)
	assert.NoError(t, err)

	err = writer.Close()
	assert.NoError(t, err)

	return writer, body, tmpFilePath
}

// createTempFile creates a temporary file with the given content and returns its path.
// The caller is responsible for deleting the file after use.
func createTempFile(t *testing.T, content string) string {
	tmpFile, err := os.CreateTemp("", "test-*.txt")
	assert.NoError(t, err)

	_, err = tmpFile.WriteString(content)
	assert.NoError(t, err)

	err = tmpFile.Close()
	assert.NoError(t, err)

	return tmpFile.Name()
}
