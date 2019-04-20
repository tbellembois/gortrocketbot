package rocket

// Config is the Rocket server configuration
type Config struct {
	// ServerHost is the Rocket server host name or IP without the scheme
	ServerHost string
	// ServerScheme is the Rocket server URL scheme such as https
	ServerScheme string
	// User is the bot Rocket username to connect with
	User string
	// Email is the bot Rocket email
	Email string
	// Password is the bot Rocket Password to connect with
	Password string
	// Debug enable logs
	Debug bool
}

// Plugin is a bot plugin
// executing "CommandFunc" with
// the given "Args"
type Plugin struct {
	Name        string
	CommandFunc func(...string) string
	Args        []string
	Help        string
}
