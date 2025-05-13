package litepage_test

import (
	"io"
	"testing"

	"github.com/man-on-box/litepage"
	"github.com/stretchr/testify/assert"
)

func TestCreateNewLitepage(t *testing.T) {
	t.Run("Errors if domain is not specified", func(t *testing.T) {
		_, err := litepage.New("")
		assert.Error(t, err)
		assert.ErrorContains(t, err, "site domain is required")
	})

	t.Run("Errors if domain is invalid", func(t *testing.T) {
		_, err := litepage.New("invalid domain")
		assert.Error(t, err)
		assert.ErrorContains(t, err, "site domain is not valid")
	})

	t.Run("Returns no error when correct domain is supplied", func(t *testing.T) {
		_, err := litepage.New("nice-domain.com")
		assert.NoError(t, err)
	})

	t.Run("Returns error if base path is supplied and not valid", func(t *testing.T) {
		_, err := litepage.New("nice-domain.com", litepage.WithBasePath("/"))
		assert.Error(t, err)
		assert.ErrorContains(t, err, "base path is not valid")

	})
}

func TestAddNewPage(t *testing.T) {
	t.Run("errors when adding an invalid path", func(t *testing.T) {
		lp, err := litepage.New("test.com")
		assert.NoError(t, err)

		err = lp.Page("/invalid path", func(w io.Writer) {})
		assert.Error(t, err)
		assert.ErrorContains(t, err, "error when validating file path")
	})

	t.Run("returns no error when adding valid page", func(t *testing.T) {
		lp, err := litepage.New("test.com")
		assert.NoError(t, err)

		err = lp.Page("/foo.html", func(w io.Writer) {})
		assert.NoError(t, err)
	})

	t.Run("Errors when trying to add the same path twice", func(t *testing.T) {
		lp, err := litepage.New("test.com")
		assert.NoError(t, err)

		err = lp.Page("/foo.html", func(w io.Writer) {})
		assert.NoError(t, err)

		err = lp.Page("/foo.html", func(w io.Writer) {})
		assert.ErrorContains(t, err, "it already exists")
	})
}
