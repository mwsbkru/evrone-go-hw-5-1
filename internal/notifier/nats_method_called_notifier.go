package notifier

import (
	"encoding/json"
	"evrone_go_hw_5_1/config"
	"fmt"
	"github.com/nats-io/nats.go"
	"log/slog"
)

// NatsMethodCalledNotifier provides functionality for notification about events via Nats
type NatsMethodCalledNotifier struct {
	conn *nats.Conn
	cfg  *config.Config
}

type notificationType struct {
	Method string            `json:"method"`
	Params map[string]string `json:"params"`
}

// NewNatsMethodCalledNotifier returns new NatsMethodCalledNotifier
func NewNatsMethodCalledNotifier(conn *nats.Conn, cfg *config.Config) *NatsMethodCalledNotifier {
	return &NatsMethodCalledNotifier{conn: conn, cfg: cfg}
}

// NotifyMethodCalled sends notification about method was called to Nats
func (n *NatsMethodCalledNotifier) NotifyMethodCalled(method string, params map[string]string) error {
	notification := notificationType{
		Method: method,
		Params: params,
	}

	msg, err := json.Marshal(&notification)
	if err != nil {
		return fmt.Errorf("не удалось сериализовать сообщение для отправки в Nats: %w", err)
	}

	if err := n.conn.Publish(n.cfg.NatsMethodCalledSubject, msg); err != nil {
		return fmt.Errorf("не удалось отправить сообщение в Nats: %w", err)
	}

	slog.Info("Сообщение отправлено в Nats", slog.String("message", string(msg)))
	return nil
}
