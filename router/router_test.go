package router_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brewinski/systems-design/router"
)

func TestRouter_ServeHttp(t *testing.T) {
	type args struct {
		w      *httptest.ResponseRecorder
		r      *http.Request
		routes []router.RouteEntry
	}
	tests := []struct {
		name string
		sr   *router.Router
		args args
		want int
	}{
		{
			name: "Return 404 for unmatched route",
			sr:   &router.Router{},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "http://testing.test", nil),
			},
			want: 404,
		},
		{
			name: "Return 200 for matched route",
			sr:   &router.Router{},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "http://testing.test/", nil),
				routes: []router.RouteEntry{
					{
						Path:   "/",
						Method: "GET",
						Handler: func(w http.ResponseWriter, r *http.Request) {
							w.Write([]byte("testing"))
						},
					},
				},
			},
			want: 200,
		},
		{
			name: "Return 200 for different route",
			sr:   &router.Router{},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "http://testing.test/different", nil),
				routes: []router.RouteEntry{
					{
						Path:   "/different",
						Method: "GET",
						Handler: func(w http.ResponseWriter, r *http.Request) {
							w.Write([]byte("testing"))
						},
					},
				},
			},
			want: 200,
		},
		{
			name: "Return 200 for POST route",
			sr:   &router.Router{},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "http://testing.test/different", nil),
				routes: []router.RouteEntry{
					{
						Path:   "/different",
						Method: "POST",
						Handler: func(w http.ResponseWriter, r *http.Request) {
							w.Write([]byte("testing"))
						},
					},
				},
			},
			want: 200,
		},
		{
			name: "Return 404 for unregistered POST route",
			sr:   &router.Router{},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "http://testing.test/different", nil),
			},
			want: 404,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, re := range tt.args.routes {
				tt.sr.Route(re.Method, re.Path, re.Handler)
			}

			tt.sr.ServeHttp(tt.args.w, tt.args.r)
			resp := tt.args.w.Result()

			if resp.StatusCode != tt.want {
				t.Errorf("Response status: want %d, got %d", tt.want, resp.StatusCode)
			}
		})
	}
}
