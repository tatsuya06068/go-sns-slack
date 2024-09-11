package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/slack-go/slack"
)

// SNSMessage はSNSメッセージの構造体です。
type SNSMessage struct {
	ChannelID string   `json:"channel_id"`
	Message   string   `json:"message"`
	Mentions  []string `json:"mentions"` // メンションするユーザーやグループのIDリスト
}

// handleRequest はSNSからのメッセージを処理してSlackに通知する関数です。
func handleRequest(ctx context.Context, snsEvent events.SNSEvent) error {
	// Slack API トークンを環境変数から取得
	slackToken := os.Getenv("SLACK_BOT_TOKEN")
	if slackToken == "" {
		return fmt.Errorf("Slack API token is not set")
	}

	// Slack API クライアントを作成
	client := slack.New(slackToken)

	for _, record := range snsEvent.Records {
		snsRecord := record.SNS
		var snsMessage SNSMessage

		// SNSメッセージをパース
		if err := json.Unmarshal([]byte(snsRecord.Message), &snsMessage); err != nil {
			log.Printf("Failed to unmarshal SNS message: %v", err)
			return err
		}

		// チャンネルIDとメッセージが正しく取得できたか確認
		if snsMessage.ChannelID == "" || snsMessage.Message == "" {
			return fmt.Errorf("channel_id or message is missing in SNS message")
		}

		// メッセージにメンションを追加
		mentions := ""
		for _, mention := range snsMessage.Mentions {
			mentions += fmt.Sprintf("<@%s> ", mention)
		}

		// メッセージ本文を作成
		message := fmt.Sprintf("%s%s", mentions, snsMessage.Message)

		// Slackにメッセージを送信
		_, _, err := client.PostMessage(
			snsMessage.ChannelID,
			slack.MsgOptionText(message, false),
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



// {
// 	"Records": [
// 	  {
// 		"EventSource": "aws:sns",
// 		"EventVersion": "1.0",
// 		"EventSubscriptionArn": "arn:aws:sns:region:account-id:example-topic",
// 		"Sns": {
// 		  "Type": "Notification",
// 		  "MessageId": "11111111-2222-3333-4444-555555555555",
// 		  "TopicArn": "arn:aws:sns:region:account-id:example-topic",
// 		  "Subject": "Test SNS Subject",
// 		  "Message": "{\"channel_id\": \"C1234567890\", \"message\": \"This is a test message\", \"mentions\": [\"U1234567890\", \"U0987654321\"]}",
// 		  "Timestamp": "2024-09-11T12:00:00.000Z",
// 		  "SignatureVersion": "1",
// 		  "Signature": "EXAMPLE_SIGNATURE",
// 		  "SigningCertUrl": "EXAMPLE_CERT_URL",
// 		  "UnsubscribeUrl": "EXAMPLE_UNSUBSCRIBE_URL",
// 		  "MessageAttributes": {}
// 		}
// 	  }
// 	]
//   }
