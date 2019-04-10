package sms_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/martinsirbe/go-sms/pkg/sms"
	"github.com/martinsirbe/go-sms/pkg/sms/mocks"
)

type testSuite struct {
	mockController *gomock.Controller
	mockedSNS      *mocks.MockSNSAPI
	sender         *sms.Sender
}

func setupTest(t *testing.T) *testSuite {
	mockController := gomock.NewController(t)
	mockedSNS := mocks.NewMockSNSAPI(mockController)

	return &testSuite{
		mockController: mockController,
		mockedSNS:      mockedSNS,
		sender:         sms.NewWithSNS(mockedSNS),
	}
}

func TestSuccessfullyPublishedSMS(t *testing.T) {
	ts := setupTest(t)

	expectedMessageID := "test"
	expectedResponse := sns.PublishOutput{
		MessageId: &expectedMessageID,
	}
	ts.mockedSNS.EXPECT().Publish(gomock.Any()).Return(&expectedResponse, nil).Times(1)

	actualMessageID, err := ts.sender.Send("test-msg", "test-receiver")
	assert.Nil(t, err)
	assert.Equal(t, expectedMessageID, *actualMessageID)
}

func TestOnFailedSMSPublishReturnError(t *testing.T) {
	ts := setupTest(t)

	ts.mockedSNS.EXPECT().Publish(gomock.Any()).Return(nil, errors.New("bad")).Times(1)

	_, err := ts.sender.Send("test-msg", "test-receiver")
	assert.NotNil(t, err)

	assert.Equal(t, "failed to publish a text message to test-receiver: bad", err.Error())
}

func TestSuccessfullySetSenderID(t *testing.T) {
	ts := setupTest(t)

	assert.Nil(t, ts.sender.MessageAttributes[sms.SenderIDSMSAttribute])

	expectedSenderID := "test"
	ts.sender.WithSenderID(expectedSenderID)

	assert.NotNil(t, ts.sender.MessageAttributes[sms.SenderIDSMSAttribute].StringValue)
	assert.Equal(t, expectedSenderID, *ts.sender.MessageAttributes[sms.SenderIDSMSAttribute].StringValue)
}

func TestSuccessfullySetMaxPrice(t *testing.T) {
	ts := setupTest(t)

	assert.Nil(t, ts.sender.MessageAttributes[sms.MaxPriceSMSAttribute])

	var expectedMaxPrice float32 = 0.05
	ts.sender.WithMaxPrice(expectedMaxPrice)

	assert.NotNil(t, ts.sender.MessageAttributes[sms.MaxPriceSMSAttribute].StringValue)
	assert.Equal(t, fmt.Sprintf("%.2f", expectedMaxPrice),
		*ts.sender.MessageAttributes[sms.MaxPriceSMSAttribute].StringValue)
}

func TestMaxPriceNotSetIfGivenValueIsBellowOneCent(t *testing.T) {
	ts := setupTest(t)

	assert.Nil(t, ts.sender.MessageAttributes[sms.MaxPriceSMSAttribute])

	expectedMaxPrice := float32(0.00)
	ts.sender.WithMaxPrice(expectedMaxPrice)

	assert.Nil(t, ts.sender.MessageAttributes[sms.MaxPriceSMSAttribute])
}

func TestSenderIDNotSetIfGivenValueIsEmptyString(t *testing.T) {
	ts := setupTest(t)

	assert.Nil(t, ts.sender.MessageAttributes[sms.SenderIDSMSAttribute])

	ts.sender.WithSenderID("")

	assert.Nil(t, ts.sender.MessageAttributes[sms.SenderIDSMSAttribute])
}

func TestSuccessfullySetMessageTypeAsPromotional(t *testing.T) {
	ts := setupTest(t)

	assert.Nil(t, ts.sender.MessageAttributes[sms.MessageTypeSMSAttribute])

	expectedMessageType := sms.Promotional
	ts.sender.WithMessageType(expectedMessageType)

	assert.NotNil(t, ts.sender.MessageAttributes[sms.MessageTypeSMSAttribute].StringValue)
	assert.Equal(t, string(expectedMessageType), *ts.sender.MessageAttributes[sms.MessageTypeSMSAttribute].StringValue)
}

func TestSuccessfullySetMessageTypeAsTransactional(t *testing.T) {
	ts := setupTest(t)

	assert.Nil(t, ts.sender.MessageAttributes[sms.MessageTypeSMSAttribute])

	expectedMessageType := sms.Transactional
	ts.sender.WithMessageType(expectedMessageType)

	assert.NotNil(t, ts.sender.MessageAttributes[sms.MessageTypeSMSAttribute].StringValue)
	assert.Equal(t, string(expectedMessageType), *ts.sender.MessageAttributes[sms.MessageTypeSMSAttribute].StringValue)
}

func TestSuccessfullyCreatedMessageAttributesFromConfig(t *testing.T) {
	expectedSenderID := "test-sender-id"
	expectedMaxPrice := float32(0.07)
	expectedMessageType := string(sms.Promotional)

	attrs := sms.GetMessageAttributes(&sms.Config{
		SenderID:    &expectedSenderID,
		MaxPrice:    &expectedMaxPrice,
		MessageType: &expectedMessageType,
	})

	assert.Equal(t, expectedSenderID, *attrs[sms.SenderIDSMSAttribute].StringValue)
	assert.Equal(t, fmt.Sprintf("%.2f", expectedMaxPrice), *attrs[sms.MaxPriceSMSAttribute].StringValue)
	assert.Equal(t, expectedMessageType, *attrs[sms.MessageTypeSMSAttribute].StringValue)
}

func TestConfigValuesAreOptionalWhenCreatingMessageAttributesFromConfig(t *testing.T) {
	attrs := sms.GetMessageAttributes(&sms.Config{})

	assert.Nil(t, attrs[sms.SenderIDSMSAttribute])
	assert.Nil(t, attrs[sms.MaxPriceSMSAttribute])
	assert.Nil(t, attrs[sms.MessageTypeSMSAttribute])
}
