package logical

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/license"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/helper/wrapping"
)

// SystemView exposes system configuration information in a safe way
// for logical backends to consume
type SystemView interface {
	// DefaultLeaseTTL returns the default lease TTL set in Vault configuration
	DefaultLeaseTTL() time.Duration

	// MaxLeaseTTL returns the max lease TTL set in Vault configuration; backend
	// authors should take care not to issue credentials that last longer than
	// this value, as Vault will revoke them
	MaxLeaseTTL() time.Duration

	// Returns true if the mount is tainted. A mount is tainted if it is in the
	// process of being unmounted. This should only be used in special
	// circumstances; a primary use-case is as a guard in revocation functions.
	// If revocation of a backend's leases fails it can keep the unmounting
	// process from being successful. If the reason for this failure is not
	// relevant when the mount is tainted (for instance, saving a CRL to disk
	// when the stored CRL will be removed during the unmounting process
	// anyways), we can ignore the errors to allow unmounting to complete.
	Tainted() bool

	// Returns true if caching is disabled. If true, no caches should be used,
	// despite known slowdowns.
	CachingDisabled() bool

	// When run from a system view attached to a request, indicates whether the
	// request is affecting a local mount or not
	LocalMount() bool

	// ReplicationState indicates the state of cluster replication
	ReplicationState() consts.ReplicationState

	// HasFeature returns true if the feature is currently enabled
	HasFeature(feature license.Features) bool

	// ResponseWrapData wraps the given data in a cubbyhole and returns the
	// token used to unwrap.
	ResponseWrapData(ctx context.Context, data map[string]interface{}, ttl time.Duration, jwt bool) (*wrapping.ResponseWrapInfo, error)

	// LookupPlugin looks into the plugin catalog for a plugin with the given
	// name. Returns a PluginRunner or an error if a plugin can not be found.
	LookupPlugin(context.Context, string, consts.PluginType) (*pluginutil.PluginRunner, error)

	NewPluginClient(ctx context.Context, pluginRunner *pluginutil.PluginRunner, logger log.Logger) (pluginutil.PluginClient, error)

	// MlockEnabled returns the configuration setting for enabling mlock on
	// plugins.
	MlockEnabled() bool

	// EntityInfo returns a subset of information related to the identity entity
	// for the given entity id
	EntityInfo(entityID string) (*Entity, error)

	// GroupsForEntity returns the group membership information for the provided
	// entity id
	GroupsForEntity(entityID string) ([]*Group, error)

	// PluginEnv returns Vault environment information used by plugins
	PluginEnv(context.Context) (*PluginEnvironment, error)

	// GeneratePasswordFromPolicy generates a password from the policy referenced.
	// If the policy does not exist, this will return an error.
	GeneratePasswordFromPolicy(ctx context.Context, policyName string) (password string, err error)
}

type PasswordPolicy interface {
	// Generate a random password
	Generate(context.Context, io.Reader) (string, error)
}

type ExtendedSystemView interface {
	Auditor() Auditor
	ForwardGenericRequest(context.Context, *Request) (*Response, error)
}

type PasswordGenerator func() (password string, err error)

type StaticSystemView struct {
	DefaultLeaseTTLVal  time.Duration
	MaxLeaseTTLVal      time.Duration
	SudoPrivilegeVal    bool
	TaintedVal          bool
	CachingDisabledVal  bool
	Primary             bool
	EnableMlock         bool
	LocalMountVal       bool
	ReplicationStateVal consts.ReplicationState
	EntityVal           *Entity
	GroupsVal           []*Group
	Features            license.Features
	VaultVersion        string
	PluginEnvironment   *PluginEnvironment
	PasswordPolicies    map[string]PasswordGenerator
}

type noopAuditor struct{}

func (a noopAuditor) AuditRequest(ctx context.Context, input *LogInput) error {
	return nil
}

func (a noopAuditor) AuditResponse(ctx context.Context, input *LogInput) error {
	return nil
}

func (d StaticSystemView) Auditor() Auditor {
	return noopAuditor{}
}

func (d StaticSystemView) ForwardGenericRequest(ctx context.Context, req *Request) (*Response, error) {
	return nil, errors.New("ForwardGenericRequest is not implemented in StaticSystemView")
}

func (d StaticSystemView) DefaultLeaseTTL() time.Duration {
	return d.DefaultLeaseTTLVal
}

func (d StaticSystemView) MaxLeaseTTL() time.Duration {
	return d.MaxLeaseTTLVal
}

func (d StaticSystemView) SudoPrivilege(_ context.Context, path string, token string) bool {
	return d.SudoPrivilegeVal
}

func (d StaticSystemView) Tainted() bool {
	return d.TaintedVal
}

func (d StaticSystemView) CachingDisabled() bool {
	return d.CachingDisabledVal
}

func (d StaticSystemView) LocalMount() bool {
	return d.LocalMountVal
}

func (d StaticSystemView) ReplicationState() consts.ReplicationState {
	return d.ReplicationStateVal
}

func (d StaticSystemView) NewPluginClient(ctx context.Context, pluginRunner *pluginutil.PluginRunner, logger log.Logger) (pluginutil.PluginClient, error) {
	return nil, errors.New("NewPluginClient is not implemented in StaticSystemView")
}

func (d StaticSystemView) ResponseWrapData(_ context.Context, data map[string]interface{}, ttl time.Duration, jwt bool) (*wrapping.ResponseWrapInfo, error) {
	return nil, errors.New("ResponseWrapData is not implemented in StaticSystemView")
}

func (d StaticSystemView) LookupPlugin(_ context.Context, _ string, _ consts.PluginType) (*pluginutil.PluginRunner, error) {
	return nil, errors.New("LookupPlugin is not implemented in StaticSystemView")
}

func (d StaticSystemView) MlockEnabled() bool {
	return d.EnableMlock
}

func (d StaticSystemView) EntityInfo(entityID string) (*Entity, error) {
	return d.EntityVal, nil
}

func (d StaticSystemView) GroupsForEntity(entityID string) ([]*Group, error) {
	return d.GroupsVal, nil
}

func (d StaticSystemView) HasFeature(feature license.Features) bool {
	return d.Features.HasFeature(feature)
}

func (d StaticSystemView) PluginEnv(_ context.Context) (*PluginEnvironment, error) {
	return d.PluginEnvironment, nil
}

func (d StaticSystemView) GeneratePasswordFromPolicy(ctx context.Context, policyName string) (password string, err error) {
	select {
	case <-ctx.Done():
		return "", fmt.Errorf("context timed out")
	default:
	}

	if d.PasswordPolicies == nil {
		return "", fmt.Errorf("password policy not found")
	}
	policy, exists := d.PasswordPolicies[policyName]
	if !exists {
		return "", fmt.Errorf("password policy not found")
	}
	return policy()
}

func (d *StaticSystemView) SetPasswordPolicy(name string, generator PasswordGenerator) {
	if d.PasswordPolicies == nil {
		d.PasswordPolicies = map[string]PasswordGenerator{}
	}
	d.PasswordPolicies[name] = generator
}

func (d *StaticSystemView) DeletePasswordPolicy(name string) (existed bool) {
	_, existed = d.PasswordPolicies[name]
	delete(d.PasswordPolicies, name)
	return existed
}
