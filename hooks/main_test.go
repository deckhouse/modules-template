package main_test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	hook "hook"

	"github.com/deckhouse/module-sdk/pkg/jq"
)

func Test_JQFilterApplyGolangVersion(t *testing.T) {
	t.Run("apply golang", func(t *testing.T) {
		const rawGolang = `
		{
	  "apiVersion": "example.io/v1",
	  "kind": "Golang",
	  "metadata": {
		"name": "some-pytnon",
		"namespace": "some-ns"
	  },
	  "spec": {
		"version": {
			"major":1,
			"minor":23,
			"patch":8
		}
	  }
	}`

		q, err := jq.NewQuery(hook.ApplyNodeJQFilter)
		assert.NoError(t, err)

		res, err := q.FilterStringObject(context.Background(), rawGolang)
		assert.NoError(t, err)

		golangVersion := new(hook.VersionInfoMetadata)
		err = json.NewDecoder(bytes.NewBufferString(res.String())).Decode(golangVersion)
		assert.NoError(t, err)

		assert.Equal(t, 1, golangVersion.Major)
		assert.Equal(t, 23, golangVersion.Minor)
		assert.Equal(t, 8, golangVersion.Patch)
	})
}
