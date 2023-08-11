/*
   Copyright The containerd Authors.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package config

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/containerd/containerd/log"
	"github.com/containerd/containerd/plugin"
)

type SandboxControllerMode string

const (
	// ModePodSandbox means use Controller implementation from sbserver podsandbox package.
	// We take this one as a default mode.
	ModePodSandbox SandboxControllerMode = "podsandbox"
	// ModeShim means use whatever Controller implementation provided by shim.
	ModeShim SandboxControllerMode = "shim"
)

// Runtime struct to contain the type(ID), engine, and root variables for a default runtime
// and a runtime for untrusted workload.
type Runtime struct {
	// Type is the runtime type to use in containerd e.g. io.containerd.runtime.v1.linux
	Type string `toml:"runtime_type" json:"runtimeType"`
	// Path is an optional field that can be used to overwrite path to a shim runtime binary.
	// When specified, containerd will ignore runtime name field when resolving shim location.
	// Path must be abs.
	Path string `toml:"runtime_path" json:"runtimePath"`
	// Engine is the name of the runtime engine used by containerd.
	// This only works for runtime type "io.containerd.runtime.v1.linux".
	// DEPRECATED: use Options instead. Remove when shim v1 is deprecated.
	Engine string `toml:"runtime_engine" json:"runtimeEngine"`
	// PodAnnotations is a list of pod annotations passed to both pod sandbox as well as
	// container OCI annotations.
	PodAnnotations []string `toml:"pod_annotations" json:"PodAnnotations"`
	// ContainerAnnotations is a list of container annotations passed through to the OCI config of the containers.
	// Container annotations in CRI are usually generated by other Kubernetes node components (i.e., not users).
	// Currently, only device plugins populate the annotations.
	ContainerAnnotations []string `toml:"container_annotations" json:"ContainerAnnotations"`
	// Root is the directory used by containerd for runtime state.
	// DEPRECATED: use Options instead. Remove when shim v1 is deprecated.
	// This only works for runtime type "io.containerd.runtime.v1.linux".
	Root string `toml:"runtime_root" json:"runtimeRoot"`
	// Options are config options for the runtime.
	// If options is loaded from toml config, it will be map[string]interface{}.
	// Options can be converted into toml.Tree using toml.TreeFromMap().
	// Using options type as map[string]interface{} helps in correctly marshaling options from Go to JSON.
	Options map[string]interface{} `toml:"options" json:"options"`
	// PrivilegedWithoutHostDevices overloads the default behaviour for adding host devices to the
	// runtime spec when the container is privileged. Defaults to false.
	PrivilegedWithoutHostDevices bool `toml:"privileged_without_host_devices" json:"privileged_without_host_devices"`
	// PrivilegedWithoutHostDevicesAllDevicesAllowed overloads the default behaviour device allowlisting when
	// to the runtime spec when the container when PrivilegedWithoutHostDevices is already enabled. Requires
	// PrivilegedWithoutHostDevices to be enabled. Defaults to false.
	PrivilegedWithoutHostDevicesAllDevicesAllowed bool `toml:"privileged_without_host_devices_all_devices_allowed" json:"privileged_without_host_devices_all_devices_allowed"`
	// BaseRuntimeSpec is a json file with OCI spec to use as base spec that all container's will be created from.
	BaseRuntimeSpec string `toml:"base_runtime_spec" json:"baseRuntimeSpec"`
	// NetworkPluginConfDir is a directory containing the CNI network information for the runtime class.
	NetworkPluginConfDir string `toml:"cni_conf_dir" json:"cniConfDir"`
	// NetworkPluginMaxConfNum is the max number of plugin config files that will
	// be loaded from the cni config directory by go-cni. Set the value to 0 to
	// load all config files (no arbitrary limit). The legacy default value is 1.
	NetworkPluginMaxConfNum int `toml:"cni_max_conf_num" json:"cniMaxConfNum"`
	// Snapshotter setting snapshotter at runtime level instead of making it as a global configuration.
	// An example use case is to use devmapper or other snapshotters in Kata containers for performance and security
	// while using default snapshotters for operational simplicity.
	// See https://github.com/containerd/containerd/issues/6657 for details.
	Snapshotter string `toml:"snapshotter" json:"snapshotter"`
	// SandboxMode defines which sandbox runtime to use when scheduling pods
	// This features requires experimental CRI server to be enabled (use ENABLE_CRI_SANDBOXES=1)
	// shim - means use whatever Controller implementation provided by shim (e.g. use RemoteController).
	// podsandbox - means use Controller implementation from sbserver podsandbox package.
	SandboxMode string `toml:"sandbox_mode" json:"sandboxMode"`
}

// ContainerdConfig contains toml config related to containerd
type ContainerdConfig struct {
	// Snapshotter is the snapshotter used by containerd.
	Snapshotter string `toml:"snapshotter" json:"snapshotter"`
	// DefaultRuntimeName is the default runtime name to use from the runtimes table.
	DefaultRuntimeName string `toml:"default_runtime_name" json:"defaultRuntimeName"`
	// DefaultRuntime is the default runtime to use in containerd.
	// This runtime is used when no runtime handler (or the empty string) is provided.
	// DEPRECATED: use DefaultRuntimeName instead. Remove in containerd 1.4.
	DefaultRuntime Runtime `toml:"default_runtime" json:"defaultRuntime"`
	// UntrustedWorkloadRuntime is a runtime to run untrusted workloads on it.
	// DEPRECATED: use `untrusted` runtime in Runtimes instead. Remove in containerd 1.4.
	UntrustedWorkloadRuntime Runtime `toml:"untrusted_workload_runtime" json:"untrustedWorkloadRuntime"`
	// Runtimes is a map from CRI RuntimeHandler strings, which specify types of runtime
	// configurations, to the matching configurations.
	Runtimes map[string]Runtime `toml:"runtimes" json:"runtimes"`
	// NoPivot disables pivot-root (linux only), required when running a container in a RamDisk with runc
	// This only works for runtime type "io.containerd.runtime.v1.linux".
	NoPivot bool `toml:"no_pivot" json:"noPivot"`

	// DisableSnapshotAnnotations disables to pass additional annotations (image
	// related information) to snapshotters. These annotations are required by
	// stargz snapshotter (https://github.com/containerd/stargz-snapshotter).
	DisableSnapshotAnnotations bool `toml:"disable_snapshot_annotations" json:"disableSnapshotAnnotations"`

	// DiscardUnpackedLayers is a boolean flag to specify whether to allow GC to
	// remove layers from the content store after successfully unpacking these
	// layers to the snapshotter.
	DiscardUnpackedLayers bool `toml:"discard_unpacked_layers" json:"discardUnpackedLayers"`

	// IgnoreBlockIONotEnabledErrors is a boolean flag to ignore
	// blockio related errors when blockio support has not been
	// enabled.
	IgnoreBlockIONotEnabledErrors bool `toml:"ignore_blockio_not_enabled_errors" json:"ignoreBlockIONotEnabledErrors"`

	// IgnoreRdtNotEnabledErrors is a boolean flag to ignore RDT related errors
	// when RDT support has not been enabled.
	IgnoreRdtNotEnabledErrors bool `toml:"ignore_rdt_not_enabled_errors" json:"ignoreRdtNotEnabledErrors"`
}

// CniConfig contains toml config related to cni
type CniConfig struct {
	// NetworkPluginBinDir is the directory in which the binaries for the plugin is kept.
	NetworkPluginBinDir string `toml:"bin_dir" json:"binDir"`
	// NetworkPluginConfDir is the directory in which the admin places a CNI conf.
	NetworkPluginConfDir string `toml:"conf_dir" json:"confDir"`
	// NetworkPluginMaxConfNum is the max number of plugin config files that will
	// be loaded from the cni config directory by go-cni. Set the value to 0 to
	// load all config files (no arbitrary limit). The legacy default value is 1.
	NetworkPluginMaxConfNum int `toml:"max_conf_num" json:"maxConfNum"`
	// NetworkPluginSetupSerially is a boolean flag to specify whether containerd sets up networks serially
	// if there are multiple CNI plugin config files existing and NetworkPluginMaxConfNum is larger than 1.
	NetworkPluginSetupSerially bool `toml:"setup_serially" json:"setupSerially"`
	// NetworkPluginConfTemplate is the file path of golang template used to generate cni config.
	// When it is set, containerd will get cidr(s) from kubelet to replace {{.PodCIDR}},
	// {{.PodCIDRRanges}} or {{.Routes}} in the template, and write the config into
	// NetworkPluginConfDir.
	// Ideally the cni config should be placed by system admin or cni daemon like calico,
	// weaveworks etc. However, this is useful for the cases when there is no cni daemonset to place cni config.
	// This allowed for very simple generic networking using the Kubernetes built in node pod CIDR IPAM, avoiding the
	// need to fetch the node object through some external process (which has scalability, auth, complexity issues).
	// It is currently heavily used in kubernetes-containerd CI testing
	// NetworkPluginConfTemplate was once deprecated in containerd v1.7.0,
	// but its deprecation was cancelled in v1.7.3.
	NetworkPluginConfTemplate string `toml:"conf_template" json:"confTemplate"`
	// IPPreference specifies the strategy to use when selecting the main IP address for a pod.
	//
	// Options include:
	// * ipv4, "" - (default) select the first ipv4 address
	// * ipv6 - select the first ipv6 address
	// * cni - use the order returned by the CNI plugins, returning the first IP address from the results
	IPPreference string `toml:"ip_pref" json:"ipPref"`
}

// Mirror contains the config related to the registry mirror
type Mirror struct {
	// Endpoints are endpoints for a namespace. CRI plugin will try the endpoints
	// one by one until a working one is found. The endpoint must be a valid url
	// with host specified.
	// The scheme, host and path from the endpoint URL will be used.
	Endpoints []string `toml:"endpoint" json:"endpoint"`
}

// AuthConfig contains the config related to authentication to a specific registry
type AuthConfig struct {
	// Username is the username to login the registry.
	Username string `toml:"username" json:"username"`
	// Password is the password to login the registry.
	Password string `toml:"password" json:"password"`
	// Auth is a base64 encoded string from the concatenation of the username,
	// a colon, and the password.
	Auth string `toml:"auth" json:"auth"`
	// IdentityToken is used to authenticate the user and get
	// an access token for the registry.
	IdentityToken string `toml:"identitytoken" json:"identitytoken"`
}

// TLSConfig contains the CA/Cert/Key used for a registry
type TLSConfig struct {
	InsecureSkipVerify bool   `toml:"insecure_skip_verify" json:"insecure_skip_verify"`
	CAFile             string `toml:"ca_file" json:"caFile"`
	CertFile           string `toml:"cert_file" json:"certFile"`
	KeyFile            string `toml:"key_file" json:"keyFile"`
}

// Registry is registry settings configured
type Registry struct {
	// ConfigPath is a path to the root directory containing registry-specific
	// configurations.
	// If ConfigPath is set, the rest of the registry specific options are ignored.
	ConfigPath string `toml:"config_path" json:"configPath"`
	// Mirrors are namespace to mirror mapping for all namespaces.
	// This option will not be used when ConfigPath is provided.
	// DEPRECATED: Use ConfigPath instead. Remove in containerd 2.0.
	Mirrors map[string]Mirror `toml:"mirrors" json:"mirrors"`
	// Configs are configs for each registry.
	// The key is the domain name or IP of the registry.
	// DEPRECATED: Use ConfigPath instead.
	Configs map[string]RegistryConfig `toml:"configs" json:"configs"`
	// Auths are registry endpoint to auth config mapping. The registry endpoint must
	// be a valid url with host specified.
	// DEPRECATED: Use ConfigPath instead. Remove in containerd 1.6.
	Auths map[string]AuthConfig `toml:"auths" json:"auths"`
	// Headers adds additional HTTP headers that get sent to all registries
	Headers map[string][]string `toml:"headers" json:"headers"`
}

// RegistryConfig contains configuration used to communicate with the registry.
type RegistryConfig struct {
	// Auth contains information to authenticate to the registry.
	Auth *AuthConfig `toml:"auth" json:"auth"`
	// TLS is a pair of CA/Cert/Key which then are used when creating the transport
	// that communicates with the registry.
	// This field will not be used when ConfigPath is provided.
	// DEPRECATED: Use ConfigPath instead. Remove in containerd 1.7.
	TLS *TLSConfig `toml:"tls" json:"tls"`
}

// ImageDecryption contains configuration to handling decryption of encrypted container images.
type ImageDecryption struct {
	// KeyModel specifies the trust model of where keys should reside.
	//
	// Details of field usage can be found in:
	// https://github.com/containerd/containerd/tree/main/docs/cri/config.md
	//
	// Details of key models can be found in:
	// https://github.com/containerd/containerd/tree/main/docs/cri/decryption.md
	KeyModel string `toml:"key_model" json:"keyModel"`
}

// PluginConfig contains toml config related to CRI plugin,
// it is a subset of Config.
type PluginConfig struct {
	// ContainerdConfig contains config related to containerd
	ContainerdConfig `toml:"containerd" json:"containerd"`
	// CniConfig contains config related to cni
	CniConfig `toml:"cni" json:"cni"`
	// Registry contains config related to the registry
	Registry Registry `toml:"registry" json:"registry"`
	// ImageDecryption contains config related to handling decryption of encrypted container images
	ImageDecryption `toml:"image_decryption" json:"imageDecryption"`
	// DisableTCPService disables serving CRI on the TCP server.
	DisableTCPService bool `toml:"disable_tcp_service" json:"disableTCPService"`
	// StreamServerAddress is the ip address streaming server is listening on.
	StreamServerAddress string `toml:"stream_server_address" json:"streamServerAddress"`
	// StreamServerPort is the port streaming server is listening on.
	StreamServerPort string `toml:"stream_server_port" json:"streamServerPort"`
	// StreamIdleTimeout is the maximum time a streaming connection
	// can be idle before the connection is automatically closed.
	// The string is in the golang duration format, see:
	//   https://golang.org/pkg/time/#ParseDuration
	StreamIdleTimeout string `toml:"stream_idle_timeout" json:"streamIdleTimeout"`
	// EnableSelinux indicates to enable the selinux support.
	EnableSelinux bool `toml:"enable_selinux" json:"enableSelinux"`
	// SelinuxCategoryRange allows the upper bound on the category range to be set.
	// If not specified or set to 0, defaults to 1024 from the selinux package.
	SelinuxCategoryRange int `toml:"selinux_category_range" json:"selinuxCategoryRange"`
	// SandboxImage is the image used by sandbox container.
	SandboxImage string `toml:"sandbox_image" json:"sandboxImage"`
	// StatsCollectPeriod is the period (in seconds) of snapshots stats collection.
	StatsCollectPeriod int `toml:"stats_collect_period" json:"statsCollectPeriod"`
	// SystemdCgroup enables systemd cgroup support.
	// This only works for runtime type "io.containerd.runtime.v1.linux".
	// DEPRECATED: config runc runtime handler instead. Remove when shim v1 is deprecated.
	SystemdCgroup bool `toml:"systemd_cgroup" json:"systemdCgroup"`
	// EnableTLSStreaming indicates to enable the TLS streaming support.
	EnableTLSStreaming bool `toml:"enable_tls_streaming" json:"enableTLSStreaming"`
	// X509KeyPairStreaming is a x509 key pair used for TLS streaming
	X509KeyPairStreaming `toml:"x509_key_pair_streaming" json:"x509KeyPairStreaming"`
	// MaxContainerLogLineSize is the maximum log line size in bytes for a container.
	// Log line longer than the limit will be split into multiple lines. Non-positive
	// value means no limit.
	MaxContainerLogLineSize int `toml:"max_container_log_line_size" json:"maxContainerLogSize"`
	// DisableCgroup indicates to disable the cgroup support.
	// This is useful when the containerd does not have permission to access cgroup.
	DisableCgroup bool `toml:"disable_cgroup" json:"disableCgroup"`
	// DisableApparmor indicates to disable the apparmor support.
	// This is useful when the containerd does not have permission to access Apparmor.
	DisableApparmor bool `toml:"disable_apparmor" json:"disableApparmor"`
	// RestrictOOMScoreAdj indicates to limit the lower bound of OOMScoreAdj to the containerd's
	// current OOMScoreADj.
	// This is useful when the containerd does not have permission to decrease OOMScoreAdj.
	RestrictOOMScoreAdj bool `toml:"restrict_oom_score_adj" json:"restrictOOMScoreAdj"`
	// MaxConcurrentDownloads restricts the number of concurrent downloads for each image.
	MaxConcurrentDownloads int `toml:"max_concurrent_downloads" json:"maxConcurrentDownloads"`
	// DisableProcMount disables Kubernetes ProcMount support. This MUST be set to `true`
	// when using containerd with Kubernetes <=1.11.
	DisableProcMount bool `toml:"disable_proc_mount" json:"disableProcMount"`
	// UnsetSeccompProfile is the profile containerd/cri will use If the provided seccomp profile is
	// unset (`""`) for a container (default is `unconfined`)
	UnsetSeccompProfile string `toml:"unset_seccomp_profile" json:"unsetSeccompProfile"`
	// TolerateMissingHugetlbController if set to false will error out on create/update
	// container requests with huge page limits if the cgroup controller for hugepages is not present.
	// This helps with supporting Kubernetes <=1.18 out of the box. (default is `true`)
	TolerateMissingHugetlbController bool `toml:"tolerate_missing_hugetlb_controller" json:"tolerateMissingHugetlbController"`
	// DisableHugetlbController indicates to silently disable the hugetlb controller, even when it is
	// present in /sys/fs/cgroup/cgroup.controllers.
	// This helps with running rootless mode + cgroup v2 + systemd but without hugetlb delegation.
	DisableHugetlbController bool `toml:"disable_hugetlb_controller" json:"disableHugetlbController"`
	// DeviceOwnershipFromSecurityContext changes the default behavior of setting container devices uid/gid
	// from CRI's SecurityContext (RunAsUser/RunAsGroup) instead of taking host's uid/gid. Defaults to false.
	DeviceOwnershipFromSecurityContext bool `toml:"device_ownership_from_security_context" json:"device_ownership_from_security_context"`
	// IgnoreImageDefinedVolumes ignores volumes defined by the image. Useful for better resource
	// isolation, security and early detection of issues in the mount configuration when using
	// ReadOnlyRootFilesystem since containers won't silently mount a temporary volume.
	IgnoreImageDefinedVolumes bool `toml:"ignore_image_defined_volumes" json:"ignoreImageDefinedVolumes"`
	// NetNSMountsUnderStateDir places all mounts for network namespaces under StateDir/netns instead
	// of being placed under the hardcoded directory /var/run/netns. Changing this setting requires
	// that all containers are deleted.
	NetNSMountsUnderStateDir bool `toml:"netns_mounts_under_state_dir" json:"netnsMountsUnderStateDir"`
	// EnableUnprivilegedPorts configures net.ipv4.ip_unprivileged_port_start=0
	// for all containers which are not using host network
	// and if it is not overwritten by PodSandboxConfig
	// Note that currently default is set to disabled but target change it in future, see:
	//   https://github.com/kubernetes/kubernetes/issues/102612
	EnableUnprivilegedPorts bool `toml:"enable_unprivileged_ports" json:"enableUnprivilegedPorts"`
	// EnableUnprivilegedICMP configures net.ipv4.ping_group_range="0 2147483647"
	// for all containers which are not using host network, are not running in user namespace
	// and if it is not overwritten by PodSandboxConfig
	// Note that currently default is set to disabled but target change it in future together with EnableUnprivilegedPorts
	EnableUnprivilegedICMP bool `toml:"enable_unprivileged_icmp" json:"enableUnprivilegedICMP"`
	// EnableCDI indicates to enable injection of the Container Device Interface Specifications
	// into the OCI config
	// For more details about CDI and the syntax of CDI Spec files please refer to
	// https://github.com/container-orchestrated-devices/container-device-interface.
	EnableCDI bool `toml:"enable_cdi" json:"enableCDI"`
	// CDISpecDirs is the list of directories to scan for Container Device Interface Specifications
	// For more details about CDI configuration please refer to
	// https://github.com/container-orchestrated-devices/container-device-interface#containerd-configuration
	CDISpecDirs []string `toml:"cdi_spec_dirs" json:"cdiSpecDirs"`
	// ImagePullProgressTimeout is the maximum duration that there is no
	// image data read from image registry in the open connection. It will
	// be reset whatever a new byte has been read. If timeout, the image
	// pulling will be cancelled. A zero value means there is no timeout.
	//
	// The string is in the golang duration format, see:
	//   https://golang.org/pkg/time/#ParseDuration
	ImagePullProgressTimeout string `toml:"image_pull_progress_timeout" json:"imagePullProgressTimeout"`
	// DrainExecSyncIOTimeout is the maximum duration to wait for ExecSync
	// API' IO EOF event after exec init process exits. A zero value means
	// there is no timeout.
	//
	// The string is in the golang duration format, see:
	//   https://golang.org/pkg/time/#ParseDuration
	//
	// For example, the value can be '5h', '2h30m', '10s'.
	DrainExecSyncIOTimeout string `toml:"drain_exec_sync_io_timeout" json:"drainExecSyncIOTimeout"`
}

// X509KeyPairStreaming contains the x509 configuration for streaming
type X509KeyPairStreaming struct {
	// TLSCertFile is the path to a certificate file
	TLSCertFile string `toml:"tls_cert_file" json:"tlsCertFile"`
	// TLSKeyFile is the path to a private key file
	TLSKeyFile string `toml:"tls_key_file" json:"tlsKeyFile"`
}

// Config contains all configurations for cri server.
type Config struct {
	// PluginConfig is the config for CRI plugin.
	PluginConfig
	// ContainerdRootDir is the root directory path for containerd.
	ContainerdRootDir string `json:"containerdRootDir"`
	// ContainerdEndpoint is the containerd endpoint path.
	ContainerdEndpoint string `json:"containerdEndpoint"`
	// RootDir is the root directory path for managing cri plugin files
	// (metadata checkpoint etc.)
	RootDir string `json:"rootDir"`
	// StateDir is the root directory path for managing volatile pod/container data
	StateDir string `json:"stateDir"`
}

const (
	// RuntimeUntrusted is the implicit runtime defined for ContainerdConfig.UntrustedWorkloadRuntime
	RuntimeUntrusted = "untrusted"
	// RuntimeDefault is the implicit runtime defined for ContainerdConfig.DefaultRuntime
	RuntimeDefault = "default"
	// KeyModelNode is the key model where key for encrypted images reside
	// on the worker nodes
	KeyModelNode = "node"
)

// ValidatePluginConfig validates the given plugin configuration.
func ValidatePluginConfig(ctx context.Context, c *PluginConfig) error {
	if c.ContainerdConfig.Runtimes == nil {
		c.ContainerdConfig.Runtimes = make(map[string]Runtime)
	}

	// Validation for deprecated untrusted_workload_runtime.
	if c.ContainerdConfig.UntrustedWorkloadRuntime.Type != "" {
		log.G(ctx).Warning("`untrusted_workload_runtime` is deprecated, please use `untrusted` runtime in `runtimes` instead")
		if _, ok := c.ContainerdConfig.Runtimes[RuntimeUntrusted]; ok {
			return fmt.Errorf("conflicting definitions: configuration includes both `untrusted_workload_runtime` and `runtimes[%q]`", RuntimeUntrusted)
		}
		c.ContainerdConfig.Runtimes[RuntimeUntrusted] = c.ContainerdConfig.UntrustedWorkloadRuntime
	}

	// Validation for deprecated default_runtime field.
	if c.ContainerdConfig.DefaultRuntime.Type != "" {
		log.G(ctx).Warning("`default_runtime` is deprecated, please use `default_runtime_name` to reference the default configuration you have defined in `runtimes`")
		c.ContainerdConfig.DefaultRuntimeName = RuntimeDefault
		c.ContainerdConfig.Runtimes[RuntimeDefault] = c.ContainerdConfig.DefaultRuntime
	}

	// Validation for default_runtime_name
	if c.ContainerdConfig.DefaultRuntimeName == "" {
		return errors.New("`default_runtime_name` is empty")
	}
	if _, ok := c.ContainerdConfig.Runtimes[c.ContainerdConfig.DefaultRuntimeName]; !ok {
		return fmt.Errorf("no corresponding runtime configured in `containerd.runtimes` for `containerd` `default_runtime_name = \"%s\"", c.ContainerdConfig.DefaultRuntimeName)
	}

	// Validation for deprecated runtime options.
	if c.SystemdCgroup {
		if c.ContainerdConfig.Runtimes[c.ContainerdConfig.DefaultRuntimeName].Type != plugin.RuntimeLinuxV1 {
			return fmt.Errorf("`systemd_cgroup` only works for runtime %s", plugin.RuntimeLinuxV1)
		}
		log.G(ctx).Warning("`systemd_cgroup` is deprecated, please use runtime `options` instead")
	}
	if c.NoPivot {
		if c.ContainerdConfig.Runtimes[c.ContainerdConfig.DefaultRuntimeName].Type != plugin.RuntimeLinuxV1 {
			return fmt.Errorf("`no_pivot` only works for runtime %s", plugin.RuntimeLinuxV1)
		}
		// NoPivot can't be deprecated yet, because there is no alternative config option
		// for `io.containerd.runtime.v1.linux`.
	}
	for k, r := range c.ContainerdConfig.Runtimes {
		if r.Engine != "" {
			if r.Type != plugin.RuntimeLinuxV1 {
				return fmt.Errorf("`runtime_engine` only works for runtime %s", plugin.RuntimeLinuxV1)
			}
			log.G(ctx).Warning("`runtime_engine` is deprecated, please use runtime `options` instead")
		}
		if r.Root != "" {
			if r.Type != plugin.RuntimeLinuxV1 {
				return fmt.Errorf("`runtime_root` only works for runtime %s", plugin.RuntimeLinuxV1)
			}
			log.G(ctx).Warning("`runtime_root` is deprecated, please use runtime `options` instead")
		}
		if !r.PrivilegedWithoutHostDevices && r.PrivilegedWithoutHostDevicesAllDevicesAllowed {
			return errors.New("`privileged_without_host_devices_all_devices_allowed` requires `privileged_without_host_devices` to be enabled")
		}
		// If empty, use default podSandbox mode
		if len(r.SandboxMode) == 0 {
			r.SandboxMode = string(ModePodSandbox)
			c.ContainerdConfig.Runtimes[k] = r
		}
	}

	useConfigPath := c.Registry.ConfigPath != ""
	if len(c.Registry.Mirrors) > 0 {
		if useConfigPath {
			return errors.New("`mirrors` cannot be set when `config_path` is provided")
		}
		log.G(ctx).Warning("`mirrors` is deprecated, please use `config_path` instead")
	}
	var hasDeprecatedTLS bool
	for _, r := range c.Registry.Configs {
		if r.TLS != nil {
			hasDeprecatedTLS = true
			break
		}
	}
	if hasDeprecatedTLS {
		if useConfigPath {
			return errors.New("`configs.tls` cannot be set when `config_path` is provided")
		}
		log.G(ctx).Warning("`configs.tls` is deprecated, please use `config_path` instead")
	}

	// Validation for deprecated auths options and mapping it to configs.
	if len(c.Registry.Auths) != 0 {
		if c.Registry.Configs == nil {
			c.Registry.Configs = make(map[string]RegistryConfig)
		}
		for endpoint, auth := range c.Registry.Auths {
			auth := auth
			u, err := url.Parse(endpoint)
			if err != nil {
				return fmt.Errorf("failed to parse registry url %q from `registry.auths`: %w", endpoint, err)
			}
			if u.Scheme != "" {
				// Do not include the scheme in the new registry config.
				endpoint = u.Host
			}
			config := c.Registry.Configs[endpoint]
			config.Auth = &auth
			c.Registry.Configs[endpoint] = config
		}
		log.G(ctx).Warning("`auths` is deprecated, please use `configs` instead")
	}

	// Validation for stream_idle_timeout
	if c.StreamIdleTimeout != "" {
		if _, err := time.ParseDuration(c.StreamIdleTimeout); err != nil {
			return fmt.Errorf("invalid stream idle timeout: %w", err)
		}
	}

	// Validation for image_pull_progress_timeout
	if c.ImagePullProgressTimeout != "" {
		if _, err := time.ParseDuration(c.ImagePullProgressTimeout); err != nil {
			return fmt.Errorf("invalid image pull progress timeout: %w", err)
		}
	}

	// Validation for drain_exec_sync_io_timeout
	if c.DrainExecSyncIOTimeout != "" {
		if _, err := time.ParseDuration(c.DrainExecSyncIOTimeout); err != nil {
			return fmt.Errorf("invalid `drain_exec_sync_io_timeout`: %w", err)
		}
	}
	return nil
}
