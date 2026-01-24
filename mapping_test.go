package ocl

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func testServer[T any](responseCode int, responseBody T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(responseCode)

		json.NewEncoder(w).Encode(responseBody)
	}))
}

func TestClient_FetchMappings(t *testing.T) {
	type fields struct {
		token string
		HTTP  *http.Client
	}
	type args struct {
		ctx          context.Context
		searchParams map[string]string
		headers      *Headers
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		want         any
		responseCode int
		wantErr      bool
	}{
		{
			name: "Happy case: successfully searched for mapping",
			fields: fields{
				token: "test-token",
				HTTP:  http.DefaultClient,
			},
			args: args{
				ctx: t.Context(),
				searchParams: map[string]string{
					"toConceptCode": "REG-10001",
				},
				headers: &Headers{
					Organisation: "org-001",
					Source:       "test-source",
					MappingID:    "2434",
				},
			},
			wantErr: false,
			want: []Mapping{
				{
					Retired: false,
				},
			},
			responseCode: 200,
		},
		{
			name: "Sad case: error occurred while searching for mapping",
			fields: fields{
				token: "test-token",
				HTTP:  http.DefaultClient,
			},
			args: args{
				ctx: t.Context(),
				searchParams: map[string]string{
					"toConceptCode": "REG-10001",
				},
				headers: &Headers{
					Organisation: "org-001",
					Source:       "test-source",
					MappingID:    "2434",
				},
			},
			wantErr:      true,
			want:         "Internal server error",
			responseCode: 500,
		},
		{
			name: "Sad case: error occurred while searching for mapping - incorrect source or org",
			fields: fields{
				token: "test-token",
				HTTP:  http.DefaultClient,
			},
			args: args{
				ctx: t.Context(),
				searchParams: map[string]string{
					"toConceptCode": "REG-10001",
				},
				headers: &Headers{
					Organisation: "non-existent",
					Source:       "test-source",
				},
			},
			wantErr:      true,
			want:         "Not found",
			responseCode: 400,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := testServer(tt.responseCode, tt.want)
			defer server.Close()

			c := &Client{
				baseURL: server.URL,
				token:   tt.fields.token,
				HTTP:    tt.fields.HTTP,
			}

			got, err := c.FetchMappings(tt.args.ctx, tt.args.searchParams, tt.args.headers)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Client.FetchMappings() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.FetchMappings() = %v, want %v", got, tt.want)
			}
		})
	}
}
