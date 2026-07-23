package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/RedditUclaista/community-service/internal/usecases"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type UserRegisterPayload struct {
	UserID    string `json:"user_id"`
	Name      string `json:"name"`
	Lastname  string `json:"lastname"`
	Role      string `json:"role"`
	Timestamp int64  `json:"timestamp"`
	Email     string `json:"email"`
}

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	useCase *usecases.UserUseCase
	wg      sync.WaitGroup
}

func NewConsumer(mqURL string, vhost string, uc *usecases.UserUseCase) (*Consumer, error) {
	config := amqp.Config{
		Vhost: vhost,
	}

	conn, err := amqp.DialConfig(mqURL, config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to LavinMQ cluster: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open execution channel: %w", err)
	}

	return &Consumer{
		conn:    conn,
		channel: ch,
		useCase: uc,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) error {
	exchangeName := "topic"
	queueName := "community_service.core.user.registered"
	routingKeys := []string{
		"core.user.registered.#",
		"core.user.registered",
	}

	err := c.channel.ExchangeDeclare(
		exchangeName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare target exchange: %w", err)
	}

	err = c.channel.Qos(
		100,
		0,
		false,
	)
	if err != nil {
		return fmt.Errorf("failed to enforce qos prefetch constraints: %w", err)
	}

	q, err := c.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare processing queue: %w", err)
	}

	for _, key := range routingKeys {
		err = c.channel.QueueBind(
			q.Name,
			key,
			exchangeName,
			false,
			nil,
		)
		if err != nil {
			return fmt.Errorf("failed to bind queue topology to exchange for key %s: %w", key, err)
		}
	}

	msgs, err := c.channel.Consume(
		q.Name,
		"community_service_consumer",
		false, // manual ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to initialize consumer stream: %w", err)
	}

	go func() {
		log.Println("======> COMMUNITY CONSUMER GOROUTINE STARTED! WAITING FOR MESSAGES <======")
		for {
			select {
			case d, ok := <-msgs:
				if !ok {
					log.Println("======> LavinMQ channel stream closed abruptly <======")
					return
				}

				c.wg.Add(1)
				go func(delivery amqp.Delivery) {
					defer c.wg.Done()
					workerCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
					defer cancel()

					var payload UserRegisterPayload
					if err := json.Unmarshal(delivery.Body, &payload); err != nil {
						log.Printf("!!! Malformed json payload detected: %s. Dropping message.\n", err)
						_ = delivery.Nack(false, false)
						return
					}

					log.Printf(">>> Community processing registration payload for UserID: %s\n", payload.UserID)

					id, err := uuid.Parse(payload.UserID)
					if err != nil {
						log.Printf("!!! Invalid UUID %s: %s\n", payload.UserID, err)
						_ = delivery.Nack(false, false)
						return
					}

					err = c.useCase.CreateUser(workerCtx, id)
					if err != nil {
						log.Printf("!!! UseCase execution failed for user %s: %s\n", payload.UserID, err)
						delivery.Nack(false, true) // Requeue
					} else {
						delivery.Ack(false)
						log.Printf(">>> Successfully saved user %s to community db\n", payload.UserID)
					}
				}(d)

			case <-ctx.Done():
				log.Println("======> Lifecycle context termination signaled <======")
				return
			}
		}
	}()

	return nil
}

func (c *Consumer) Close() {
	c.wg.Wait()
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
	log.Println("Community LavinMQ client connection released successfully")
}
