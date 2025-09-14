package operations

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTape_Write(t *testing.T) {
	file, teardown := setupFile(t, []byte("12345"))
	defer teardown()

	tape := &tape{file}

	tape.Write([]byte("abc"))

	file.Seek(0, io.SeekStart)
	newFileContents, err := io.ReadAll(file)
	assert.NoError(t, err)

	assert.Equal(t, "abc", string(newFileContents))
}
