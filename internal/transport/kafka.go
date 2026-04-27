package transport

import (
	"context"
	"encoding/json"
	"job-aggregator/internal/models"
	"github.com/segmentio/kafka-go"
)

func SendVacancy(v models.Vacancy) error{

	writer := &kafka.Writer{
		Addr: kafka.TCP("localhost:9092"),
		Topic: "jobs_raw",
		Balancer: &kafka.LeastBytes{},
	}

	payload, err := json.Marshal(v)
	if err != nil{
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