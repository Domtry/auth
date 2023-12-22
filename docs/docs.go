// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/users": {
            "get": {
                "description": "Récupère la liste de tous les utilisateurs.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Récupère tous les utilisateurs",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.User"
                            }
                        }
                    }
                }
            },
            "post": {
                "description": "Crée un nouvel utilisateur avec les détails fournis.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Crée un nouvel utilisateur",
                "parameters": [
                    {
                        "description": "Détails de l'utilisateur",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.User"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/model.User"
                        }
                    }
                }
            }
        },
        "/users/{id}": {
            "get": {
                "description": "Récupère un utilisateur en fonction de son ID.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Récupère un utilisateur par ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID de l'utilisateur",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.User"
                        }
                    }
                }
            },
            "put": {
                "description": "Met à jour un utilisateur en fonction de son ID avec les détails fournis.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Met à jour un utilisateur",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID de l'utilisateur",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Détails de l'utilisateur à mettre à jour",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.User"
                        }
                    }
                }
            },
            "delete": {
                "description": "Supprime un utilisateur en fonction de son ID.",
                "tags": [
                    "users"
                ],
                "summary": "Supprime un utilisateur",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID de l'utilisateur",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "Aucun contenu"
                    }
                }
            }
        },
        "/users/{id}/assign-role/{role}": {
            "post": {
                "description": "Assigner un rôle spécifié à un utilisateur en fonction de son ID.",
                "tags": [
                    "users"
                ],
                "summary": "Assigner un rôle à un utilisateur",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID de l'utilisateur",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Rôle à assigner",
                        "name": "role",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "Aucun contenu"
                    }
                }
            }
        },
        "/users/{id}/remove-role/{role}": {
            "post": {
                "description": "Supprimer un rôle spécifié d'un utilisateur en fonction de son ID.",
                "tags": [
                    "users"
                ],
                "summary": "Supprimer un rôle d'un utilisateur",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID de l'utilisateur",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Rôle à supprimer",
                        "name": "role",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "Aucun contenu"
                    }
                }
            }
        }
    },
    "definitions": {
        "model.User": {
            "type": "object",
            "required": [
                "email",
                "name",
                "password",
                "role_id",
                "sername",
                "username"
            ],
            "properties": {
                "create_at": {
                    "type": "string"
                },
                "deleted_at": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "is_visible": {
                    "type": "boolean"
                },
                "name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "role_id": {
                    "type": "string"
                },
                "sername": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "BasicAuth": {
            "type": "basic"
        }
    },
    "externalDocs": {
        "description": "OpenAPI",
        "url": "https://swagger.io/resources/open-api/"
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0.0",
	Host:             "localhost:8000",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "CMagic Auth",
	Description:      "Service d'authentification.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}