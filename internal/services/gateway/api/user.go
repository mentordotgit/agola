// Copyright 2019 Sorint.lab
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied
// See the License for the specific language governing permissions and
// limitations under the License.

package api

import (
	"context"
	"encoding/json"
	"net/http"
	"sort"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/sorintlab/errors"

	"agola.io/agola/internal/services/gateway/action"
	"agola.io/agola/internal/services/gateway/common"
	"agola.io/agola/internal/util"
	csapitypes "agola.io/agola/services/configstore/api/types"
	cstypes "agola.io/agola/services/configstore/types"
	gwapitypes "agola.io/agola/services/gateway/api/types"
)

type CreateUserHandler struct {
	log zerolog.Logger
	ah  *action.ActionHandler
}

func NewCreateUserHandler(log zerolog.Logger, ah *action.ActionHandler) *CreateUserHandler {
	return &CreateUserHandler{log: log, ah: ah}
}

func (h *CreateUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req gwapitypes.CreateUserRequest
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&req); err != nil {
		util.HTTPError(w, util.NewAPIError(util.ErrBadRequest, err))
		return
	}

	creq := &action.CreateUserRequest{
		UserName: req.UserName,
	}

	u, err := h.ah.CreateUser(ctx, creq)
	if util.HTTPError(w, err) {
		h.log.Err(err).Send()
		return
	}

	res := createUserResponse(u)
	if err := util.HTTPResponse(w, http.StatusCreated, res); err != nil {
		h.log.Err(err).Send()
	}
}

type DeleteUserHandler struct {
	log zerolog.Logger
	ah  *action.ActionHandler
}

func NewDeleteUserHandler(log zerolog.Logger, ah *action.ActionHandler) *DeleteUserHandler {
	return &DeleteUserHandler{log: log, ah: ah}
}

func (h *DeleteUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	userRef := vars["userref"]

	err := h.ah.DeleteUser(ctx, userRef)
	if util.HTTPError(w, err) {
		h.log.Err(err).Send()
		return
	}

	if err := util.HTTPResponse(w, http.StatusNoContent, nil); err != nil {
		h.log.Err(err).Send()
	}
}

type CurrentUserHandler struct {
	log zerolog.Logger
	ah  *action.ActionHandler
}

func NewCurrentUserHandler(log zerolog.Logger, ah *action.ActionHandler) *CurrentUserHandler {
	return &CurrentUserHandler{log: log, ah: ah}
}

func (h *CurrentUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := common.CurrentUserID(ctx)
	if userID == "" {
		util.HTTPError(w, util.NewAPIError(util.ErrBadRequest, errors.Errorf("user not authenticated")))
		return
	}

	user, err := h.ah.GetCurrentUser(ctx, userID)
	if util.HTTPError(w, err) {
		h.log.Err(err).Send()
		return
	}

	res := createPrivateUserResponse(user.User, user.Tokens, user.LinkedAccounts)
	if err := util.HTTPResponse(w, http.StatusOK, res); err != nil {
		h.log.Err(err).Send()
	}
}

type UserHandler struct {
	log zerolog.Logger
	ah  *action.ActionHandler
}

func NewUserHandler(log zerolog.Logger, ah *action.ActionHandler) *UserHandler {
	return &UserHandler{log: log, ah: ah}
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	userRef := vars["userref"]

	user, err := h.ah.GetUser(ctx, userRef)
	if util.HTTPError(w, err) {
		h.log.Err(err).Send()
		return
	}

	res := createUserResponse(user)
	if err := util.HTTPResponse(w, http.StatusOK, res); err != nil {
		h.log.Err(err).Send()
	}
}

func createPrivateUserResponse(u *cstypes.User, tokens []*cstypes.UserToken, linkedAccounts []*cstypes.LinkedAccount) *gwapitypes.PrivateUserResponse {
	user := &gwapitypes.PrivateUserResponse{
		ID:             u.ID,
		UserName:       u.Name,
		Tokens:         make([]string, 0, len(tokens)),
		LinkedAccounts: make([]*gwapitypes.LinkedAccountResponse, 0, len(linkedAccounts)),
	}
	for _, token := range tokens {
		user.Tokens = append(user.Tokens, token.Name)
	}
	sort.Strings(user.Tokens)

	for _, la := range linkedAccounts {
		user.LinkedAccounts = append(user.LinkedAccounts, &gwapitypes.LinkedAccountResponse{
			ID:                  la.ID,
			RemoteSourceID:      la.RemoteSourceID,
			RemoteUserName:      la.RemoteUserName,
			RemoteUserAvatarURL: la.RemoteUserAvatarURL,
		})
	}

	return user
}

func createUserResponse(u *cstypes.User) *gwapitypes.UserResponse {
	user := &gwapitypes.UserResponse{
		ID:       u.ID,
		UserName: u.Name,
	}

	return user
}

type UsersHandler struct {
	log zerolog.Logger
	ah  *action.ActionHandler
}

func NewUsersHandler(log zerolog.Logger, ah *action.ActionHandler) *UsersHandler {
	return &UsersHandler{log: log, ah: ah}
}

func (h *UsersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res, err := h.do(w, r)
	if util.HTTPError(w, err) {
		h.log.Err(err).Send()
		return
	}

	if err := util.HTTPResponse(w, http.StatusOK, res); err != nil {
		h.log.Err(err).Send()
	}
}

func (h *UsersHandler) do(w http.ResponseWriter, r *http.Request) ([]*gwapitypes.PrivateUserResponse, error) {
	ctx := r.Context()

	query := r.URL.Query()

	ropts, err := parseRequestOptions(r)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	queryType := query.Get("query_type")

	var ausers []*action.PrivateUserResponse
	switch queryType {
	case "byremoteuser":
		remoteUserID := query.Get("remoteuserid")
		rsRef := query.Get("remotesourceref")

		user, err := h.ah.GetUserByLinkedAccountRemoteUserAndSource(ctx, remoteUserID, rsRef)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		ausers = []*action.PrivateUserResponse{user}
	case "":
		areq := &action.GetUsersRequest{Cursor: ropts.Cursor, Limit: ropts.Limit, SortDirection: action.SortDirection(ropts.SortDirection)}
		ares, err := h.ah.GetUsers(ctx, areq)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		ausers = ares.Users
		addCursorHeader(w, ares.Cursor)
	default:
		return nil, util.NewAPIError(util.ErrBadRequest, errors.Errorf("unknown query_type: %q", queryType))
	}

	users := make([]*gwapitypes.PrivateUserResponse, len(ausers))
	for i, p := range ausers {
		users[i] = createPrivateUserResponse(p.User, p.Tokens, p.LinkedAccounts)
	}

	return users, nil
}

type CreateUserLAHandler struct {
	log zerolog.Logger
	ah  *action.ActionHandler
}

func NewCreateUserLAHandler(log zerolog.Logger, ah *action.ActionHandler) *CreateUserLAHandler {
	return &CreateUserLAHandler{log: log, ah: ah}
}

func (h *CreateUserLAHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	userRef := vars["userref"]

	var req *gwapitypes.CreateUserLARequest
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&req); err != nil {
		util.HTTPError(w, util.NewAPIError(util.ErrBadRequest, err))
		return
	}

	res, err := h.createUserLA(ctx, userRef, req)
	if util.HTTPError(w, err) {
		h.log.Err(err).Send()
		return
	}

	if err := util.HTTPResponse(w, http.StatusCreated, res); err != nil {
		h.log.Err(err).Send()
	}
}

func (h *CreateUserLAHandler) createUserLA(ctx context.Context, userRef string, req *gwapitypes.CreateUserLARequest) (*gwapitypes.CreateUserLAResponse, error) {
	creq := &action.CreateUserLARequest{
		UserRef:          userRef,
		RemoteSourceName: req.RemoteSourceName,
	}

	h.log.Info().Msgf("creating linked account")
	cresp, err := h.ah.HandleRemoteSourceAuth(ctx, req.RemoteSourceName, req.RemoteSourceLoginName, req.RemoteSourceLoginPassword, action.RemoteSourceRequestTypeCreateUserLA, creq)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if cresp.Oauth2Redirect != "" {
		return &gwapitypes.CreateUserLAResponse{
			Oauth2Redirect: cresp.Oauth2Redirect,
		}, nil
	}
	authresp := cresp.Response.(*action.CreateUserLAResponse)

	resp := &gwapitypes.CreateUserLAResponse{
		LinkedAccount: &gwapitypes.LinkedAccount{
			ID:                  authresp.LinkedAccount.ID,
			RemoteUserID:        authresp.LinkedAccount.RemoteUserID,
			RemoteUserName:      authresp.LinkedAccount.RemoteUserName,
			RemoteUserAvatarURL: authresp.LinkedAccount.RemoteUserAvatarURL,
			RemoteSourceID:      authresp.LinkedAccount.RemoteUserID,
		},
	}
	h.log.Info().Msgf("linked account %q for user %q created", resp.LinkedAccount.ID, userRef)
	return resp, nil
}

type DeleteUserLAHandler struct {
	log zerolog.Logger
	ah  *action.ActionHandler
}

func NewDeleteUserLAHandler(log zerolog.Logger, ah *action.ActionHandler) *DeleteUserLAHandler {
	return &DeleteUserLAHandler{log: log, ah: ah}
}

func (h *DeleteUserLAHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	userRef := vars["userref"]
	laID := vars["laid"]

	err := h.ah.DeleteUserLA(ctx, userRef, laID)
	if util.HTTPError(w, err) {
		h.log.Err(err).Send()
		return
	}

	if err := util.HTTPResponse(w, http.StatusNoContent, nil); err != nil {
		h.log.Err(err).Send()
	}
}

type CreateUserTokenHandler struct {
	log zerolog.Logger
	ah  *action.ActionHandler
}

func NewCreateUserTokenHandler(log zerolog.Logger, ah *action.ActionHandler) *CreateUserTokenHandler {
	return &CreateUserTokenHandler{log: log, ah: ah}
}

func (h *CreateUserTokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	userRef := vars["userref"]

	var req gwapitypes.CreateUserTokenRequest
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&req); err != nil {
		util.HTTPError(w, util.NewAPIError(util.ErrBadRequest, err))
		return
	}

	creq := &action.CreateUserTokenRequest{
		UserRef:   userRef,
		TokenName: req.TokenName,
	}
	h.log.Info().Msgf("creating user %q token", userRef)
	token, err := h.ah.CreateUserToken(ctx, creq)
	if util.HTTPError(w, err) {
		h.log.Err(err).Send()
		return
	}

	res := &gwapitypes.CreateUserTokenResponse{
		Token: token,
	}

	if err := util.HTTPResponse(w, http.StatusCreated, res); err != nil {
		h.log.Err(err).Send()
	}
}

type DeleteUserTokenHandler struct {
	log zerolog.Logger
	ah  *action.ActionHandler
}

func NewDeleteUserTokenHandler(log zerolog.Logger, ah *action.ActionHandler) *DeleteUserTokenHandler {
	return &DeleteUserTokenHandler{log: log, ah: ah}
}

func (h *DeleteUserTokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	userRef := vars["userref"]
	tokenName := vars["tokenname"]

	h.log.Info().Msgf("deleting user %q token %q", userRef, tokenName)
	err := h.ah.DeleteUserToken(ctx, userRef, tokenName)
	if util.HTTPError(w, err) {
		h.log.Err(err).Send()
		return
	}

	if err := util.HTTPResponse(w, http.StatusNoContent, nil); err != nil {
		h.log.Err(err).Send()
	}
}

type RegisterUserHandler struct {
	log zerolog.Logger
	ah  *action.ActionHandler
}

func NewRegisterUserHandler(log zerolog.Logger, ah *action.ActionHandler) *RegisterUserHandler {
	return &RegisterUserHandler{log: log, ah: ah}
}

func (h *RegisterUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req *gwapitypes.RegisterUserRequest
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&req); err != nil {
		util.HTTPError(w, util.NewAPIError(util.ErrBadRequest, err))
		return
	}

	res, err := h.registerUser(ctx, req)
	if util.HTTPError(w, err) {
		h.log.Err(err).Send()
		return
	}

	if err := util.HTTPResponse(w, http.StatusCreated, res); err != nil {
		h.log.Err(err).Send()
	}
}

func (h *RegisterUserHandler) registerUser(ctx context.Context, req *gwapitypes.RegisterUserRequest) (*gwapitypes.RegisterUserResponse, error) {
	creq := &action.RegisterUserRequest{
		UserName:         req.CreateUserRequest.UserName,
		RemoteUserName:   req.CreateUserLARequest.RemoteSourceLoginName,
		RemotePassword:   req.CreateUserLARequest.RemoteSourceLoginPassword,
		RemoteSourceName: req.CreateUserLARequest.RemoteSourceName,
	}

	cresp, err := h.ah.HandleRemoteSourceAuth(ctx, req.CreateUserLARequest.RemoteSourceName, req.CreateUserLARequest.RemoteSourceLoginName, req.CreateUserLARequest.RemoteSourceLoginPassword, action.RemoteSourceRequestTypeRegisterUser, creq)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if cresp.Oauth2Redirect != "" {
		return &gwapitypes.RegisterUserResponse{
			Oauth2Redirect: cresp.Oauth2Redirect,
		}, nil
	}
	//authresp := cresp.Response.(*action.RegisterUserResponse)

	resp := &gwapitypes.RegisterUserResponse{}
	return resp, nil
}

type AuthorizeHandler struct {
	log zerolog.Logger
	ah  *action.ActionHandler
}

func NewAuthorizeHandler(log zerolog.Logger, ah *action.ActionHandler) *AuthorizeHandler {
	return &AuthorizeHandler{log: log, ah: ah}
}

func (h *AuthorizeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req *gwapitypes.LoginUserRequest
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&req); err != nil {
		util.HTTPError(w, util.NewAPIError(util.ErrBadRequest, err))
		return
	}

	res, err := h.authorize(ctx, req)
	if util.HTTPError(w, err) {
		h.log.Err(err).Send()
		return
	}

	if err := util.HTTPResponse(w, http.StatusCreated, res); err != nil {
		h.log.Err(err).Send()
	}
}

func (h *AuthorizeHandler) authorize(ctx context.Context, req *gwapitypes.LoginUserRequest) (*gwapitypes.AuthorizeResponse, error) {
	creq := &action.LoginUserRequest{
		RemoteSourceName: req.RemoteSourceName,
	}

	cresp, err := h.ah.HandleRemoteSourceAuth(ctx, req.RemoteSourceName, req.LoginName, req.LoginPassword, action.RemoteSourceRequestTypeAuthorize, creq)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if cresp.Oauth2Redirect != "" {
		return &gwapitypes.AuthorizeResponse{
			Oauth2Redirect: cresp.Oauth2Redirect,
		}, nil
	}
	authresp := cresp.Response.(*action.AuthorizeResponse)

	resp := &gwapitypes.AuthorizeResponse{
		RemoteUserInfo: &gwapitypes.UserInfo{
			ID:        authresp.RemoteUserInfo.ID,
			LoginName: authresp.RemoteUserInfo.LoginName,
			Email:     authresp.RemoteUserInfo.Email,
		},
		RemoteSourceName: authresp.RemoteSourceName,
	}
	return resp, nil
}

type LoginUserHandler struct {
	log zerolog.Logger
	ah  *action.ActionHandler
}

func NewLoginUserHandler(log zerolog.Logger, ah *action.ActionHandler) *LoginUserHandler {
	return &LoginUserHandler{log: log, ah: ah}
}

func (h *LoginUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req *gwapitypes.LoginUserRequest
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&req); err != nil {
		util.HTTPError(w, util.NewAPIError(util.ErrBadRequest, err))
		return
	}

	res, err := h.loginUser(ctx, w, req)
	if util.HTTPError(w, err) {
		h.log.Err(err).Send()
		return
	}

	if err := util.HTTPResponse(w, http.StatusCreated, res); err != nil {
		h.log.Err(err).Send()
	}
}

func (h *LoginUserHandler) loginUser(ctx context.Context, w http.ResponseWriter, req *gwapitypes.LoginUserRequest) (*gwapitypes.LoginUserResponse, error) {
	creq := &action.LoginUserRequest{
		RemoteSourceName: req.RemoteSourceName,
	}

	h.log.Info().Msgf("logging in user")
	cresp, err := h.ah.HandleRemoteSourceAuth(ctx, req.RemoteSourceName, req.LoginName, req.LoginPassword, action.RemoteSourceRequestTypeLoginUser, creq)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if cresp.Oauth2Redirect != "" {
		return &gwapitypes.LoginUserResponse{
			Oauth2Redirect: cresp.Oauth2Redirect,
		}, nil
	}
	authresp := cresp.Response.(*action.LoginUserResponse)

	resp := &gwapitypes.LoginUserResponse{
		User: createUserResponse(authresp.User),
	}

	http.SetCookie(w, authresp.Cookie)
	http.SetCookie(w, authresp.SecondaryCookie)

	return resp, nil
}

type UserCreateRunHandler struct {
	log zerolog.Logger
	ah  *action.ActionHandler
}

func NewUserCreateRunHandler(log zerolog.Logger, ah *action.ActionHandler) *UserCreateRunHandler {
	return &UserCreateRunHandler{log: log, ah: ah}
}

func (h *UserCreateRunHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req gwapitypes.UserCreateRunRequest
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&req); err != nil {
		util.HTTPError(w, util.NewAPIError(util.ErrBadRequest, err))
		return
	}

	creq := &action.UserCreateRunRequest{
		RepoUUID:              req.RepoUUID,
		RepoPath:              req.RepoPath,
		Branch:                req.Branch,
		Tag:                   req.Tag,
		Ref:                   req.Ref,
		CommitSHA:             req.CommitSHA,
		Message:               req.Message,
		PullRequestRefRegexes: req.PullRequestRefRegexes,
		Variables:             req.Variables,
	}
	err := h.ah.UserCreateRun(ctx, creq)
	if util.HTTPError(w, err) {
		h.log.Err(err).Send()
		return
	}

	if err := util.HTTPResponse(w, http.StatusCreated, nil); err != nil {
		h.log.Err(err).Send()
	}
}

type UserOrgsHandler struct {
	log zerolog.Logger
	ah  *action.ActionHandler
}

func NewUserOrgsHandler(log zerolog.Logger, ah *action.ActionHandler) *UserOrgsHandler {
	return &UserOrgsHandler{log: log, ah: ah}
}

func (h *UserOrgsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res, err := h.do(w, r)
	if util.HTTPError(w, err) {
		h.log.Err(err).Send()
		return
	}

	if err := util.HTTPResponse(w, http.StatusOK, res); err != nil {
		h.log.Err(err).Send()
	}
}

func (h *UserOrgsHandler) do(w http.ResponseWriter, r *http.Request) ([]*gwapitypes.UserOrgResponse, error) {
	ctx := r.Context()

	userID := common.CurrentUserID(ctx)
	if userID == "" {
		return nil, util.NewAPIError(util.ErrBadRequest, errors.Errorf("user not authenticated"))
	}

	ropts, err := parseRequestOptions(r)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	ares, err := h.ah.GetUserOrgs(ctx, &action.GetUserOrgsRequest{UserRef: userID, Cursor: ropts.Cursor, Limit: ropts.Limit, SortDirection: action.SortDirection(ropts.SortDirection)})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res := make([]*gwapitypes.UserOrgResponse, len(ares.Orgs))
	for i, userOrg := range ares.Orgs {
		res[i] = createUserOrgsResponse(userOrg)
	}

	addCursorHeader(w, ares.Cursor)

	return res, nil
}

func createUserOrgsResponse(o *csapitypes.UserOrgResponse) *gwapitypes.UserOrgResponse {
	userOrgs := &gwapitypes.UserOrgResponse{
		Organization: createOrgResponse(o.Organization),
		Role:         gwapitypes.MemberRole(o.Role),
	}

	return userOrgs
}

type UserOrgInvitationsHandler struct {
	log zerolog.Logger
	ah  *action.ActionHandler
}

func NewUserOrgInvitationsHandler(log zerolog.Logger, ah *action.ActionHandler) *UserOrgInvitationsHandler {
	return &UserOrgInvitationsHandler{log: log, ah: ah}
}

func (h *UserOrgInvitationsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := common.CurrentUserID(ctx)
	if userID == "" {
		util.HTTPError(w, util.NewAPIError(util.ErrBadRequest, errors.Errorf("user not authenticated")))
		return
	}

	user, err := h.ah.GetUser(ctx, userID)
	if util.HTTPError(w, err) {
		h.log.Err(err).Send()
		return
	}

	query := r.URL.Query()

	limitS := query.Get("limit")
	limit := DefaultOrgInvitationsLimit
	if limitS != "" {
		var err error
		limit, err = strconv.Atoi(limitS)
		if err != nil {
			util.HTTPError(w, util.NewAPIError(util.ErrBadRequest, errors.Wrapf(err, "cannot parse limit")))
			return
		}
	}
	if limit < 0 {
		util.HTTPError(w, util.NewAPIError(util.ErrBadRequest, errors.Errorf("limit must be greater or equal than 0")))
		return
	}
	if limit > MaxOrgInvitationsLimit {
		limit = MaxOrgInvitationsLimit
	}

	userInvitations, err := h.ah.GetUserOrgInvitations(ctx, user.ID, limit)
	if util.HTTPError(w, err) {
		h.log.Err(err).Send()
		return
	}

	orgInvitations := make([]*gwapitypes.OrgInvitationResponse, len(userInvitations))
	for i, p := range userInvitations {
		orgInvitations[i] = createOrgInvitationResponse(p.OrgInvitation, p.Organization)
	}

	if err := util.HTTPResponse(w, http.StatusOK, orgInvitations); err != nil {
		h.log.Err(err).Send()
	}
}
