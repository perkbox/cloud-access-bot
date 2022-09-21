package main

import (
	"context"
	"errors"
	"log"
	"os"
	"strings"

	"github.com/perkbox/cloud-access-bot/commands"
	"github.com/perkbox/cloud-access-bot/internal"
	"github.com/perkbox/cloud-access-bot/internal/awsproviderv2"
	"github.com/perkbox/cloud-access-bot/internal/identitydata"
	"github.com/perkbox/cloud-access-bot/internal/messenger"
	"github.com/perkbox/cloud-access-bot/internal/policy"
	"github.com/perkbox/cloud-access-bot/internal/repository"
	"github.com/perkbox/cloud-access-bot/internal/settings"

	"github.com/aws/aws-sdk-go-v2/config"
	runtime "github.com/banzaicloud/logrus-runtime-formatter"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

func init() {
	formatter := runtime.Formatter{ChildFormatter: &logrus.JSONFormatter{}}
	formatter.Line = true
	logrus.SetFormatter(&formatter)
}

func main() {
	//Load config from .env when running, works both locally and in production settings
	_ = godotenv.Load()
	//
	if os.Getenv("BOT_CONFIG_S3_BUCKET") == "" || os.Getenv("BOT_CONFIG_S3_KEY") == "" {
		logrus.Errorf("Missing Env Vars Err: Ensure both BOT_CONFIG_S3_BUCKET & BOT_CONFIG_S3_KEY are set. ")
		os.Exit(1)
	}

	cfg, _ := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("eu-west-1"),
	)

	settings, err := settings.NewS3Config(cfg, os.Getenv("BOT_CONFIG_S3_KEY"), os.Getenv("BOT_CONFIG_S3_BUCKET"))
	if err != nil {
		logrus.Errorf("Unable get Config. Err %s", err)
		os.Exit(1)
	}

	client, err := connectToSlackViaSocketmode()
	if err != nil {
		logrus.Errorf("Unable to connect to slack. Err: %s", err)
		os.Exit(1)
	}

	service := internal.NewService(
		awsproviderv2.NewAwsResourceFinder(cfg, settings),
		repository.NewDynamoDBRRepo(cfg, settings.GetDynamodbTable()),
		policy.NewPolicyManager(cfg, settings, nil, nil),
		identitydata.NewIamDefinitions(),
		messenger.NewMessenger(&client.Client),
	)

	socketmodeHandler := socketmode.NewSocketmodeHandler(client)

	commands.NewRequestCommandHandler(settings, service, socketmodeHandler)

	socketmodeHandler.RunEventLoop()
}

func connectToSlackViaSocketmode() (*socketmode.Client, error) {
	appToken := os.Getenv("SLACK_APP_TOKEN")
	if appToken == "" {
		return nil, errors.New("SLACK_APP_TOKEN must be set")
	}

	if !strings.HasPrefix(appToken, "xapp-") {
		return nil, errors.New("SLACK_APP_TOKEN must have the prefix \"xapp-\".")
	}

	botToken := os.Getenv("SLACK_BOT_TOKEN")
	if botToken == "" {
		return nil, errors.New("SLACK_BOT_TOKEN must be set.")
	}

	if !strings.HasPrefix(botToken, "xoxb-") {
		return nil, errors.New("SLACK_BOT_TOKEN must have the prefix \"xoxb-\".")
	}

	api := slack.New(
		botToken,
		//slack.OptionDebug(true),
		slack.OptionAppLevelToken(appToken),
		slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
	)

	client := socketmode.New(
		api,
		//socketmode.OptionDebug(true),
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)

	return client, nil
}
