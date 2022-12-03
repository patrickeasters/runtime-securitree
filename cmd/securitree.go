package main

import (
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/cenkalti/backoff/v4"
	"github.com/spf13/viper"

	"github.com/patrickeasters/runtime-securitree/decorator"
	"github.com/patrickeasters/runtime-securitree/queue"
	"github.com/patrickeasters/runtime-securitree/wled"
)

func main() {
	// read config
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/securitree/")
	viper.AddConfigPath("$HOME/.securitree")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("TREE")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	viper.ReadInConfig()

	// setup SQS client
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		panic(err)
	}
	q := queue.NewQueue(sess, viper.GetString("aws.sqsQueueURL"))

	// setup WLED client
	wledClient := wled.NewClient(viper.GetString("wled.host"))

	// setup decorator and initial colors
	dec := decorator.NewDecorator(viper.GetInt("wled.stripLength"), viper.GetInt("wled.brightness"))

	wledClient.SetLEDs(dec.Brightness, dec.LEDState)

	// start listening to queue
	b := backoff.NewExponentialBackOff()
	notify := func(err error, d time.Duration) {
		log.Printf("Encountered error: %s. Retrying in %s", err, d)
	}
	for {
		backoff.RetryNotify(func() error {
			return q.Poll(dec.MessageHandler)
		}, b, notify)
		backoff.RetryNotify(func() error {
			return wledClient.SetLEDs(dec.Brightness, dec.LEDState)
		}, b, notify)
	}

}
