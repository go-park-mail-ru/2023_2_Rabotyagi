// Package docs Code generated by swaggo/swag. DO NOT EDIT
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
        "/logout": {
            "post": {
                "description": "logout in app",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "logout",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    },
                    "405": {
                        "description": "Method Not Allowed",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/post/add": {
            "post": {
                "description": "add post by data",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "add post",
                "parameters": [
                    {
                        "description": "post data for adding",
                        "name": "post",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/storage.PrePost"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    },
                    "405": {
                        "description": "Method Not Allowed",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/post/get/{id}": {
            "get": {
                "description": "get post by id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "get post",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "post id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    },
                    "405": {
                        "description": "Method Not Allowed",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/post/get_list": {
            "get": {
                "description": "get posts by count",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "get posts",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "count posts",
                        "name": "count",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    },
                    "405": {
                        "description": "Method Not Allowed",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/signin": {
            "post": {
                "description": "signin in app",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "signin",
                "parameters": [
                    {
                        "description": "user data for signin",
                        "name": "preUser",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/storage.PreUser"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    },
                    "405": {
                        "description": "Method Not Allowed",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/signup": {
            "post": {
                "description": "signup in app",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "signup",
                "parameters": [
                    {
                        "description": "user data for signup",
                        "name": "preUser",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/storage.PreUser"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler.Response"
                        }
                    },
                    "222": {
                        "description": "error",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    },
                    "405": {
                        "description": "Method Not Allowed",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "handler.ErrorResponse": {
            "type": "object",
            "properties": {
                "body": {
                    "$ref": "#/definitions/handler.ResponseBodyError"
                },
                "status": {
                    "type": "integer"
                }
            }
        },
        "handler.PostsListResponse": {
            "type": "object",
            "properties": {
                "body": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/storage.Post"
                    }
                },
                "status": {
                    "type": "integer"
                }
            }
        },
        "handler.Response": {
            "type": "object",
            "properties": {
                "body": {
                    "$ref": "#/definitions/handler.ResponseBody"
                },
                "status": {
                    "type": "integer"
                }
            }
        },
        "handler.ResponseBody": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "handler.ResponseBodyError": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "storage.Post": {
            "type": "object",
            "properties": {
                "author": {
                    "type": "integer"
                },
                "city": {
                    "type": "string"
                },
                "delivery": {
                    "type": "boolean"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "price": {
                    "type": "integer"
                },
                "safe": {
                    "type": "boolean"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "storage.PrePost": {
            "type": "object",
            "properties": {
                "author": {
                    "type": "integer"
                },
                "city": {
                    "type": "string"
                },
                "delivery": {
                    "type": "boolean"
                },
                "description": {
                    "type": "string"
                },
                "price": {
                    "type": "integer"
                },
                "safe": {
                    "type": "boolean"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "storage.PreUser": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "84.23.53.28:8080",
	BasePath:         "/api/v1",
	Schemes:          []string{"http"},
	Title:            "YULA project API",
	Description:      "This is a server of YULA server.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
