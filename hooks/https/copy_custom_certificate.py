#!/usr/bin/env python3

from deckhouse import hook
from dotmap import DotMap


# External module mechanism makes YAML files execulable in hooks, so we use in-line configs
config = """
configVersion: v1
beforeHelm: 10
kubernetes:
- name: custom_certificates
  apiVersion: v1
  kind: Secret
  namespace:
    nameSelector:
      matchNames: ["d8-system"]
  labelSelector:
    matchExpressions:
    - key: "owner"
      operator: "NotIn"
      values: ["helm"]
  jqFilter: |
    {
      name: .metadata.name,
      data: .data
    }
"""


def main(ctx: hook.Context):
    snaps = ctx.snapshots.get("custom_certificates", [])
    certs = {c["name"]: c["data"] for c in [s["filterResult"] for s in snaps]}

    # module values and global values
    mv = DotMap(ctx.values["echoserver"])
    gv = DotMap(ctx.values["global"])

    # we have to copy DotMap instances FOR READING to avoid initializing absent fields that will be
    # stored in values
    https_mode = first_non_empty(
        mv.copy().https.mode,
        gv.copy().modules.https.mode,
        default="",
    )
    if https_mode != "CustomCertificate":
        if mv.copy().internal.customCertificateData:
            del mv.internal.customCertificateData
        ctx.values["echoserver"] = mv.toDict()
        return

    # Setting custom certificate into values
    secret_name = first_non_empty(
        mv.copy().https.customCertificate.secretName,
        gv.copy().modules.https.customCertificate.secretName,
        default="",
    )

    if secret_name not in certs:
        return
    mv.internal.customCertificateData = certs[secret_name]

    ctx.values["echoserver"] = mv.toDict()
    return


def first_non_empty(*args, default=None):
    for a in args:
        if a:
            return a
    return default


if __name__ == "__main__":
    hook.run(main, config=config)
