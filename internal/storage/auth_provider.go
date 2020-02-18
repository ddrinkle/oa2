package storage

import (
	"github.com/ddrinkle/oa2"
	"github.com/ddrinkle/platform/query"

	uuid "github.com/satori/go.uuid"
)

// GetAuthProviders gets a list of auth_providers
func (s Crate) GetAuthProviders(sqlCondition query.SQLCondition, page int, limit int) (oa2.AuthProviders, error) {

	qry := `SELECT uuid,name,oa2_service_name,client_key,secret_key from AuthProvider `

	qry += sqlCondition.ToSQL()

	qry += query.Limit(page, limit)

	rows, err := s.DB.Query(qry, sqlCondition.Values()...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	aa := oa2.AuthProviders{}
	for rows.Next() {
		a := oa2.AuthProvider{}
		err = rows.Scan(
			&a.UUID,
			&a.Name,
			&a.OA2ServiceName,
			&a.ClientKey,
			&a.ClientSecret,
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

// GetAuthProvidersForApp gets a list of auth_providers
func (s Crate) GetAuthProvidersForApp(uid uuid.UUID) (oa2.AuthProviders, error) {

	condition := query.NewSQLConditionFromCriteria([]query.Criteria{{
		Field: "a.uuid",
		Value: uid,
	}})

	qry := `SELECT ap.uuid,ap.name,ap.oa2_service_name from AuthProvider ap
		join AppAuthProviders aps on aps.auth_provider_id = ap.id
		join App a on a.id = aps.app_id`

	qry += condition.ToSQL()

	rows, err := s.DB.Query(qry, condition.Values()...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	aa := oa2.AuthProviders{}
	for rows.Next() {
		a := oa2.AuthProvider{}
		err = rows.Scan(
			&a.UUID,
			&a.Name,
			&a.OA2ServiceName,
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

// GetAuthProvider gets a list of auth_providers
func (s Crate) GetAuthProvider(uid uuid.UUID) (oa2.AuthProvider, error) {

	condition := query.NewSQLConditionFromCriteria([]query.Criteria{{
		Field: "uuid",
		Value: uid,
	}})

	authProviders, err := s.GetAuthProviders(condition, 1, 1)

	if err != nil {
		return oa2.AuthProvider{}, err
	}

	if len(authProviders) == 0 {
		return oa2.AuthProvider{}, query.ErrRecordNotFound
	}
	return authProviders[0], nil

}

// AddAuthProvider adds a AuthProvider
func (s Crate) AddAuthProvider(a *oa2.AuthProvider) (uuid.UUID, error) {

	if a.UUID == uuid.Nil {
		a.UUID = uuid.NewV4()
	}

	_, err := s.DB.Exec("INSERT into AuthProvider (uuid,name,oa2_service_name,client_key,secret_key,authorization_method) VALUES (?,?,?,?,?,?)",
		a.UUID,
		a.Name,
		a.OA2ServiceName,
		a.ClientKey,
		a.ClientSecret,
		a.AuthorizationMethod,
	)

	if err != nil {
		return uuid.Nil, err
	}
	return a.UUID, nil
}
