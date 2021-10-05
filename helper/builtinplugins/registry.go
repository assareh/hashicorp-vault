package builtinplugins

import (
	credAliCloud "github.com/hashicorp/vault-plugin-auth-alicloud"
	credAzure "github.com/hashicorp/vault-plugin-auth-azure"
	credCentrify "github.com/hashicorp/vault-plugin-auth-centrify"
	credCF "github.com/hashicorp/vault-plugin-auth-cf"
	credGcp "github.com/hashicorp/vault-plugin-auth-gcp/plugin"
	credJWT "github.com/hashicorp/vault-plugin-auth-jwt"
	credKerb "github.com/hashicorp/vault-plugin-auth-kerberos"
	credKube "github.com/hashicorp/vault-plugin-auth-kubernetes"
	credOCI "github.com/hashicorp/vault-plugin-auth-oci"
	logicalAd "github.com/hashicorp/vault-plugin-secrets-ad/plugin"
	logicalAlicloud "github.com/hashicorp/vault-plugin-secrets-alicloud"
	logicalAzure "github.com/hashicorp/vault-plugin-secrets-azure"
	logicalGcp "github.com/hashicorp/vault-plugin-secrets-gcp/plugin"
	logicalGcpKms "github.com/hashicorp/vault-plugin-secrets-gcpkms"
	logicalKv "github.com/hashicorp/vault-plugin-secrets-kv"
	logicalMongoAtlas "github.com/hashicorp/vault-plugin-secrets-mongodbatlas"
	logicalOpenLDAP "github.com/hashicorp/vault-plugin-secrets-openldap"
	logicalTerraform "github.com/hashicorp/vault-plugin-secrets-terraform"
	credAppId "github.com/hashicorp/vault/builtin/credential/app-id"
	credAppRole "github.com/hashicorp/vault/builtin/credential/approle"
	credAws "github.com/hashicorp/vault/builtin/credential/aws"
	credCert "github.com/hashicorp/vault/builtin/credential/cert"
	credGitHub "github.com/hashicorp/vault/builtin/credential/github"
	credLdap "github.com/hashicorp/vault/builtin/credential/ldap"
	credOkta "github.com/hashicorp/vault/builtin/credential/okta"
	credRadius "github.com/hashicorp/vault/builtin/credential/radius"
	credUserpass "github.com/hashicorp/vault/builtin/credential/userpass"
	logicalAws "github.com/hashicorp/vault/builtin/logical/aws"
	logicalCass "github.com/hashicorp/vault/builtin/logical/cassandra"
	logicalConsul "github.com/hashicorp/vault/builtin/logical/consul"
	logicalMongo "github.com/hashicorp/vault/builtin/logical/mongodb"
	logicalMssql "github.com/hashicorp/vault/builtin/logical/mssql"
	logicalMysql "github.com/hashicorp/vault/builtin/logical/mysql"
	logicalNomad "github.com/hashicorp/vault/builtin/logical/nomad"
	logicalPki "github.com/hashicorp/vault/builtin/logical/pki"
	logicalPostgres "github.com/hashicorp/vault/builtin/logical/postgresql"
	logicalRabbit "github.com/hashicorp/vault/builtin/logical/rabbitmq"
	logicalSsh "github.com/hashicorp/vault/builtin/logical/ssh"
	logicalTotp "github.com/hashicorp/vault/builtin/logical/totp"
	logicalTransit "github.com/hashicorp/vault/builtin/logical/transit"
	dbCass "github.com/hashicorp/vault/plugins/database/cassandra"
	dbHana "github.com/hashicorp/vault/plugins/database/hana"
	dbInflux "github.com/hashicorp/vault/plugins/database/influxdb"
	dbMongo "github.com/hashicorp/vault/plugins/database/mongodb"
	dbMssql "github.com/hashicorp/vault/plugins/database/mssql"
	dbMysql "github.com/hashicorp/vault/plugins/database/mysql"
	dbPostgres "github.com/hashicorp/vault/plugins/database/postgresql"
	dbplugin "github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
)

// Registry is inherently thread-safe because it's immutable.
// Thus, rather than creating multiple instances of it, we only need one.
var Registry = newRegistry()

var addExternalPlugins = addExtPluginsImpl

func newRegistry() *registry {
	reg := &registry{
		credentialBackends: map[string]logical.Factory{
			"alicloud":   credAliCloud.Factory,
			"app-id":     credAppId.Factory,
			"approle":    credAppRole.Factory,
			"aws":        credAws.Factory,
			"azure":      credAzure.Factory,
			"centrify":   credCentrify.Factory,
			"cert":       credCert.Factory,
			"cf":         credCF.Factory,
			"gcp":        credGcp.Factory,
			"github":     credGitHub.Factory,
			"jwt":        credJWT.Factory,
			"kerberos":   credKerb.Factory,
			"kubernetes": credKube.Factory,
			"ldap":       credLdap.Factory,
			"oci":        credOCI.Factory,
			"oidc":       credJWT.Factory,
			"okta":       credOkta.Factory,
			"pcf":        credCF.Factory, // Deprecated.
			"radius":     credRadius.Factory,
			"userpass":   credUserpass.Factory,
		},
		databasePlugins: map[string]dbplugin.Factory{
			// These four plugins all use the same mysql implementation but with
			// different username settings passed by the constructor.
			"mysql-database-plugin":        dbMysql.New(dbMysql.DefaultUserNameTemplate),
			"mysql-aurora-database-plugin": dbMysql.New(dbMysql.DefaultLegacyUserNameTemplate),
			"mysql-rds-database-plugin":    dbMysql.New(dbMysql.DefaultLegacyUserNameTemplate),
			"mysql-legacy-database-plugin": dbMysql.New(dbMysql.DefaultLegacyUserNameTemplate),

			"cassandra-database-plugin": dbCass.New,
			// JASON: TODO
			//"couchbase-database-plugin":     dbCouchbase.New,
			// JASON: TODO
			//"elasticsearch-database-plugin": dbElastic.New,
			"hana-database-plugin":     dbHana.New,
			"influxdb-database-plugin": dbInflux.New,
			"mongodb-database-plugin":  dbMongo.New,
			// JASON: TODO
			//"mongodbatlas-database-plugin":  dbMongoAtlas.New,
			"mssql-database-plugin":      dbMssql.New,
			"postgresql-database-plugin": dbPostgres.New,
			// JASON: TODO
			//"redshift-database-plugin":      dbRedshift.New,
			// JASON: TODO
			//"snowflake-database-plugin":     dbSnowflake.New,
		},
		logicalBackends: map[string]logical.Factory{
			"ad":           logicalAd.Factory,
			"alicloud":     logicalAlicloud.Factory,
			"aws":          logicalAws.Factory,
			"azure":        logicalAzure.Factory,
			"cassandra":    logicalCass.Factory, // Deprecated
			"consul":       logicalConsul.Factory,
			"gcp":          logicalGcp.Factory,
			"gcpkms":       logicalGcpKms.Factory,
			"kv":           logicalKv.Factory,
			"mongodb":      logicalMongo.Factory, // Deprecated
			"mongodbatlas": logicalMongoAtlas.Factory,
			"mssql":        logicalMssql.Factory, // Deprecated
			"mysql":        logicalMysql.Factory, // Deprecated
			"nomad":        logicalNomad.Factory,
			"openldap":     logicalOpenLDAP.Factory,
			"pki":          logicalPki.Factory,
			"postgresql":   logicalPostgres.Factory, // Deprecated
			"rabbitmq":     logicalRabbit.Factory,
			"ssh":          logicalSsh.Factory,
			"terraform":    logicalTerraform.Factory,
			"totp":         logicalTotp.Factory,
			"transit":      logicalTransit.Factory,
		},
	}

	addExternalPlugins(reg)

	return reg
}

func addExtPluginsImpl(r *registry) {}

type registry struct {
	credentialBackends map[string]logical.Factory
	databasePlugins    map[string]dbplugin.Factory
	logicalBackends    map[string]logical.Factory
}

// Get returns the BuiltinFactory func for a particular backend plugin
// from the plugins map.
func (r *registry) Get(name string, pluginType consts.PluginType) (func() (interface{}, error), bool) {
	switch pluginType {
	case consts.PluginTypeCredential:
		f, ok := r.credentialBackends[name]
		return toFunc(f), ok
	case consts.PluginTypeSecrets:
		f, ok := r.logicalBackends[name]
		return toFunc(f), ok
	case consts.PluginTypeDatabase:
		f, ok := r.databasePlugins[name]
		return toFunc(f), ok
	default:
		return nil, false
	}
}

// Keys returns the list of plugin names that are considered builtin plugins.
func (r *registry) Keys(pluginType consts.PluginType) []string {
	var keys []string
	switch pluginType {
	case consts.PluginTypeDatabase:
		for key := range r.databasePlugins {
			keys = append(keys, key)
		}
	case consts.PluginTypeCredential:
		for key := range r.credentialBackends {
			keys = append(keys, key)
		}
	case consts.PluginTypeSecrets:
		for key := range r.logicalBackends {
			keys = append(keys, key)
		}
	}
	return keys
}

func (r *registry) Contains(name string, pluginType consts.PluginType) bool {
	for _, key := range r.Keys(pluginType) {
		if key == name {
			return true
		}
	}
	return false
}

func toFunc(ifc interface{}) func() (interface{}, error) {
	return func() (interface{}, error) {
		return ifc, nil
	}
}
