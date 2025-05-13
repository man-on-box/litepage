package validate

import (
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
)

var ErrDomainEmpty = errors.New("domain cannot be empty")
var ErrDomainInvalid = errors.New("domain is invalid, check it does not include spaces or any illegal characters")

func IsValidDomain(domain string) error {
	if domain == "" {
		return ErrDomainEmpty
	}

	parsedUrl, err := url.Parse(domain)
	if err != nil || parsedUrl.String() != domain {
		return ErrDomainInvalid
	}

	return nil
}

var ErrPathMustStartWithSlash = errors.New("path must start with '/'")
var ErrPathContainsInvalidCharacters = errors.New("path contains invalid characters")
var ErrPathContainsDirTraversal = errors.New("path contains illegal '..' for directory traversal")
var ErrPathMustIncludeFileExt = errors.New("path must end with a file extension e.g. '.html'")

func IsValidFilePath(filePath string) error {
	if !strings.HasPrefix(filePath, "/") {
		return ErrPathMustStartWithSlash
	}

	parsedURL, err := url.Parse(filePath)
	if err != nil {
		return fmt.Errorf("failed to parse path: %v", err)
	}

	if parsedURL.String() != filePath {
		return ErrPathContainsInvalidCharacters
	}

	if strings.Contains(parsedURL.Path, "..") {
		return ErrPathContainsDirTraversal
	}

	ext := filepath.Ext(parsedURL.Path)
	if ext == "" {
		return ErrPathMustIncludeFileExt
	}

	return nil
}

var ErrBasePathInvalid = errors.New("base path is invalid, check it does not include spaces or any illegal characters")
var ErrBasePathTrailingSlash = errors.New("base path should not have a trailing slash '/'")

func IsValidBasePath(basePath string) error {
	if !strings.HasPrefix(basePath, "/") {
		return ErrPathMustStartWithSlash
	}

	parsedUrl, err := url.Parse(basePath)
	if err != nil || parsedUrl.String() != basePath {
		return ErrBasePathInvalid
	}

	if strings.HasSuffix(basePath, "/") {
		return ErrBasePathTrailingSlash
	}

	return nil
}
