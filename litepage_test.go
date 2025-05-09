package litepage_test

import (
	"fmt"
	"io"
	"testing"

	"github.com/man-on-box/litepage"
	"github.com/stretchr/testify/assert"
)

func TestCreateNewLitepage(t *testing.T) {
	t.Run("Errors if domain is not specified", func(t *testing.T) {
		_, err := litepage.New("")
		assert.ErrorContains(t, err, "site domain is required")
	})

	t.Run("Errors if domain contains spaces", func(t *testing.T) {
		_, err := litepage.New("invalid domain")
		assert.ErrorContains(t, err, "site domain is not valid")
	})

	t.Run("Errors if domain contains illegal characters", func(t *testing.T) {
		_, err := litepage.New("invaliddomain%=?")
		assert.ErrorContains(t, err, "site domain is not valid")
	})

	t.Run("Returns no error when correct domain is supplied", func(t *testing.T) {
		_, err := litepage.New("nice-domain.com")
		assert.NoError(t, err)
	})
}

func TestAddNewPage(t *testing.T) {
	tests := []struct {
		filePath      string
		expectError   bool
		errorContains string
	}{
		{
			filePath:    "/foo.html",
			expectError: false,
		},
		{
			filePath:    "/foo.htm",
			expectError: false,
		},
		{
			filePath:    "/nested/foo.htm",
			expectError: false,
		},
		{
			filePath:    "/nested/foo.txt",
			expectError: false,
		}, {
			filePath:    "/file-hyphen_underscore.html",
			expectError: false,
		},
		{
			filePath:      "",
			expectError:   true,
			errorContains: "must start with '/'",
		},
		{
			filePath:      "index.html",
			expectError:   true,
			errorContains: "must start with '/'",
		},
		{
			filePath:      "/../index.html",
			expectError:   true,
			errorContains: "contains illegal '..'",
		},
		{
			filePath:      "/index",
			expectError:   true,
			errorContains: "must end with a file extension",
		},
		{
			filePath:      "/index file.html",
			expectError:   true,
			errorContains: "contains invalid character",
		}, {
			filePath:      "/index<.html",
			expectError:   true,
			errorContains: "contains invalid character",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Expect error '%t' for path '%s'", tt.expectError, tt.filePath), func(t *testing.T) {
			lp, err := litepage.New("test.com")
			assert.NoError(t, err)

			err = lp.Page(tt.filePath, func(w io.Writer) {})

			if tt.expectError {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tt.errorContains)
			} else {
				assert.NoError(t, err)
			}

		})
	}

	t.Run("Errors when trying to add the same path twice", func(t *testing.T) {
		lp, err := litepage.New("test.com")
		assert.NoError(t, err)

		err = lp.Page("/foo.html", func(w io.Writer) {})
		assert.NoError(t, err)

		err = lp.Page("/foo.html", func(w io.Writer) {})
		assert.ErrorContains(t, err, "it already exists")
	})
}
