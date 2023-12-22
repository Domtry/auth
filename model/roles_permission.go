package model

type RolePermission struct {
	AbstractModel
	RoleId       string `json:"role_id" gorm:"size:120; index"`
	PermissionId string `json:"permission_id" gorm:"size:120; index"`
	Describe     string `json:"describe" gorm:"size:500"`
}
