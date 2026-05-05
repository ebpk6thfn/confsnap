package ssh

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

// Client wraps an SSH connection to a remote host.
type Client struct {
	host   string
	port   int
	client *ssh.Client
}

// Config holds the parameters needed to establish an SSH connection.
type Config struct {
	Host       string
	Port       int
	User       string
	KeyPath    string
	Timeout    time.Duration
}

// NewClient dials an SSH connection using the provided Config.
func NewClient(cfg Config) (*Client, error) {
	key, err := os.ReadFile(cfg.KeyPath)
	if err != nil {
		return nil, fmt.Errorf("read identity file %q: %w", cfg.KeyPath, err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("parse private key: %w", err)
	}

	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 15 * time.Second
	}

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	sshCfg := &ssh.ClientConfig{
		User:            cfg.User,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // TODO: use known_hosts
		Timeout:         timeout,
	}

	c, err := ssh.Dial("tcp", addr, sshCfg)
	if err != nil {
		return nil, fmt.Errorf("dial %s: %w", addr, err)
	}

	return &Client{host: cfg.Host, port: cfg.Port, client: c}, nil
}

// ReadFile retrieves the contents of a remote file over SSH.
func (c *Client) ReadFile(path string) ([]byte, error) {
	sess, err := c.client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("new session: %w", err)
	}
	defer sess.Close()

	out, err := sess.Output(fmt.Sprintf("cat -- %q", path))
	if err != nil {
		return nil, fmt.Errorf("read remote file %q on %s: %w", path, c.host, err)
	}
	return out, nil
}

// RunCommand executes a command on the remote host and returns its combined
// stdout output. The session is closed automatically when the command finishes.
func (c *Client) RunCommand(cmd string) ([]byte, error) {
	sess, err := c.client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("new session: %w", err)
	}
	defer sess.Close()

	out, err := sess.Output(cmd)
	if err != nil {
		return nil, fmt.Errorf("run command %q on %s: %w", cmd, c.host, err)
	}
	return out, nil
}

// Close terminates the underlying SSH connection.
func (c *Client) Close() error {
	return c.client.Close()
}

// Host returns the target hostname.
func (c *Client) Host() string { return c.host }
