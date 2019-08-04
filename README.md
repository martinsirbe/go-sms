# Go SMS
A simple app for sending Short Message Service (SMS) text messages using AWS Simple Notification Service.  

![json_ast_badge](https://img.shields.io/badge/SNS-green.svg?logo=amazon-aws&style=flat) 
[![Go Report Card](https://goreportcard.com/badge/github.com/martinsirbe/go-sms)](https://goreportcard.com/report/github.com/martinsirbe/go-sms) 
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fmartinsirbe%2Fgo-sms.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fmartinsirbe%2Fgo-sms?ref=badge_shield) 
[![codecov](https://codecov.io/gh/martinsirbe/go-sms/branch/master/graph/badge.svg)](https://codecov.io/gh/martinsirbe/go-sms) 
[![CircleCI](https://circleci.com/gh/martinsirbe/go-sms/tree/master.svg?style=svg)](https://circleci.com/gh/martinsirbe/go-sms/tree/master) 


## Build
Run `make build`, and you should see `go-sms` binary in `bin` directory. Alternatively you can build a docker image 
by running `make docker-build` and run it by `make run`. Note that `config.yaml` should be located in the root of 
this project which should be based on `config_sample.yaml`.  

You also can use go-sms pre-built docker image by `docker pull martinsirbe/go-sms`.  

## Configuration
### Mandatory
* `aws_access_key` - AWS account access key id string.  
* `aws_secret_access_key` - AWS account secret access key string.  
* `aws_region` - AWS account secret access key string. (Only certain AWS regions are 
supported, check [AWS documentation][1].)  

### Optional
* `sender_id` - Sender ID which will be visible on the receiver's device. Can be up to 11 alphanumeric characters which 
must contain at least one latter. When not set will default to `NOTICE`. This configuration value will be overridden by 
the CLI `sender-id` argument. (Note that only certain countries support sender ID, check [AWS documentation][1].)  
* `max_price` - The maximum price in USD that you are willing to pay to send the message. Note that 
your message won't be sent if the cost to send the message exceeds the set maximum price. This attribute will have 
no effect if the limit set for the `MonthlySpendLimit` attribute is exceeded. Check [AWS documentation][2] for SMS prices, 
based on this you can determine the possible `max_price`.  
* `sms_type` - Signifies SMS type which is being sent. It can be either `Promotional` (default) or 
`Transactional`.  
  * `Promotional` - Noncritical messages with optimised message delivery to incur the lowest cost, e.g. marketing messages.  
  * `Transactional` - Critical messages with optimised message delivery to achieve the highest reliability, e.g. authentication messages.  

See `config_sample.yaml` for an example configuration file.  

### CLI Options
You can provide options as environment variables, or pass options as CLI arguments.  
```bash
Usage: go-sms [OPTIONS]

Short Message Service (SMS) text message sender using AWS Simple Notification Service.
                      
Options:              
  --sender-id     The sender ID which will appear on the receiver's device. (Optional, if provided will override sender ID provided via configuration file.) (env $SENDER_ID)
  --receiver      The receiver mobile phone number. (Mandatory) (env $RECEIVER)
  --message       The text message you wish to send. (Mandatory) (env $MESSAGE)
  --config-path   The path to the configurations file. (Mandatory) (env $GO_SMS_CONFIG_PATH)
```

* `message` - Can be 160 GSM, 140 ASCII or 70 UCS-2 characters long with a total size limit of 1600 bytes per SMS publish action.  
* `config-path` - Should point to the configurations file. You can use `config_sample.yaml` as a reference.  

## Examples
### CLI
```bash
go-sms --sender-id=<sender_id> --receiver=<mobile_phone_number> --message=<your_message> --config-path=<path_to_config_file>
```
  
### Go
```golang
package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"

	log "github.com/sirupsen/logrus"

	"github.com/martinsirbe/go-sms/pkg/sms"
)

func main() {
	configFile, err := ioutil.ReadFile("path/to/config.yaml")
	if err != nil {
		log.WithError(err).Fatal("failed to load go-sms config.yaml")
	}

	var config sms.Config
	if err = yaml.Unmarshal(configFile, &config); err != nil {
		log.WithError(err).Fatal("failed to unmarshal go-sms config.yaml")
	}

	sender := sms.New(config)
	if _, err := sender.Send("Hello world!", "+44xxx"); err != nil {
		log.WithError(err).Fatal("failed to send the text message")
	}
}
```

## Tests
To run tests, just run `make test`.
```bash
?       github.com/martinsirbe/go-sms/cmd/go-sms        [no test files]
=== RUN   TestSuccessfullyPublishedSMS
--- PASS: TestSuccessfullyPublishedSMS (0.00s)
=== RUN   TestOnFailedSMSPublishReturnError
--- PASS: TestOnFailedSMSPublishReturnError (0.00s)
=== RUN   TestSuccessfullySetSenderID
--- PASS: TestSuccessfullySetSenderID (0.00s)
=== RUN   TestSuccessfullySetMaxPrice
--- PASS: TestSuccessfullySetMaxPrice (0.00s)
=== RUN   TestMaxPriceNotSetIfGivenValueIsBellowOneCent
--- PASS: TestMaxPriceNotSetIfGivenValueIsBellowOneCent (0.00s)
=== RUN   TestSenderIDNotSetIfGivenValueIsEmptyString
--- PASS: TestSenderIDNotSetIfGivenValueIsEmptyString (0.00s)
=== RUN   TestSuccessfullySetMessageTypeAsPromotional
--- PASS: TestSuccessfullySetMessageTypeAsPromotional (0.00s)
=== RUN   TestSuccessfullySetMessageTypeAsTransactional
--- PASS: TestSuccessfullySetMessageTypeAsTransactional (0.00s)
=== RUN   TestSuccessfullyCreatedMessageAttributesFromConfig
--- PASS: TestSuccessfullyCreatedMessageAttributesFromConfig (0.00s)
=== RUN   TestConfigValuesAreOptionalWhenCreatingMessageAttributesFromConfig
--- PASS: TestConfigValuesAreOptionalWhenCreatingMessageAttributesFromConfig (0.00s)
PASS
coverage: 84.8% of statements
ok      github.com/martinsirbe/go-sms/pkg/sms   0.020s  coverage: 84.8% of statements
?       github.com/martinsirbe/go-sms/pkg/sms/mocks     [no test files]
```

## License
This project is licensed under the MIT License - see the [LICENSE.md](LICENCE.md) file for details.  

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fmartinsirbe%2Fgo-sms.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fmartinsirbe%2Fgo-sms?ref=badge_large)

## Contributing
1. Go get it! `go get github.com/martinsirbe/go-sms`  
2. Create your feature branch. (`git checkout -b my-feature-branch`)  
3. Commit your changes. (`git commit -m 'Add ...'`)  
4. Push to the branch. (`git push origin my-feature-branch`)  
5. Create a new pull request.  

[1]: https://docs.aws.amazon.com/sns/latest/dg/sms_supported-countries.html
[2]: https://aws.amazon.com/sns/sms-pricing/
[3]: https://github.com/golangci/golangci-lint
