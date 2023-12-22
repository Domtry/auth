package config

import (
	"auth/model"
	"fmt"

	"github.com/google/uuid"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func GetDB(config DbConfig) (*gorm.DB, error) {

	dbDns := ""
	var db *gorm.DB
	var err error

	switch config.DbProvider {
	case "pg":
		dbDns = fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable",
			config.DbHost, config.DbUser, config.DbPassword, config.DbName, config.DbPort)
		db, err = gorm.Open(postgres.Open(dbDns), &gorm.Config{
			TranslateError: true,
		})
	case "mysql":
		dbDns = fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
			config.DbUser, config.DbPassword, config.DbHost, config.DbPort, config.DbName)
		db, err = gorm.Open(mysql.Open(dbDns), &gorm.Config{
			TranslateError: true,
		})
	case "sqlite":
		dbDns = fmt.Sprintf("%v.sqlite", config.DbName)
		db, err = gorm.Open(sqlite.Open(dbDns), &gorm.Config{})
	}

	return db, err
}

func CreateUpdateTable(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Permission{},
		&model.RolePermission{},
	)
}

func DropTable(db *gorm.DB) error {
	return db.Migrator().DropTable(
		&model.User{},
		&model.Role{},
		&model.Permission{},
		&model.RolePermission{},
	)
}

func InitSystem(db *gorm.DB) error {
	/* create basic permission */
	permissions := []model.Permission{
		{Name: "create_roles", Describe: "insert role table"},
		{Name: "update_roles", Describe: "update role"},
		{Name: "assign_role", Describe: "Assign new role to user"},

		//permission to search on database
		{Name: "filter_all", Describe: "filter all"},
		{Name: "find_all", Describe: "find all"},
		{Name: "get_all", Describe: "get all"},

		//user permissios
		{Name: "create_user", Describe: "insert user"},
		{Name: "update_user_profile", Describe: "update user profile"},
		{Name: "remove_user", Describe: "remove user"},
		{Name: "chang_password", Describe: "Chang user password"},
		{Name: "create_new_password", Describe: "Create new password"},
	}

	/*create system all roles*/
	roles := []model.Role{
		{Name: "admin", Describe: "Control all app"},
		{Name: "client", Describe: "Client role "},
		{Name: "manager", Describe: "Manager role"},
		{Name: "comptable", Describe: "Comptable role"},
	}

	var count int64
	if err := db.Model(model.Role{}).Count(&count).Error; err != nil {
	}

	if count == 0 {
		db.Create(roles)
		db.Create(permissions)
		db.Commit()
	}

	var adminRoleModel model.Role
	if err := db.Where("name = ?", "admin").Find(&adminRoleModel).Error; err != nil {
	}

	if _, err := uuid.Parse(adminRoleModel.Id); err != nil {
	}

	count = 0
	if err := db.Model(model.RolePermission{}).Count(&count).Error; err != nil {
	}

	if count == 0 {
		var permissionModel []model.Permission
		var rolePermissions []model.RolePermission

		db.Find(&permissionModel)

		for i := 0; i < len(permissionModel); i++ {
			currentItem := model.RolePermission{
				RoleId:       adminRoleModel.Id,
				PermissionId: permissionModel[i].Id,
				Describe:     "",
			}

			rolePermissions = append(rolePermissions, currentItem)
		}

		db.Create(rolePermissions)
		db.Commit()
	}

	if err := db.Model(model.User{}).Where("").Error; err != nil {
	}
	return nil
}
