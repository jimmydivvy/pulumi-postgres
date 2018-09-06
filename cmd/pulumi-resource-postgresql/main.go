package main

import (
	"github.com/pulumi/pulumi-terraform/pkg/tfbridge"
	"github.com/jimmydivvy/pulumi-postgres"
	"github.com/jimmydivvy/pulumi-postgres/pkg/version"
)

func main() {
	tfbridge.Main("postgresql", version.Version, postgresql.Provider())
}