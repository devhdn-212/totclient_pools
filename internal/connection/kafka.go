package connection

import (
	"fmt"
	"strings"
	"time"

	"github.com/devhdn-212/totclient_pools/internal/config"
	"github.com/segmentio/kafka-go"
)

// kafkaLogf routes kafka-go's own internal diagnostics (coordinator
// discovery, group join/sync, rebalances, per-broker connection errors,
// ...) to stdout — plain fmt, not connection.Log, so this never triggers
// the Telegram alert hook (see checkoutconsumer.go's Run) for what's often
// just noisy protocol-level chatter. Wired in mainly so that a generic
// "i/o timeout" surfaced by Run() has an accompanying trail just above it
// in the container logs showing WHICH internal step (metadata? coordinator?
// join group? fetch?) actually stalled.
func kafkaLogf(prefix string) func(string, ...interface{}) {
	return func(msg string, args ...interface{}) {
		fmt.Printf("[kafka %s] "+msg+"\n", append([]interface{}{prefix}, args...)...)
	}
}

// NewKafkaReader builds the checkout event bus consumer for this worker.
// GroupID makes this a real consumer group: the broker tracks committed
// offsets per group, so if totclient_pools is ever scaled to more than one
// instance, each message still only goes to one of them.
//
// Dialer.Timeout is widened from kafka-go's default (10s) — the broker
// (34.128.79.54:7070) sits across the public internet rather than a
// low-latency VPC link, so coordinator/group round trips that legitimately
// take a while easily trip a too-tight connection deadline, surfacing as a
// spurious "read tcp ...: i/o timeout" even though the broker is perfectly
// reachable. MaxWait is left at kafka-go's own default (10s, i.e. how long
// the broker holds a Fetch open waiting for new data before replying
// empty) — a shorter MaxWait was tried here at one point to "leave more
// headroom", but all it did was make the reader poll (and log an entirely
// normal "no messages received within the allocated time") 3x more often
// while idle, for no actual benefit once Dialer.Timeout had already been
// widened.
func NewKafkaReader(conf config.Kafka) *kafka.Reader {
	brokers := strings.Split(conf.Brokers, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   conf.Topic,
		GroupID: conf.GroupID,
		Dialer: &kafka.Dialer{
			Timeout:   30 * time.Second,
			DualStack: true,
		},
		ReadBackoffMin: 1 * time.Second,
		ReadBackoffMax: 10 * time.Second,
		// Without this, a consumer that joins the group before the topic
		// exists (or before it gets more partitions) just sits on an empty
		// assignment forever — kafka-go won't notice the topic/partition
		// count changed until something forces a fresh rebalance. Watching
		// for changes means a topic created after this process started
		// still gets picked up without a restart.
		WatchPartitionChanges:  true,
		PartitionWatchInterval: 30 * time.Second,
		Logger:                 kafka.LoggerFunc(kafkaLogf("info")),
		ErrorLogger:            kafka.LoggerFunc(kafkaLogf("error")),
	})
}
