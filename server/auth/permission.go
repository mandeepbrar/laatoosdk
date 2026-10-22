package auth

// PermissionDeclaration is one resource access permission a plugin declares: the identifier a
// channel or service guard names and a role permission lists.
//
// The declaration is read from the plugin's `permissions:` map (identifier to name) and its
// `externalpermissions:` map; the server owns the catalogue of them.
type PermissionDeclaration struct {
	// Id is the stable identifier a guard names, for example "students.read". It is never localized.
	Id string `json:"Id"`

	// Name is what a person reads: a key into the declaring plugin's properties, or literal text when
	// no such key exists.
	Name string `json:"Name"`

	// Plugin is the plugin that authors the permission. For a permission a plugin borrows from a
	// dependency this is the dependency, not the borrower.
	Plugin string `json:"Plugin"`
}

// RolePermission is a named bundle of resource access permissions and service access rows, authored
// by a namespace. A role lists role permissions; the authorization back end that defines a role
// permission is the one that enforces it.
type RolePermission struct {
	// Id is the role permission's identifier, unique across the namespace tree. It is never localized.
	Id string `json:"Id"`

	// Name is what a person reads: literal text, or a property key a user interface resolves.
	Name string `json:"Name"`

	// Namespace is the address of the namespace whose back end defines the role permission.
	Namespace string `json:"Namespace"`
}

// RolePermissionHolder is implemented by a Role that carries role-permission ids. A role built in
// role management is a list of role permissions and nothing else; the back end reads the list from
// here when it loads the role.
type RolePermissionHolder interface {
	Role

	// GetRolePermissions returns the ids of the role permissions the role holds. A nil slice is an
	// empty set.
	GetRolePermissions() []string

	// SetRolePermissions replaces the role's role-permission ids.
	SetRolePermissions([]string)
}
