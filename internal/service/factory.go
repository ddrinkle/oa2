package service

import (
	"github.com/ddrinkle/oa2"
	"errors"
)

type OA2Factory struct {
	EnvironmentName string
}

func NewOA2Factory(env string) *OA2Factory {
	return &OA2Factory{
		EnvironmentName: env,
	}
}

func (f *OA2Factory) GetServiceFromAuthProvider(provider oa2.AuthProvider) (oa2.OA2ServiceI, error) {

	clientKey, err := provider.ClientKey.Get()
	if err != nil {
		return nil, err
	}
	clientSecret, err := provider.ClientSecret.Get()
	if err != nil {
		return nil, err
	}

	service, err := f.getServiceFromName(provider.OA2ServiceName)
	if err != nil {
		return nil, err
	}

	service.SetClientKey(clientKey)
	service.SetClientSecret(clientSecret)

	return service, nil
}

func (f *OA2Factory) getServiceFromName(name string) (oa2.OA2ServiceI, error) {
	var svc oa2.OA2ServiceI
	var err error

	switch name {
	case "facebook":
		svc, err = NewFacebookService(f.EnvironmentName)
	case "google":
		svc, err = NewGoogleService(f.EnvironmentName)
	case "eveonline":
		svc, err = NewEveService(f.EnvironmentName)
	default:
		return nil, errors.New("Invalid Service Name:" + name)
	}

	return svc, err
}
