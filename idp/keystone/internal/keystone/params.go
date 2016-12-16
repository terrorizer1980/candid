// Copyright 2015 Canonical Ltd.

// Package keystone implements a keystone client.
package keystone

import (
	"net/http"
	"time"

	"github.com/juju/httprequest"
)

// TokensRequest is the request sent to /v2.0/tokens to perform a login.
// See
// http://developer.openstack.org/api-ref-identity-v2.html#authenticate-v2.0
// for more information.
type TokensRequest struct {
	httprequest.Route `httprequest:"POST /v2.0/tokens"`
	Body              TokensBody `httprequest:",body"`
}

// TokensBody represents the JSON body sent in a login request.
type TokensBody struct {
	Auth Auth `json:"auth"`
}

// TokensResponse is the response from /v2.0/tokens on success.
type TokensResponse struct {
	Access Access `json:"access"`
}

// Auth is the authentication information sent in a login request.
type Auth struct {
	TenantName          string               `json:"tenantName,omitempty"`
	TenantID            string               `json:"tenantId,omitempty"`
	PasswordCredentials *PasswordCredentials `json:"passwordCredentials,omitempty"`
	Token               *Token               `json:"token,omitempty"`
}

// PasswordCredentials holds the credentials for a username/password
// authentication.
type PasswordCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Token contains the details of a token generated by keystone.
type Token struct {
	ID       string  `json:"id,omitempty"`
	IssuedAt *Time   `json:"issued_at,omitempty"`
	Expires  *Time   `json:"expires,omitempty"`
	Tenant   *Tenant `json:"tenant,omitempty"`
}

// Tenant contains details of a tenant in the openstack environment.
type Tenant struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

// Access contains the access granted in the login attempt.
type Access struct {
	Token Token `json:"token"`
	User  User  `json:"user"`
}

// User contains details of a user in the openstack environment.
type User struct {
	ID       string  `json:"id,omitempty"`
	Name     string  `json:"name,omitempty"`
	Username string  `json:"username,omitemtpy"`
	Domain   *Domain `json:"domain,omitempty"`
	Password string  `json:"password,omitempty"`
}

// TenantsRequest is the request sent to /v2.0/tenants to list tenants a
// token has access to. See
// http://developer.openstack.org/api-ref-identity-v2.html#listTenants
// for more information.
type TenantsRequest struct {
	httprequest.Route `httprequest:"GET /v2.0/tenants"`
	AuthToken         string `httprequest:"X-Auth-Token,header"`
}

// TenantsResponse is the list of tenants a token has access to.
type TenantsResponse struct {
	Tenants []Tenant `json:"tenants"`
}

// Time is a time.Time that provides a custom UnmarshalJSON method.
type Time struct {
	time.Time
}

func (t *Time) UnmarshalJSON(data []byte) error {
	if err := t.Time.UnmarshalJSON(data); err == nil {
		return nil
	}
	var err error
	t.Time, err = time.Parse(`"2006-01-02T15:04:05"`, string(data))
	return err
}

// AuthTokensRequest is the request sent to /v3/auth/tokens to perform a
// login. See
// http://developer.openstack.org/api-ref/identity/v3/index.html?expanded=password-authentication-with-unscoped-authorization-detail
// for more information.
type AuthTokensRequest struct {
	httprequest.Route `httprequest:"POST /v3/auth/tokens"`
	Body              AuthTokensBody `httprequest:",body"`
}

// AuthTokensBody represents the JSON body sent in a v3 login request.
type AuthTokensBody struct {
	Auth AuthV3 `json:"auth"`
}

// AuthV3 is the authentication information sent in a v3 login request.
type AuthV3 struct {
	Identity Identity `json:"identity"`
}

// Identity contains the identity information sent in a v3 login request.
type Identity struct {
	Methods  []string       `json:"methods"`
	Password *Password      `json:"password,omitempty"`
	Token    *IdentityToken `json:"token,omitempty"`
}

// Password contains the password based identity information sent in a
// v3 login request.
type Password struct {
	User User `json:"user"`
}

// IdentityToken contains the token based identity information sent in a
// v3 login request.
type IdentityToken struct {
	ID string `json:"id"`
}

// Domain contains the domain of a user in the v3 API.
type Domain struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name:omitempty"`
}

// AuthTokensResponse is the reponse sent by /v3/auth/tokens when there
// has been a successful login. See
// http://developer.openstack.org/api-ref/identity/v3/index.html?expanded=password-authentication-with-unscoped-authorization-detail
// for more information.
type AuthTokensResponse struct {
	SubjectToken string
	Token        TokenV3 `json:"token"`
}

// SetHeader implements httprequest.HeaderSetter by setting the
// appropriate X-Subject-Token header for the response.
func (resp AuthTokensResponse) SetHeader(h http.Header) {
	h.Set(subjectTokenHeader, resp.SubjectToken)
}

// TokenV3 represents the token returned from /v3/auth/tokens after a
// successful login.
type TokenV3 struct {
	IssuedAt  *Time    `json:"issued_at,omitempty"`
	Methods   []string `json:"methods,omitempty"`
	ExpiresAt *Time    `json:"expires_at,omitempty"`
	User      User     `json:"user"`
}

// UserGroupsRequest represents a request to the /v3/users/:id/groups
// endpoint. See
// http://developer.openstack.org/api-ref/identity/v3/index.html?expanded=list-groups-to-which-a-user-belongs-detail
// for more information.
type UserGroupsRequest struct {
	httprequest.Route `httprequest:"GET /v3/users/:UserID/groups"`
	UserID            string `httprequest:",path"`
	AuthToken         string `httprequest:"X-Auth-Token,header"`
}

// UserGroupsResponse represents a response to the /v3/users/:id/groups
// endpoint. See
// http://developer.openstack.org/api-ref/identity/v3/index.html?expanded=list-groups-to-which-a-user-belongs-detail
// for more information.
type UserGroupsResponse struct {
	Groups []Group `json:"groups"`
}

// Group contains information on a keystone group.
type Group struct {
	ID          string `json:"id"`
	DomainID    string `json:"domain_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
