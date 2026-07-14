package redis

import (
	"context"
	"fmt"
	"legacy-messenger-control-plane/configs"
	"legacy-messenger-control-plane/internal/adapters/ssh"
	"legacy-messenger-control-plane/internal/domain"
	"net"
	"strconv"

	goredis "github.com/redis/go-redis/v9"
)

type redisClient struct {
	client *goredis.Client
}

func NewRedisClient(ctx context.Context, redisCfg *configs.RedisConfig, sshClient *ssh.SSHClient) (*redisClient, error) {

	address := net.JoinHostPort(
		redisCfg.Host,
		strconv.Itoa(redisCfg.Port),
	)

	options := &goredis.Options{
		Addr:     address,
		Password: redisCfg.Password,
		DB:       redisCfg.DB,
	}

	// SSH Client가 존재하면 SSH 서버를 경유해 Redis에 연결한다.
	if sshClient != nil {
		fmt.Println("sshClient used")
		options.Dialer = sshClient.DialContext
		// golang.org/x/crypto/ssh의 tcpChan은
		// SetReadDeadline / SetWriteDeadline을 지원하지 않는다.
		options.ReadTimeout = -2
		options.WriteTimeout = -2
	}

	client := goredis.NewClient(options)

	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()

		connectionType := "direct"
		if sshClient != nil {
			connectionType = "ssh tunnel"
		}

		return nil, fmt.Errorf(
			"failed to connect to Redis via %s: %w",
			connectionType,
			err,
		)
	}

	return &redisClient{
		client: client,
	}, nil
}

func (c *redisClient) SaveTaskSessionReport(ctx context.Context, report domain.TaskSessionReport) error {

	fmt.Printf("[SaveTaskSessionReport] serviceName : %s, taskID : %s, sessionCount : %d \n", report.ServiceName, report.TaskID, report.SessionCount)

	return nil
}

func (c *redisClient) Close() error {
	if c == nil || c.client == nil {
		return nil
	}

	return c.client.Close()
}
