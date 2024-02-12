package rate_limiter_test

import (
	"testing"
	"time"

	rate_limiter "github.com/brewinski/systems-design/rate-limter"
)

func TestTokenBucket_Request(t *testing.T) {
	type fields struct {
		maxTokens  float64
		refillRate float64
	}
	type args struct {
		tokens           float64
		numberOfRequests int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Allows a requests when enough tokens are available.",
			want: true,
			fields: fields{
				maxTokens:  10,
				refillRate: 1,
			},
			args: args{
				tokens:           1,
				numberOfRequests: 1,
			},
		},
		{
			name: "Denies a requests when enough tokens are available.",
			want: false,
			fields: fields{
				maxTokens:  10,
				refillRate: 1,
			},
			args: args{
				tokens:           15,
				numberOfRequests: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			for i := 0; i < tt.args.numberOfRequests; i++ {
				tb := rate_limiter.NewTokenBucket(tt.fields.maxTokens, tt.fields.refillRate)

				if got := tb.Request(tt.args.tokens); got != tt.want {
					t.Errorf("TokenBucket.Request() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestTokenBucket_RequestRateLimitingAfterXCalls(t *testing.T) {
	type fields struct {
		maxTokens  float64
		refillRate float64
	}
	type args struct {
		tokens           float64
		numberOfRequests int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Allows 10 requests per second and rate limits on the last request",
			want: false,
			fields: fields{
				maxTokens:  10,
				refillRate: 1,
			},
			args: args{
				tokens:           1,
				numberOfRequests: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tb := rate_limiter.NewTokenBucket(tt.fields.maxTokens, tt.fields.refillRate)

			for i := 0; i < int(tt.fields.maxTokens); i++ {
				canRequest := tb.Request(1)
				if !canRequest {
					t.Errorf("TokenBucket.Request() = %v, want %v", canRequest, true)
				}
			}

			if got := tb.Request(tt.args.tokens); got != tt.want {
				t.Errorf("TokenBucket.Request() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenBucket_RequestBucketRefilling(t *testing.T) {
	type fields struct {
		maxTokens  float64
		refillRate float64
		sleepTime  time.Duration
	}

	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Refills the bucket after 10 seconds and allows 20 requests",
			fields: fields{
				maxTokens:  10,
				refillRate: 1,
				sleepTime:  10 * time.Second,
			},
			want: true,
		},
		{
			name: "Refills the bucket after 1 second and rate limits requests",
			fields: fields{
				maxTokens:  10,
				refillRate: 1,
				sleepTime:  1 * time.Second,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tb := rate_limiter.NewTokenBucket(tt.fields.maxTokens, tt.fields.refillRate)

			for i := 0; i < int(tt.fields.maxTokens); i++ {
				canRequest := tb.Request(1)
				if !canRequest && tt.want {
					t.Errorf("TokenBucket.Request() = %v, want %v", canRequest, tt.want)
					return
				}
			}

			time.Sleep(tt.fields.sleepTime)

			for i := 0; i < int(tt.fields.maxTokens); i++ {
				canRequest := tb.Request(1)
				if !canRequest && tt.want {
					t.Errorf("TokenBucket.Request() = %v, want %v", canRequest, tt.want)
					return
				}
			}

			if !tt.want {
				t.Errorf("TokenBucket.Request(), expected failure but got success")
			}
		})
	}
}
