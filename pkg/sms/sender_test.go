package sms_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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

func TestSendMessage(t *testing.T) {
	t.Parallel()

	expectedMessageID := "test"
	for name, tc := range map[string]struct {
		description  string
		mockResponse *sns.PublishOutput
		mockError    error
		assertError  assert.ErrorAssertionFunc
		expectedID   *string
	}{
		"SuccessfullySentMessage": {
			description:  "Successfully sent SMS message.",
			mockResponse: &sns.PublishOutput{MessageId: &expectedMessageID},
			assertError:  assert.NoError,
			expectedID:   &expectedMessageID,
		},
		"FailOnSNSPublish": {
			description: "An error returned when failed to send the message.",
			mockError:   errors.New("broken"),
			assertError: assert.Error,
		},
	} {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Log(tc.description)
			t.Parallel()

			ts := setupTest(t)

			ts.mockedSNS.EXPECT().
				Publish(gomock.Any()).
				Return(tc.mockResponse, tc.mockError).
				Times(1)

			messageID, err := ts.sender.Send("test-msg", "test-receiver")
			tc.assertError(t, err)
			assert.Equal(t, tc.expectedID, messageID)
		})
	}

}

func TestMessageAttributes(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		description      string
		attribute        string
		expectedValue    string
		assertValueIsSet assert.ValueAssertionFunc
		setAttribute     func(s *sms.Sender)
	}{
		"SetSenderID": {
			description:      "Successfully set the sender ID SMS attribute.",
			attribute:        sms.SenderIDSMSAttribute,
			expectedValue:    "test-sender-id",
			assertValueIsSet: assert.NotEmpty,
			setAttribute: func(sender *sms.Sender) {
				sender.WithSenderID("test-sender-id")
			},
		},
		"FailSetSenderIDOnEmptyString": {
			description:      "Fail to set sender ID SMS attribute when an empty string provided.",
			attribute:        sms.SenderIDSMSAttribute,
			assertValueIsSet: assert.Empty,
			setAttribute: func(sender *sms.Sender) {
				sender.WithSenderID("")
			},
		},
		"FailSetSenderIDOnEmptyStringWithWhitespaces": {
			description:      "Fail to set sender ID SMS attribute when an empty string provided.",
			attribute:        sms.SenderIDSMSAttribute,
			assertValueIsSet: assert.Empty,
			setAttribute: func(sender *sms.Sender) {
				sender.WithSenderID("\n\t\t\t\t  \t\n")
			},
		},
		"SetMaxPrice": {
			description:      "Successfully set the max price SMS attribute when over 1 USD cent.",
			attribute:        sms.MaxPriceSMSAttribute,
			expectedValue:    "0.05",
			assertValueIsSet: assert.NotEmpty,
			setAttribute: func(sender *sms.Sender) {
				sender.WithMaxPrice(0.05)
			},
		},
		"FailSetMaxPriceBellowOneCent": {
			description:      "Fail to set max price SMS attribute when the provided value is bellow 1 USD cent.",
			attribute:        sms.MaxPriceSMSAttribute,
			assertValueIsSet: assert.Empty,
			setAttribute: func(sender *sms.Sender) {
				sender.WithMaxPrice(0)
			},
		},
		"SetMessageTypeAsPromotional": {
			description:      "Successfully set the message type SMS attribute as promotional.",
			attribute:        sms.MessageTypeSMSAttribute,
			expectedValue:    sms.Promotional.String(),
			assertValueIsSet: assert.NotEmpty,
			setAttribute: func(sender *sms.Sender) {
				sender.WithMessageType(sms.Promotional)
			},
		},
		"SetMessageTypeAsTransactional": {
			description:      "Successfully set the message type SMS attribute as transactional.",
			attribute:        sms.MessageTypeSMSAttribute,
			expectedValue:    sms.Transactional.String(),
			assertValueIsSet: assert.NotEmpty,
			setAttribute: func(sender *sms.Sender) {
				sender.WithMessageType(sms.Transactional)
			},
		},
		"FailSetMessageType": {
			description:      "When message type isn't valid, message type SMS attribute isn't set.",
			attribute:        sms.MessageTypeSMSAttribute,
			assertValueIsSet: assert.Empty,
			setAttribute: func(sender *sms.Sender) {
				sender.WithMessageType("FooBar")
			},
		},
	} {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Log(tc.description)
			t.Parallel()

			ts := setupTest(t)

			require.Nil(t, ts.sender.MessageAttributes[tc.attribute])

			tc.setAttribute(ts.sender)
			tc.assertValueIsSet(t, ts.sender.MessageAttributes[tc.attribute])

			// proceed with further attribute tests only when the attribute value was successfully set
			if ts.sender.MessageAttributes[tc.attribute] == nil {
				return
			}

			switch *ts.sender.MessageAttributes[tc.attribute].DataType {
			case "String":
				strVal := ts.sender.MessageAttributes[tc.attribute].StringValue
				require.NotNil(t, strVal)
				assert.Equal(t, tc.expectedValue, *strVal)
			}
		})
	}
}

func TestMessageAttributesCreatedFromConfig(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		description      string
		config           *sms.Config
		expectedValues   map[string]string
		assertValueIsSet assert.ValueAssertionFunc
	}{
		"Success": {
			description: "",
			expectedValues: map[string]string{
				sms.SenderIDSMSAttribute:    "test-sender-id",
				sms.MaxPriceSMSAttribute:    "0.07",
				sms.MessageTypeSMSAttribute: sms.Promotional.String(),
			},
			config: &sms.Config{
				SenderID:    strp("test-sender-id"),
				MaxPrice:    floatp(0.07),
				MessageType: strp(sms.Promotional.String()),
			},
			assertValueIsSet: assert.NotEmpty,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Log(tc.description)
			t.Parallel()

			attrs := sms.GetMessageAttributes(tc.config)

			for _, a := range []string{
				sms.SenderIDSMSAttribute,
				sms.MaxPriceSMSAttribute,
				sms.MessageTypeSMSAttribute,
			} {
				attr := attrs[a]
				tc.assertValueIsSet(t, attr)

				if attr == nil || attr.StringValue == nil {
					continue
				}

				assert.Equal(t, tc.expectedValues[a], *attr.StringValue)
			}
		})
	}

}

func strp(s string) *string {
	return &s
}

func floatp(f float32) *float32 {
	return &f
}
