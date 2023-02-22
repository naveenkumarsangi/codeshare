package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListFiles(t *testing.T) {
	files, err := getFileList()
	if err != nil {
		t.Error(err)
		return
	}

	assert.EqualValues(t, files, []File{})
	assert.NotEqualValues(t, len(files), 0)
}
