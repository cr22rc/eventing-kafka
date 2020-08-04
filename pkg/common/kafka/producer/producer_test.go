package producer

import (
	"crypto/tls"
	"github.com/Shopify/sarama"
	"github.com/stretchr/testify/assert"
	"knative.dev/eventing-kafka/pkg/common/kafka/constants"
	"testing"
)

// Test Constants
const (
	ClientId      = "TestClientId"
	KafkaBrokers  = "TestBrokers"
	KafkaUsername = "TestUsername"
	KafkaPassword = "TestPassword"
)

// Test The CreateSyncProducer() Functionality
func TestCreateSyncProducer(t *testing.T) {
	performCreateSyncProducerTest(t, "", "")
}

// Test The CreateSyncProducer() Functionality With SASL PLAIN Auth
func TestCreateSyncProducerWithSasl(t *testing.T) {
	performCreateSyncProducerTest(t, KafkaUsername, KafkaPassword)
}

// Perform A Single Instance Of The CreateProducer() Test With Specified Config
func performCreateSyncProducerTest(t *testing.T, username string, password string) {

	// Create A Mock SyncProducer
	mockSyncProducer := &MockSyncProducer{}

	// Stub The Kafka SyncProducer Creation Wrapper With Test Version Returning Mock SyncProducer
	newSyncProducerWrapperPlaceholder := newSyncProducerWrapper
	newSyncProducerWrapper = func(brokers []string, config *sarama.Config) (sarama.SyncProducer, error) {
		assert.NotNil(t, brokers)
		assert.Len(t, brokers, 1)
		assert.Equal(t, brokers[0], KafkaBrokers)
		verifySaramaConfig(t, config, ClientId, username, password)
		return mockSyncProducer, nil
	}
	defer func() { newSyncProducerWrapper = newSyncProducerWrapperPlaceholder }()

	// Perform The Test
	producer, registry, err := CreateSyncProducer(ClientId, []string{KafkaBrokers}, username, password)

	// Verify The Results
	assert.Nil(t, err)
	assert.NotNil(t, producer)
	assert.Equal(t, mockSyncProducer, producer)
	assert.NotNil(t, registry)
}

// Verify The Sarama Config Is As Expected
func verifySaramaConfig(t *testing.T, config *sarama.Config, clientId string, username string, password string) {
	assert.NotNil(t, config)
	assert.Equal(t, clientId, config.ClientID)
	assert.Equal(t, constants.ConfigKafkaVersion, config.Version)
	assert.Equal(t, constants.ConfigNetKeepAlive, config.Net.KeepAlive)
	assert.Equal(t, constants.ConfigProducerIdempotent, config.Producer.Idempotent)
	assert.Equal(t, constants.ConfigProducerRequiredAcks, config.Producer.RequiredAcks)
	assert.True(t, config.Producer.Return.Successes)

	if len(username) > 0 && len(password) > 0 {
		assert.Equal(t, constants.ConfigNetSaslVersion, config.Net.SASL.Version)
		assert.True(t, config.Net.SASL.Enable)
		assert.Equal(t, sarama.SASLMechanism(sarama.SASLTypePlaintext), config.Net.SASL.Mechanism)
		assert.Equal(t, username, config.Net.SASL.User)
		assert.Equal(t, password, config.Net.SASL.Password)
		assert.True(t, config.Net.TLS.Enable)
		assert.NotNil(t, config.Net.TLS.Config)
		assert.False(t, config.Net.TLS.Config.InsecureSkipVerify)
		assert.Equal(t, tls.NoClientCert, config.Net.TLS.Config.ClientAuth)
	} else {
		assert.NotNil(t, config.Net)
		assert.False(t, config.Net.TLS.Enable)
		assert.Nil(t, config.Net.TLS.Config)
		assert.Equal(t, sarama.SASLHandshakeV0, config.Net.SASL.Version)
		assert.NotNil(t, config.Net.SASL)
		assert.False(t, config.Net.SASL.Enable)
		assert.Equal(t, sarama.SASLMechanism(""), config.Net.SASL.Mechanism)
		assert.Equal(t, "", config.Net.SASL.User)
		assert.Equal(t, "", config.Net.SASL.Password)
	}
	assert.Equal(t, constants.ConfigMetadataRefreshFrequency, config.Metadata.RefreshFrequency)
}

//
// Mock Sarama SyncProducer Implementation
//

var _ sarama.SyncProducer = &MockSyncProducer{}

type MockSyncProducer struct {
}

func (p *MockSyncProducer) SendMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	return 0, 0, nil
}

func (p *MockSyncProducer) SendMessages(msgs []*sarama.ProducerMessage) error {
	return nil
}

func (p *MockSyncProducer) Close() error {
	return nil
}