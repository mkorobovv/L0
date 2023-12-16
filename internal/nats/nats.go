package nats

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/mkorobovv/L0/config"
	"github.com/mkorobovv/L0/internal/models"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

type Nats struct {
	cfg *config.Nats
	nc  *nats.Conn
	sc  stan.Conn
}

func NewNats(cfg *config.Nats) *Nats {

	natsUrl := fmt.Sprintf("nats://%s:%s", cfg.Host, cfg.Port)
	natsConn, err := nats.Connect(natsUrl)
	if err != nil {
		fmt.Printf("Connection lost: %v\n", err)
		return nil
	}
	stanConn, err := stan.Connect(cfg.Cluster, cfg.Client,
		stan.Pings(10, 3),
		stan.SetConnectionLostHandler(func(_ stan.Conn, err error) {
			fmt.Printf("Connection lost: %v\n", err)
		}))

	if err != nil {
		fmt.Printf("Connection lost: %v\n", err)
		return nil
	}

	return &Nats{cfg, natsConn, stanConn}
}

func (nats *Nats) Publish(msg models.OrderJSON) error {

	order, err := json.MarshalIndent(msg, "", "\t")
	if err != nil {
		return fmt.Errorf("error marshalling order %v", err)
	}

	return nats.sc.Publish(nats.cfg.Topic, order)
}

func (nats *Nats) Subscribe() (*models.OrderJSON, error) {
	var order models.OrderJSON
	ch := make(chan *models.OrderJSON)

	_, err := nats.sc.Subscribe(nats.cfg.Topic, func(msg *stan.Msg) {
		if err := json.Unmarshal(msg.Data, &order); err != nil {
			fmt.Printf("Error unmarshalling: %v", err)
		}
		ch <- &order
	})

	if err != nil {
		fmt.Printf("Error on subscribe: %v", err)
	}

	select {
	case order := <-ch:
		return order, nil
	case <-time.After(60 * time.Second):
		return nil, fmt.Errorf("timeout error")
	}
}
