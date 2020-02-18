package service_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/ddrinkle/oa2"
	"github.com/ddrinkle/oa2/internal/service"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

func TestNewFacebookService(t *testing.T) {

	tests := []struct {
		name     string
		input    string
		expected oa2.OA2ServiceI
	}{
		{
			name:  "Positive NewFacebookService",
			input: "dev",
			expected: &service.FacebookOA2Service{
				oa2.OA2Service{
					UserInfoEndpoint: "https://graph.facebook.com/v2.5/me?fields=name,first_name,last_name,email,verified",
					Config: oauth2.Config{
						Endpoint: facebook.Endpoint,
						Scopes:   []string{"email", "public_profile"},
					},
				},
			},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			got, _ := service.NewFacebookService(tt.input)
			if !reflect.DeepEqual(got, tt.expected) {
				expected, _ := json.Marshal(tt.expected)
				gotstr, _ := json.Marshal(got)
				t.Errorf("Expected:" + string(expected) + " but got:" + string(gotstr))
			}
		})
	}
}

func TestParseUserInfoResponse(t *testing.T) {

	tests := []struct {
		name         string
		input        string
		expected     oa2.AuthProviderUser
		expect_error bool
	}{
		{
			name:  "Positive NewFacebookService",
			input: `{"id":"id","name":"username","first_name":"foo","last_name":"bar","email":"foo@bar.com"}`,
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
			name:         "Negative NewFacebookService",
			input:        `"id":"id","name":"username","first_name":"foo","last_name":"bar","email":"foo@bar.com"}`,
			expect_error: true,
		},
	}
	for _, tt := range tests {
		fb, _ := service.NewFacebookService("dev")
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

func TestFacebookGetAuthorizationMethod(t *testing.T) {
	t.Run("Positive Facebook:GetAuthorizationMethod Test", func(t *testing.T) {
		fb, _ := service.NewFacebookService("dev")
		got := fb.GetAuthorizationMethod()
		if got != oa2.AUTHORIZATION_METHOD_HEADER_OAUTH {
			t.Errorf("Incorrect Authorization Method. Expected AUTHORIZATION_METHOD_HEADER_OAUTH")
		}
	})
}
func TestFacebookIsTrustedLoginService(t *testing.T) {
	t.Run("Positive Facebook:IsTrustedLoginService Test", func(t *testing.T) {
		fb, _ := service.NewFacebookService("dev")
		got := fb.IsTrustedLoginService()
		if got != false {
			t.Errorf("Incorrect Trusted Status. Expected false but got true")
		}
	})
}
