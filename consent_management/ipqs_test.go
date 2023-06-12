package main

import (
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	v3 "influencemobile.com/libs/v3_helpers"
)

type ClientMock struct {
}

func (c *ClientMock) Do(req *http.Request) (*http.Response, error) {
	stringReadCloser := io.NopCloser(strings.NewReader(`{"success":true,"message":"Success","fraud_score":100,"country_code":"GB","region":"England","city":"London","ISP":"OVH SAS","ASN":16276,"organization":"Gecko VPN","is_crawler":false,"timezone":"Europe\/London","mobile":false,"host":"127.0.0.1","proxy":true,"vpn":true,"tor":false,"active_vpn":true,"active_tor":false,"recent_abuse":true,"bot_status":true,"connection_type":"Premium required.","abuse_velocity":"Premium required.","zip_code":"N\/A","latitude":51.51,"longitude":-0.09,"request_id":"D01cFLBQIG"}`))
	stringReadCloser1 := io.NopCloser(strings.NewReader(`{"success":true,"message":"Success","fraud_score":100,"country_code":"GB","region":"England","city":"Belgium","ISP":"OVH SAS","ASN":16276,"organization":"Gecko VPN","is_crawler":false,"timezone":"Europe\/Belgium","mobile":false,"host":"127.0.0.10","proxy":true,"vpn":true,"tor":false,"active_vpn":true,"active_tor":false,"recent_abuse":true,"bot_status":true,"connection_type":"Premium required.","abuse_velocity":"Premium required.","zip_code":"N\/A","latitude":51.51,"longitude":-0.09,"request_id":"D01cFLBQIG"}`))
	stringReadCloser2 := io.NopCloser(strings.NewReader(`{"success":true,"message":"Success","fraud_score":100,"country_code":"GB","region":"England","city":"Dublin","ISP":"OVH SAS","ASN":16276,"organization":"Gecko VPN","is_crawler":false,"timezone":"Europe\/Dublin","mobile":false,"host":"127.0.0.20","proxy":true,"vpn":true,"tor":false,"active_vpn":true,"active_tor":false,"recent_abuse":true,"bot_status":true,"connection_type":"Premium required.","abuse_velocity":"Premium required.","zip_code":"N\/A","latitude":51.51,"longitude":-0.09,"request_id":"D01cFLBQIG"}`))

	switch req.URL.String() {
	case "https://ipqualityscore.com/api/json/ip/test-api-key/127.0.0.1":
		return &http.Response{Body: stringReadCloser, StatusCode: 200, Status: "200 Ok"}, nil
	case "https://ipqualityscore.com/api/json/ip/test-api-key/127.0.0.10":
		return &http.Response{Body: stringReadCloser1, StatusCode: 200, Status: "200 Ok"}, nil
	case "https://ipqualityscore.com/api/json/ip/test-api-key/127.0.0.20":
		return &http.Response{Body: stringReadCloser2, StatusCode: 200, Status: "200 Ok"}, nil
	default:
		return &http.Response{Body: io.NopCloser(strings.NewReader(`{"success":false,"message":"Unexpected error","timezone":""}`)), StatusCode: 200, Status: "200 Ok"}, nil
	}
}

func TestGetTimezoneFromIPQS(t *testing.T) {
	type args struct {
		apiKey string
		ip     string
	}
	tests := []struct {
		name     string
		args     args
		want     *IPQSResponse
		wantCode int
		wantErr  bool
	}{
		{
			name: "#1: Get Request success [IP: 127.0.0.1]",
			args: args{
				apiKey: "test-api-key",
				ip:     "127.0.0.1",
			},
			want: &IPQSResponse{
				Timezone: "Europe/London",
				Success:  true,
				Message:  "Success",
			},
			wantCode: int(v3.ResponseSuccess),
			wantErr:  false,
		},
		{
			name: "#2: Get Request success [IP: 127.0.0.10]",
			args: args{
				apiKey: "test-api-key",
				ip:     "127.0.0.10",
			},
			want: &IPQSResponse{
				Timezone: "Europe/Belgium",
				Success:  true,
				Message:  "Success",
			},
			wantCode: int(v3.ResponseSuccess),
			wantErr:  false,
		},
		{
			name: "#3: Get Request success [IP: 127.0.0.20]",
			args: args{
				apiKey: "test-api-key",
				ip:     "127.0.0.20",
			},
			want: &IPQSResponse{
				Timezone: "Europe/Dublin",
				Success:  true,
				Message:  "Success",
			},
			wantErr:  false,
			wantCode: int(v3.ResponseSuccess),
		},
		{
			name: "#4: Get Request Error [IP: 127.0.0.30]",
			args: args{
				apiKey: "test-api-key",
				ip:     "127.0.0.30",
			},
			want:     nil,
			wantCode: int(v3.ResponseInternalServerError),
			wantErr:  true,
		},
	}

	httpClient = &ClientMock{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, code, err := GetTimezoneFromIPQS(tt.args.apiKey, tt.args.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTimezoneFromIPQS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTimezoneFromIPQS() = %v, want %v", got, tt.want)
			}
			if code != tt.wantCode {
				t.Errorf("GetTimezoneFromIPQS() code = %v, wantCode %v", code, tt.wantCode)
			}

		})
	}
}
