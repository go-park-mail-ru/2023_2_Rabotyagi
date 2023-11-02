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
                            "$ref": "#/definitions/delivery.Response"
                        }
                    },
                    "222": {
                        "description": "Error",
                        "schema": {
                            "$ref": "#/definitions/delivery.ErrorResponse"
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
        "/product/add": {
            "post": {
                "description": "add product by data",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "add product",
                "parameters": [
                    {
                        "description": "product data for adding",
                        "name": "product",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.PreProduct"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/delivery.Response"
                        }
                    },
                    "222": {
                        "description": "Error",
                        "schema": {
                            "$ref": "#/definitions/delivery.ErrorResponse"
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
        "/product/get/{id}": {
            "get": {
                "description": "get product by id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "get product",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "product id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/delivery.PostResponse"
                        }
                    },
                    "222": {
                        "description": "Error",
                        "schema": {
                            "$ref": "#/definitions/delivery.ErrorResponse"
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
        "/product/get_list": {
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
                            "$ref": "#/definitions/delivery.PostsListResponse"
                        }
                    },
                    "222": {
                        "description": "Error",
                        "schema": {
                            "$ref": "#/definitions/delivery.ErrorResponse"
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
                            "$ref": "#/definitions/models.UserWithoutID"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/delivery.Response"
                        }
                    },
                    "222": {
                        "description": "Error",
                        "schema": {
                            "$ref": "#/definitions/delivery.ErrorResponse"
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
                "description": "signup in app\nError.status can be:\nStatusErrBadRequest      = 400\nStatusErrInternalServer  = 500",
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
                            "$ref": "#/definitions/models.UserWithoutID"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/delivery.Response"
                        }
                    },
                    "222": {
                        "description": "Error",
                        "schema": {
                            "$ref": "#/definitions/delivery.ErrorResponse"
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
        "delivery.ErrorResponse": {
            "type": "object",
            "properties": {
                "body": {
                    "$ref": "#/definitions/delivery.ResponseBodyError"
                },
                "status": {
                    "type": "integer"
                }
            }
        },
        "delivery.PostResponse": {
            "type": "object",
            "properties": {
                "body": {
                    "$ref": "#/definitions/models.Product"
                },
                "status": {
                    "type": "integer"
                }
            }
        },
        "delivery.PostsListResponse": {
            "type": "object",
            "properties": {
                "body": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.ProductInFeed"
                    }
                },
                "status": {
                    "type": "integer"
                }
            }
        },
        "delivery.Response": {
            "type": "object",
            "properties": {
                "body": {
                    "$ref": "#/definitions/delivery.ResponseBody"
                },
                "status": {
                    "type": "integer"
                }
            }
        },
        "delivery.ResponseBody": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "delivery.ResponseBodyError": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "models.Images": {
            "type": "object",
            "properties": {
                "alt": {
                    "type": "string"
                },
                "url": {
                    "type": "string"
                }
            }
        },
        "models.PreProduct": {
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
                "image": {
                    "$ref": "#/definitions/models.Images"
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
        "models.Product": {
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
                "image": {
                    "$ref": "#/definitions/models.Images"
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
        "models.ProductInFeed": {
            "type": "object",
            "properties": {
                "city": {
                    "type": "string"
                },
                "delivery": {
                    "type": "boolean"
                },
                "id": {
                    "type": "integer"
                },
                "image": {
                    "$ref": "#/definitions/models.Images"
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
        "models.UserWithoutID": {
            "type": "object",
            "properties": {
                "birthday": {
                    "description": "nolint",
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "name": {
                    "description": "nolint",
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "phone": {
                    "description": "nolint",
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
