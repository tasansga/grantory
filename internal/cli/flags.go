package cli

const (
	FlagBackend   = "backend"
	FlagServerURL = "server-url"
	FlagToken     = "token"
	FlagUser      = "user"
	FlagPassword  = "password"
	EnvBackend    = "BACKEND"
	EnvServerURL  = "SERVER"
	EnvToken      = "TOKEN"
	EnvUser       = "USER"
	EnvPassword   = "PASSWORD"
	FlagNamespace = "namespace"
	EnvNamespace  = "NAMESPACE"
)

type backendMode string

const (
	backendModeDirect backendMode = "direct"
	backendModeAPI    backendMode = "api"
)

const (
	BackendModeDirect = backendModeDirect
	BackendModeAPI    = backendModeAPI
)
