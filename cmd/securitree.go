package main

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/spf13/viper"

	"github.com/patrickeasters/runtime-securitree/decorator"
	"github.com/patrickeasters/runtime-securitree/queue"
	"github.com/patrickeasters/runtime-securitree/wled"
)

func main() {
	// read config
	viper.SetConfigName("config")            // name of config file (without extension)
	viper.AddConfigPath("/etc/securitree/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.securitree") // call multiple times to add many search paths
	viper.AddConfigPath(".")                 // optionally look for config in the working directory
	viper.SetEnvPrefix("TREE")
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
	for {
		q.Poll(dec.MessageHandler)
		wledClient.SetLEDs(dec.Brightness, dec.LEDState)
	}

}
