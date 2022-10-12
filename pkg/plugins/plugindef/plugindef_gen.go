// This file is autogenerated. DO NOT EDIT.
//
// Generated by "thema lineage gen" from lineage defined in pkg/plugins/plugindef/plugindef.cue

package plugindef

// Defines values for PlugindefCategory.
const (
	Cloud      PlugindefCategory = "cloud"
	Enterprise PlugindefCategory = "enterprise"
	Logging    PlugindefCategory = "logging"
	Other      PlugindefCategory = "other"
	Sql        PlugindefCategory = "sql"
	Tracing    PlugindefCategory = "tracing"
	Tsdb       PlugindefCategory = "tsdb"
)

// Defines values for PlugindefType.
const (
	PlugindefTypeApp            PlugindefType = "app"
	PlugindefTypeDatasource     PlugindefType = "datasource"
	PlugindefTypePanel          PlugindefType = "panel"
	PlugindefTypeRenderer       PlugindefType = "renderer"
	PlugindefTypeSecretsmanager PlugindefType = "secretsmanager"
)

// Defines values for PlugindefDependencyType.
const (
	PlugindefDependencyTypeApp        PlugindefDependencyType = "app"
	PlugindefDependencyTypeDatasource PlugindefDependencyType = "datasource"
	PlugindefDependencyTypePanel      PlugindefDependencyType = "panel"
)

// Defines values for PlugindefIncludeRole.
const (
	Admin  PlugindefIncludeRole = "Admin"
	Editor PlugindefIncludeRole = "Editor"
	Viewer PlugindefIncludeRole = "Viewer"
)

// Defines values for PlugindefIncludeType.
const (
	App            PlugindefIncludeType = "app"
	Dashboard      PlugindefIncludeType = "dashboard"
	Datasource     PlugindefIncludeType = "datasource"
	Page           PlugindefIncludeType = "page"
	Panel          PlugindefIncludeType = "panel"
	Renderer       PlugindefIncludeType = "renderer"
	Secretsmanager PlugindefIncludeType = "secretsmanager"
)

// Defines values for PlugindefReleaseState.
const (
	Alpha      PlugindefReleaseState = "alpha"
	Beta       PlugindefReleaseState = "beta"
	Deprecated PlugindefReleaseState = "deprecated"
	Stable     PlugindefReleaseState = "stable"
)

// Plugindef defines model for plugindef.
type Plugindef struct {
	// For data source plugins, if the plugin supports alerting.
	Alerting *bool `json:"alerting,omitempty"`

	// For data source plugins, if the plugin supports annotation
	// queries.
	Annotations *bool `json:"annotations,omitempty"`

	// Set to true for app plugins that should be enabled by default
	// in all orgs
	AutoEnabled *bool `json:"autoEnabled,omitempty"`

	// If the plugin has a backend component.
	Backend *bool `json:"backend,omitempty"`

	// builtin indicates whether the plugin is developed and shipped as part
	// of Grafana. Also known as a "core plugin."
	BuiltIn bool `json:"builtIn"`

	// Plugin category used on the Add data source page.
	Category *PlugindefCategory `json:"category,omitempty"`

	// Dependencies needed by the plugin.
	Dependencies struct {
		// Required Grafana version for this plugin. Validated using
		// https://github.com/npm/node-semver.
		GrafanaDependency string `json:"grafanaDependency"`

		// (Deprecated) Required Grafana version for this plugin, e.g.
		// `6.x.x 7.x.x` to denote plugin requires Grafana v6.x.x or
		// v7.x.x.
		GrafanaVersion *string `json:"grafanaVersion,omitempty"`

		// An array of required plugins on which this plugin depends.
		Plugins *[]PlugindefDependency `json:"plugins,omitempty"`
	} `json:"dependencies"`

	// Grafana Enerprise specific features.
	EnterpriseFeatures *struct {
		// Enable/Disable health diagnostics errors. Requires Grafana
		// >=7.5.5.
		HealthDiagnosticsErrors *bool `json:"healthDiagnosticsErrors,omitempty"`
	} `json:"enterpriseFeatures,omitempty"`

	// The first part of the file name of the backend component
	// executable. There can be multiple executables built for
	// different operating system and architecture. Grafana will
	// check for executables named `<executable>_<$GOOS>_<lower case
	// $GOARCH><.exe for Windows>`, e.g. `plugin_linux_amd64`.
	// Combination of $GOOS and $GOARCH can be found here:
	// https://golang.org/doc/install/source#environment.
	Executable *string `json:"executable,omitempty"`

	// For data source plugins, include hidden queries in the data
	// request.
	HiddenQueries *bool `json:"hiddenQueries,omitempty"`

	// hideFromList excludes the plugin from listings in Grafana's UI. Only
	// allowed for builtin plugins.
	HideFromList bool `json:"hideFromList"`

	// Unique name of the plugin. If the plugin is published on
	// grafana.com, then the plugin id has to follow the naming
	// conventions.
	Id string `json:"id"`

	// Resources to include in plugin.
	Includes *[]PlugindefInclude `json:"includes,omitempty"`

	// Metadata about the plugin.
	Info struct {
		// Information about the plugin author.
		Author *struct {
			// Author's name.
			Email *string `json:"email,omitempty"`

			// Author's name.
			Name *string `json:"name,omitempty"`

			// Link to author's website.
			Url *string `json:"url,omitempty"`
		} `json:"author,omitempty"`
		Build *PlugindefBuildInfo `json:"build,omitempty"`

		// Description of plugin. Used on the plugins page in Grafana and
		// for search on grafana.com.
		Description *string `json:"description,omitempty"`

		// Array of plugin keywords. Used for search on grafana.com.
		Keywords []string `json:"keywords"`

		// An array of link objects to be displayed on this plugin's
		// project page in the form `{name: 'foo', url:
		// 'http://example.com'}`
		Links *[]struct {
			Name *string `json:"name,omitempty"`
			Url  *string `json:"url,omitempty"`
		} `json:"links,omitempty"`

		// SVG images that are used as plugin icons.
		Logos *struct {
			// Link to the "large" version of the plugin logo, which must be
			// an SVG image. "Large" and "small" logos can be the same image.
			Large string `json:"large"`

			// Link to the "small" version of the plugin logo, which must be
			// an SVG image. "Large" and "small" logos can be the same image.
			Small string `json:"small"`
		} `json:"logos,omitempty"`

		// An array of screenshot objects in the form `{name: 'bar', path:
		// 'img/screenshot.png'}`
		Screenshots *[]struct {
			Name *string `json:"name,omitempty"`
			Path *string `json:"path,omitempty"`
		} `json:"screenshots,omitempty"`

		// Date when this plugin was built.
		Updated *string `json:"updated,omitempty"`

		// Project version of this commit, e.g. `6.7.x`.
		Version *string `json:"version,omitempty"`
	} `json:"info"`

	// For data source plugins, if the plugin supports logs.
	Logs *bool `json:"logs,omitempty"`

	// For data source plugins, if the plugin supports metric queries.
	// Used in Explore.
	Metrics *bool `json:"metrics,omitempty"`

	// Human-readable name of the plugin that is shown to the user in
	// the UI.
	Name string `json:"name"`

	// Initialize plugin on startup. By default, the plugin
	// initializes on first use.
	Preload *bool `json:"preload,omitempty"`

	// For data source plugins. There is a query options section in
	// the plugin's query editor and these options can be turned on
	// if needed.
	QueryOptions *struct {
		// For data source plugins. If the `cache timeout` option should
		// be shown in the query options section in the query editor.
		CacheTimeout *bool `json:"cacheTimeout,omitempty"`

		// For data source plugins. If the `max data points` option should
		// be shown in the query options section in the query editor.
		MaxDataPoints *bool `json:"maxDataPoints,omitempty"`

		// For data source plugins. If the `min interval` option should be
		// shown in the query options section in the query editor.
		MinInterval *bool `json:"minInterval,omitempty"`
	} `json:"queryOptions,omitempty"`

	// Routes is a list of proxy routes, if any. For datasource plugins only.
	Routes *[]PlugindefRoute `json:"routes,omitempty"`

	// For panel plugins. Hides the query editor.
	SkipDataQuery *bool `json:"skipDataQuery,omitempty"`

	// ReleaseState indicates release maturity state of a plugin.
	State *PlugindefReleaseState `json:"state,omitempty"`

	// For data source plugins, if the plugin supports streaming.
	Streaming *bool `json:"streaming,omitempty"`

	// This is an undocumented feature.
	Tables *bool `json:"tables,omitempty"`

	// For data source plugins, if the plugin supports tracing.
	Tracing *bool `json:"tracing,omitempty"`

	// type indicates which type of Grafana plugin this is, of the defined
	// set of Grafana plugin types.
	Type PlugindefType `json:"type"`
}

// Plugin category used on the Add data source page.
type PlugindefCategory string

// type indicates which type of Grafana plugin this is, of the defined
// set of Grafana plugin types.
type PlugindefType string

// PlugindefBuildInfo defines model for plugindef.BuildInfo.
type PlugindefBuildInfo struct {
	// Git branch the plugin was built from.
	Branch *string `json:"branch,omitempty"`

	// Git hash of the commit the plugin was built from
	Hash   *string `json:"hash,omitempty"`
	Number *int64  `json:"number,omitempty"`

	// GitHub pull request the plugin was built from
	Pr   *int32  `json:"pr,omitempty"`
	Repo *string `json:"repo,omitempty"`

	// Time when the plugin was built, as a Unix timestamp.
	Time *int64 `json:"time,omitempty"`
}

// PlugindefDependencies defines model for plugindef.Dependencies.
type PlugindefDependencies struct {
	// Required Grafana version for this plugin. Validated using
	// https://github.com/npm/node-semver.
	GrafanaDependency string `json:"grafanaDependency"`

	// (Deprecated) Required Grafana version for this plugin, e.g.
	// `6.x.x 7.x.x` to denote plugin requires Grafana v6.x.x or
	// v7.x.x.
	GrafanaVersion *string `json:"grafanaVersion,omitempty"`

	// An array of required plugins on which this plugin depends.
	Plugins *[]PlugindefDependency `json:"plugins,omitempty"`
}

// Dependency describes another plugin on which a plugin depends.
// The id refers to the plugin package identifier, as given on
// the grafana.com plugin marketplace.
type PlugindefDependency struct {
	Id      string                  `json:"id"`
	Name    string                  `json:"name"`
	Type    PlugindefDependencyType `json:"type"`
	Version string                  `json:"version"`
}

// PlugindefDependencyType defines model for PlugindefDependency.Type.
type PlugindefDependencyType string

// Header describes an HTTP header that is forwarded with a proxied request for
// a plugin route.
type PlugindefHeader struct {
	Content string `json:"content"`
	Name    string `json:"name"`
}

// A resource to be included in a plugin.
type PlugindefInclude struct {
	// Add the include to the side menu.
	AddToNav *bool `json:"addToNav,omitempty"`

	// (Legacy) The Angular component to use for a page.
	Component *string `json:"component,omitempty"`

	// Page or dashboard when user clicks the icon in the side menu.
	DefaultNav *bool `json:"defaultNav,omitempty"`

	// Icon to use in the side menu. For information on available
	// icon, refer to [Icons
	// Overview](https://developers.grafana.com/ui/latest/index.html?path=/story/docs-overview-icon--icons-overview).
	Icon *string `json:"icon,omitempty"`
	Name *string `json:"name,omitempty"`

	// Used for app plugins.
	Path *string               `json:"path,omitempty"`
	Role *PlugindefIncludeRole `json:"role,omitempty"`
	Type PlugindefIncludeType  `json:"type"`

	// Unique identifier of the included resource
	Uid *string `json:"uid,omitempty"`
}

// PlugindefIncludeRole defines model for PlugindefInclude.Role.
type PlugindefIncludeRole string

// PlugindefIncludeType defines model for PlugindefInclude.Type.
type PlugindefIncludeType string

// Metadata about a Grafana plugin. Some fields are used on the plugins
// page in Grafana and others on grafana.com, if the plugin is published.
type PlugindefInfo struct {
	// Information about the plugin author.
	Author *struct {
		// Author's name.
		Email *string `json:"email,omitempty"`

		// Author's name.
		Name *string `json:"name,omitempty"`

		// Link to author's website.
		Url *string `json:"url,omitempty"`
	} `json:"author,omitempty"`
	Build *PlugindefBuildInfo `json:"build,omitempty"`

	// Description of plugin. Used on the plugins page in Grafana and
	// for search on grafana.com.
	Description *string `json:"description,omitempty"`

	// Array of plugin keywords. Used for search on grafana.com.
	Keywords []string `json:"keywords"`

	// An array of link objects to be displayed on this plugin's
	// project page in the form `{name: 'foo', url:
	// 'http://example.com'}`
	Links *[]struct {
		Name *string `json:"name,omitempty"`
		Url  *string `json:"url,omitempty"`
	} `json:"links,omitempty"`

	// SVG images that are used as plugin icons.
	Logos *struct {
		// Link to the "large" version of the plugin logo, which must be
		// an SVG image. "Large" and "small" logos can be the same image.
		Large string `json:"large"`

		// Link to the "small" version of the plugin logo, which must be
		// an SVG image. "Large" and "small" logos can be the same image.
		Small string `json:"small"`
	} `json:"logos,omitempty"`

	// An array of screenshot objects in the form `{name: 'bar', path:
	// 'img/screenshot.png'}`
	Screenshots *[]struct {
		Name *string `json:"name,omitempty"`
		Path *string `json:"path,omitempty"`
	} `json:"screenshots,omitempty"`

	// Date when this plugin was built.
	Updated *string `json:"updated,omitempty"`

	// Project version of this commit, e.g. `6.7.x`.
	Version *string `json:"version,omitempty"`
}

// TODO docs
// TODO should this really be separate from TokenAuth?
type PlugindefJWTTokenAuth struct {
	// Parameters for the JWT token authentication request.
	Params map[string]interface{} `json:"params"`

	// The list of scopes that your application should be granted
	// access to.
	Scopes []string `json:"scopes"`

	// URL to fetch the JWT token.
	Url string `json:"url"`
}

// ReleaseState indicates release maturity state of a plugin.
type PlugindefReleaseState string

// A proxy route used in datasource plugins for plugin authentication
// and adding headers to HTTP requests made by the plugin.
// For more information, refer to [Authentication for data source
// plugins](https://grafana.com/docs/grafana/latest/developers/plugins/authentication/).
type PlugindefRoute struct {
	// For data source plugins. Route headers set the body content and
	// length to the proxied request.
	Body *map[string]interface{} `json:"body,omitempty"`

	// For data source plugins. Route headers adds HTTP headers to the
	// proxied request.
	Headers *[]PlugindefHeader `json:"headers,omitempty"`

	// TODO docs
	// TODO should this really be separate from TokenAuth?
	JwtTokenAuth *PlugindefJWTTokenAuth `json:"jwtTokenAuth,omitempty"`

	// For data source plugins. Route method matches the HTTP verb
	// like GET or POST. Multiple methods can be provided as a
	// comma-separated list.
	Method *string `json:"method,omitempty"`

	// For data source plugins. The route path that is replaced by the
	// route URL field when proxying the call.
	Path        *string `json:"path,omitempty"`
	ReqRole     *string `json:"reqRole,omitempty"`
	ReqSignedIn *bool   `json:"reqSignedIn,omitempty"`

	// TODO docs
	TokenAuth *PlugindefTokenAuth `json:"tokenAuth,omitempty"`

	// For data source plugins. Route URL is where the request is
	// proxied to.
	Url       *string              `json:"url,omitempty"`
	UrlParams *[]PlugindefURLParam `json:"urlParams,omitempty"`
}

// TODO docs
type PlugindefTokenAuth struct {
	// Parameters for the token authentication request.
	Params map[string]interface{} `json:"params"`

	// The list of scopes that your application should be granted
	// access to.
	Scopes *[]string `json:"scopes,omitempty"`

	// URL to fetch the authentication token.
	Url *string `json:"url,omitempty"`
}

// URLParam describes query string parameters for
// a url in a plugin route
type PlugindefURLParam struct {
	Content string `json:"content"`
	Name    string `json:"name"`
}
