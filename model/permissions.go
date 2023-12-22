package model

type Permission struct {
	AbstractModel
	Name             string           `json:"name" gorm:"size:80"`
	Describe         string           `json:"describe" gorm:"size:500"`
	RolesPermissions []RolePermission `json:"roles_permissions" gorm:"foreignkey:PermissionId"`
}
