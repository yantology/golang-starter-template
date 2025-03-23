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
        "/auth/forget-password": {
            "post": {
                "description": "Reset user password using activation code",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Reset password",
                "parameters": [
                    {
                        "description": "Password reset details",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.ForgetPasswordRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success response with message",
                        "schema": {
                            "$ref": "#/definitions/auth.MessageResponse"
                        }
                    },
                    "400": {
                        "description": "Bad request response",
                        "schema": {
                            "$ref": "#/definitions/auth.MessageResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized response",
                        "schema": {
                            "$ref": "#/definitions/auth.MessageResponse"
                        }
                    }
                }
            }
        },
        "/auth/login": {
            "post": {
                "description": "Authenticate user and return JWT tokens",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "User login",
                "parameters": [
                    {
                        "description": "Login credentials",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success response with JWT tokens",
                        "schema": {
                            "$ref": "#/definitions/auth.DataResponse-auth_jwtResponseData"
                        }
                    },
                    "400": {
                        "description": "Bad request response",
                        "schema": {
                            "$ref": "#/definitions/auth.MessageResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized response",
                        "schema": {
                            "$ref": "#/definitions/auth.MessageResponse"
                        }
                    }
                }
            }
        },
        "/auth/refresh-token": {
            "post": {
                "description": "Get new access token using refresh token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Refresh token",
                "parameters": [
                    {
                        "description": "Refresh token details",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.RefreshTokenRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success response with new JWT tokens",
                        "schema": {
                            "$ref": "#/definitions/auth.DataResponse-auth_jwtResponseData"
                        }
                    },
                    "400": {
                        "description": "Bad request response",
                        "schema": {
                            "$ref": "#/definitions/auth.MessageResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized response",
                        "schema": {
                            "$ref": "#/definitions/auth.MessageResponse"
                        }
                    }
                }
            }
        },
        "/auth/register": {
            "post": {
                "description": "Register a new user with activation code",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Register new user",
                "parameters": [
                    {
                        "description": "Registration details",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.RegisterRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Success response with message",
                        "schema": {
                            "$ref": "#/definitions/auth.MessageResponse"
                        }
                    },
                    "400": {
                        "description": "Bad request response",
                        "schema": {
                            "$ref": "#/definitions/auth.MessageResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized response",
                        "schema": {
                            "$ref": "#/definitions/auth.MessageResponse"
                        }
                    }
                }
            }
        },
        "/auth/token/{type}": {
            "post": {
                "description": "Request a token for registration or password reset",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Request activation token",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Token type (registration or forget-password)",
                        "name": "type",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Token request parameters",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.TokenRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success response with message",
                        "schema": {
                            "$ref": "#/definitions/auth.MessageResponse"
                        }
                    },
                    "400": {
                        "description": "Bad request response",
                        "schema": {
                            "$ref": "#/definitions/auth.MessageResponse"
                        }
                    },
                    "404": {
                        "description": "Not found response",
                        "schema": {
                            "$ref": "#/definitions/auth.MessageResponse"
                        }
                    },
                    "409": {
                        "description": "Conflict response",
                        "schema": {
                            "$ref": "#/definitions/auth.MessageResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "auth.DataResponse-auth_jwtResponseData": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/auth.jwtResponseData"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "auth.ForgetPasswordRequest": {
            "type": "object",
            "required": [
                "activation_code",
                "email",
                "new_password",
                "new_password_confirmation"
            ],
            "properties": {
                "activation_code": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "new_password": {
                    "type": "string"
                },
                "new_password_confirmation": {
                    "type": "string"
                }
            }
        },
        "auth.LoginRequest": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "auth.MessageResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "auth.RefreshTokenRequest": {
            "type": "object",
            "required": [
                "refresh_token"
            ],
            "properties": {
                "refresh_token": {
                    "type": "string"
                }
            }
        },
        "auth.RegisterRequest": {
            "type": "object",
            "required": [
                "activation_code",
                "email",
                "fullname",
                "password",
                "password_confirmation"
            ],
            "properties": {
                "activation_code": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "fullname": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "password_confirmation": {
                    "type": "string"
                }
            }
        },
        "auth.TokenRequest": {
            "type": "object",
            "required": [
                "email"
            ],
            "properties": {
                "email": {
                    "type": "string"
                }
            }
        },
        "auth.jwtResponseData": {
            "type": "object",
            "properties": {
                "access_token": {
                    "type": "string"
                },
                "expires_in": {
                    "type": "integer"
                },
                "refresh_token": {
                    "type": "string"
                },
                "token_type": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8000",
	BasePath:         "/api/v1",
	Schemes:          []string{"http", "https"},
	Title:            "Retail Pro API",
	Description:      "This is a retail management system server.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
