# Cloud Access Bot

The cloud access bot is a  Slack bot developed within Perkbox which enables developers to request elevated permission's to AWS accounts. The entire process is contained within slack and is mostly automated..



## :magic_wand:	Features

- Automated Temporary Permissions management 
- Time limited elevated permissions
- Simple approval workflow
- Audit logging of all requests



## Slack App Setup

The setup of the slack application provides you with the 
- SLACK_APP_TOKEN
- SLACK_BOT_TOKEN

Which allows your Slack workspace and the bot to communicate with each other..

The Cloud Access Bot uses slacks [Socket mode]("https://api.slack.com/apis/connections/socket) to communicate with slack which means no public endpoint(Request URL)
is needed, making it more secure..


For the setup you can follow the video link bellow which shows how to setup and create the Slack application 

To start you will need to be an admin in your workspace to create applications.
Also navigate to https://api.slack.com/apps


**Useful Video Timestamps**
 - Slack app Token Creation  (00:09 -> 00:28)
 - Slack bot Token Creation  (00:50 -> 00:56)


For the list of OAuth Scopes a written list of them are 
- channels:history
- channels:read
- chat:write
- commands
- im:history
- im:write
- usergroups:read
- users:read
- users:read:email



## :nut_and_bolt: Installation

The preferred installation of the Cloud Access Bot is to run the container within an container orchestrion platform, within Perkbox are using ECS Fargate to host the cloud access bot.

Please see the terraform module within the repo as a quick start to demo the Slack bot in your own environment.

[Terraform](./terrraform/)


The cloud access bot can also be ran locally on your machine quite easily.


Please follow the Running Bot Locally Guide which explains everything needed.


## :computer: Running Bot Locally

### Pre Requisites
- Go 1.17+
- AWS Access
- S3 Bucket
- Dynamo Table
- Slack Credentials (App and Bot Token)


### 1. Building Compiling code

**Requires Go to be installed**

To build the bot all you need to run is the command 

```shell
go build . 
```
From the working DIR where the main.go is located..
In Later versions the binary will be released as part of the package.


### 2. Setting up environment

**Requires S3 Bucket and Dynamo Table**

Even running locally you still require the S3 Bucket and Dynamo Table for the bot to function.
You local environment will also have access to your AWS account with the required permissions to edit roles and read/write from s3 and dynamoDb.


The [Environment Variables](#environment-variables) can be set in your local environment using 

**eg.**
```shell
EXPORT SLACK_BOT_TOKEN="xoxb-25-TokenHere"
```

Or they can be place into a file called `.env` next to the binary. The binary will read these in when it first starts without the need to specif where the file is.

**eg.**
```
SLACK_BOT_TOKEN=xoxb-25-TokenHere
```


You will also need to have the [Configuration file](#configuration-file) in s3 with the correct parameters for the bot to work.. 



### 3. Running 

Once the above steps are completed you can run the bot locally

```
./cloud-access-bot 
```


## :page_facing_up: Configuration

There are 2 main configuration types used by the bot to get it up and running. 

- Environment Variables
- Configuration file(S3)

Both are required for the bots proper operation. 



### Environment Variables


| ENV Variable         | Description                                                   | Example                     |
|----------------------|---------------------------------------------------------------|-----------------------------|
| BOT_CONFIG_S3_BUCKET | Name of the S3 Bucket Containing the Configuration            | `request-access-bot-config` |
| BOT_CONFIG_S3_KEY    | The path or file name of the configuration file in the bucket | `config.yml`                |
| SLACK_APP_TOKEN      | Slack App Token, Required for authentication                  | `xapp-1-A02K...`            |
| SLACK_BOT_TOKEN      | Slack BOT Token, Required for authentication                  | `xoxb-2557...`              |



### Configuration file 

A configuration file is used to store the primary configuration of the bot which does effect some of the selectable options within the modal prompts.


| Var Name       | Type         | Description                                                                       |
|----------------|--------------|-----------------------------------------------------------------------------------|
| loginRoles     | `list`       | The Roles that users can use to login                                             |
| approvalGroups | `list`       | The names of the slack user groups which approval messages will be sent to        |
| regions        | `list`       | Regions to fetch resources from when the auto discovery is being used             |
| accounts       | `dictionary` | Dictionary of accounts which the slack bot can search and create policies within  |



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


### DynamoDB 

The Cloud Access Bot uses a dynamodb table to store all audit data and some metadata from messages which allows it to update there content thru-out the request process.

As the table isn't used heavily in most cases a provisioned capacity table can be used, thou using an on-demand table may save on costs if the bot is heavily used.


| Table fields    | Name | Type   |
|-----------------|------|--------|
| Partition Key   | PK   | string |
| Sort(Range) Key | SK   | string |

*All other table felids are added as needed by the bot*



### IAM Policy

The Cloud Access Bot Requires an IAM role to access the required cloud resources which it requires to use.

This also includes cross-account access create and edit IAM roles for users in other accounts, tho cross account access isn't required if its only being used within a single account.

Please note this is a working sample permissions can be tuned and tweaked to your preferences. 

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

- Only Supports SAML based AWS logins 
- Limited auto-loading aws services (Services where ARN's are generated)
- Single Service policy's ( At the moment )



# Contributing

All contributions to this project are welcome please refer to the Contributing file
[CONTRIBUTING](./CONTRIBUTING.md)
