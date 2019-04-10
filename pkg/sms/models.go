package sms

const (
	// SenderIDSMSAttribute a custom name that's displayed as the message sender on the receiving device.
	SenderIDSMSAttribute = "AWS.SNS.SMS.SenderID"
	// MaxPriceSMSAttribute a maximum price in USD that you are willing to pay to send the message.
	MaxPriceSMSAttribute = "AWS.SNS.SMS.MaxPrice"
	// MessageTypeSMSAttribute a SMS type, can be either Promotional or Transactional.
	MessageTypeSMSAttribute = "AWS.SNS.SMS.MessageType"

	// Promotional message type used for promotional purposes that are noncritical, won't be delivered
	// to DND (Do Not Disturb) numbers.
	Promotional MessageType = "Promotional"
	// Transactional message type used for transactional purposes which includes critical messages as multi-factor
	// authentication. This message type might be more expensive than Promotional message type. Will be delivered to
	// to DND numbers.
	Transactional MessageType = "Transactional"
)

// MessageType SNS SMS type
type MessageType string

// Config sender configuration
type Config struct {
	AWSAccessKey       string   `yaml:"aws_access_key"`
	AWSSecretAccessKey string   `yaml:"aws_secret_access_key"`
	AWSRegion          string   `yaml:"aws_region"`
	SenderID           *string  `yaml:"sender_id,omitempty"`
	MaxPrice           *float32 `yaml:"max_price,omitempty"`
	MessageType        *string  `yaml:"message_type,omitempty"`
}
