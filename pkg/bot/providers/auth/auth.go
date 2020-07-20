package auth

// Config defines auth configuration for github and docker.
type Config struct {
	GithubToken    string
	GithubUsername string
	DockerToken    string
	DockerUsername string
}
