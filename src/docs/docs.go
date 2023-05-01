// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/challenge/": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "summary": "Creates challenge based in a given user",
                "parameters": [
                    {
                        "description": "Challenge",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api_endpoints.Challenge"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api_endpoints.Challenge"
                        }
                    }
                }
            },
            "delete": {
                "produces": [
                    "application/json"
                ],
                "summary": "Deletes challenge in a user",
                "parameters": [
                    {
                        "description": "Challenge",
                        "name": "challenge",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api_endpoints.DelChallenge"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api_endpoints.DelRespChallenge"
                        }
                    }
                }
            }
        },
        "/kali/": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "summary": "Creates Kali based on given user",
                "parameters": [
                    {
                        "description": "User",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api_endpoints.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api_endpoints.Kali"
                        }
                    }
                }
            }
        },
        "/user/": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "Retrieves all challenges, as well as Kalis or wireguards running for a user",
                "parameters": [
                    {
                        "description": "User",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api_endpoints.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api_endpoints.UserInfo"
                        }
                    }
                }
            },
            "post": {
                "produces": [
                    "application/json"
                ],
                "summary": "Creates user based on given name",
                "parameters": [
                    {
                        "description": "User",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api_endpoints.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api_endpoints.User"
                        }
                    }
                }
            },
            "delete": {
                "produces": [
                    "application/json"
                ],
                "summary": "Deletes user based on given name",
                "parameters": [
                    {
                        "description": "User",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api_endpoints.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            }
        },
        "/users": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "Retrieves all users",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api_endpoints.Users"
                        }
                    }
                }
            }
        },
        "/wireguard/": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "summary": "Sends a public key and starts Wireguard",
                "parameters": [
                    {
                        "description": "Wireguard",
                        "name": "publicKey",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api_endpoints.Wireguard"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api_endpoints.ConfigFile"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "api_endpoints.Challenge": {
            "type": "object",
            "properties": {
                "challengeName": {
                    "type": "string"
                },
                "ports": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "user": {
                    "type": "string"
                }
            }
        },
        "api_endpoints.ConfigFile": {
            "type": "object",
            "properties": {
                "file": {
                    "type": "string"
                }
            }
        },
        "api_endpoints.DelChallenge": {
            "type": "object",
            "properties": {
                "challengeName": {
                    "type": "string"
                },
                "user": {
                    "type": "string"
                }
            }
        },
        "api_endpoints.DelRespChallenge": {
            "type": "object",
            "properties": {
                "challengeName": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                },
                "user": {
                    "type": "string"
                }
            }
        },
        "api_endpoints.Kali": {
            "type": "object",
            "properties": {
                "message": {
                    "description": "Message m\nin: string",
                    "type": "string"
                },
                "name": {
                    "description": "Name\nin: string",
                    "type": "string"
                }
            }
        },
        "api_endpoints.User": {
            "type": "object",
            "properties": {
                "name": {
                    "description": "Username\nin: string",
                    "type": "string"
                }
            }
        },
        "api_endpoints.UserInfo": {
            "type": "object",
            "properties": {
                "pods": {
                    "description": "UserInfo pods\nin: array",
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "api_endpoints.Users": {
            "type": "object",
            "properties": {
                "names": {
                    "description": "Users names\nin: array",
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "api_endpoints.Wireguard": {
            "type": "object",
            "properties": {
                "key": {
                    "type": "string"
                },
                "namespace": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
