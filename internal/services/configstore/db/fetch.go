
// Code generated by go generate; DO NOT EDIT.
package db

import (
	stdsql "database/sql"
	"encoding/json"

	"github.com/sorintlab/errors"
	sq "github.com/huandu/go-sqlbuilder"

	"agola.io/agola/internal/sqlg"
	"agola.io/agola/internal/sqlg/sql"

	types "agola.io/agola/services/configstore/types"

	"time"
)

func (d *DB) fetchRemoteSources(tx *sql.Tx, q sq.Builder) ([]*types.RemoteSource, []string, error) {
	rows, err := d.query(tx, q)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	defer rows.Close()

	return d.scanRemoteSources(rows, tx.ID(), 0)
}

func (d *DB) fetchRemoteSourcesSkipLastFields(tx *sql.Tx, q sq.Builder, skipFieldsCount uint) ([]*types.RemoteSource, []string, error) {
	rows, err := d.query(tx, q)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	defer rows.Close()

	return d.scanRemoteSources(rows, tx.ID(), skipFieldsCount)
}

func (d *DB) scanRemoteSource(rows *stdsql.Rows, skipFieldsCount uint) (*types.RemoteSource, string, error) {

	v := &types.RemoteSource{}

	var vi any = v
	if x, ok := vi.(sqlg.Initer); ok {
		x.Init()
	}

	fields := append([]any{&v.ID, &v.Revision, &v.CreationTime, &v.UpdateTime, &v.Name, &v.APIURL, &v.SkipVerify, &v.Type, &v.AuthType, &v.Oauth2ClientID, &v.Oauth2ClientSecret, &v.SSHHostKey, &v.SkipSSHHostKeyCheck, &v.RegistrationEnabled, &v.LoginEnabled})

	for i := uint(0); i < skipFieldsCount; i++ {
		fields = append(fields, new(any))
	}

	if err := rows.Scan(fields...); err != nil {
		return nil, "", errors.Wrap(err, "failed to scan row")
	}

	if x, ok := vi.(sqlg.PreJSONSetupper); ok {
		if err := x.PreJSON(); err != nil {
			return nil, "", errors.Wrap(err, "prejson error")
		}
	}

	return v, v.ID, nil
}

func (d *DB) scanRemoteSources(rows *stdsql.Rows, txID string, skipFieldsCount uint) ([]*types.RemoteSource, []string, error) {
	vs := []*types.RemoteSource{}
	ids := []string{}
	for rows.Next() {
		v, id, err := d.scanRemoteSource(rows, skipFieldsCount)
		if err != nil {
			rows.Close()
			return nil, nil, errors.WithStack(err)
		}
		v.TxID = txID
		vs = append(vs, v)
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, errors.WithStack(err)
	}
	return vs, ids, nil
}

func (d *DB) RemoteSourceArray() []any {
	a := []any{}
	a = append(a, new(string))
	a = append(a, new(uint64))
	a = append(a, new(time.Time))
	a = append(a, new(time.Time))
	a = append(a, new(string))
	a = append(a, new(string))
	a = append(a, new(bool))
	a = append(a, new(types.RemoteSourceType))
	a = append(a, new(types.RemoteSourceAuthType))
	a = append(a, new(string))
	a = append(a, new(string))
	a = append(a, new(string))
	a = append(a, new(bool))
	a = append(a, new(bool))
	a = append(a, new(bool))

	return a
}

func (d *DB) RemoteSourceFromArray(a []any, txID string) (*types.RemoteSource, string, error) {
	v := &types.RemoteSource{}

	var vi any = v
	if x, ok := vi.(sqlg.Initer); ok {
		x.Init()
	}
	v.ID = *a[0].(*string)
	v.Revision = *a[1].(*uint64)
	v.CreationTime = *a[2].(*time.Time)
	v.UpdateTime = *a[3].(*time.Time)
	v.Name = *a[4].(*string)
	v.APIURL = *a[5].(*string)
	v.SkipVerify = *a[6].(*bool)
	v.Type = *a[7].(*types.RemoteSourceType)
	v.AuthType = *a[8].(*types.RemoteSourceAuthType)
	v.Oauth2ClientID = *a[9].(*string)
	v.Oauth2ClientSecret = *a[10].(*string)
	v.SSHHostKey = *a[11].(*string)
	v.SkipSSHHostKeyCheck = *a[12].(*bool)
	v.RegistrationEnabled = *a[13].(*bool)
	v.LoginEnabled = *a[14].(*bool)

	if x, ok := vi.(sqlg.PreJSONSetupper); ok {
		if err := x.PreJSON(); err != nil {
			return nil, "", errors.Wrap(err, "prejson error")
		}
	}

	v.TxID = txID

	return v, v.ID, nil
}

func (d *DB) fetchUsers(tx *sql.Tx, q sq.Builder) ([]*types.User, []string, error) {
	rows, err := d.query(tx, q)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	defer rows.Close()

	return d.scanUsers(rows, tx.ID(), 0)
}

func (d *DB) fetchUsersSkipLastFields(tx *sql.Tx, q sq.Builder, skipFieldsCount uint) ([]*types.User, []string, error) {
	rows, err := d.query(tx, q)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	defer rows.Close()

	return d.scanUsers(rows, tx.ID(), skipFieldsCount)
}

func (d *DB) scanUser(rows *stdsql.Rows, skipFieldsCount uint) (*types.User, string, error) {

	v := &types.User{}

	var vi any = v
	if x, ok := vi.(sqlg.Initer); ok {
		x.Init()
	}

	fields := append([]any{&v.ID, &v.Revision, &v.CreationTime, &v.UpdateTime, &v.Name, &v.Secret, &v.Admin})

	for i := uint(0); i < skipFieldsCount; i++ {
		fields = append(fields, new(any))
	}

	if err := rows.Scan(fields...); err != nil {
		return nil, "", errors.Wrap(err, "failed to scan row")
	}

	if x, ok := vi.(sqlg.PreJSONSetupper); ok {
		if err := x.PreJSON(); err != nil {
			return nil, "", errors.Wrap(err, "prejson error")
		}
	}

	return v, v.ID, nil
}

func (d *DB) scanUsers(rows *stdsql.Rows, txID string, skipFieldsCount uint) ([]*types.User, []string, error) {
	vs := []*types.User{}
	ids := []string{}
	for rows.Next() {
		v, id, err := d.scanUser(rows, skipFieldsCount)
		if err != nil {
			rows.Close()
			return nil, nil, errors.WithStack(err)
		}
		v.TxID = txID
		vs = append(vs, v)
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, errors.WithStack(err)
	}
	return vs, ids, nil
}

func (d *DB) UserArray() []any {
	a := []any{}
	a = append(a, new(string))
	a = append(a, new(uint64))
	a = append(a, new(time.Time))
	a = append(a, new(time.Time))
	a = append(a, new(string))
	a = append(a, new(string))
	a = append(a, new(bool))

	return a
}

func (d *DB) UserFromArray(a []any, txID string) (*types.User, string, error) {
	v := &types.User{}

	var vi any = v
	if x, ok := vi.(sqlg.Initer); ok {
		x.Init()
	}
	v.ID = *a[0].(*string)
	v.Revision = *a[1].(*uint64)
	v.CreationTime = *a[2].(*time.Time)
	v.UpdateTime = *a[3].(*time.Time)
	v.Name = *a[4].(*string)
	v.Secret = *a[5].(*string)
	v.Admin = *a[6].(*bool)

	if x, ok := vi.(sqlg.PreJSONSetupper); ok {
		if err := x.PreJSON(); err != nil {
			return nil, "", errors.Wrap(err, "prejson error")
		}
	}

	v.TxID = txID

	return v, v.ID, nil
}

func (d *DB) fetchUserTokens(tx *sql.Tx, q sq.Builder) ([]*types.UserToken, []string, error) {
	rows, err := d.query(tx, q)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	defer rows.Close()

	return d.scanUserTokens(rows, tx.ID(), 0)
}

func (d *DB) fetchUserTokensSkipLastFields(tx *sql.Tx, q sq.Builder, skipFieldsCount uint) ([]*types.UserToken, []string, error) {
	rows, err := d.query(tx, q)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	defer rows.Close()

	return d.scanUserTokens(rows, tx.ID(), skipFieldsCount)
}

func (d *DB) scanUserToken(rows *stdsql.Rows, skipFieldsCount uint) (*types.UserToken, string, error) {

	v := &types.UserToken{}

	var vi any = v
	if x, ok := vi.(sqlg.Initer); ok {
		x.Init()
	}

	fields := append([]any{&v.ID, &v.Revision, &v.CreationTime, &v.UpdateTime, &v.UserID, &v.Name, &v.Value})

	for i := uint(0); i < skipFieldsCount; i++ {
		fields = append(fields, new(any))
	}

	if err := rows.Scan(fields...); err != nil {
		return nil, "", errors.Wrap(err, "failed to scan row")
	}

	if x, ok := vi.(sqlg.PreJSONSetupper); ok {
		if err := x.PreJSON(); err != nil {
			return nil, "", errors.Wrap(err, "prejson error")
		}
	}

	return v, v.ID, nil
}

func (d *DB) scanUserTokens(rows *stdsql.Rows, txID string, skipFieldsCount uint) ([]*types.UserToken, []string, error) {
	vs := []*types.UserToken{}
	ids := []string{}
	for rows.Next() {
		v, id, err := d.scanUserToken(rows, skipFieldsCount)
		if err != nil {
			rows.Close()
			return nil, nil, errors.WithStack(err)
		}
		v.TxID = txID
		vs = append(vs, v)
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, errors.WithStack(err)
	}
	return vs, ids, nil
}

func (d *DB) UserTokenArray() []any {
	a := []any{}
	a = append(a, new(string))
	a = append(a, new(uint64))
	a = append(a, new(time.Time))
	a = append(a, new(time.Time))
	a = append(a, new(string))
	a = append(a, new(string))
	a = append(a, new(string))

	return a
}

func (d *DB) UserTokenFromArray(a []any, txID string) (*types.UserToken, string, error) {
	v := &types.UserToken{}

	var vi any = v
	if x, ok := vi.(sqlg.Initer); ok {
		x.Init()
	}
	v.ID = *a[0].(*string)
	v.Revision = *a[1].(*uint64)
	v.CreationTime = *a[2].(*time.Time)
	v.UpdateTime = *a[3].(*time.Time)
	v.UserID = *a[4].(*string)
	v.Name = *a[5].(*string)
	v.Value = *a[6].(*string)

	if x, ok := vi.(sqlg.PreJSONSetupper); ok {
		if err := x.PreJSON(); err != nil {
			return nil, "", errors.Wrap(err, "prejson error")
		}
	}

	v.TxID = txID

	return v, v.ID, nil
}

func (d *DB) fetchLinkedAccounts(tx *sql.Tx, q sq.Builder) ([]*types.LinkedAccount, []string, error) {
	rows, err := d.query(tx, q)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	defer rows.Close()

	return d.scanLinkedAccounts(rows, tx.ID(), 0)
}

func (d *DB) fetchLinkedAccountsSkipLastFields(tx *sql.Tx, q sq.Builder, skipFieldsCount uint) ([]*types.LinkedAccount, []string, error) {
	rows, err := d.query(tx, q)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	defer rows.Close()

	return d.scanLinkedAccounts(rows, tx.ID(), skipFieldsCount)
}

func (d *DB) scanLinkedAccount(rows *stdsql.Rows, skipFieldsCount uint) (*types.LinkedAccount, string, error) {

	v := &types.LinkedAccount{}

	var vi any = v
	if x, ok := vi.(sqlg.Initer); ok {
		x.Init()
	}

	fields := append([]any{&v.ID, &v.Revision, &v.CreationTime, &v.UpdateTime, &v.UserID, &v.RemoteUserID, &v.RemoteUserName, &v.RemoteUserAvatarURL, &v.RemoteSourceID, &v.UserAccessToken, &v.Oauth2AccessToken, &v.Oauth2RefreshToken, &v.Oauth2AccessTokenExpiresAt})

	for i := uint(0); i < skipFieldsCount; i++ {
		fields = append(fields, new(any))
	}

	if err := rows.Scan(fields...); err != nil {
		return nil, "", errors.Wrap(err, "failed to scan row")
	}

	if x, ok := vi.(sqlg.PreJSONSetupper); ok {
		if err := x.PreJSON(); err != nil {
			return nil, "", errors.Wrap(err, "prejson error")
		}
	}

	return v, v.ID, nil
}

func (d *DB) scanLinkedAccounts(rows *stdsql.Rows, txID string, skipFieldsCount uint) ([]*types.LinkedAccount, []string, error) {
	vs := []*types.LinkedAccount{}
	ids := []string{}
	for rows.Next() {
		v, id, err := d.scanLinkedAccount(rows, skipFieldsCount)
		if err != nil {
			rows.Close()
			return nil, nil, errors.WithStack(err)
		}
		v.TxID = txID
		vs = append(vs, v)
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, errors.WithStack(err)
	}
	return vs, ids, nil
}

func (d *DB) LinkedAccountArray() []any {
	a := []any{}
	a = append(a, new(string))
	a = append(a, new(uint64))
	a = append(a, new(time.Time))
	a = append(a, new(time.Time))
	a = append(a, new(string))
	a = append(a, new(string))
	a = append(a, new(string))
	a = append(a, new(string))
	a = append(a, new(string))
	a = append(a, new(string))
	a = append(a, new(string))
	a = append(a, new(string))
	a = append(a, new(time.Time))

	return a
}

func (d *DB) LinkedAccountFromArray(a []any, txID string) (*types.LinkedAccount, string, error) {
	v := &types.LinkedAccount{}

	var vi any = v
	if x, ok := vi.(sqlg.Initer); ok {
		x.Init()
	}
	v.ID = *a[0].(*string)
	v.Revision = *a[1].(*uint64)
	v.CreationTime = *a[2].(*time.Time)
	v.UpdateTime = *a[3].(*time.Time)
	v.UserID = *a[4].(*string)
	v.RemoteUserID = *a[5].(*string)
	v.RemoteUserName = *a[6].(*string)
	v.RemoteUserAvatarURL = *a[7].(*string)
	v.RemoteSourceID = *a[8].(*string)
	v.UserAccessToken = *a[9].(*string)
	v.Oauth2AccessToken = *a[10].(*string)
	v.Oauth2RefreshToken = *a[11].(*string)
	v.Oauth2AccessTokenExpiresAt = *a[12].(*time.Time)

	if x, ok := vi.(sqlg.PreJSONSetupper); ok {
		if err := x.PreJSON(); err != nil {
			return nil, "", errors.Wrap(err, "prejson error")
		}
	}

	v.TxID = txID

	return v, v.ID, nil
}

func (d *DB) fetchOrganizations(tx *sql.Tx, q sq.Builder) ([]*types.Organization, []string, error) {
	rows, err := d.query(tx, q)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	defer rows.Close()

	return d.scanOrganizations(rows, tx.ID(), 0)
}

func (d *DB) fetchOrganizationsSkipLastFields(tx *sql.Tx, q sq.Builder, skipFieldsCount uint) ([]*types.Organization, []string, error) {
	rows, err := d.query(tx, q)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	defer rows.Close()

	return d.scanOrganizations(rows, tx.ID(), skipFieldsCount)
}

func (d *DB) scanOrganization(rows *stdsql.Rows, skipFieldsCount uint) (*types.Organization, string, error) {

	v := &types.Organization{}

	var vi any = v
	if x, ok := vi.(sqlg.Initer); ok {
		x.Init()
	}

	fields := append([]any{&v.ID, &v.Revision, &v.CreationTime, &v.UpdateTime, &v.Name, &v.Visibility, &v.CreatorUserID})

	for i := uint(0); i < skipFieldsCount; i++ {
		fields = append(fields, new(any))
	}

	if err := rows.Scan(fields...); err != nil {
		return nil, "", errors.Wrap(err, "failed to scan row")
	}

	if x, ok := vi.(sqlg.PreJSONSetupper); ok {
		if err := x.PreJSON(); err != nil {
			return nil, "", errors.Wrap(err, "prejson error")
		}
	}

	return v, v.ID, nil
}

func (d *DB) scanOrganizations(rows *stdsql.Rows, txID string, skipFieldsCount uint) ([]*types.Organization, []string, error) {
	vs := []*types.Organization{}
	ids := []string{}
	for rows.Next() {
		v, id, err := d.scanOrganization(rows, skipFieldsCount)
		if err != nil {
			rows.Close()
			return nil, nil, errors.WithStack(err)
		}
		v.TxID = txID
		vs = append(vs, v)
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, errors.WithStack(err)
	}
	return vs, ids, nil
}

func (d *DB) OrganizationArray() []any {
	a := []any{}
	a = append(a, new(string))
	a = append(a, new(uint64))
	a = append(a, new(time.Time))
	a = append(a, new(time.Time))
	a = append(a, new(string))
	a = append(a, new(types.Visibility))
	a = append(a, new(string))

	return a
}

func (d *DB) OrganizationFromArray(a []any, txID string) (*types.Organization, string, error) {
	v := &types.Organization{}

	var vi any = v
	if x, ok := vi.(sqlg.Initer); ok {
		x.Init()
	}
	v.ID = *a[0].(*string)
	v.Revision = *a[1].(*uint64)
	v.CreationTime = *a[2].(*time.Time)
	v.UpdateTime = *a[3].(*time.Time)
	v.Name = *a[4].(*string)
	v.Visibility = *a[5].(*types.Visibility)
	v.CreatorUserID = *a[6].(*string)

	if x, ok := vi.(sqlg.PreJSONSetupper); ok {
		if err := x.PreJSON(); err != nil {
			return nil, "", errors.Wrap(err, "prejson error")
		}
	}

	v.TxID = txID

	return v, v.ID, nil
}

func (d *DB) fetchOrganizationMembers(tx *sql.Tx, q sq.Builder) ([]*types.OrganizationMember, []string, error) {
	rows, err := d.query(tx, q)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	defer rows.Close()

	return d.scanOrganizationMembers(rows, tx.ID(), 0)
}

func (d *DB) fetchOrganizationMembersSkipLastFields(tx *sql.Tx, q sq.Builder, skipFieldsCount uint) ([]*types.OrganizationMember, []string, error) {
	rows, err := d.query(tx, q)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	defer rows.Close()

	return d.scanOrganizationMembers(rows, tx.ID(), skipFieldsCount)
}

func (d *DB) scanOrganizationMember(rows *stdsql.Rows, skipFieldsCount uint) (*types.OrganizationMember, string, error) {

	v := &types.OrganizationMember{}

	var vi any = v
	if x, ok := vi.(sqlg.Initer); ok {
		x.Init()
	}

	fields := append([]any{&v.ID, &v.Revision, &v.CreationTime, &v.UpdateTime, &v.OrganizationID, &v.UserID, &v.MemberRole})

	for i := uint(0); i < skipFieldsCount; i++ {
		fields = append(fields, new(any))
	}

	if err := rows.Scan(fields...); err != nil {
		return nil, "", errors.Wrap(err, "failed to scan row")
	}

	if x, ok := vi.(sqlg.PreJSONSetupper); ok {
		if err := x.PreJSON(); err != nil {
			return nil, "", errors.Wrap(err, "prejson error")
		}
	}

	return v, v.ID, nil
}

func (d *DB) scanOrganizationMembers(rows *stdsql.Rows, txID string, skipFieldsCount uint) ([]*types.OrganizationMember, []string, error) {
	vs := []*types.OrganizationMember{}
	ids := []string{}
	for rows.Next() {
		v, id, err := d.scanOrganizationMember(rows, skipFieldsCount)
		if err != nil {
			rows.Close()
			return nil, nil, errors.WithStack(err)
		}
		v.TxID = txID
		vs = append(vs, v)
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, errors.WithStack(err)
	}
	return vs, ids, nil
}

func (d *DB) OrganizationMemberArray() []any {
	a := []any{}
	a = append(a, new(string))
	a = append(a, new(uint64))
	a = append(a, new(time.Time))
	a = append(a, new(time.Time))
	a = append(a, new(string))
	a = append(a, new(string))
	a = append(a, new(types.MemberRole))

	return a
}

func (d *DB) OrganizationMemberFromArray(a []any, txID string) (*types.OrganizationMember, string, error) {
	v := &types.OrganizationMember{}

	var vi any = v
	if x, ok := vi.(sqlg.Initer); ok {
		x.Init()
	}
	v.ID = *a[0].(*string)
	v.Revision = *a[1].(*uint64)
	v.CreationTime = *a[2].(*time.Time)
	v.UpdateTime = *a[3].(*time.Time)
	v.OrganizationID = *a[4].(*string)
	v.UserID = *a[5].(*string)
	v.MemberRole = *a[6].(*types.MemberRole)

	if x, ok := vi.(sqlg.PreJSONSetupper); ok {
		if err := x.PreJSON(); err != nil {
			return nil, "", errors.Wrap(err, "prejson error")
		}
	}

	v.TxID = txID

	return v, v.ID, nil
}

func (d *DB) fetchProjectGroups(tx *sql.Tx, q sq.Builder) ([]*types.ProjectGroup, []string, error) {
	rows, err := d.query(tx, q)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	defer rows.Close()

	return d.scanProjectGroups(rows, tx.ID(), 0)
}

func (d *DB) fetchProjectGroupsSkipLastFields(tx *sql.Tx, q sq.Builder, skipFieldsCount uint) ([]*types.ProjectGroup, []string, error) {
	rows, err := d.query(tx, q)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	defer rows.Close()

	return d.scanProjectGroups(rows, tx.ID(), skipFieldsCount)
}

func (d *DB) scanProjectGroup(rows *stdsql.Rows, skipFieldsCount uint) (*types.ProjectGroup, string, error) {

	v := &types.ProjectGroup{}

	var vi any = v
	if x, ok := vi.(sqlg.Initer); ok {
		x.Init()
	}

	fields := append([]any{&v.ID, &v.Revision, &v.CreationTime, &v.UpdateTime, &v.Name, &v.Parent.Kind, &v.Parent.ID, &v.Visibility})

	for i := uint(0); i < skipFieldsCount; i++ {
		fields = append(fields, new(any))
	}

	if err := rows.Scan(fields...); err != nil {
		return nil, "", errors.Wrap(err, "failed to scan row")
	}

	if x, ok := vi.(sqlg.PreJSONSetupper); ok {
		if err := x.PreJSON(); err != nil {
			return nil, "", errors.Wrap(err, "prejson error")
		}
	}

	return v, v.ID, nil
}

func (d *DB) scanProjectGroups(rows *stdsql.Rows, txID string, skipFieldsCount uint) ([]*types.ProjectGroup, []string, error) {
	vs := []*types.ProjectGroup{}
	ids := []string{}
	for rows.Next() {
		v, id, err := d.scanProjectGroup(rows, skipFieldsCount)
		if err != nil {
			rows.Close()
			return nil, nil, errors.WithStack(err)
		}
		v.TxID = txID
		vs = append(vs, v)
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, errors.WithStack(err)
	}
	return vs, ids, nil
}

func (d *DB) ProjectGroupArray() []any {
	a := []any{}
	a = append(a, new(string))
	a = append(a, new(uint64))
	a = append(a, new(time.Time))
	a = append(a, new(time.Time))
	a = append(a, new(string))
	a = append(a, new(types.ObjectKind))
	a = append(a, new(string))
	a = append(a, new(types.Visibility))

	return a
}

func (d *DB) ProjectGroupFromArray(a []any, txID string) (*types.ProjectGroup, string, error) {
	v := &types.ProjectGroup{}

	var vi any = v
	if x, ok := vi.(sqlg.Initer); ok {
		x.Init()
	}
	v.ID = *a[0].(*string)
	v.Revision = *a[1].(*uint64)
	v.CreationTime = *a[2].(*time.Time)
	v.UpdateTime = *a[3].(*time.Time)
	v.Name = *a[4].(*string)
	v.Parent.Kind = *a[5].(*types.ObjectKind)
	v.Parent.ID = *a[6].(*string)
	v.Visibility = *a[7].(*types.Visibility)

	if x, ok := vi.(sqlg.PreJSONSetupper); ok {
		if err := x.PreJSON(); err != nil {
			return nil, "", errors.Wrap(err, "prejson error")
		}
	}

	v.TxID = txID

	return v, v.ID, nil
}

func (d *DB) fetchProjects(tx *sql.Tx, q sq.Builder) ([]*types.Project, []string, error) {
	rows, err := d.query(tx, q)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	defer rows.Close()

	return d.scanProjects(rows, tx.ID(), 0)
}

func (d *DB) fetchProjectsSkipLastFields(tx *sql.Tx, q sq.Builder, skipFieldsCount uint) ([]*types.Project, []string, error) {
	rows, err := d.query(tx, q)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	defer rows.Close()

	return d.scanProjects(rows, tx.ID(), skipFieldsCount)
}

func (d *DB) scanProject(rows *stdsql.Rows, skipFieldsCount uint) (*types.Project, string, error) {

	v := &types.Project{}

	var vi any = v
	if x, ok := vi.(sqlg.Initer); ok {
		x.Init()
	}

	fields := append([]any{&v.ID, &v.Revision, &v.CreationTime, &v.UpdateTime, &v.Name, &v.Parent.Kind, &v.Parent.ID, &v.Secret, &v.Visibility, &v.RemoteRepositoryConfigType, &v.RemoteSourceID, &v.LinkedAccountID, &v.RepositoryID, &v.RepositoryPath, &v.SSHPrivateKey, &v.SkipSSHHostKeyCheck, &v.WebhookSecret, &v.PassVarsToForkedPR, &v.DefaultBranch, &v.MembersCanPerformRunActions})

	for i := uint(0); i < skipFieldsCount; i++ {
		fields = append(fields, new(any))
	}

	if err := rows.Scan(fields...); err != nil {
		return nil, "", errors.Wrap(err, "failed to scan row")
	}

	if x, ok := vi.(sqlg.PreJSONSetupper); ok {
		if err := x.PreJSON(); err != nil {
			return nil, "", errors.Wrap(err, "prejson error")
		}
	}

	return v, v.ID, nil
}

func (d *DB) scanProjects(rows *stdsql.Rows, txID string, skipFieldsCount uint) ([]*types.Project, []string, error) {
	vs := []*types.Project{}
	ids := []string{}
	for rows.Next() {
		v, id, err := d.scanProject(rows, skipFieldsCount)
		if err != nil {
			rows.Close()
			return nil, nil, errors.WithStack(err)
		}
		v.TxID = txID
		vs = append(vs, v)
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, errors.WithStack(err)
	}
	return vs, ids, nil
}

func (d *DB) ProjectArray() []any {
	a := []any{}
	a = append(a, new(string))
	a = append(a, new(uint64))
	a = append(a, new(time.Time))
	a = append(a, new(time.Time))
	a = append(a, new(string))
	a = append(a, new(types.ObjectKind))
	a = append(a, new(string))
	a = append(a, new(string))
	a = append(a, new(types.Visibility))
	a = append(a, new(types.RemoteRepositoryConfigType))
	a = append(a, new(string))
	a = append(a, new(string))
	a = append(a, new(string))
	a = append(a, new(string))
	a = append(a, new(string))
	a = append(a, new(bool))
	a = append(a, new(string))
	a = append(a, new(bool))
	a = append(a, new(string))
	a = append(a, new(bool))

	return a
}

func (d *DB) ProjectFromArray(a []any, txID string) (*types.Project, string, error) {
	v := &types.Project{}

	var vi any = v
	if x, ok := vi.(sqlg.Initer); ok {
		x.Init()
	}
	v.ID = *a[0].(*string)
	v.Revision = *a[1].(*uint64)
	v.CreationTime = *a[2].(*time.Time)
	v.UpdateTime = *a[3].(*time.Time)
	v.Name = *a[4].(*string)
	v.Parent.Kind = *a[5].(*types.ObjectKind)
	v.Parent.ID = *a[6].(*string)
	v.Secret = *a[7].(*string)
	v.Visibility = *a[8].(*types.Visibility)
	v.RemoteRepositoryConfigType = *a[9].(*types.RemoteRepositoryConfigType)
	v.RemoteSourceID = *a[10].(*string)
	v.LinkedAccountID = *a[11].(*string)
	v.RepositoryID = *a[12].(*string)
	v.RepositoryPath = *a[13].(*string)
	v.SSHPrivateKey = *a[14].(*string)
	v.SkipSSHHostKeyCheck = *a[15].(*bool)
	v.WebhookSecret = *a[16].(*string)
	v.PassVarsToForkedPR = *a[17].(*bool)
	v.DefaultBranch = *a[18].(*string)
	v.MembersCanPerformRunActions = *a[19].(*bool)

	if x, ok := vi.(sqlg.PreJSONSetupper); ok {
		if err := x.PreJSON(); err != nil {
			return nil, "", errors.Wrap(err, "prejson error")
		}
	}

	v.TxID = txID

	return v, v.ID, nil
}

func (d *DB) fetchSecrets(tx *sql.Tx, q sq.Builder) ([]*types.Secret, []string, error) {
	rows, err := d.query(tx, q)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	defer rows.Close()

	return d.scanSecrets(rows, tx.ID(), 0)
}

func (d *DB) fetchSecretsSkipLastFields(tx *sql.Tx, q sq.Builder, skipFieldsCount uint) ([]*types.Secret, []string, error) {
	rows, err := d.query(tx, q)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	defer rows.Close()

	return d.scanSecrets(rows, tx.ID(), skipFieldsCount)
}

func (d *DB) scanSecret(rows *stdsql.Rows, skipFieldsCount uint) (*types.Secret, string, error) {
	var inDataJSON []byte

	v := &types.Secret{}

	var vi any = v
	if x, ok := vi.(sqlg.Initer); ok {
		x.Init()
	}

	fields := append([]any{&v.ID, &v.Revision, &v.CreationTime, &v.UpdateTime, &v.Name, &v.Parent.Kind, &v.Parent.ID, &v.Type, &inDataJSON, &v.SecretProviderID, &v.Path})

	for i := uint(0); i < skipFieldsCount; i++ {
		fields = append(fields, new(any))
	}

	if err := rows.Scan(fields...); err != nil {
		return nil, "", errors.Wrap(err, "failed to scan row")
	}

	if x, ok := vi.(sqlg.PreJSONSetupper); ok {
		if err := x.PreJSON(); err != nil {
			return nil, "", errors.Wrap(err, "prejson error")
		}
	}
	if err := json.Unmarshal(inDataJSON, &v.Data); err != nil {
		return nil, "", errors.Wrap(err, "failed to unmarshal v.Data")
	}

	return v, v.ID, nil
}

func (d *DB) scanSecrets(rows *stdsql.Rows, txID string, skipFieldsCount uint) ([]*types.Secret, []string, error) {
	vs := []*types.Secret{}
	ids := []string{}
	for rows.Next() {
		v, id, err := d.scanSecret(rows, skipFieldsCount)
		if err != nil {
			rows.Close()
			return nil, nil, errors.WithStack(err)
		}
		v.TxID = txID
		vs = append(vs, v)
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, errors.WithStack(err)
	}
	return vs, ids, nil
}

func (d *DB) SecretArray() []any {
	a := []any{}
	a = append(a, new(string))
	a = append(a, new(uint64))
	a = append(a, new(time.Time))
	a = append(a, new(time.Time))
	a = append(a, new(string))
	a = append(a, new(types.ObjectKind))
	a = append(a, new(string))
	a = append(a, new(types.SecretType))
	a = append(a, new([]byte))
	a = append(a, new(string))
	a = append(a, new(string))

	return a
}

func (d *DB) SecretFromArray(a []any, txID string) (*types.Secret, string, error) {
	v := &types.Secret{}

	var vi any = v
	if x, ok := vi.(sqlg.Initer); ok {
		x.Init()
	}
	v.ID = *a[0].(*string)
	v.Revision = *a[1].(*uint64)
	v.CreationTime = *a[2].(*time.Time)
	v.UpdateTime = *a[3].(*time.Time)
	v.Name = *a[4].(*string)
	v.Parent.Kind = *a[5].(*types.ObjectKind)
	v.Parent.ID = *a[6].(*string)
	v.Type = *a[7].(*types.SecretType)
	v.SecretProviderID = *a[9].(*string)
	v.Path = *a[10].(*string)

	if x, ok := vi.(sqlg.PreJSONSetupper); ok {
		if err := x.PreJSON(); err != nil {
			return nil, "", errors.Wrap(err, "prejson error")
		}
	}
	if err := json.Unmarshal(a[8].([]byte), &v.Data); err != nil {
		return nil, "", errors.Wrap(err, "failed to unmarshal v.v.Data")
	}

	v.TxID = txID

	return v, v.ID, nil
}

func (d *DB) fetchVariables(tx *sql.Tx, q sq.Builder) ([]*types.Variable, []string, error) {
	rows, err := d.query(tx, q)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	defer rows.Close()

	return d.scanVariables(rows, tx.ID(), 0)
}

func (d *DB) fetchVariablesSkipLastFields(tx *sql.Tx, q sq.Builder, skipFieldsCount uint) ([]*types.Variable, []string, error) {
	rows, err := d.query(tx, q)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	defer rows.Close()

	return d.scanVariables(rows, tx.ID(), skipFieldsCount)
}

func (d *DB) scanVariable(rows *stdsql.Rows, skipFieldsCount uint) (*types.Variable, string, error) {
	var inValuesJSON []byte

	v := &types.Variable{}

	var vi any = v
	if x, ok := vi.(sqlg.Initer); ok {
		x.Init()
	}

	fields := append([]any{&v.ID, &v.Revision, &v.CreationTime, &v.UpdateTime, &v.Name, &v.Parent.Kind, &v.Parent.ID, &inValuesJSON})

	for i := uint(0); i < skipFieldsCount; i++ {
		fields = append(fields, new(any))
	}

	if err := rows.Scan(fields...); err != nil {
		return nil, "", errors.Wrap(err, "failed to scan row")
	}

	if x, ok := vi.(sqlg.PreJSONSetupper); ok {
		if err := x.PreJSON(); err != nil {
			return nil, "", errors.Wrap(err, "prejson error")
		}
	}
	if err := json.Unmarshal(inValuesJSON, &v.Values); err != nil {
		return nil, "", errors.Wrap(err, "failed to unmarshal v.Values")
	}

	return v, v.ID, nil
}

func (d *DB) scanVariables(rows *stdsql.Rows, txID string, skipFieldsCount uint) ([]*types.Variable, []string, error) {
	vs := []*types.Variable{}
	ids := []string{}
	for rows.Next() {
		v, id, err := d.scanVariable(rows, skipFieldsCount)
		if err != nil {
			rows.Close()
			return nil, nil, errors.WithStack(err)
		}
		v.TxID = txID
		vs = append(vs, v)
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, errors.WithStack(err)
	}
	return vs, ids, nil
}

func (d *DB) VariableArray() []any {
	a := []any{}
	a = append(a, new(string))
	a = append(a, new(uint64))
	a = append(a, new(time.Time))
	a = append(a, new(time.Time))
	a = append(a, new(string))
	a = append(a, new(types.ObjectKind))
	a = append(a, new(string))
	a = append(a, new([]byte))

	return a
}

func (d *DB) VariableFromArray(a []any, txID string) (*types.Variable, string, error) {
	v := &types.Variable{}

	var vi any = v
	if x, ok := vi.(sqlg.Initer); ok {
		x.Init()
	}
	v.ID = *a[0].(*string)
	v.Revision = *a[1].(*uint64)
	v.CreationTime = *a[2].(*time.Time)
	v.UpdateTime = *a[3].(*time.Time)
	v.Name = *a[4].(*string)
	v.Parent.Kind = *a[5].(*types.ObjectKind)
	v.Parent.ID = *a[6].(*string)

	if x, ok := vi.(sqlg.PreJSONSetupper); ok {
		if err := x.PreJSON(); err != nil {
			return nil, "", errors.Wrap(err, "prejson error")
		}
	}
	if err := json.Unmarshal(a[7].([]byte), &v.Values); err != nil {
		return nil, "", errors.Wrap(err, "failed to unmarshal v.v.Values")
	}

	v.TxID = txID

	return v, v.ID, nil
}

func (d *DB) fetchOrgInvitations(tx *sql.Tx, q sq.Builder) ([]*types.OrgInvitation, []string, error) {
	rows, err := d.query(tx, q)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	defer rows.Close()

	return d.scanOrgInvitations(rows, tx.ID(), 0)
}

func (d *DB) fetchOrgInvitationsSkipLastFields(tx *sql.Tx, q sq.Builder, skipFieldsCount uint) ([]*types.OrgInvitation, []string, error) {
	rows, err := d.query(tx, q)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	defer rows.Close()

	return d.scanOrgInvitations(rows, tx.ID(), skipFieldsCount)
}

func (d *DB) scanOrgInvitation(rows *stdsql.Rows, skipFieldsCount uint) (*types.OrgInvitation, string, error) {

	v := &types.OrgInvitation{}

	var vi any = v
	if x, ok := vi.(sqlg.Initer); ok {
		x.Init()
	}

	fields := append([]any{&v.ID, &v.Revision, &v.CreationTime, &v.UpdateTime, &v.UserID, &v.OrganizationID, &v.Role})

	for i := uint(0); i < skipFieldsCount; i++ {
		fields = append(fields, new(any))
	}

	if err := rows.Scan(fields...); err != nil {
		return nil, "", errors.Wrap(err, "failed to scan row")
	}

	if x, ok := vi.(sqlg.PreJSONSetupper); ok {
		if err := x.PreJSON(); err != nil {
			return nil, "", errors.Wrap(err, "prejson error")
		}
	}

	return v, v.ID, nil
}

func (d *DB) scanOrgInvitations(rows *stdsql.Rows, txID string, skipFieldsCount uint) ([]*types.OrgInvitation, []string, error) {
	vs := []*types.OrgInvitation{}
	ids := []string{}
	for rows.Next() {
		v, id, err := d.scanOrgInvitation(rows, skipFieldsCount)
		if err != nil {
			rows.Close()
			return nil, nil, errors.WithStack(err)
		}
		v.TxID = txID
		vs = append(vs, v)
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, errors.WithStack(err)
	}
	return vs, ids, nil
}

func (d *DB) OrgInvitationArray() []any {
	a := []any{}
	a = append(a, new(string))
	a = append(a, new(uint64))
	a = append(a, new(time.Time))
	a = append(a, new(time.Time))
	a = append(a, new(string))
	a = append(a, new(string))
	a = append(a, new(types.MemberRole))

	return a
}

func (d *DB) OrgInvitationFromArray(a []any, txID string) (*types.OrgInvitation, string, error) {
	v := &types.OrgInvitation{}

	var vi any = v
	if x, ok := vi.(sqlg.Initer); ok {
		x.Init()
	}
	v.ID = *a[0].(*string)
	v.Revision = *a[1].(*uint64)
	v.CreationTime = *a[2].(*time.Time)
	v.UpdateTime = *a[3].(*time.Time)
	v.UserID = *a[4].(*string)
	v.OrganizationID = *a[5].(*string)
	v.Role = *a[6].(*types.MemberRole)

	if x, ok := vi.(sqlg.PreJSONSetupper); ok {
		if err := x.PreJSON(); err != nil {
			return nil, "", errors.Wrap(err, "prejson error")
		}
	}

	v.TxID = txID

	return v, v.ID, nil
}
