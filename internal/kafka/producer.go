package kafka

import (
	"encoding/json"
	"github.com/Alex1472/ozon-film-service/internal/model"
	"github.com/Shopify/sarama"
)

type Producer struct {
	p     sarama.SyncProducer
	topic string
}

const (
	topicName = "film_events"
)

func NewProducer() (*Producer, error) {
	// TODO add to config
	// address outside docker
	brokers := []string{"127.0.0.1:9094"}
	p, err := NewSyncProducer(brokers)
	if err != nil {
		return nil, err
	}

	return &Producer{
		p:     p,
		topic: topicName,
	}, nil
}

func (p *Producer) SendCreated(film *model.Film) error {
	event := NewKafkaEvent(film.ID, Created, film)
	return p.send(event)
}

func (p *Producer) SendUpdated(film *model.Film) error {
	event := NewKafkaEvent(film.ID, Updated, film)
	return p.send(event)
}

func (p *Producer) SendRemoved(filmId uint64) error {
	event := NewKafkaEvent(filmId, Removed, nil)
	return p.send(event)
}

func (p *Producer) send(event *FilmEvent) error {
	bytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = SendMessage(p.p, p.topic, bytes)
	return err
}
