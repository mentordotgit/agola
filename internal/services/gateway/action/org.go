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

package action

import (
	"context"

	"github.com/sorintlab/errors"

	"agola.io/agola/internal/services/gateway/common"
	"agola.io/agola/internal/util"
	csapitypes "agola.io/agola/services/configstore/api/types"
	"agola.io/agola/services/configstore/client"
	cstypes "agola.io/agola/services/configstore/types"
)

func (h *ActionHandler) GetOrg(ctx context.Context, orgRef string) (*cstypes.Organization, error) {
	org, _, err := h.configstoreClient.GetOrg(ctx, orgRef)
	if err != nil {
		return nil, util.NewAPIError(util.KindFromRemoteError(err), err)
	}

	if org.Visibility == cstypes.VisibilityPublic {
		return org, nil
	}

	isMember, err := h.IsAuthUserMember(ctx, cstypes.ObjectKindOrg, org.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to determine ownership")
	}
	if !isMember {
		return nil, util.NewAPIError(util.ErrForbidden, errors.Errorf("user not authorized"))
	}

	return org, nil
}

type GetOrgsRequest struct {
	Cursor string

	Limit         int
	SortDirection SortDirection

	Public  bool
	Private bool
}

type GetOrgsResponse struct {
	Orgs   []*cstypes.Organization
	Cursor string
}

func (h *ActionHandler) GetOrgs(ctx context.Context, req *GetOrgsRequest) (*GetOrgsResponse, error) {
	inCursor := &StartCursor{}
	sortDirection := req.SortDirection
	if req.Cursor != "" {
		if err := UnmarshalCursor(req.Cursor, inCursor); err != nil {
			return nil, errors.WithStack(err)
		}
		sortDirection = inCursor.SortDirection
	}
	if sortDirection == "" {
		sortDirection = SortDirectionAsc
	}

	isAdmin := common.IsUserAdmin(ctx)

	// only admin will also get private organizations
	// normal user will only get public organizations. Private organizations where they are members won't be returned.
	visibilites := []cstypes.Visibility{cstypes.VisibilityPublic}
	if isAdmin {
		visibilites = append(visibilites, cstypes.VisibilityPrivate)
	}

	orgs, resp, err := h.configstoreClient.GetOrgs(ctx, &client.GetOrgsOptions{ListOptions: &client.ListOptions{Limit: req.Limit, SortDirection: cstypes.SortDirection(sortDirection)}, StartOrgName: inCursor.Start, Visibilities: visibilites})
	if err != nil {
		return nil, util.NewAPIError(util.KindFromRemoteError(err), err)
	}

	var outCursor string
	if resp.HasMore && len(orgs) > 0 {
		lastRemoteSourceName := orgs[len(orgs)-1].Name
		outCursor, err = MarshalCursor(&StartCursor{
			Start:         lastRemoteSourceName,
			SortDirection: sortDirection,
		})
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	res := &GetOrgsResponse{
		Orgs:   orgs,
		Cursor: outCursor,
	}

	return res, nil
}

type GetOrgMembersRequest struct {
	OrgRef string

	Cursor string

	Limit         int
	SortDirection SortDirection
}

type OrgMemberResponse struct {
	User *cstypes.User
	Role cstypes.MemberRole
}

type GetOrgMembersResponse struct {
	Organization *cstypes.Organization
	Members      []*OrgMemberResponse
	Cursor       string
}

func (h *ActionHandler) GetOrgMembers(ctx context.Context, req *GetOrgMembersRequest) (*GetOrgMembersResponse, error) {
	inCursor := &StartCursor{}
	sortDirection := req.SortDirection
	if req.Cursor != "" {
		if err := UnmarshalCursor(req.Cursor, inCursor); err != nil {
			return nil, errors.WithStack(err)
		}
		sortDirection = inCursor.SortDirection
	}
	if sortDirection == "" {
		sortDirection = SortDirectionAsc
	}

	org, _, err := h.configstoreClient.GetOrg(ctx, req.OrgRef)
	if err != nil {
		return nil, util.NewAPIError(util.KindFromRemoteError(err), err)
	}

	orgMembers, resp, err := h.configstoreClient.GetOrgMembers(ctx, req.OrgRef, &client.GetOrgMembersOptions{ListOptions: &client.ListOptions{Limit: req.Limit, SortDirection: cstypes.SortDirection(sortDirection)}, StartUserName: inCursor.Start})
	if err != nil {
		return nil, util.NewAPIError(util.KindFromRemoteError(err), err)
	}

	var outCursor string
	if resp.HasMore && len(orgMembers) > 0 {
		lastUserName := orgMembers[len(orgMembers)-1].User.Name
		outCursor, err = MarshalCursor(&StartCursor{
			Start:         lastUserName,
			SortDirection: sortDirection,
		})
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	res := &GetOrgMembersResponse{
		Organization: org,
		Members:      make([]*OrgMemberResponse, len(orgMembers)),
		Cursor:       outCursor,
	}

	for i, orgMember := range orgMembers {
		res.Members[i] = &OrgMemberResponse{
			User: orgMember.User,
			Role: orgMember.Role,
		}
	}

	return res, nil
}

type CreateOrgRequest struct {
	Name       string
	Visibility cstypes.Visibility

	CreatorUserID string
}

func (h *ActionHandler) CreateOrg(ctx context.Context, req *CreateOrgRequest) (*cstypes.Organization, error) {
	if !common.IsUserLoggedOrAdmin(ctx) {
		return nil, errors.Errorf("user not logged in")
	}

	if req.Name == "" {
		return nil, util.NewAPIError(util.ErrBadRequest, errors.Errorf("organization name required"))
	}
	if !util.ValidateName(req.Name) {
		return nil, util.NewAPIError(util.ErrBadRequest, errors.Errorf("invalid organization name %q", req.Name))
	}

	creq := &csapitypes.CreateOrgRequest{
		Name:       req.Name,
		Visibility: req.Visibility,
	}
	if req.CreatorUserID != "" {
		creq.CreatorUserID = req.CreatorUserID
	}

	h.log.Info().Msgf("creating organization")
	org, _, err := h.configstoreClient.CreateOrg(ctx, creq)
	if err != nil {
		return nil, util.NewAPIError(util.KindFromRemoteError(err), errors.Wrapf(err, "failed to create organization"))
	}
	h.log.Info().Msgf("organization %s created, ID: %s", org.Name, org.ID)

	return org, nil
}

type UpdateOrgRequest struct {
	Visibility *cstypes.Visibility
}

func (h *ActionHandler) UpdateOrg(ctx context.Context, orgRef string, req *UpdateOrgRequest) (*cstypes.Organization, error) {
	org, _, err := h.configstoreClient.GetOrg(ctx, orgRef)
	if err != nil {
		return nil, util.NewAPIError(util.KindFromRemoteError(err), err)
	}

	isOrgOwner, err := h.IsAuthUserOrgOwner(ctx, org.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to determine ownership")
	}
	if !isOrgOwner {
		return nil, util.NewAPIError(util.ErrForbidden, errors.Errorf("user not authorized"))
	}

	if req.Visibility != nil {
		org.Visibility = *req.Visibility
	}

	creq := &csapitypes.UpdateOrgRequest{
		Visibility: org.Visibility,
	}

	h.log.Info().Msgf("updating organization")
	org, _, err = h.configstoreClient.UpdateOrg(ctx, orgRef, creq)
	if err != nil {
		return nil, util.NewAPIError(util.KindFromRemoteError(err), errors.Wrapf(err, "failed to update organization"))
	}
	h.log.Info().Msgf("organization %s updated, ID: %s", org.Name, org.ID)

	return org, nil
}

func (h *ActionHandler) DeleteOrg(ctx context.Context, orgRef string) error {
	org, _, err := h.configstoreClient.GetOrg(ctx, orgRef)
	if err != nil {
		return util.NewAPIError(util.KindFromRemoteError(err), err)
	}

	isOrgOwner, err := h.IsAuthUserOrgOwner(ctx, org.ID)
	if err != nil {
		return errors.Wrapf(err, "failed to determine ownership")
	}
	if !isOrgOwner {
		return util.NewAPIError(util.ErrForbidden, errors.Errorf("user not authorized"))
	}

	if _, err := h.configstoreClient.DeleteOrg(ctx, orgRef); err != nil {
		return util.NewAPIError(util.KindFromRemoteError(err), errors.Wrapf(err, "failed to delete org"))
	}
	return nil
}

type AddOrgMemberResponse struct {
	OrganizationMember *cstypes.OrganizationMember
	Org                *cstypes.Organization
	User               *cstypes.User
}

func (h *ActionHandler) AddOrgMember(ctx context.Context, orgRef, userRef string, role cstypes.MemberRole) (*AddOrgMemberResponse, error) {
	if h.organizationMemberAddingMode != OrganizationMemberAddingModeDirect && !common.IsUserAdmin(ctx) {
		return nil, util.NewAPIError(util.ErrBadRequest, errors.Errorf("cannot directly add user to organization"))
	}

	org, _, err := h.configstoreClient.GetOrg(ctx, orgRef)
	if err != nil {
		return nil, util.NewAPIError(util.KindFromRemoteError(err), err)
	}
	user, _, err := h.configstoreClient.GetUser(ctx, userRef)
	if err != nil {
		return nil, util.NewAPIError(util.KindFromRemoteError(err), err)
	}

	isOrgOwner, err := h.IsAuthUserOrgOwner(ctx, org.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to determine ownership")
	}
	if !isOrgOwner {
		return nil, util.NewAPIError(util.ErrForbidden, errors.Errorf("user not authorized"))
	}

	orgmember, _, err := h.configstoreClient.AddOrgMember(ctx, orgRef, userRef, role)
	if err != nil {
		return nil, util.NewAPIError(util.KindFromRemoteError(err), errors.Wrapf(err, "failed to add/update organization member"))
	}

	return &AddOrgMemberResponse{
		OrganizationMember: orgmember,
		Org:                org,
		User:               user,
	}, nil
}

func (h *ActionHandler) RemoveOrgMember(ctx context.Context, orgRef, userRef string) error {
	org, _, err := h.configstoreClient.GetOrg(ctx, orgRef)
	if err != nil {
		return util.NewAPIError(util.KindFromRemoteError(err), err)
	}

	isOrgOwner, err := h.IsAuthUserOrgOwner(ctx, org.ID)
	if err != nil {
		return errors.Wrapf(err, "failed to determine ownership")
	}
	if !isOrgOwner {
		return util.NewAPIError(util.ErrForbidden, errors.Errorf("user not authorized"))
	}

	if _, err = h.configstoreClient.RemoveOrgMember(ctx, orgRef, userRef); err != nil {
		return util.NewAPIError(util.KindFromRemoteError(err), errors.Wrapf(err, "failed to remove organization member"))
	}

	return nil
}

type OrgInvitationResponse struct {
	Organization  *cstypes.Organization
	OrgInvitation *cstypes.OrgInvitation
}

func (h *ActionHandler) GetOrgInvitations(ctx context.Context, orgRef string, limit int) ([]*cstypes.OrgInvitation, error) {
	if !common.IsUserLogged(ctx) {
		return nil, errors.Errorf("user not logged in")
	}

	org, _, err := h.configstoreClient.GetOrg(ctx, orgRef)
	if err != nil {
		return nil, util.NewAPIError(util.KindFromRemoteError(err), errors.Wrapf(err, "failed to get org %s", orgRef))
	}

	isOrgOwner, err := h.IsAuthUserOrgOwner(ctx, org.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to determine ownership")
	}
	if !isOrgOwner {
		return nil, util.NewAPIError(util.ErrForbidden, errors.Errorf("user not authorized"))
	}

	orgInvitations, _, err := h.configstoreClient.GetOrgInvitations(ctx, orgRef, limit)
	if err != nil {
		return nil, util.NewAPIError(util.KindFromRemoteError(err), err)
	}
	return orgInvitations, nil
}

type CreateOrgInvitationRequest struct {
	UserRef         string
	OrganizationRef string
	Role            cstypes.MemberRole
}

func (h *ActionHandler) CreateOrgInvitation(ctx context.Context, req *CreateOrgInvitationRequest) (*OrgInvitationResponse, error) {
	if !common.IsUserLogged(ctx) {
		return nil, errors.Errorf("user not logged in")
	}

	if h.organizationMemberAddingMode != OrganizationMemberAddingModeInvitation {
		return nil, util.NewAPIError(util.ErrBadRequest, errors.Errorf("user members can not added by invitation"))
	}

	if req.UserRef == "" {
		return nil, util.NewAPIError(util.ErrBadRequest, errors.Errorf("user id required"))
	}
	if req.OrganizationRef == "" {
		return nil, util.NewAPIError(util.ErrBadRequest, errors.Errorf("organization id required"))
	}
	if req.Role == "" {
		return nil, util.NewAPIError(util.ErrBadRequest, errors.Errorf("role is required"))
	}

	org, _, err := h.configstoreClient.GetOrg(ctx, req.OrganizationRef)
	if err != nil {
		return nil, util.NewAPIError(util.KindFromRemoteError(err), errors.Wrapf(err, "failed to get org %s", req.OrganizationRef))
	}

	isOrgOwner, err := h.IsAuthUserOrgOwner(ctx, org.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to determine ownership")
	}
	if !isOrgOwner {
		return nil, util.NewAPIError(util.ErrForbidden, errors.Errorf("user not authorized"))
	}

	isOrgMember, err := h.IsUserOrgMember(ctx, req.UserRef, req.OrganizationRef)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to determine membership")
	}
	if isOrgMember {
		return nil, errors.Errorf("user is already an org member")
	}

	_, _, err = h.configstoreClient.GetOrgInvitation(ctx, req.OrganizationRef, req.UserRef)
	if err != nil {
		if !util.RemoteErrorIs(err, util.ErrNotExist) {
			return nil, errors.Wrapf(err, "failed to determine if org invitation exists")
		}
	} else {
		return nil, errors.Errorf("invitation already exists")
	}

	creq := &csapitypes.CreateOrgInvitationRequest{
		UserRef: req.UserRef,
		Role:    req.Role,
	}

	h.log.Info().Msgf("creating org invitation")
	orgInvitation, _, err := h.configstoreClient.CreateOrgInvitation(ctx, org.ID, creq)
	if err != nil {
		return nil, util.NewAPIError(util.KindFromRemoteError(err), errors.Wrapf(err, "failed to create org invitation"))
	}
	h.log.Info().Msgf("org invitation created, ID: %s", orgInvitation.ID)

	return &OrgInvitationResponse{
		OrgInvitation: orgInvitation,
		Organization:  org,
	}, nil
}

type OrgInvitationActionRequest struct {
	OrgRef string
	Action csapitypes.OrgInvitationActionType `json:"action_type"`
}

func (h *ActionHandler) OrgInvitationAction(ctx context.Context, req *OrgInvitationActionRequest) error {
	if !req.Action.IsValid() {
		return errors.Errorf("action is not valid")
	}

	userID := common.CurrentUserID(ctx)
	if userID == "" {
		return util.NewAPIError(util.ErrBadRequest, errors.Errorf("user not authenticated"))
	}

	orgInvitation, _, err := h.configstoreClient.GetOrgInvitation(ctx, req.OrgRef, userID)
	if err != nil {
		return util.NewAPIError(util.KindFromRemoteError(err), errors.Wrapf(err, "failed to get org invitation"))
	}
	if orgInvitation == nil {
		return util.NewAPIError(util.ErrBadRequest, errors.Errorf("invitation for org %s user %s not found", req.OrgRef, userID))
	}

	if userID != orgInvitation.UserID {
		return errors.Errorf("user not authorized")
	}

	org, _, err := h.configstoreClient.GetOrg(ctx, req.OrgRef)
	if err != nil {
		return util.NewAPIError(util.KindFromRemoteError(err), errors.Wrapf(err, "failed to get org %s", req.OrgRef))
	}

	if org == nil {
		return util.NewAPIError(util.ErrBadRequest, errors.Errorf("org %s not found", req.OrgRef))
	}

	creq := &csapitypes.OrgInvitationActionRequest{Action: req.Action}
	_, err = h.configstoreClient.UserOrgInvitationAction(ctx, userID, req.OrgRef, creq)
	if err != nil {
		return util.NewAPIError(util.KindFromRemoteError(err), err)
	}

	return nil
}

func (h *ActionHandler) DeleteOrgInvitation(ctx context.Context, orgRef string, userRef string) error {
	userID := common.CurrentUserID(ctx)
	if userID == "" {
		return errors.Errorf("user not authenticated")
	}

	orgInvitation, _, err := h.configstoreClient.GetOrgInvitation(ctx, orgRef, userRef)
	if err != nil {
		return util.NewAPIError(util.KindFromRemoteError(err), err)
	}

	org, _, err := h.configstoreClient.GetOrg(ctx, orgInvitation.OrganizationID)
	if err != nil {
		return util.NewAPIError(util.KindFromRemoteError(err), errors.Wrapf(err, "failed to get org %s", orgInvitation.OrganizationID))
	}

	isOrgOwner, err := h.IsAuthUserOrgOwner(ctx, org.ID)
	if err != nil {
		return util.NewAPIError(util.ErrBadRequest, errors.Wrapf(err, "failed to determine ownership"))
	}
	if !isOrgOwner {
		return util.NewAPIError(util.ErrForbidden, errors.Errorf("user is not owner"))
	}

	_, err = h.configstoreClient.DeleteOrgInvitation(ctx, orgRef, userRef)
	if err != nil {
		return util.NewAPIError(util.KindFromRemoteError(err), errors.Wrapf(err, "failed to delete org invitation"))
	}
	return nil
}

func (h *ActionHandler) GetOrgInvitation(ctx context.Context, orgRef string, userRef string) (*OrgInvitationResponse, error) {
	cOrgInvitation, _, err := h.configstoreClient.GetOrgInvitation(ctx, orgRef, userRef)
	if err != nil {
		return nil, util.NewAPIError(util.KindFromRemoteError(err), err)
	}

	org, _, err := h.configstoreClient.GetOrg(ctx, cOrgInvitation.OrganizationID)
	if err != nil {
		return nil, util.NewAPIError(util.KindFromRemoteError(err), err)
	}

	res := OrgInvitationResponse{
		OrgInvitation: cOrgInvitation,
		Organization:  org,
	}

	return &res, nil
}
