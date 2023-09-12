package internal

// Config struct for SonarCloud Token and Exporter
type Config struct {
	ListenAddress string
	ListenPath    string
	Organization  string
	Token         string
	Metrics       string
}
