package testing

import (
	"context"
	"testing"

	"github.com/go-gotop/gotop/broker/kafka"
	"github.com/go-kratos/kratos/v2/log"
)

func Test_Kafka_Publish(t *testing.T) {
	publisher, err := kafka.NewPublisher(
		log.NewHelper(log.DefaultLogger),
		kafka.WithAddrs("10.0.0.141:9092"),
	)
	if err != nil {
		t.Fatal(err)
	}

	err = publisher.Publish(context.Background(), "munal_test", []byte("test"), []byte("test"), nil)
	if err != nil {
		t.Fatal(err)
	}
}
