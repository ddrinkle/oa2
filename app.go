package oa2

import (
	"github.com/ddrinkle/platform/crypt"
	"errors"

	"github.com/manyminds/api2go/jsonapi"
	uuid "github.com/satori/go.uuid"
)

type App struct {
	UUID      uuid.UUID     `json:"-" jsonapi:"primary,app"`
	Name      string        `json:"name,omitempty" jsonapi:"attr,name"`
	AppKey    *crypt.String `json:"app_key,omitempty" jsonapi:"attr:app_key"`
	SecretKey *crypt.String `json:"secret_key,omitempty" jsonapi:"attr:secret_key"`

	AuthProviders []AuthProvider `json:"-" jsonapi:"relation,auth_providers"`
}

func (a App) GetName() string {
	return "app"
}

func (a App) GetID() string {
	return a.UUID.String()
}

func (a *App) SetID(input string) error {
	var err error

	if input != "" {
		a.UUID, err = uuid.FromString(input)
	}
	return err
}

func (a App) GetReferences() []jsonapi.Reference {

	return []jsonapi.Reference{
		{
			Type:         "AuthProviders",
			Name:         "auth_providers",
			IsNotLoaded:  false,
			Relationship: jsonapi.ToManyRelationship,
		},
	}
}

func (a App) GetReferencedIDs() []jsonapi.ReferenceID {
	refIDs := []jsonapi.ReferenceID{}
	for _, p := range a.AuthProviders {
		refIDs = append(refIDs, jsonapi.ReferenceID{
			ID:           p.GetID(),
			Type:         "auth_provider",
			Name:         "auth_providers",
			Relationship: jsonapi.ToManyRelationship,
		})
	}
	return refIDs
}

func (a App) GetReferencedStructs() []jsonapi.MarshalIdentifier {

	result := []jsonapi.MarshalIdentifier{}
	for _, p := range a.AuthProviders {
		result = append(result, p)
	}

	return result
}

func (a *App) SetToManyReferenceIDs(name string, IDs []string) error {
	if name == "auth_providers" {
		for _, id := range IDs {
			uid, err := uuid.FromString(id)
			if err != nil {
				return err
			}
			a.AuthProviders = append(a.AuthProviders, AuthProvider{
				UUID: uid,
			})
		}
		return nil
	}

	return errors.New("There is no to-many relationship with the name " + name)
}
