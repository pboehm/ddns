package shared

type Config struct {
	Verbose            bool
	Domain             string
	SOAFqdn            string
	HostExpirationDays int
	FrontendListen     string
	BackendListen      string
	RedisHost          string
}
