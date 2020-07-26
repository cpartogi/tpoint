package consumer

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/Shopify/sarama"
	"github.com/cpartogi/tpoint/models"
	"github.com/sirupsen/logrus"
)

// KafkaConsumer hold sarama consumer
type KafkaConsumer struct {
	Consumer sarama.Consumer
}

var err error
var pesan string

// Consume function to consume message from apache kafka
func (c *KafkaConsumer) Consume(topics []string, signals chan os.Signal) {
	chanMessage := make(chan *sarama.ConsumerMessage, 10000)
	for _, topic := range topics {
		partitionList, err := c.Consumer.Partitions(topic)
		if err != nil {
			logrus.Errorf("Unable to get partition got error %v", err)
			continue
		}
		for _, partition := range partitionList {
			go consumeMessage(c.Consumer, topic, partition, chanMessage)
		}
	}
	logrus.Infof("Kafka is consuming....")

ConsumerLoop:
	for {
		select {
		case msg := <-chanMessage:
			logrus.Infof("New Message from kafka, message: %v", string(msg.Value))

			//proccess message
			pesan = string(msg.Value)
			processMessage(pesan)
		case sig := <-signals:
			if sig == os.Interrupt {
				break ConsumerLoop
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func consumeMessage(consumer sarama.Consumer, topic string, partition int32, c chan *sarama.ConsumerMessage) {
	msg, err := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
	if err != nil {
		logrus.Errorf("Unable to consume partition %v got error %v", partition, err)
		return
	}

	defer func() {
		if err := msg.Close(); err != nil {
			logrus.Errorf("Unable to close partition %v: %v", partition, err)
		}
	}()

	for {
		msg := <-msg.Messages()
		c <- msg
	}

}

func processMessage(c string) error {
	p := new(models.Pointlog)
	json.Unmarshal([]byte(c), &p)

	// send to another api
	url := "http://127.0.0.1:1234/point"

	var jsonStr = []byte(c)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	return nil
}
