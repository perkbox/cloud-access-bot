package main

import (
	"context"
	"os"

	"github.com/perkbox/cloud-access-bot/internal/settings"

	"github.com/perkbox/cloud-access-bot/internal/messenger"

	"github.com/perkbox/cloud-access-bot/controllers"
	"github.com/slack-go/slack/socketmode"

	"github.com/perkbox/cloud-access-bot/drivers"

	"github.com/perkbox/cloud-access-bot/internal/identitydata"

	"github.com/perkbox/cloud-access-bot/internal/policy"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/joho/godotenv"
	"github.com/perkbox/cloud-access-bot/internal"
	"github.com/perkbox/cloud-access-bot/internal/awsproviderv2"
	"github.com/perkbox/cloud-access-bot/internal/repository"
	"github.com/sirupsen/logrus"
)

func main() {
	//Load config from .env when running, works both locally and in production settings
	_ = godotenv.Load()
	//
	if os.Getenv("BOT_CONFIG_S3_BUCKET") == "" || os.Getenv("BOT_CONFIG_S3_KEY") == "" {
		logrus.WithField("err", "Ensure both BOT_CONFIG_S3_BUCKET & BOT_CONFIG_S3_KEY are set.").Errorf("Missing Env Vars")
		os.Exit(1)
	}

	cfg, _ := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("eu-west-1"),
	)

	settings, err := settings.NewS3Config(cfg, os.Getenv("BOT_CONFIG_S3_KEY"), os.Getenv("BOT_CONFIG_S3_BUCKET"))
	if err != nil {
		logrus.WithField("err", err.Error()).Errorf("Unable get Config")
		os.Exit(1)
	}

	client, err := drivers.ConnectToSlackViaSocketmode()
	if err != nil {
		logrus.WithField("err", err.Error()).Errorf("Unable to connect to slack")
		os.Exit(1)
	}

	service := internal.NewService(
		awsproviderv2.NewAwsResourceFinder(cfg, settings),
		repository.NewDynamoDBRRepo(cfg, settings.GetDynamodbTable()),
		policy.NewPolicyManager(cfg, settings, nil, nil),
		identitydata.NewIamDefinitions(),
		messenger.NewMessenger(client.GetApiClient()),
	)

	socketmodeHandler := socketmode.NewsSocketmodeHandler(client)

	controllers.NewSlashCommandController(settings, service, socketmodeHandler)

	socketmodeHandler.RunEventLoop()
}
