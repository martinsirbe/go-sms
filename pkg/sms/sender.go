package sms

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

//go:generate mockgen -package=mocks -destination=mocks/sns.go github.com/aws/aws-sdk-go/service/sns/snsiface SNSAPI

// Sender used to send text messages using AWS SNS
type Sender struct {
	awsSNS            snsiface.SNSAPI
	MessageAttributes map[string]*sns.MessageAttributeValue
}

// New initialise sender with the given configuration
func New(config Config) *Sender {
	awsConfig := aws.NewConfig().
		WithCredentials(credentials.NewStaticCredentials(
			config.AWSAccessKey,
			config.AWSSecretAccessKey, "")).
		WithRegion(config.AWSRegion)

	awsSession, err := session.NewSession(awsConfig)
	if err != nil {
		log.WithError(err).Fatal("failed to initialise a new AWS session")
	}

	return &Sender{
		awsSNS:            sns.New(awsSession, awsConfig),
		MessageAttributes: GetMessageAttributes(&config),
	}
}

// NewWithSNS initialise sender with custom SNS
func NewWithSNS(awsSNS snsiface.SNSAPI) *Sender {
	return &Sender{
		awsSNS:            awsSNS,
		MessageAttributes: make(map[string]*sns.MessageAttributeValue),
	}
}

// Send message to the given message to the receiver mobile phone number
func (s *Sender) Send(msg, receiver string) (*string, error) {
	params := &sns.PublishInput{
		Message:           aws.String(msg),
		PhoneNumber:       aws.String(receiver),
		MessageAttributes: s.MessageAttributes,
	}

	resp, err := s.awsSNS.Publish(params)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to publish a text message to %s", receiver)
	}

	return resp.MessageId, nil
}

// WithSenderID will include a sender ID in the sent text message
func (s *Sender) WithSenderID(id string) *Sender {
	if id == "" {
		return s
	}

	s.MessageAttributes[SenderIDSMSAttribute] = &sns.MessageAttributeValue{
		DataType:    aws.String("String"),
		StringValue: &id,
	}
	return s
}

// WithMaxPrice sets the max price in USD that you are willing to pay to send a single text message
func (s *Sender) WithMaxPrice(p float32) *Sender {
	if p < 0.01 {
		return s
	}

	mp := fmt.Sprintf("%.2f", p)
	s.MessageAttributes[MaxPriceSMSAttribute] = &sns.MessageAttributeValue{
		DataType:    aws.String("String"),
		StringValue: &mp,
	}
	return s
}

// WithMessageType sets the message type
func (s *Sender) WithMessageType(t MessageType) *Sender {
	mt := string(t)
	s.MessageAttributes[MessageTypeSMSAttribute] = &sns.MessageAttributeValue{
		DataType:    aws.String("String"),
		StringValue: &mt,
	}
	return s
}

// GetMessageAttributes returns message attributes map based on the given config
func GetMessageAttributes(c *Config) map[string]*sns.MessageAttributeValue {
	attrs := make(map[string]*sns.MessageAttributeValue)

	if c.SenderID != nil {
		attrs[SenderIDSMSAttribute] = &sns.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: c.SenderID,
		}
	}

	if c.MaxPrice != nil {
		mp := fmt.Sprintf("%.2f", *c.MaxPrice)
		attrs[MaxPriceSMSAttribute] = &sns.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: &mp,
		}
	}

	mt := c.MessageType
	if mt != nil && (*mt == string(Promotional) || *mt == string(Transactional)) {
		attrs[MessageTypeSMSAttribute] = &sns.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: mt,
		}
	}

	return attrs
}
