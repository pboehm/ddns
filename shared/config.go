package shared

type Config struct {
	Verbose            bool
	Domain             string
	SOAFqdn            string
	HostExpirationDays int
	Listen             string
	RedisHost          string
}
