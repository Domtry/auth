package model

type Role struct {
	AbstractModel
	Name             string           `json:"name" gorm:"size:80; index"`
	Describe         string           `json:"describe" gorm:"size:500"`
	Users            []User           `json:"users" gorm:"foreignKey:RoleId"`
	RolesPermissions []RolePermission `json:"roles_permissions" gorm:"foreignKey:RoleId"`
}
