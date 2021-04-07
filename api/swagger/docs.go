// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package swagger

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/post": {
            "get": {
                "description": "Get all published with pagination",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Post"
                ],
                "summary": "Get all published",
                "operationId": "post-get-all-published",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Number of current page",
                        "name": "current_page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Number of posts count",
                        "name": "count_per_page",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Posts with user id",
                        "name": "user_id",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.PostPaginationResponseDto"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponseDto"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Create new post",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Post"
                ],
                "summary": "Create post",
                "operationId": "post-create",
                "parameters": [
                    {
                        "description": "Post details",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.CreatePostRequestDto"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/response.PostResponseDto"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponseDto"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponseDto"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponseDto"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponseDto"
                        }
                    }
                }
            }
        },
        "/post/self": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get all self post with pagination",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Post"
                ],
                "summary": "Get all self posts",
                "operationId": "post-get-all-self",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Number of current page",
                        "name": "current_page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Number of posts count",
                        "name": "count_per_page",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.PostPaginationResponseDto"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponseDto"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponseDto"
                        }
                    }
                }
            }
        },
        "/post/{id}": {
            "get": {
                "description": "Get single post by id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Post"
                ],
                "summary": "Get single post",
                "operationId": "post-get-single",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Post with id",
                        "name": "id",
                        "in": "path"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.PostResponseDto"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponseDto"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponseDto"
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Update post with id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Post"
                ],
                "summary": "Update post",
                "operationId": "post-update",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Post with id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Post details",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.UpdatePostRequestDto"
                        }
                    }
                ],
                "responses": {
                    "202": {
                        "description": "Accepted",
                        "schema": {
                            "$ref": "#/definitions/response.PostResponseDto"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponseDto"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponseDto"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponseDto"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponseDto"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Delete post with id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Post"
                ],
                "summary": "Delete post",
                "operationId": "post-delete",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Post with id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": ""
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponseDto"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponseDto"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponseDto"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponseDto"
                        }
                    }
                }
            }
        },
        "/post/{id}/publish": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Publish post with id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Post"
                ],
                "summary": "Publish post",
                "operationId": "post-publish",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Post with id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "202": {
                        "description": "Accepted",
                        "schema": {
                            "$ref": "#/definitions/response.PostResponseDto"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponseDto"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponseDto"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponseDto"
                        }
                    }
                }
            }
        },
        "/user/self": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get self user by authorized information",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Get self user",
                "operationId": "user-get-self",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.UserResponseDto"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponseDto"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponseDto"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponseDto"
                        }
                    }
                }
            }
        },
        "/user/sign-in": {
            "post": {
                "description": "Sign in with account details",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Sign in",
                "operationId": "user-sign-in",
                "parameters": [
                    {
                        "description": "Sign in with account details",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.SignInUserRequestDto"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.TokenResponseDto"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponseDto"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponseDto"
                        }
                    }
                }
            }
        },
        "/user/sign-up": {
            "post": {
                "description": "Sign up with account details",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Sign up",
                "operationId": "user-sign-up",
                "parameters": [
                    {
                        "description": "Sign up with account details",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.SignUpUserRequestDto"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/response.UserResponseDto"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponseDto"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponseDto"
                        }
                    }
                }
            }
        },
        "/user/{id}": {
            "get": {
                "description": "Get single user by id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Get single user",
                "operationId": "user-get-single",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User with id",
                        "name": "id",
                        "in": "path"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.UserResponseDto"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponseDto"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponseDto"
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Update user with id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Update user",
                "operationId": "user-update",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User with id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "User details",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.UpdateUserRequestDto"
                        }
                    }
                ],
                "responses": {
                    "202": {
                        "description": "Accepted",
                        "schema": {
                            "$ref": "#/definitions/response.UserResponseDto"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponseDto"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponseDto"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponseDto"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponseDto"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponseDto"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "request.CreatePostRequestDto": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "string"
                },
                "is_published": {
                    "type": "boolean"
                },
                "slug": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "request.SignInUserRequestDto": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "request.SignUpUserRequestDto": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "request.UpdatePostRequestDto": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "string"
                },
                "is_published": {
                    "type": "boolean"
                },
                "slug": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "request.UpdateUserRequestDto": {
            "type": "object",
            "properties": {
                "username": {
                    "type": "string"
                }
            }
        },
        "response.ErrorResponseDto": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "information": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "response.PaginationResponseDto": {
            "type": "object",
            "properties": {
                "count_per_page": {
                    "type": "integer"
                },
                "current_page": {
                    "type": "integer"
                },
                "next_page": {
                    "type": "integer"
                },
                "previous_page": {
                    "type": "integer"
                },
                "total": {
                    "type": "integer"
                }
            }
        },
        "response.PostPaginationResponseDto": {
            "type": "object",
            "properties": {
                "pagination": {
                    "$ref": "#/definitions/response.PaginationResponseDto"
                },
                "posts": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/response.PostResponseDto"
                    }
                }
            }
        },
        "response.PostResponseDto": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "deleted_at": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "is_deleted": {
                    "type": "boolean"
                },
                "is_published": {
                    "description": "User        *UserResponseDto ` + "`" + `json:\"user,omitempty\"` + "`" + `",
                    "type": "boolean"
                },
                "published_at": {
                    "type": "string"
                },
                "slug": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                },
                "user_id": {
                    "type": "string"
                }
            }
        },
        "response.TokenResponseDto": {
            "type": "object",
            "properties": {
                "access_token": {
                    "type": "string"
                }
            }
        },
        "response.UserResponseDto": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "id": {
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
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0.0",
	Host:        "",
	BasePath:    "/api/v1",
	Schemes:     []string{},
	Title:       "Go Simple Blog API",
	Description: "",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
