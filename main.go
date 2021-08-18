package main

import (
    "context"
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
    "github.com/aws/aws-sdk-go/service/s3"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/aws"
)
var sess *session.Session

func writeFile(body string,filePath string) error{

    d1 := []byte(body)
    err := ioutil.WriteFile(filePath, d1, 0644)
    return err
}

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {

     // Setup AWS S3 Session (build once use every function)
    sess = session.Must(session.NewSession(&aws.Config{
        Region: aws.String("ap-southeast-1"),
    }))


    for _, message := range sqsEvent.Records {
        fmt.Printf("The message %s for event source %s = %s \n", message.MessageId, message.EventSource, message.Body)
        file := "/tmp/test.txt"

        error:=writeFile(message.Body,file)
        if error != nil {
            fmt.Println("writeFile error", error.Error())
        }

        fileName, err := os.Open(file)
        if err != nil {
            log.Println("os.Open - filename: %s, err: %v", file, err)
            return err
        }
        defer fileName.Close()

        _, err = s3.New(sess).PutObject(&s3.PutObjectInput{
            Bucket: aws.String("s3healthians-dev/"),
            Key:    aws.String("shweta_test.txt"),
            Body:  fileName,
            ContentType:aws.String("text/plain"),
        })
    }
    




    return nil
}

func main() {
    lambda.Start(handler)
}