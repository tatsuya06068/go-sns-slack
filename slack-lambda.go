package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/slack-go/slack"
	"os"
)

// handleRequest はSNSからのメッセージを処理してSlackに通知する関数です。
func handleRequest(ctx context.Context, snsEvent events.SNSEvent) error {
	// Slack API トークンを環境変数から取得
	slackToken := os.Getenv("SLACK_BOT_TOKEN")
	if slackToken == "" {
		return fmt.Errorf("Slack API token is not set")
	}

	// Slack API クライアントを作成
	client := slack.New(slackToken)

	// Slack チャンネルID（例: #general）
	channelID := os.Getenv("SLACK_CHANNEL_ID")
	if channelID == "" {
		return fmt.Errorf("Slack channel ID is not set")
	}

	// SNSのメッセージをSlackに投稿
	for _, record := range snsEvent.Records {
		snsRecord := record.SNS
		message := snsRecord.Message

		// Slackにメッセージを送信
		_, _, err := client.PostMessage(
			channelID,
			slack.MsgOptionText(fmt.Sprintf("New SNS message: %s", message), false),
		)
		if err != nil {
			log.Printf("Failed to send message to Slack: %v", err)
			return err
		}
	}

	return nil
}

func main() {
	// Lambdaのエントリーポイント
	lambda.Start(handleRequest)
}
