package service_test

import (
	"github.com/ddrinkle/oa2"
	"github.com/ddrinkle/oa2/internal/service"
	"github.com/ddrinkle/platform/crypt"
	"encoding/json"
	"reflect"
	"testing"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
)

// servicekey encrypted = L3d0Q3HczAMTZDhDJ8ZufnEAwpRecbiRdi0=
// servicesecret encrypted = IGYL7huGaojyvruzYCRdAN2y5QAP6WNPS0rb3J4=

var encryptKey = []byte("12345678901234567890123456789012")

var clientKey = "servicekey"
var clientSecret = "servicesecret"

func TestGetServiceFromAuthProvider(t *testing.T) {

	crypt.SetCryptKeyFunction(func() []byte {
		return encryptKey
	})

	tests := []struct {
		name         string
		input        oa2.AuthProvider
		expected     oa2.OA2ServiceI
		expect_error bool
		env          string
	}{
		{
			name: "Positive Facebook Auth Provider",
			env:  "dev",
			input: oa2.AuthProvider{
				ClientKey:      &crypt.String{},
				ClientSecret:   &crypt.String{},
				OA2ServiceName: "facebook",
			},
			expected: &service.FacebookOA2Service{
				oa2.OA2Service{
					UserInfoEndpoint: "https://graph.facebook.com/v2.5/me?fields=name,first_name,last_name,email,verified",
					Config: oauth2.Config{
						Endpoint:     facebook.Endpoint,
						ClientID:     "servicekey",
						ClientSecret: "servicesecret",
						Scopes:       []string{"email", "public_profile"},
					},
				},
			},
		},
		{
			name: "Positive Google Auth Provider",
			env:  "dev",
			input: oa2.AuthProvider{
				ClientKey:      &crypt.String{},
				ClientSecret:   &crypt.String{},
				OA2ServiceName: "google",
			},
			expected: &service.GoogleOA2Service{
				oa2.OA2Service{
					UserInfoEndpoint: "https://www.googleapis.com/oauth2/v1/userinfo",
					Config: oauth2.Config{
						Endpoint:     google.Endpoint,
						ClientID:     "servicekey",
						ClientSecret: "servicesecret",
						Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := service.NewOA2Factory(tt.env)

			tt.input.ClientKey.Set(clientKey)
			tt.input.ClientSecret.Set(clientSecret)

			got, err := factory.GetServiceFromAuthProvider(tt.input)

			if err != nil && !tt.expect_error {
				t.Errorf("Unexpected Error:" + err.Error())
			}

			if !reflect.DeepEqual(got, tt.expected) {
				expected, _ := json.Marshal(tt.expected)
				gotstr, _ := json.Marshal(got)
				t.Errorf("Expected:" + string(expected) + "\nbut got:              " + string(gotstr))
			}
		})
	}
}
