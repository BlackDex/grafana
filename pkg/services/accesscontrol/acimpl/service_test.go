package acimpl

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grafana/grafana/pkg/api/routing"
	"github.com/grafana/grafana/pkg/infra/db"
	"github.com/grafana/grafana/pkg/infra/localcache"
	"github.com/grafana/grafana/pkg/infra/log"
	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/models/roletype"
	"github.com/grafana/grafana/pkg/services/accesscontrol"
	"github.com/grafana/grafana/pkg/services/accesscontrol/actest"
	"github.com/grafana/grafana/pkg/services/accesscontrol/database"
	"github.com/grafana/grafana/pkg/services/user"
	"github.com/grafana/grafana/pkg/setting"
)

func setupTestEnv(t testing.TB) *Service {
	t.Helper()
	cfg := setting.NewCfg()
	cfg.RBACEnabled = true

	ac := &Service{
		cfg:           cfg,
		log:           log.New("accesscontrol"),
		registrations: accesscontrol.RegistrationList{},
		store:         database.ProvideService(db.InitTestDB(t)),
		roles:         accesscontrol.BuildBasicRoleDefinitions(),
	}
	require.NoError(t, ac.RegisterFixedRoles(context.Background()))
	return ac
}

func TestUsageMetrics(t *testing.T) {
	tests := []struct {
		name          string
		enabled       bool
		expectedValue int
	}{
		{
			name:          "Expecting metric with value 0",
			enabled:       false,
			expectedValue: 0,
		},
		{
			name:          "Expecting metric with value 1",
			enabled:       true,
			expectedValue: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := setting.NewCfg()
			cfg.RBACEnabled = tt.enabled

			s, errInitAc := ProvideService(
				cfg,
				db.InitTestDB(t),
				routing.NewRouteRegister(),
				localcache.ProvideService(),
				actest.FakeAccessControl{},
			)
			require.NoError(t, errInitAc)
			assert.Equal(t, tt.expectedValue, s.GetUsageStats(context.Background())["stats.oss.accesscontrol.enabled.count"])
		})
	}
}

func TestService_DeclareFixedRoles(t *testing.T) {
	tests := []struct {
		name          string
		registrations []accesscontrol.RoleRegistration
		wantErr       bool
		err           error
	}{
		{
			name:    "should work with empty list",
			wantErr: false,
		},
		{
			name: "should add registration",
			registrations: []accesscontrol.RoleRegistration{
				{
					Role: accesscontrol.RoleDTO{
						Name: "fixed:test:test",
					},
					Grants: []string{"Admin"},
				},
			},
			wantErr: false,
		},
		{
			name: "should fail registration invalid role name",
			registrations: []accesscontrol.RoleRegistration{
				{
					Role: accesscontrol.RoleDTO{
						Name: "custom:test:test",
					},
					Grants: []string{"Admin"},
				},
			},
			wantErr: true,
			err:     accesscontrol.ErrFixedRolePrefixMissing,
		},
		{
			name: "should fail registration invalid builtin role assignment",
			registrations: []accesscontrol.RoleRegistration{
				{
					Role: accesscontrol.RoleDTO{
						Name: "fixed:test:test",
					},
					Grants: []string{"WrongAdmin"},
				},
			},
			wantErr: true,
			err:     accesscontrol.ErrInvalidBuiltinRole,
		},
		{
			name: "should add multiple registrations at once",
			registrations: []accesscontrol.RoleRegistration{
				{
					Role: accesscontrol.RoleDTO{
						Name: "fixed:test:test",
					},
					Grants: []string{"Admin"},
				},
				{
					Role: accesscontrol.RoleDTO{
						Name: "fixed:test2:test2",
					},
					Grants: []string{"Admin"},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ac := setupTestEnv(t)

			// Reset the registations
			ac.registrations = accesscontrol.RegistrationList{}

			// Test
			err := ac.DeclareFixedRoles(tt.registrations...)
			if tt.wantErr {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)

			registrationCnt := 0
			ac.registrations.Range(func(registration accesscontrol.RoleRegistration) bool {
				registrationCnt++
				return true
			})
			assert.Equal(t, len(tt.registrations), registrationCnt,
				"expected service registration list to contain all test registrations")
		})
	}
}

func TestService_RegisterFixedRoles(t *testing.T) {
	tests := []struct {
		name          string
		token         models.Licensing
		registrations []accesscontrol.RoleRegistration
		wantErr       bool
	}{
		{
			name: "should work with empty list",
		},
		{
			name: "should register and assign role",
			registrations: []accesscontrol.RoleRegistration{
				{
					Role: accesscontrol.RoleDTO{
						Name:        "fixed:test:test",
						Permissions: []accesscontrol.Permission{{Action: "test:test"}},
					},
					Grants: []string{"Editor"},
				},
			},
			wantErr: false,
		},
		{
			name: "should register and assign multiple roles",
			registrations: []accesscontrol.RoleRegistration{
				{
					Role: accesscontrol.RoleDTO{
						Name:        "fixed:test:test",
						Permissions: []accesscontrol.Permission{{Action: "test:test"}},
					},
					Grants: []string{"Editor"},
				},
				{
					Role: accesscontrol.RoleDTO{
						Name: "fixed:test2:test2",
						Permissions: []accesscontrol.Permission{
							{Action: "test:test2"},
							{Action: "test:test3", Scope: "test:*"},
						},
					},
					Grants: []string{"Viewer"},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ac := setupTestEnv(t)

			ac.registrations.Append(tt.registrations...)

			// Test
			err := ac.RegisterFixedRoles(context.Background())
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Check
			for _, registration := range tt.registrations {
				// Check builtin roles (parents included) have been granted with the permissions
				for br := range accesscontrol.BuiltInRolesWithParents(registration.Grants) {
					builtinRole, ok := ac.roles[br]
					assert.True(t, ok)
					for _, expectedPermission := range registration.Role.Permissions {
						assert.Contains(t, builtinRole.Permissions, expectedPermission)
					}
				}
			}
		})
	}
}

func TestService_GetSimplifiedUsersPermissions(t *testing.T) {
	actionPrefix := "teams"
	ctx := context.Background()
	listAllPerms := map[string][]string{accesscontrol.ActionUsersPermissionsRead: {"users:*"}}
	listSomePerms := map[string][]string{accesscontrol.ActionUsersPermissionsRead: {"users:id:2"}}
	tests := []struct {
		name           string
		siuPermissions map[string][]string
		ramRoles       map[string]*accesscontrol.RoleDTO    // BasicRole => RBAC BasicRole
		storedPerms    map[int64][]accesscontrol.Permission // UserID => Permissions
		storedRoles    map[int64][]string                   // UserID => Roles
		want           map[int64][]accesscontrol.SimplifiedUserPermissionDTO
		wantErr        bool
	}{
		{
			name:           "ram only",
			siuPermissions: listAllPerms,
			ramRoles: map[string]*accesscontrol.RoleDTO{
				string(roletype.RoleAdmin): {Permissions: []accesscontrol.Permission{
					{Action: accesscontrol.ActionTeamsRead, Scope: "teams:*"},
				}},
				accesscontrol.RoleGrafanaAdmin: {Permissions: []accesscontrol.Permission{
					{Action: accesscontrol.ActionTeamsPermissionsRead, Scope: "teams:*"},
				}},
			},
			storedRoles: map[int64][]string{
				1: {string(roletype.RoleEditor)},
				2: {string(roletype.RoleAdmin), accesscontrol.RoleGrafanaAdmin},
			},
			want: map[int64][]accesscontrol.SimplifiedUserPermissionDTO{
				2: {{Action: accesscontrol.ActionTeamsRead, All: true},
					{Action: accesscontrol.ActionTeamsPermissionsRead, All: true}},
			},
		},
		{
			name:           "stored only",
			siuPermissions: listAllPerms,
			storedPerms: map[int64][]accesscontrol.Permission{
				1: {{Action: accesscontrol.ActionTeamsRead, Scope: "teams:id:1"}},
				2: {{Action: accesscontrol.ActionTeamsRead, Scope: "teams:*"},
					{Action: accesscontrol.ActionTeamsPermissionsRead, Scope: "teams:*"}},
			},
			storedRoles: map[int64][]string{
				1: {string(roletype.RoleEditor)},
				2: {string(roletype.RoleAdmin), accesscontrol.RoleGrafanaAdmin},
			},
			want: map[int64][]accesscontrol.SimplifiedUserPermissionDTO{
				1: {{Action: accesscontrol.ActionTeamsRead, UIDs: []string{"1"}}},
				2: {{Action: accesscontrol.ActionTeamsRead, All: true},
					{Action: accesscontrol.ActionTeamsPermissionsRead, All: true}},
			},
		},
		{
			name:           "ram and stored",
			siuPermissions: listAllPerms,
			ramRoles: map[string]*accesscontrol.RoleDTO{
				string(roletype.RoleAdmin): {Permissions: []accesscontrol.Permission{
					{Action: accesscontrol.ActionTeamsRead, Scope: "teams:*"},
				}},
				accesscontrol.RoleGrafanaAdmin: {Permissions: []accesscontrol.Permission{
					{Action: accesscontrol.ActionTeamsPermissionsRead, Scope: "teams:*"},
				}},
			},
			storedPerms: map[int64][]accesscontrol.Permission{
				1: {{Action: accesscontrol.ActionTeamsRead, Scope: "teams:id:1"}},
				2: {{Action: accesscontrol.ActionTeamsRead, Scope: "teams:id:1"},
					{Action: accesscontrol.ActionTeamsPermissionsRead, Scope: "teams:id:1"}},
			},
			storedRoles: map[int64][]string{
				1: {string(roletype.RoleEditor)},
				2: {string(roletype.RoleAdmin), accesscontrol.RoleGrafanaAdmin},
			},
			want: map[int64][]accesscontrol.SimplifiedUserPermissionDTO{
				1: {{Action: accesscontrol.ActionTeamsRead, UIDs: []string{"1"}}},
				2: {{Action: accesscontrol.ActionTeamsRead, All: true},
					{Action: accesscontrol.ActionTeamsPermissionsRead, All: true}},
			},
		},
		{
			name:           "view permission on subset of users only",
			siuPermissions: listSomePerms,
			ramRoles: map[string]*accesscontrol.RoleDTO{
				accesscontrol.RoleGrafanaAdmin: {Permissions: []accesscontrol.Permission{
					{Action: accesscontrol.ActionTeamsPermissionsRead, Scope: "teams:*"},
				}},
			},
			storedPerms: map[int64][]accesscontrol.Permission{
				1: {{Action: accesscontrol.ActionTeamsRead, Scope: "teams:id:1"}},
				2: {{Action: accesscontrol.ActionTeamsRead, Scope: "teams:id:1"},
					{Action: accesscontrol.ActionTeamsPermissionsRead, Scope: "teams:id:1"}},
			},
			storedRoles: map[int64][]string{
				1: {string(roletype.RoleEditor)},
				2: {accesscontrol.RoleGrafanaAdmin},
			},
			want: map[int64][]accesscontrol.SimplifiedUserPermissionDTO{
				2: {{Action: accesscontrol.ActionTeamsRead, UIDs: []string{"1"}},
					{Action: accesscontrol.ActionTeamsPermissionsRead, All: true}},
			},
		},
		{
			name:           "check action filter on RAM permissions works correctly",
			siuPermissions: listAllPerms,
			ramRoles: map[string]*accesscontrol.RoleDTO{
				accesscontrol.RoleGrafanaAdmin: {Permissions: []accesscontrol.Permission{
					{Action: accesscontrol.ActionUsersCreate},
					{Action: accesscontrol.ActionTeamsPermissionsRead, Scope: "teams:*"},
				}},
			},
			storedRoles: map[int64][]string{1: {accesscontrol.RoleGrafanaAdmin}},
			want: map[int64][]accesscontrol.SimplifiedUserPermissionDTO{
				1: {{Action: accesscontrol.ActionTeamsPermissionsRead, All: true}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ac := setupTestEnv(t)

			ac.roles = tt.ramRoles
			ac.store = actest.FakeStore{
				ExpectedUsersPermissions: tt.storedPerms,
				ExpectedUsersRoles:       tt.storedRoles,
			}

			siu := &user.SignedInUser{OrgID: 2, Permissions: map[int64]map[string][]string{2: tt.siuPermissions}}
			got, err := ac.GetSimplifiedUsersPermissions(ctx, siu, 2, actionPrefix)
			if tt.wantErr {
				require.NotNil(t, err)
				return
			}
			require.Nil(t, err)

			require.Len(t, got, len(tt.want), "expected more users permissions")
			for userID, wantPerm := range tt.want {
				gotPerm, ok := got[userID]
				require.True(t, ok, "expected permissions for user", userID)

				require.ElementsMatch(t, gotPerm, wantPerm)
			}
		})
	}
}
