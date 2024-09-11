package infrastructure

import (
    "context"
    "encoding/json"
    "fmt"
    "myapp/domain"
    "myapp/repository"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/sns"
)

// snsRepositoryはSNSRepositoryインターフェースを実装した構造体です。
type snsRepository struct {
    client *sns.SNS
    topic  string
}

// NewSNSRepositoryは新しいSNSリポジトリを作成します。
func NewSNSRepository(topic string) repository.SNSRepository {
    sess := session.Must(session.NewSession())
    return &snsRepository{
        client: sns.New(sess),
        topic:  topic,
    }
}

// PublishはSNSにメッセージを送信します。
func (s *snsRepository) Publish(ctx context.Context, message domain.SNSMessage) error {
    msgJSON, err := json.Marshal(message)
    if err != nil {
        return fmt.Errorf("failed to marshal message: %v", err)
    }

    input := &sns.PublishInput{
        Message:  aws.String(string(msgJSON)),
        TopicArn: aws.String(s.topic),
    }

    _, err = s.client.Publish(input)
    if err != nil {
        return fmt.Errorf("failed to publish to SNS: %v", err)
    }

    return nil
}



// package domain

// // SNSMessageはSNSに送信するメッセージの構造体です。
// type SNSMessage struct {
//     ChannelID string
//     Message   string
//     Mentions  []string
// }

// package main

// import (
//     "context"
//     "log"
//     "myapp/domain"
//     "myapp/infrastructure"
//     "myapp/usecase"
// )

// func main() {
//     // インフラ（SNSリポジトリ）のセットアップ
//     snsRepo := infrastructure.NewSNSRepository("arn:aws:sns:region:account-id:topic-name")

//     // ユースケースのセットアップ
//     snsUseCase := usecase.NewSNSUseCase(snsRepo)

//     // メッセージデータ
//     message := domain.SNSMessage{
//         ChannelID: "CXXXXXXXXX",
//         Message:   "This is a notification with a mention.",
//         Mentions:  []string{"U12345678", "U87654321"},
//     }

//     // メッセージの送信
//     err := snsUseCase.PublishMessage(context.Background(), message)
//     if err != nil {
//         log.Fatalf("Failed to publish message: %v", err)
//     }

//     log.Println("Message successfully published to SNS.")
// }
