package copycustomcertificate

import (
	copyCustomCertificate "github.com/deckhouse/module-sdk/common-hooks/copy-custom-certificate"
)

const moduleName = "my-module"

var _ = copyCustomCertificate.RegisterHook(moduleName)
