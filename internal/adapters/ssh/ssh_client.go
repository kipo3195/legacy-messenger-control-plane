package ssh

import (
	"context"
	"fmt"
	"legacy-messenger-control-plane/configs"
	"net"
	"strconv"
	"time"

	gossh "golang.org/x/crypto/ssh"
)

type SSHClient struct {
	client *gossh.Client
}

func NewSSHClient(cfg *configs.SSHConfig) (*SSHClient, error) {
	if cfg.Host == "" {
		return nil, fmt.Errorf("SSH host is required")
	}

	if cfg.Port <= 0 {
		cfg.Port = 22
	}

	if cfg.User == "" {
		return nil, fmt.Errorf("SSH user is required")
	}

	if cfg.Password == "" {
		return nil, fmt.Errorf("SSH password is required")
	}

	if cfg.Timeout <= 0 {
		cfg.Timeout = 5 * time.Second
	}

	address := net.JoinHostPort(
		cfg.Host,
		strconv.Itoa(cfg.Port),
	)

	clientConfig := &gossh.ClientConfig{
		User: cfg.User,
		Auth: []gossh.AuthMethod{
			gossh.Password(cfg.Password),
		},
		HostKeyCallback: gossh.InsecureIgnoreHostKey(),
		Timeout:         cfg.Timeout,
	}

	client, err := gossh.Dial("tcp", address, clientConfig)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to connect SSH server %s: %w",
			address,
			err,
		)
	}

	return &SSHClient{
		client: client,
	}, nil
}

func (c *SSHClient) DialContext(
	ctx context.Context,
	network string,
	address string,
) (net.Conn, error) {
	if c == nil || c.client == nil {
		return nil, fmt.Errorf("SSH client is not initialized")
	}

	conn, err := c.client.DialContext(ctx, network, address)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to connect through SSH to %s: %w",
			address,
			err,
		)
	}

	return conn, nil
}

func (c *SSHClient) Close() error {
	if c == nil || c.client == nil {
		return nil
	}

	return c.client.Close()
}
