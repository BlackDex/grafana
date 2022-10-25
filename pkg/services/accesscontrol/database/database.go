package database

import (
	"context"
	"strconv"
	"strings"

	"github.com/grafana/grafana/pkg/infra/db"
	"github.com/grafana/grafana/pkg/services/accesscontrol"
)

func ProvideService(sql db.DB) *AccessControlStore {
	return &AccessControlStore{sql}
}

type AccessControlStore struct {
	sql db.DB
}

func (s *AccessControlStore) GetUserPermissions(ctx context.Context, query accesscontrol.GetUserPermissionsQuery) ([]accesscontrol.Permission, error) {
	result := make([]accesscontrol.Permission, 0)
	err := s.sql.WithDbSession(ctx, func(sess *db.Session) error {
		if query.UserID == 0 && len(query.TeamIDs) == 0 && len(query.Roles) == 0 {
			// no permission to fetch
			return nil
		}

		filter, params := accesscontrol.UserRolesFilter(query.OrgID, query.UserID, query.TeamIDs, query.Roles)

		q := `
		SELECT
			permission.action,
			permission.scope
			FROM permission
			INNER JOIN role ON role.id = permission.role_id
		` + filter

		if len(query.Actions) > 0 {
			q += " WHERE permission.action IN("
			if len(query.Actions) > 0 {
				q += "?" + strings.Repeat(",?", len(query.Actions)-1)
			}
			q += ")"
			for _, a := range query.Actions {
				params = append(params, a)
			}
		}
		if err := sess.SQL(q, params...).Find(&result); err != nil {
			return err
		}

		return nil
	})

	return result, err
}

// GetUsersPermissions returns the list of user permissions indexed by UserID
func (s *AccessControlStore) GetUsersPermissions(ctx context.Context, orgID int64, actionPrefix string) (map[int64][]accesscontrol.Permission, map[int64][]string, error) {
	type UserRBACPermission struct {
		UserID int64  `xorm:"user_id"`
		Action string `xorm:"action"`
		Scope  string `xorm:"scope"`
	}
	type UserOrgRole struct {
		UserID  int64  `xorm:"id"`
		OrgRole string `xorm:"role"`
		IsAdmin bool   `xorm:"is_admin"`
	}
	dbPerms := make([]UserRBACPermission, 0)
	dbRoles := make([]UserOrgRole, 0)
	err := s.sql.WithDbSession(ctx, func(sess *db.Session) error {
		// Find permissions
		q := `
		SELECT
			user_id,
			action,
			scope
		FROM (
			SELECT ur.user_id, ur.org_id, p.action, p.scope
				FROM permission AS p
				INNER JOIN user_role AS ur on ur.role_id = p.role_id
			UNION
				SELECT tm.user_id, tr.org_id, p.action, p.scope
					FROM permission AS p
					INNER JOIN team_role AS tr ON tr.role_id = p.role_id
					INNER JOIN team_member AS tm ON tm.team_id = tr.team_id
			UNION
				SELECT ou.user_id, br.org_id, p.action, p.scope
					FROM permission AS p
					INNER JOIN builtin_role AS br ON br.role_id = p.role_id
					INNER JOIN org_user AS ou ON ou.role = br.role
			UNION
				SELECT sa.user_id, br.org_id, p.action, p.scope
					FROM permission AS p
					INNER JOIN builtin_role AS br ON br.role_id = p.role_id
					INNER JOIN (
						SELECT user.id AS user_id
						FROM user WHERE user.is_admin
					) AS sa ON 1 = 1 
					WHERE br.role = ?
		) AS up
		WHERE (org_id = ? OR org_id = ?) AND action LIKE ?
		`

		if err := sess.SQL(q, accesscontrol.RoleGrafanaAdmin, accesscontrol.GlobalOrgID, orgID, actionPrefix+"%").
			Find(&dbPerms); err != nil {
			return err
		}

		// Find roles
		q = `
		SELECT u.id, ou.role, u.is_admin
		FROM user AS u 
		LEFT JOIN org_user AS ou ON u.id = ou.user_id
		WHERE u.is_admin OR ou.org_id = ?
		`

		if err := sess.SQL(q, orgID).Find(&dbRoles); err != nil {
			return err
		}
		return nil
	})

	mapped := map[int64][]accesscontrol.Permission{}
	for i := range dbPerms {
		mapped[dbPerms[i].UserID] = append(mapped[dbPerms[i].UserID], accesscontrol.Permission{Action: dbPerms[i].Action, Scope: dbPerms[i].Scope})
	}

	roles := map[int64][]string{}
	for i := range dbRoles {
		if dbRoles[i].OrgRole != "" {
			roles[dbRoles[i].UserID] = []string{dbRoles[i].OrgRole}
		}
		if dbRoles[i].IsAdmin {
			roles[dbRoles[i].UserID] = append(roles[dbRoles[i].UserID], accesscontrol.RoleGrafanaAdmin)
		}
	}

	return mapped, roles, err
}

func (s *AccessControlStore) DeleteUserPermissions(ctx context.Context, orgID, userID int64) error {
	err := s.sql.WithDbSession(ctx, func(sess *db.Session) error {
		roleDeleteQuery := "DELETE FROM user_role WHERE user_id = ?"
		roleDeleteParams := []interface{}{roleDeleteQuery, userID}
		if orgID != accesscontrol.GlobalOrgID {
			roleDeleteQuery += " AND org_id = ?"
			roleDeleteParams = []interface{}{roleDeleteQuery, userID, orgID}
		}

		// Delete user role assignments
		if _, err := sess.Exec(roleDeleteParams...); err != nil {
			return err
		}

		// only delete scopes to user if all permissions is removed (i.e. user is removed)
		if orgID == accesscontrol.GlobalOrgID {
			// Delete permissions that are scoped to user
			if _, err := sess.Exec("DELETE FROM permission WHERE scope = ?", accesscontrol.Scope("users", "id", strconv.FormatInt(userID, 10))); err != nil {
				return err
			}
		}

		roleQuery := "SELECT id FROM role WHERE name = ?"
		roleParams := []interface{}{accesscontrol.ManagedUserRoleName(userID)}
		if orgID != accesscontrol.GlobalOrgID {
			roleQuery += " AND org_id = ?"
			roleParams = []interface{}{accesscontrol.ManagedUserRoleName(userID), orgID}
		}

		var roleIDs []int64
		if err := sess.SQL(roleQuery, roleParams...).Find(&roleIDs); err != nil {
			return err
		}

		if len(roleIDs) == 0 {
			return nil
		}

		permissionDeleteQuery := "DELETE FROM permission WHERE role_id IN(? " + strings.Repeat(",?", len(roleIDs)-1) + ")"
		permissionDeleteParams := make([]interface{}, 0, len(roleIDs)+1)
		permissionDeleteParams = append(permissionDeleteParams, permissionDeleteQuery)
		for _, id := range roleIDs {
			permissionDeleteParams = append(permissionDeleteParams, id)
		}

		// Delete managed user permissions
		if _, err := sess.Exec(permissionDeleteParams...); err != nil {
			return err
		}

		managedRoleDeleteQuery := "DELETE FROM role WHERE id IN(? " + strings.Repeat(",?", len(roleIDs)-1) + ")"
		managedRoleDeleteParams := []interface{}{managedRoleDeleteQuery}
		for _, id := range roleIDs {
			managedRoleDeleteParams = append(managedRoleDeleteParams, id)
		}
		// Delete managed user roles
		if _, err := sess.Exec(managedRoleDeleteParams...); err != nil {
			return err
		}

		return nil
	})
	return err
}
