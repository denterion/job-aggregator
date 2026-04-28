package transport

import (
	"context"
	"encoding/json"
	"job-aggregator/internal/models"

	"github.com/segmentio/kafka-go"
)

func SendVacancy(v models.Vacancy) error {

	writer := &kafka.Writer{
		Addr:                   kafka.TCP("127.0.0.1:9092"),
		Topic:                  "jobs_raw",
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}

	payload, err := json.Marshal(v)
	if err != nil {
		return err
	}

	err = writer.WriteMessages(context.Background(),
		kafka.Message{
			Value: payload,
		},
	)

	writer.Close()
	return err
}

func GetReader() *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"127.0.0.1:9092"},
		Topic:    "jobs_raw",
		GroupID:  "processor-group",
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
}
