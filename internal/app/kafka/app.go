package kafka

import (
	"context"
	"fmt"
	"log/slog"
	"love-signal-users/internal/config"
	"love-signal-users/internal/infrastructure/kafka/handlers"
	"love-signal-users/pkg/kafka"
	"love-signal-users/pkg/logger/sl"
	"strings"
)

// App is a kafka queue application.
type App struct {
	log                      *slog.Logger
	consumers                []*kafka.Consumer
	userHasRegisteredHandler UserHasRegisteredHandler
}

// New returns new instance of kafka queue application.
func New(
	log *slog.Logger,
	cfg config.KafkaConfig,
	userHasRegisteredHandler UserHasRegisteredHandler,
) *App {
	address := strings.Split(cfg.Address, ",")
	registerNewUserConsumer := kafka.NewConsumerGroup(
		address,
		cfg.RegisterNewUserTopic.GroupID,
		cfg.RegisterNewUserTopic.Topic,
	)

	return &App{
		log:                      log,
		consumers:                []*kafka.Consumer{registerNewUserConsumer},
		userHasRegisteredHandler: userHasRegisteredHandler,
	}
}

// Start - starts the kafka queue application.
func (a *App) Start(ctx context.Context) {
	const op = "kafkaapp.Start"

	log := a.log.With(
		slog.String("op", op),
	)
	log.Info(fmt.Sprintf("running kafka %d consumers", len(a.consumers)))

	for _, c := range a.consumers {
		go c.Consume(ctx)
	}

	a.handleConsumerReceivedMessages(ctx)
	a.handleConsumerErrors()
}

// Stop - stops the kafka queue application.
func (a *App) Stop() {
	const op = "kafkaapp.Stop"

	log := a.log.With(
		slog.String("op", op),
	)
	log.Info("stopping kafka consumers")

	for _, consumer := range a.consumers {
		if err := consumer.Close(); err != nil {
			log.Error("error closing kafka consumer", sl.Err(err))
		}
	}

}

func (a *App) handleConsumerReceivedMessages(ctx context.Context) {
	const op = "kafkaapp.handleConsumerReceivedMessages"

	log := a.log.With(slog.String("op", op))
	
	for _, consumer := range a.consumers {
		go func() {
			for {
				select {
				case msg := <-consumer.Output():
					log.Info("received message from kafka", slog.String("topic", msg.Topic))

					switch msg.Topic {
					case handlers.UserHasRegisteredTopic:
						go func() {
							if err := a.userHasRegisteredHandler.Execute(ctx, msg.Data); err != nil {
								log.Error("error handling new registered user", sl.Err(err))
							}
						}()
					default:
						log.Warn("handler implementation for topic does not exist", slog.String("topic", msg.Topic))
					}
				default:
				}
			}
		}()
	}
}

func (a *App) handleConsumerErrors() {
	const op = "kafkaapp.handleConsumerErrors"

	log := a.log.With(slog.String("op", op))

	for _, c := range a.consumers {
		go func() {
			for {
				select {
				case err := <-c.Notify():
					if err != nil {
						log.Error("error reading message from kafka", sl.Err(err))
					}
				default:
				}
			}
		}()
	}
}
