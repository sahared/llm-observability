package kafka

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "time"

    "github.com/IBM/sarama"
)

// Producer handles Kafka event publishing
type Producer struct {
    producer sarama.SyncProducer
    config   *ProducerConfig
}

// ProducerConfig configures the Kafka producer
type ProducerConfig struct {
    Brokers []string
    Topic   string
}

// Event represents a Kafka event
type Event struct {
    Type      string                 `json:"type"`
    Timestamp time.Time              `json:"timestamp"`
    Data      map[string]interface{} `json:"data"`
}

// NewProducer creates a new Kafka producer
func NewProducer(config *ProducerConfig) (*Producer, error) {
    saramaConfig := sarama.NewConfig()
    saramaConfig.Producer.Return.Successes = true
    saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
    saramaConfig.Producer.Retry.Max = 3
    saramaConfig.Producer.Compression = sarama.CompressionSnappy

    producer, err := sarama.NewSyncProducer(config.Brokers, saramaConfig)
    if err != nil {
        return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
    }

    log.Printf("✅ Kafka producer connected to: %v", config.Brokers)

    return &Producer{
        producer: producer,
        config:   config,
    }, nil
}

// PublishEvent publishes an event to Kafka
func (p *Producer) PublishEvent(ctx context.Context, eventType string, data map[string]interface{}) error {
    event := Event{
        Type:      eventType,
        Timestamp: time.Now(),
        Data:      data,
    }

    eventBytes, err := json.Marshal(event)
    if err != nil {
        return fmt.Errorf("failed to marshal event: %w", err)
    }

    message := &sarama.ProducerMessage{
        Topic: p.config.Topic,
        Value: sarama.ByteEncoder(eventBytes),
        Key:   sarama.StringEncoder(eventType),
    }

    partition, offset, err := p.producer.SendMessage(message)
    if err != nil {
        return fmt.Errorf("failed to publish event: %w", err)
    }

    log.Printf("📤 Event published: %s [partition=%d, offset=%d]", eventType, partition, offset)
    return nil
}

// PublishTraceCreated publishes a trace created event
func (p *Producer) PublishTraceCreated(ctx context.Context, traceID, orgID, projectID, model, provider string, tokens int, cost float64) error {
    return p.PublishEvent(ctx, "trace.created", map[string]interface{}{
        "trace_id":        traceID,
        "organization_id": orgID,
        "project_id":      projectID,
        "model":           model,
        "provider":        provider,
        "total_tokens":    tokens,
        "total_cost_usd":  cost,
    })
}

// PublishSpanCreated publishes a span created event
func (p *Producer) PublishSpanCreated(ctx context.Context, spanID, traceID string, durationMs int64, tokens int) error {
    return p.PublishEvent(ctx, "span.created", map[string]interface{}{
        "span_id":      spanID,
        "trace_id":     traceID,
        "duration_ms":  durationMs,
        "total_tokens": tokens,
    })
}

// Close closes the Kafka producer
func (p *Producer) Close() error {
    if p.producer != nil {
        return p.producer.Close()
    }
    return nil
}
