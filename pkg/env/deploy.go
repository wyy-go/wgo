//go:generate stringer -type=Deploy -linecomment
package env

import (
	"log"
)

type Deploy int8

const (
	DeployUnknown Deploy = iota // unknown
	DeployLocal                 // local
	DeployDev                   // dev
	DeployTest                  // test
	DeployUat                   // uat
	DeployProd                  // prod
)

func ToDeploy(s string) Deploy {
	switch s {
	case DeployLocal.String():
		return DeployLocal
	case DeployDev.String():
		return DeployDev
	case DeployTest.String():
		return DeployTest
	case DeployUat.String():
		return DeployUat
	case DeployProd.String():
		return DeployProd
	}
	return DeployUnknown
}

var deploy = DeployDev

func GetDeploy() Deploy {
	return deploy
}

func SetDeploy(d Deploy) {
	deploy = d
}

func IsDeployLocal() bool {
	return deploy == DeployLocal
}

func IsDeployDev() bool {
	return deploy == DeployDev
}

func IsDeployTest() bool {
	return deploy == DeployTest
}

func IsDeployUat() bool {
	return deploy == DeployUat
}

func IsDeployProd() bool {
	return deploy == DeployProd
}

func IsDeployDebug() bool {
	return IsDeployLocal() || IsDeployDev() || IsDeployTest()
}

func IsDeployRelease() bool {
	return IsDeployUat() || IsDeployProd()
}

func MustSetDeploy(d Deploy) {
	deploy = d
	if deploy == DeployUnknown {
		log.Fatalf("Please set deploy mode first, must be one of local, dev, test, uat, prod")
	}
}
