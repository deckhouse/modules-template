// # Copyright 2023 Flant JSC
// #
// # Licensed under the Apache License, Version 2.0 (the "License");
// # you may not use this file except in compliance with the License.
// # You may obtain a copy of the License at
// #
// #     http://www.apache.org/licenses/LICENSE-2.0
// #
// # Unless required by applicable law or agreed to in writing, software
// # distributed under the License is distributed on an "AS IS" BASIS,
// # WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// # See the License for the specific language governing permissions and
// # limitations under the License.

package main

import (
	"context"
	"log/slog"

	_ "hook/https"

	"github.com/deckhouse/module-sdk/pkg"
	"github.com/deckhouse/module-sdk/pkg/app"
	objectpatch "github.com/deckhouse/module-sdk/pkg/object-patch"
	"github.com/deckhouse/module-sdk/pkg/registry"
)

const (
	SnapshotKey = "custom_certificates"
)

var _ = registry.RegisterFunc(config, HandlerHook)

// # Since we subscribed to ApiVersion example.io/v1, we get .spec.version (see jqFilter) as an
// # object with fields 'major' and 'minor'.
// v = snap["filterResult"]
// major, minor = v["major"], v["minor"]
// {"major":2,"minor":5}
type NodeInfoMetadata struct {
	Major string `json:"major"`
	Minor string `json:"minor"`
}

const ApplyNodeJQFilter = `.spec.version`

// # This hook subscribes to python.deckhouse.io/v1 CRs and puts their versions into ConfigMap
// # 'python-versions'. The 'jqFilter' expression lets us focus only on meaningful parts of resources.
// # The result of this filter will be in snapshots arral named 'python_versions'. Snapshots are in
// # sync with cluster state, because by default 'kubeternetes' subscription uses all kinds of events.
// #
// # Refer to Shell Operator doc for details https://github.com/flant/shell-operator/blob/main/HOOKS.md
var config = &pkg.HookConfig{
	Kubernetes: []pkg.KubernetesConfig{
		{
			Name:       "python_versions",
			APIVersion: "example.io/v1",
			Kind:       "Python",
			// NamespaceSelector: &pkg.NamespaceSelector{
			// 	NameSelector: &pkg.NameSelector{
			// 		MatchNames: []string{"d8-system"},
			// 	},
			// },
			// LabelSelector: &v1.LabelSelector{
			// 	MatchExpressions: []v1.LabelSelectorRequirement{
			// 		v1.LabelSelectorRequirement{
			// 			Key:      "owner",
			// 			Operator: v1.LabelSelectorOpNotIn,
			// 			Values:   []string{"helm"},
			// 		},
			// 	},
			// },
			JqFilter: ApplyNodeJQFilter,
		},
	},
}

func HandlerHook(_ context.Context, input *pkg.HookInput) error {
	// # From the hook run context we get the snapshots as we named it in the suscription. It will
	// # always be a list if it is defined in the hook config. 'versions' here contain objects of the form
	// #   [ {'filterResult': {'major': 3, 'minor': 8}} , ... ]
	// # The version dict is the result of jqFilter '.spec.version', see crds/python.yaml into version v1.
	pythonVersions, err := objectpatch.UnmarshalToStruct[NodeInfoMetadata](input.Snapshots, "python_versions")
	if err != nil {
		return err
	}

	input.Logger.Info("found python_versions", slog.Any("pythonVersions", pythonVersions))

	versions := make([]string, 0, len(pythonVersions))
	for _, version := range pythonVersions {
		versions = append(versions, parse_snap_version(version))
	}

	// # IMPORTANT: We assume that this module will be named 'echo-server' when added to Deckhouse. The
	// # name of the module is used in the values reference. For now, module name in deckhouse and
	// # values reference are tightly coupled.
	input.Values.Set("echoserver.internal.pythonVersions", versions)

	return nil
}

// # Since we subscribed to ApiVersion example.io/v1, we get .spec.version (see jqFilter) as an
// # object with fields 'major' and 'minor'.
func parse_snap_version(version NodeInfoMetadata) string {
	return version.Major + "." + version.Minor
}

func main() {
	app.Run()
}
