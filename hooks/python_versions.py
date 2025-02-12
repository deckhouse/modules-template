#!/usr/bin/env python3
#
# Copyright 2023 Flant JSC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


from deckhouse import hook
from dotmap import DotMap

# This hook subscribes to python.deckhouse.io/v1 CRs and puts their versions into ConfigMap
# 'python-versions'. The 'jqFilter' expression lets us focus only on meaningful parts of resources.
# The result of this filter will be in snapshots arral named 'python_versions'. Snapshots are in
# sync with cluster state, because by default 'kubeternetes' subscription uses all kinds of events.
#
# Refer to Shell Operator doc for details https://github.com/flant/shell-operator/blob/main/HOOKS.md
config = """
configVersion: v1
beforeHelm: 10
kubernetes:
- name: python_versions
  apiVersion: "example.io/v1"
  kind: "Python"
  queue: echoserver-python-versions
  group: main
  jqFilter: |
    .spec.version
  # We don't want to keep full custom resources in memory.
  keepFullObjectsInMemory: false
"""


def main(ctx: hook.Context):
    # Or goal is to edit values in-place (ctx.values). These values are used in template
    # templates/configmap.yaml. The schema for these values is defined in openapi/values.yaml.

    # From the hook run context we get the snapshots as we named it in the suscription. It will
    # always be a list if it is defined in the hook config. 'versions' here contain objects of the form
    #   [ {'filterResult': {'major': 3, 'minor': 8}} , ... ]
    # The version dict is the result of jqFilter '.spec.version', see crds/python.yaml into version v1.
    versions = ctx.snapshots.get("python_versions", [])

    # DotMap library simplifies access to nested fields in dicts, especially to inexisting ones.
    v = DotMap(ctx.values)

    # IMPORTANT: We assume that this module will be named 'echo-server' when added to Deckhouse. The
    # name of the module is used in the values reference. For now, module name in deckhouse and
    # values reference are tightly coupled.
    v.echoserver.internal.pythonVersions = [parse_snap_version(v) for v in versions]

    # DotMap is not JSON serializable, so we need to convert it back to dict in the end.
    ctx.values = v.toDict()


def parse_snap_version(snap):
    # Since we subscribed to ApiVersion example.io/v1, we get .spec.version (see jqFilter) as an
    # object with fields 'major' and 'minor'.
    v = snap["filterResult"]
    major, minor = v["major"], v["minor"]
    return f"{major}.{minor}"


if __name__ == "__main__":
    hook.run(main, config=config)
