package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type JsonRoleItem struct {
	Name        string   `json:"name"`
	Describe    string   `json:"describe"`
	Permissions []string `json:"permissions"`
}

type JsonRoleData struct {
	Roles []JsonRoleItem `json:"roles"`
}

func LoadPermissionByRoleName(roleName string) JsonRoleItem {

	var roles JsonRoleData
	jsonFile, err := os.Open("./ressources/config.json")
	if err != nil {
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &roles)
	if err != nil {
	}

	var selectedRole JsonRoleItem

	for i := 0; i < len(roles.Roles); i++ {
		item := roles.Roles[i]
		if item.Name == roleName {
			selectedRole = item
			break
		}
	}

	if selectedRole.Name == "" {
	}

	return selectedRole
}
