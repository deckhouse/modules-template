package copycustomcertificate

import (
	tlscertificate "github.com/deckhouse/module-sdk/common-hooks/copy-custom-certificate"
)

const moduleName = "my-module"

var _ = tlscertificate.RegisterHook(moduleName)
