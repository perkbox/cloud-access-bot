# Cloud Access Bot

[![codecov](https://codecov.io/gh/perkbox/cloud-access-bot/branch/main/graph/badge.svg?token=US5TXW4M67)](https://codecov.io/gh/perkbox/cloud-access-bot)
[![Go Report Card](https://goreportcard.com/badge/github.com/perkbox/cloud-access-bot)](https://goreportcard.com/report/github.com/perkbox/cloud-access-bot)

The cloud access bot is a Slack bot developed within Perkbox which enables developers to request elevated permissions to AWS accounts. The entire process is contained within Slack and is mostly automated.

The Slack application itself uses Slack's Socket Mode to ensure maximum security by not requiring to expose any public endpoints to the internet.

From a high level the workflow of the bot is:
1. Users runs the slack command /request (Customizable command name)
2. Modal loads prompting the user for:
 - Request reason
 - Expiry time of permissions
 - AWS SAML login role
 - AWS account
3. On selection of AWS account the modal will update and request:
 - AWS Services 
4. On selection of the AWS service the modal will update and request:
  - AWS Resources
    This will either be via a multi-select list or a free text field where AWS ARN's can be entered.


## :magic_wand:	Features

- Automated temporary permissions management
- Time limited elevated permissions
- Simple approval workflow
- Audit logging of all requests


## Video of Bot in Use

https://user-images.githubusercontent.com/26804184/149273215-3af8f8b5-f421-45c1-a9f7-20f9c458432d.mp4


# Table of Contents
- [Installation](#nut_and_bolt-installation)
  - [Slack App Setup](#slack-app-setup)
  - [Terraform Module](#terraform-module)
- [Usage](#usage)
  - [Running Bot Locally](#computer-running-bot-locally)
    - [Pre Requisites](#pre-requisites)
    - [1. Building Compiling Code](#1-building-compiling-code)
    - [2. Setting Up Environment](#2-setting-up-environment)
    - [3. Running](#3-running)
- [Configuration](#page_facing_up-configuration)
  - [Environment Variables](#environment-variables)
  - [Configuration File](#configuration-file)
- [Required AWS Resources](#required-aws-resources)
  - [S3 Bucket](#s3-bucket)
  - [DynamoDB](#dynamodb)
  - [IAM Policy](#iam-policy)
- [Limitations](#warning-limitations)
- [Contributing](#contributing)
- [Credits](#credits)


## DockerHub

New docker images are built on every release of the code. The DockerHub repo can be found [here] (https://hub.docker.com/r/perkbox/cloud-access-bot) or if using an image just by: 

```
perkbox/cloud-access-bot
```


## :nut_and_bolt: Installation

The preferred installation of the `Cloud Access Bot` is to run the container within a container orchestration platform, within Perkbox we are using ECS Fargate to host the `Cloud Access Bot`.

You will also need to run through the app setup within slack as well.


### Slack App Setup

The setup of the slack application provides you with the:
- `SLACK_APP_TOKEN`
- `SLACK_BOT_TOKEN`

Which allows your Slack workspace and the bot to communicate with each other.

The Cloud Access Bot uses Slack's [Socket mode]("https://api.slack.com/apis/connections/socket) to communicate with Slack, which means no public endpoint (Request URL) is needed, making it more secure.

For the setup you can follow the video link bellow, which shows how to setup and create the Slack application. 

To start you will need to be an admin in your workspace to create applications.
Also navigate to https://api.slack.com/apps


https://user-images.githubusercontent.com/26804184/149046380-9e915a86-178b-4997-806b-83e647711d85.mp4

**Useful Video Timestamps**
 - Slack app Token Creation  (00:09 -> 00:28)
 - Slack bot Token Creation  (00:50 -> 00:56)


For the list of OAuth Scopes a written list of them are:
- channels:history
- channels:read
- chat:write
- commands
- im:history
- im:write
- usergroups:read
- users:read
- users:read:email


### Terraform Module

Please see the Terraform module within the repo as a quick start to demo the Slack bot in your own environment.
All required documentation can be found in the module.

[Terraform](../terraform/)

The individual required resources can be found in the section [Required AWS Resources](#required-aws-resources).


The `Cloud Access Bot` can also be ran locally on your machine quite easily.

Please follow the usage guide [Running Bot Locally](#computer-running-bot-locally) which explains everything needed.


## Usage

### :computer: Running Bot Locally

#### Pre Requisites
- Go 1.17+
- AWS access
- S3 bucket
- DynamoDB table
- Slack credentials (app and bot tokens)


#### 1. Building Compiling code

**Requires Go to be installed**

To build the bot all you need to run is the command 

```shell
go build . 
```
From the working DIR where the main.go is located.
In later versions the binary will be released as part of the package.


#### 2. Setting up environment

**Requires S3 bucket and DynamoDB table**

Even running locally you still require the S3 bucket and DynamoDB table for the bot to function.
You local environment will also have access to your AWS account with the required permissions to edit roles and read/write from S3 and DynamoDB.

The [Environment Variables](#environment-variables) can be set in your local environment using 

**eg.**
```shell
$ export SLACK_BOT_TOKEN="xoxb-25-TokenHere"
```

Or they can be place into a file called `.env` next to the binary. The binary will read these in when it first starts without the need to specify where the file is.

**eg.**
```
SLACK_BOT_TOKEN=xoxb-25-TokenHere
```

You will also need to have the [configuration file](#configuration-file) in S3 with the correct parameters for the bot to work. 


#### 3. Running 

Once the above steps are completed you can run the bot locally.

```
./cloud-access-bot 
```


## :page_facing_up: Configuration

There are 2 main configuration types used by the bot to get it up and running. 

- Environment variables
- Configuration file (S3)

Both are required for the bots proper operation. 


### Environment Variables

| ENV Variable         | Description                                                   | Example                     |
|----------------------|---------------------------------------------------------------|-----------------------------|
| BOT_CONFIG_S3_BUCKET | Name of the S3 bucket containing the configuration            | `request-access-bot-config` |
| BOT_CONFIG_S3_KEY    | The path or file name of the configuration file in the bucket | `config.yml`                |
| SLACK_APP_TOKEN      | Slack app token, required for authentication                  | `xapp-1-A02K...`            |
| SLACK_BOT_TOKEN      | Slack bot token, required for authentication                  | `xoxb-2557...`              |


### Configuration file 

A configuration file is used to store the primary configuration of the bot, which does effect some of the selectable options within the modal prompts.


| Var Name        | Type         | Description                                                                       |
|-----------------|--------------|-----------------------------------------------------------------------------------|
| loginRoles      | `list`       | The Roles that users can use to login                                             |
| approvalGroups  | `list`       | The names of the Slack user groups which approval messages will be sent to        |
| regions         | `list`       | Regions to fetch resources from when the auto discovery is being used             |
| accounts        | `dictionary` | Dictionary of accounts which the Slack bot can search and create policies within  |
| dynamoDbTable   | `string`     | Name of the DynamoDB table to use                                                 |
| request_command | `string`     | Slack command name, defaults to `request`. Do not add the slash to the variable, this is done in code|



**Configuration File Example**:

```yaml
---
loginRoles:
  - SSO-Devops

approvalGroups:
  - devops

regions:
  - eu-west-1
  - eu-west-2

accounts:
  perkbox-development:
    account_number: 1234567
    iam_role: "arn:aws:iam::1234567:role/request-access-bot"
```


## Required AWS Resources 

## S3 Bucket

An S3 bucket is required to store the configuration for the Slack bot, no secrets or sensitive data are stored within the bucket.

As a recommendation, enabling versioning and server side encryption by default, along with disabling any form of public access or acl to ensure anything stored within the bucket is kept safe.


### DynamoDB 

The `Cloud Access Bot` uses a DynamoDB table to store all audit data and some metadata from messages which allows it to update there content throughout the request process.

As the table isn't used heavily in most cases a provisioned capacity table can be used, though using an on-demand table may save on costs if the bot is heavily used.


| Table fields    | Name | Type   |
|-----------------|------|--------|
| Partition Key   | PK   | string |
| Sort(Range) Key | SK   | string |

*All other table fields are added as needed by the bot*


### IAM Policy

The `Cloud Access Bot` requires an IAM role to access the required cloud resources, which it requires to use. This also includes cross-account access create and edit IAM roles for users in other accounts, thoug cross account access isn't required if its only being used within a single account.

Please note this is a working sample, permissions can be tuned and tweaked to your preferences. 

```json
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "AssumeAccessRoles",
            "Effect": "Allow",
            "Action": "sts:AssumeRole",
            "Resource": [
                "arn:aws:iam::account-num-1:role/cloud-access-bot",
                "arn:aws:iam::account-num-2:role/cloud-access-bot"
            ]
        },
        {
            "Sid": "AssumeAccessRoles",
            "Effect": "Allow",
            "Action": [
                "dynamodb:UpdateItem",
                "dynamodb:Scan",
                "dynamodb:Query",
                "dynamodb:PutItem",
                "dynamodb:ListTagsOfResource",
                "dynamodb:ListTables",
                "dynamodb:GetRecords",
                "dynamodb:GetItem",
                "dynamodb:DescribeTimeToLive",
                "dynamodb:DescribeTable",
                "dynamodb:DescribeReservedCapacityOfferings",
                "dynamodb:DescribeReservedCapacity",
                "dynamodb:DescribeLimits",
                "dynamodb:DeleteItem",
                "dynamodb:BatchWriteItem",
                "dynamodb:BatchGetItem"
            ],
            "Resource": "arn:aws:dynamodb:eu-west-1:account-num:table/cloud-access-bot"
        },
        {
            "Sid": "S3BucketAccess",
            "Effect": "Allow",
            "Action": "s3:*",
            "Resource": [
                "arn:aws:s3:::request-access-bot-config/*",
                "arn:aws:s3:::request-access-bot-config"
            ]
        }
    ]
}
```


## :warning: Limitations

- Only supports SAML based AWS logins
- Limited auto-loading AWS services (services where ARN's are generated)
  - S3
  - DynamoDB
- Single service policies (At the moment)


# Contributing

All contributions to this project are welcome, please refer to the contributing file [CONTRIBUTING](./CONTRIBUTING.md)


# Credits

There are some projects which we used generated files from or heavily used as part of the project and wanted to give mentions to them. 

- [policy_sentry](https://github.com/salesforce/policy_sentry)
This is used for getting some templating information about ARNS and AWS services available, the file can be found at `internal/identitydata/assets/iam-definition.json`

- [slack-go](https://github.com/slack-go/slack)
The slack-go work is heavily used for interacting with Slack's-APIs and orchestrating a handler for Socket Mode connections.

A massive helper was this medium post [Implement Slack Slash Command in Golang using Socket Mode](https://levelup.gitconnected.com/implement-slack-slash-command-in-golang-using-socket-mode-ac693e38148c)
