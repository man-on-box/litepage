package validate_test

import (
	"fmt"
	"testing"

	"github.com/man-on-box/litepage/internal/validate"
	"github.com/stretchr/testify/assert"
)

func TestIsValidDomain(t *testing.T) {
	invalidDomains := []struct {
		domain        string
		expectedError error
	}{
		{domain: "", expectedError: validate.ErrDomainEmpty},
		{domain: "domain with spaces", expectedError: validate.ErrDomainInvalid},
		{domain: "invalid%", expectedError: validate.ErrDomainInvalid},
	}

	for _, tt := range invalidDomains {
		t.Run(fmt.Sprintf("expect error for domain '%s'", tt.domain), func(t *testing.T) {
			err := validate.IsValidDomain(tt.domain)
			assert.Error(t, err)
			assert.ErrorIs(t, err, tt.expectedError)
		})
	}

	t.Run("returns nil for valid domain", func(t *testing.T) {
		err := validate.IsValidDomain("hello.com")
		assert.NoError(t, err)
	})
}

func TestIsValidFilePath(t *testing.T) {
	invalidPaths := []struct {
		path          string
		expectedError error
	}{
		{path: "", expectedError: validate.ErrPathMustStartWithSlash},
		{path: "index.html", expectedError: validate.ErrPathMustStartWithSlash},
		{path: "/../index.html", expectedError: validate.ErrPathContainsDirTraversal},
		{path: "/index", expectedError: validate.ErrPathMustIncludeFileExt},
		{path: "/index file.html", expectedError: validate.ErrPathContainsInvalidCharacters},
		{path: "/index>file.html", expectedError: validate.ErrPathContainsInvalidCharacters},
	}

	for _, tt := range invalidPaths {
		t.Run(fmt.Sprintf("expect error for path '%s'", tt.path), func(t *testing.T) {
			err := validate.IsValidFilePath(tt.path)
			assert.Error(t, err)
			assert.ErrorIs(t, err, tt.expectedError)
		})
	}

	validPaths := []string{
		"/foo.html",
		"/foo.htm",
		"/nested/foo.htm",
		"/nested/foo.txt",
		"/file-hyphen_underscore.html",
	}

	for _, path := range validPaths {
		t.Run(fmt.Sprintf("no error returned for valid domain '%s'", path), func(t *testing.T) {
			err := validate.IsValidFilePath(path)
			assert.NoError(t, err)
		})
	}
}

func TestIsValidBasePath(t *testing.T) {
	invalidPaths := []struct {
		path          string
		expectedError error
	}{
		{path: "", expectedError: validate.ErrPathMustStartWithSlash},
		{path: "/invalid path", expectedError: validate.ErrBasePathInvalid},
		{path: "/invalid%path", expectedError: validate.ErrBasePathInvalid},
		{path: "/invalid<>path", expectedError: validate.ErrBasePathInvalid},
		{path: "/", expectedError: validate.ErrBasePathTrailingSlash},
		{path: "/path-with-trailing-slash/", expectedError: validate.ErrBasePathTrailingSlash},
	}

	for _, tt := range invalidPaths {
		t.Run(fmt.Sprintf("expect error for path '%s'", tt.path), func(t *testing.T) {
			err := validate.IsValidBasePath(tt.path)
			assert.Error(t, err)
			assert.ErrorIs(t, err, tt.expectedError)
		})
	}

	validPaths := []string{
		"/test",
		"/nested/test",
		"/test-path_",
	}

	for _, path := range validPaths {
		t.Run(fmt.Sprintf("no error for valid path '%s'", path), func(t *testing.T) {
			err := validate.IsValidBasePath(path)
			assert.NoError(t, err)
		})

	}
}
