package databus_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/Shopify/sarama"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/zenoss/zenkit/test"

	"testing"
)

var (
	harness      test.Harness
	testProducer sarama.SyncProducer
	testConsumer sarama.Consumer
	logger       = test.TestLogger()
)

func TestDatabus(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Databus Integration Suite")
}

func ZooKeeperHealthCheck(zkaddr string) func() error {
	return func() error {
		conn, err := net.DialTimeout("tcp", zkaddr, time.Second)
		if err != nil {
			return err
		}
		defer conn.Close()
		conn.Write([]byte("ruok\n"))
		conn.SetReadDeadline(time.Now().Add(time.Second))
		data, err := ioutil.ReadAll(conn)
		if err != nil {
			return err
		}
		if bytes.Compare(data, []byte("imok")) != 0 {
			return errors.New("zk not ok")
		}
		return nil
	}
}

func KafkaHealthCheck(kafkaaddr string) func() error {
	return func() error {
		broker := sarama.NewBroker(kafkaaddr)
		if err := broker.Open(sarama.NewConfig()); err != nil {
			return err
		}
		if connected, _ := broker.Connected(); !connected {
			return errors.New("not connected")
		}
		return nil
	}
}

func SchemaRegistryHealthCheck(schemaaddr string) func() error {
	return func() error {
		resp, err := http.Get(fmt.Sprintf("http://%s/config", schemaaddr))
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		if bytes.Compare(data, []byte(`{"compatibilityLevel":"BACKWARD"}`)) != 0 {
			return errors.New("schema reg not ok")
		}
		return nil
	}
}

var _ = BeforeSuite(func() {
	var err error
	harness, err = test.NewDockerComposeHarness(fmt.Sprintf("test-%s", test.RandString(4)), "docker-compose.yml")
	Ω(err).ShouldNot(HaveOccurred())

	Ω(harness.Start()).ShouldNot(HaveOccurred())

	// ZooKeeper health check
	zkaddr, err := harness.Resolve("zk", 2181)
	Ω(err).ShouldNot(HaveOccurred())
	err = harness.Wait(ZooKeeperHealthCheck(zkaddr), 30*time.Second)
	Ω(err).ShouldNot(HaveOccurred())
	logger.WithField("address", zkaddr).Infof("ZooKeeper is ready")

	// Kafka health check
	kafka, err := harness.Resolve("kafka", 9092)
	Ω(err).ShouldNot(HaveOccurred())
	err = harness.Wait(KafkaHealthCheck(kafka), 30*time.Second)
	Ω(err).ShouldNot(HaveOccurred())
	logger.WithField("address", kafka).Infof("Kafka is ready")

	oldresolver := net.DefaultResolver
	net.DefaultResolver = &net.Resolver{
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			fmt.Println("Resolving", address)
			if address == "kafka" {
				addr, _ := harness.Resolve("kafka", 123)
				address = addr[:len(addr)-4]
			}
			return oldresolver.Dial(ctx, network, address)
		},
	}

	testProducer, err = sarama.NewSyncProducer([]string{kafka}, nil)
	Ω(err).ShouldNot(HaveOccurred())
	testConsumer, err = sarama.NewConsumer([]string{kafka}, nil)
	Ω(err).ShouldNot(HaveOccurred())

	// Schema health check
	schemareg, err := harness.Resolve("kafka-schema-registry", 8081)
	Ω(err).ShouldNot(HaveOccurred())
	err = harness.Wait(SchemaRegistryHealthCheck(schemareg), 30*time.Second)
	Ω(err).ShouldNot(HaveOccurred())
	logger.WithField("address", schemareg).Infof("Schema registry is ready")
})

var _ = AfterSuite(func() {
	//harness.Stop()
})
