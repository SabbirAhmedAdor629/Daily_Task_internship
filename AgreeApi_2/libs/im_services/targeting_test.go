package im_services

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type mockDynamoClient struct{ dynamodbiface.DynamoDBAPI }

func (svc *mockDynamoClient) Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	if v, ok := input.ExpressionAttributeNames["#Q"]; ok {
		if *v == "player_id" {
			if input.ExpressionAttributeValues != nil {
				if _, ok := input.ExpressionAttributeValues[":q"]; ok {
					switch {
					case *input.ExpressionAttributeValues[":q"].N == "15814820":
						return &dynamodb.QueryOutput{
							ConsumedCapacity: &dynamodb.ConsumedCapacity{},
							Count:            new(int64),
							Items: []map[string]*dynamodb.AttributeValue{
								{
									"player_id": &dynamodb.AttributeValue{
										N: input.ExpressionAttributeValues[":q"].N,
									},
									"prime_time": &dynamodb.AttributeValue{
										M: map[string]*dynamodb.AttributeValue{
											"afternoon": &dynamodb.AttributeValue{
												S: aws.String("16"),
											},
											"morning": &dynamodb.AttributeValue{
												S: aws.String("8"),
											},
											"evening": &dynamodb.AttributeValue{
												S: aws.String("24"),
											},
											"best": &dynamodb.AttributeValue{
												S: aws.String("22"),
											},
											"swing": &dynamodb.AttributeValue{
												S: aws.String("21"),
											},
										},
									},
									"time_zone_offset": &dynamodb.AttributeValue{
										S: aws.String("0"),
									},
								},
							},
							LastEvaluatedKey: nil,
							ScannedCount:     new(int64),
						}, nil

					case *input.ExpressionAttributeValues[":q"].N == "15826770":
						return &dynamodb.QueryOutput{
							ConsumedCapacity: &dynamodb.ConsumedCapacity{},
							Count:            new(int64),
							Items: []map[string]*dynamodb.AttributeValue{
								{
									"player_id": &dynamodb.AttributeValue{
										N: input.ExpressionAttributeValues[":q"].N,
									},
									"prime_time": &dynamodb.AttributeValue{
										M: map[string]*dynamodb.AttributeValue{
											"afternoon": &dynamodb.AttributeValue{
												S: aws.String("21"),
											},
											"morning": &dynamodb.AttributeValue{
												NULL: aws.Bool(true),
											},
											"evening": &dynamodb.AttributeValue{
												NULL: aws.Bool(true),
											},
											"best": &dynamodb.AttributeValue{
												S: aws.String("11"),
											},
											"swing": &dynamodb.AttributeValue{
												NULL: aws.Bool(true),
											},
										},
									},
									"time_zone_offset": &dynamodb.AttributeValue{
										S: aws.String("-4"),
									},
								},
							},
							LastEvaluatedKey: nil,
							ScannedCount:     new(int64),
						}, nil
					}

				}

			}
		}
	}
	return &dynamodb.QueryOutput{}, fmt.Errorf("Something went wrong")
}

type TargetingPlayerHash struct {
	PlayerId       int               `json:"player_id"`
	TimeZoneOffset int               `json:"time_zone_offset"`
	PrimeTime      map[string]string `json:"prime_time"`
}

func TestPlayer(t *testing.T) {
	type args struct {
		playerId  PlayerIdType
		tableName string
	}

	tests := []struct {
		name    string
		svc     *mockDynamoClient
		args    args
		want    TargetingPlayerHash
		wantErr bool
	}{
		{
			name: "#1: Player id present in mock dynamodb",
			svc:  &mockDynamoClient{},
			args: args{
				playerId:  15814820,
				tableName: "targets_development",
			},
			want: TargetingPlayerHash{
				PlayerId:       15814820,
				TimeZoneOffset: 0,
				PrimeTime:      map[string]string{"afternoon": "16", "best": "22", "evening": "24", "morning": "8", "swing": "21"},
			},
			wantErr: false,
		},
		{
			name: "#2: Player id present in mock dynamodb",
			svc:  &mockDynamoClient{},
			args: args{
				playerId:  15826770,
				tableName: "targets_development",
			},
			want: TargetingPlayerHash{
				PlayerId:       15826770,
				TimeZoneOffset: -4,
				PrimeTime:      map[string]string{"afternoon": "21", "best": "11", "evening": "", "morning": "", "swing": ""},
			},
			wantErr: false,
		},
		{
			name: "#3: Player id not present in mock dynamodb",
			svc:  &mockDynamoClient{},
			args: args{
				playerId:  5435334,
				tableName: "targets_development",
			},
			want:    TargetingPlayerHash{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Player(tt.svc, tt.args.tableName, tt.args.playerId)
			playerHashMap := TargetingPlayerHash{}
			_ = json.Unmarshal([]byte(got), &playerHashMap)
			if (err != nil) != tt.wantErr {
				t.Errorf("Player targeting error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(playerHashMap, tt.want) {
				t.Errorf("Player targeting error = %v, want %v", playerHashMap, tt.want)
			}
		})
	}
}
