package service_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/ddrinkle/oa2"
	"github.com/ddrinkle/oa2/internal/service"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func TestNewGoogleService(t *testing.T) {

	tests := []struct {
		name     string
		input    string
		expected oa2.OA2ServiceI
	}{
		{
			name:  "Positive NewGoogleService",
			input: "dev",
			expected: &service.GoogleOA2Service{
				oa2.OA2Service{
					UserInfoEndpoint: "https://www.googleapis.com/oauth2/v1/userinfo",
					Config: oauth2.Config{
						Endpoint: google.Endpoint,
						Scopes:   []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
					},
				},
			},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			got, _ := service.NewGoogleService(tt.input)
			if !reflect.DeepEqual(got, tt.expected) {
				expected, _ := json.Marshal(tt.expected)
				gotstr, _ := json.Marshal(got)
				t.Errorf("Expected:" + string(expected) + " but got:" + string(gotstr))
			}
		})
	}
}

func TestGoogleParseUserInfoResponse(t *testing.T) {

	tests := []struct {
		name         string
		input        string
		expected     oa2.AuthProviderUser
		expect_error bool
	}{
		{
			name:  "Positive NewGoogleService",
			input: `{"id":"id","name":"username","given_name":"foo","family_name":"bar","email":"foo@bar.com","verified_email":true}`,
			expected: oa2.AuthProviderUser{
				Id:            "id",
				UserName:      "username",
				FirstName:     "foo",
				LastName:      "bar",
				VerifiedEmail: true,
				Email:         "foo@bar.com",
			},
		},
		{
			name:         "Negative NewGoogleService",
			input:        `"id":"id","name":"username","given_name":"foo","family_name":"bar","email":"foo@bar.com","verified_email":true}`,
			expect_error: true,
		},
	}
	for _, tt := range tests {
		fb, _ := service.NewGoogleService("dev")
		t.Run(tt.name, func(t *testing.T) {
			got, err := fb.ParseUserInfoResponse([]byte(tt.input))

			if !reflect.DeepEqual(got, tt.expected) && (err == nil) {
				expected, _ := json.Marshal(tt.expected)
				gotstr, _ := json.Marshal(got)
				t.Errorf("Expected:" + string(expected) + " but got:" + string(gotstr))
			}
			if err != nil && !tt.expect_error {
				t.Errorf("Unexpected Error Condition. Did not expect Error, but received: " + err.Error())
			}
			if err == nil && tt.expect_error {
				t.Errorf("Unexpected Error Condition.  Expected Error, received nil")
			}
		})
	}
}

func TestGoogleGetAuthorizationMethod(t *testing.T) {
	t.Run("Positive Google:GetAuthorizationMethod Test", func(t *testing.T) {
		fb, _ := service.NewGoogleService("dev")
		got := fb.GetAuthorizationMethod()
		if got != oa2.AUTHORIZATION_METHOD_HEADER_OAUTH {
			t.Errorf("Incorrect Authorization Method. Expected AUTHORIZATION_METHOD_HEADER_OAUTH")
		}
	})
}
func TestGoogleIsTrustedLoginService(t *testing.T) {
	t.Run("Positive Google:IsTrustedLoginService Test", func(t *testing.T) {
		fb, _ := service.NewGoogleService("dev")
		got := fb.IsTrustedLoginService()
		if got != false {
			t.Errorf("Incorrect Trusted Status. Expected false but got true")
		}
	})
}
