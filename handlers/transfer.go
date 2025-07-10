package handlers

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Transfer struct {
    PlayerID  string    `json:"playerId"`
    TeamFrom  string    `json:"teamFrom"`
    TeamTo    string    `json:"teamTo"`
    Amount    float64   `json:"amount"`
    CreatedAt time.Time `json:"createdAt"`
}

func CreateTransfer(w http.ResponseWriter, r *http.Request) {
    var transfer Transfer
    _ = json.NewDecoder(r.Body).Decode(&transfer)
    transfer.CreatedAt = time.Now()

    sess := session.Must(session.NewSession(&aws.Config{
        Region: aws.String("us-east-1"),
    }))
    svc := dynamodb.New(sess)

    av, _ := dynamodbattribute.MarshalMap(transfer)
    input := &dynamodb.PutItemInput{
        TableName: aws.String("Transfers"),
        Item:      av,
    }

    _, err := svc.PutItem(input)
    if err != nil {
        http.Error(w, "Failed to save transfer", http.StatusInternalServerError)
        return
    }

    fmt.Printf("EVENT: transfer.created -> %+v\n", transfer)
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(transfer)
}

func GetTransfers(w http.ResponseWriter, r *http.Request) {
    sess := session.Must(session.NewSession(&aws.Config{
        Region: aws.String("us-east-1"),
    }))
    svc := dynamodb.New(sess)

    input := &dynamodb.ScanInput{
        TableName: aws.String("Transfers"),
    }

    result, err := svc.Scan(input)
    if err != nil {
        http.Error(w, "Failed to retrieve transfers", http.StatusInternalServerError)
        return
    }

    var transfers []Transfer
    _ = dynamodbattribute.UnmarshalListOfMaps(result.Items, &transfers)
    json.NewEncoder(w).Encode(transfers)
}
