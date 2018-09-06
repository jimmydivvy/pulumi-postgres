package postgresql

import (
	"unicode"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/pulumi/pulumi-terraform/pkg/tfbridge"
	"github.com/pulumi/pulumi/pkg/resource"
	"github.com/pulumi/pulumi/pkg/tokens"
	"github.com/terraform-providers/terraform-provider-postgresql/postgresql"
)

// all of the AWS token components used below.
const (
	// packages:
	pgPkg = "postgresql"
	// modules:
	pgMod = "index" // the root index.
)

// pgMember manufactures a type token for the package and the given module and type.
func pgMember(mod string, mem string) tokens.ModuleMember {
	return tokens.ModuleMember(pgPkg + ":" + mod + ":" + mem)
}

// pgType manufactures a type token for the package and the given module and type.
func pgType(mod string, typ string) tokens.Type {
	return tokens.Type(pgMember(mod, typ))
}

// pgDataSource manufactures a standard resource token given a module and resource name.
// It automatically uses the package and names the file by simply lower casing the data
// source's first character.
func pgDataSource(mod string, res string) tokens.ModuleMember {
	fn := string(unicode.ToLower(rune(res[0]))) + res[1:]
	return pgMember(mod+"/"+fn, res)
}

// pgResource manufactures a standard resource token given a module and resource name.  It automatically uses the AWS
// package and names the file by simply lower casing the resource's first character.
func pgResource(mod string, res string) tokens.Type {
	fn := string(unicode.ToLower(rune(res[0]))) + res[1:]
	return pgType(mod+"/"+fn, res)
}

// stringValue gets a string value from a property map if present, else ""
func stringValue(vars resource.PropertyMap, prop resource.PropertyKey) string {
	val, ok := vars[prop]
	if ok && val.IsString() {
		return val.StringValue()
	}
	return ""
}

func preConfigureCallback(vars resource.PropertyMap, c *terraform.ResourceConfig) error {
	return nil
}

// Provider returns additional overlaid schema and metadata associated with the  package.
func Provider() tfbridge.ProviderInfo {
	p := postgresql.Provider().(*schema.Provider)
	prov := tfbridge.ProviderInfo{
		P:                    p,
		Name:                 "postgresql",
		Description:          "A Pulumi package for creating and managing Postgresql databases",
		Keywords:             []string{"pulumi", "postgresql"},
		License:              "Apache-2.0",
		Homepage:             "https://pulumi.io",
		Repository:           "https://github.com/jimmydivvy/pulimi-postgres",
		PreConfigureCallback: preConfigureCallback,

		Resources: map[string]*tfbridge.ResourceInfo{
			"postgresql_database":  {Tok: pgResource(pgMod, "Database")},
			"postgresql_extension":       {Tok: pgResource(pgMod, "Extension")},
			"postgresql_role":      {Tok: pgResource(pgMod, "Role")},
			"postgresql_schema":     {Tok: pgResource(pgMod, "Schema")},
		},
		JavaScript: &tfbridge.JavaScriptInfo{
			Dependencies: map[string]string{
				"@pulumi/pulumi":    "^0.15.0",
				"builtin-modules":   "3.0.0",
				"read-package-tree": "^5.2.1",
				"resolve":           "^1.7.1",
			},
			DevDependencies: map[string]string{
				"@types/node": "^8.0.25", // so we can access strongly typed node definitions.
			},
			Overlay: &tfbridge.OverlayInfo{
				Files:   []string{},
				Modules: map[string]*tfbridge.OverlayInfo{},
			},
		},
		Python: &tfbridge.PythonInfo{
			Requires: map[string]string{
				"pulumi": ">=0.14.2,<0.15.0",
			},
		},
	}

	// For all resources with name properties, we will add an auto-name property.  Make sure to skip those that
	// already have a name mapping entry, since those may have custom overrides set above (e.g., for length).
	const nameField = "name"
	for resname, res := range prov.Resources {
		if schema := p.ResourcesMap[resname]; schema != nil {
			// Only apply auto-name to input properties (Optional || Required) named `name`
			if tfs, has := schema.Schema[nameField]; has && (tfs.Optional || tfs.Required) {
				if _, hasfield := res.Fields[nameField]; !hasfield {
					if res.Fields == nil {
						res.Fields = make(map[string]*tfbridge.SchemaInfo)
					}

					res.Fields[nameField] = tfbridge.AutoName(nameField, 255)
				}
			}
		}
	}

	return prov
}
