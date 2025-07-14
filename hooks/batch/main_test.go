/*
Copyright 2025 Flant JSC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
		"name": "some-golang",
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
