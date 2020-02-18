package storage

import (
	"github.com/ddrinkle/oa2"
	"github.com/ddrinkle/platform/query"

	uuid "github.com/satori/go.uuid"
)

// GetApps returns all apps
func (s Crate) GetApps(sqlCondition query.SQLCondition, page int, limit int) ([]oa2.App, error) {

	qry := `SELECT uuid,name,app_key,secret_key from App `

	qry += sqlCondition.ToSQL()

	qry += query.Limit(page, limit)

	rows, err := s.DB.Query(qry, sqlCondition.Values()...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	aa := []oa2.App{}

	for rows.Next() {
		a := oa2.App{}
		err = rows.Scan(
			&a.UUID,
			&a.Name,
			&a.AppKey,
			&a.SecretKey,
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
func (s Crate) GetApp(uid uuid.UUID) (oa2.App, error) {

	condition := query.NewSQLConditionFromCriteria([]query.Criteria{{
		Field: "uuid",
		Value: uid,
	}})

	apps, err := s.GetApps(condition, 1, 1)
	if err != nil {
		return oa2.App{}, err
	}

	if len(apps) == 0 {
		return oa2.App{}, query.ErrRecordNotFound
	}

	authProviders, err := s.GetAuthProvidersForApp(apps[0].UUID)
	if err != nil {
		return oa2.App{}, err
	}

	apps[0].AuthProviders = authProviders

	return apps[0], nil
}

// AddApp adds a App
func (s Crate) AddApp(a *oa2.App) (uuid.UUID, error) {

	if a.UUID == uuid.Nil {
		a.UUID = uuid.NewV4()
	}

	_, err := s.DB.Exec("INSERT into App (uuid,name,app_key,secret_key) VALUES (?,?,?,?)",
		a.UUID,
		a.Name,
		a.AppKey,
		a.SecretKey,
	)

	if err != nil {
		return uuid.Nil, err
	}
	return a.UUID, nil
}
