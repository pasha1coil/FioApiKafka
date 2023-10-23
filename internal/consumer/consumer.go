package consumer

import (
	"FioapiKafka/internal/models"
	"context"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
	"os"
)

func InitReader(broker []string, topic string) *kafka.Reader {
	log.Infof("Init reader... Topic:%s", topic)
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: broker,
		Topic:   topic,
	})
	return reader
}

func InitConsumerFIO(broker []string, topic string, aproovedFIO chan *models.PersonOut) context.CancelFunc {
	reader := InitReader(broker, topic)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		for {
			msg, err := reader.ReadMessage(ctx)
			if err != nil {
				if err == context.Canceled {
					log.Infoln("Consumer stopped")
					reader.Close()
					return
				}
				log.Fatalf("failed to read message: %s", err.Error())
			}
			log.Infof("Received message from %s: %v", msg.Topic, string(msg.Value))
			data, err := GetData(msg.Value)
			if err != nil {
				err := WriterToField(msg.Value)
				if err != nil {
					log.Errorf("error writing a message to KAFKA FIO_FAILED_TOPIC: %s", err.Error())
				}
			} else {
				aproovedFIO <- data
			}
		}
	}()

	return cancel
}

func InitConsumerFailed(broker []string, topic string) context.CancelFunc {
	reader := InitReader(broker, topic)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		for {
			msg, err := reader.ReadMessage(ctx)
			if err != nil {
				if err == context.Canceled {
					log.Infoln("Consumer stopped")
					reader.Close()
					return
				}
				log.Fatalf("failed to read message: %s", err.Error())
			}
			log.Infof("new field message from %s: %v", msg.Topic, string(msg.Value))
		}
	}()

	return cancel
}

func WriterToField(data []byte) error {
	log.Infoln("Init KAFKA FIO_FAILED_TOPIC")
	writer := &kafka.Writer{
		Addr:     kafka.TCP(os.Getenv("KAFKA_BROKER")),
		Topic:    os.Getenv("KAFKA_FIO_FAILED_TOPIC"),
		Balancer: &kafka.LeastBytes{},
	}
	log.Infoln("Recording message to KAFKA FIO_FAILED_TOPIC")
	err := writer.WriteMessages(context.Background(), kafka.Message{Value: data})
	if err != nil {
		log.Errorf("error writing a message to KAFKA FIO_FAILED_TOPIC: %s", err.Error())
		return err
	}
	writer.Close()
	log.Infoln("Writing the message to KAFKA FIO_FAILED_TOPIC was successful")
	return nil
}
