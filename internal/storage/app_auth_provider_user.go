package storage

import (
	"github.com/ddrinkle/oa2"
	"github.com/ddrinkle/platform/query"
	"errors"

	uuid "github.com/satori/go.uuid"
)

// GetAuthProviders gets a list of auth_providers
func (s Crate) GetAppAuthProviderUsers(sqlCondition query.SQLCondition, page int, limit int) (oa2.AppAuthProviderUsers, error) {

	qry := `SELECT aapu.uuid,a.uuid,ap.uuid,provider_user_id, provider_user_name, access_token,refresh_token
					FROM AppAuthProviderUser aapu
					JOIN App a ON a.id = aapu.app_id
					JOIN AuthProvider ap ON ap.id = aapu.auth_provider_id`

	qry += sqlCondition.ToSQL()

	qry += query.Limit(page, limit)

	rows, err := s.DB.Query(qry, sqlCondition.Values()...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	aa := oa2.AppAuthProviderUsers{}
	for rows.Next() {
		a := oa2.AppAuthProviderUser{
			AuthProvider: &oa2.AuthProvider{},
			App:          &oa2.App{},
		}
		err = rows.Scan(
			&a.UUID,
			&a.App.UUID,
			&a.AuthProvider.UUID,
			&a.UserID,
			&a.UserName,
			&a.AccessToken,
			&a.RefreshToken,
		)
		if err != nil {
			return nil, err
		}
		aa = append(aa, a)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return aa, nil

}

// GetAppAuthProviderUser gets a list of auth_providers
func (s Crate) GetAppAuthProviderUser(uid uuid.UUID) (oa2.AppAuthProviderUser, error) {

	condition := query.NewSQLConditionFromCriteria([]query.Criteria{{
		Field: "aapu.uuid",
		Value: uid,
	}})

	authProviderUsers, err := s.GetAppAuthProviderUsers(condition, 1, 1)

	if err != nil {
		return oa2.AppAuthProviderUser{}, err
	}

	if len(authProviderUsers) == 0 {
		return oa2.AppAuthProviderUser{}, query.ErrRecordNotFound
	}
	return authProviderUsers[0], nil

}

// AddAppAuthProviderUser adds a AuthProvider
func (s Crate) AddAppAuthProviderUser(a *oa2.AppAuthProviderUser) (uuid.UUID, error) {

	if a.UUID == uuid.Nil {
		a.UUID = uuid.NewV4()
	}

	if a.App == nil || a.App.UUID == uuid.Nil {
		return uuid.Nil, errors.New("Missing App UUID")
	}
	if a.AuthProvider == nil || a.AuthProvider.UUID == uuid.Nil {
		return uuid.Nil, errors.New("Missing AuthProvider UUID")
	}

	_, err := s.DB.Exec(`INSERT into AppAuthProviderUser
				(uuid,app_id,auth_provider_id,provider_user_id,provider_user_name,access_token,refresh_token)
				VALUES (?,
					(SELECT id FROM App where uuid=? LIMIT 1),
					(SELECT id FROM AuthProvider where uuid=? LIMIT 1),
					?,?,?,?)`,
		a.UUID,
		a.App.UUID,
		a.AuthProvider.UUID,
		a.UserID,
		a.UserName,
		a.AccessToken,
		a.RefreshToken,
	)

	if err != nil {
		return uuid.Nil, err
	}
	return a.UUID, nil
}

// UpdateAppAuthProviderUser adds a AuthProvider
func (s Crate) UpdateAppAuthProviderUser(a oa2.AppAuthProviderUser) (uuid.UUID, error) {

	if a.UUID == uuid.Nil {
		return uuid.Nil, errors.New("Missing UUID")
	}

	_, err := s.DB.Exec(`Update AppAuthProviderUser set provider_user_name=?, access_token=?, refresh_token=? where
				uuid = ?`,
		a.UserName,
		a.AccessToken,
		a.RefreshToken,
		a.UUID,
	)

	if err != nil {
		return uuid.Nil, err
	}
	return a.UUID, nil
}
