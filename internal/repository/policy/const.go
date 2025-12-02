package policy

const (
	SuperAdminRole = 1
	MemberRole     = 2
)

const (
	queryCreateRole = `
		INSERT INTO roles (name, description)
		VALUES ($1, $2) RETURNING id;
	`
	queryCreatePermission = `
		INSERT INTO permissions (name, description)
		VALUES ($1, $2) RETURNING id;
	`
	queryCreateResource = `
		INSERT INTO resources ("type", method, path)
		VALUES ($1, $2, $3) RETURNING id;
	`
	queryBulkAssignResource = `
		INSERT INTO resource_permissions (resource_id, permission_id)
	`
	queryAssignRole = `
		INSERT INTO user_roles (user_id, role_id)
		VALUES ($1, $2);
	`
	queryBulkAssignPermission = `
		INSERT INTO role_permissions (role_id, permission_id)
	`
	queryGetRoleByUserID = `
		SELECT r.name, r.description
		FROM user_roles ur
		LEFT JOIN roles r ON ur.role_id = r.id
		WHERE ur.user_id = $1;
	`

	queryLoadPolicy = `
        SELECT 
            r.name AS role,
			res.type AS "type",
            res.path AS path,
            res.method AS method
        FROM role_permissions rp
        JOIN roles r ON r.id = rp.role_id
        JOIN resources_permissions resp ON resp.permission_id = rp.permission_id
        JOIN resources res ON res.id = resp.resource_id WHERE res.type = $1;
    `
)
