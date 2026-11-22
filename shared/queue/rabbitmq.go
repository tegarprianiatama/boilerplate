package queue

import (
	"boilerplate/shared/config"
	"boilerplate/shared/logger"
	"fmt"
	"net"
	"net/url"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

const rabbitMQDialTimeout = 5 * time.Second

// OpenRabbitMQ opens a connection to RabbitMQ using the configured credentials and vhost.
func OpenRabbitMQ() (*amqp.Connection, error) {
	cfg := config.GetConfig()
	log := logger.Log()

	vhost := cfg.RabbitMQVHost
	if vhost == "" {
		vhost = "/"
	}

	uri := url.URL{
		Scheme: "amqp",
		User:   url.UserPassword(cfg.RabbitMQUser, cfg.RabbitMQPassword),
		Host:   net.JoinHostPort(cfg.RabbitMQHost, cfg.RabbitMQPort),
		Path:   "/",
	}

	conn, err := amqp.DialConfig(uri.String(), amqp.Config{
		Vhost:     vhost,
		Heartbeat: 10 * time.Second,
		Locale:    "en_US",
		Dial:      amqp.DefaultDial(rabbitMQDialTimeout),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to rabbitmq at %s: %w", uri.Host, err)
	}

	log.Info("Connected to RabbitMQ",
		zap.String("host", uri.Host),
		zap.String("vhost", vhost),
	)
	return conn, nil
}
