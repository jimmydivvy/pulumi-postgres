package main

import (
	"github.com/jimmydivvy/pulumi-postgres"
	"github.com/jimmydivvy/pulumi-postgres/pkg/version"
	"github.com/pulumi/pulumi-terraform/pkg/tfgen"
)

func main() {
	tfgen.Main("postgresql", version.Version, postgresql.Provider())
}