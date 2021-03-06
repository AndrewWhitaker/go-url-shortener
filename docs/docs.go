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
        "/shorturls": {
            "get": {
                "description": "List all short URLs",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "shorturls"
                ],
                "summary": "List all short URLs",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.ShortUrlReadFields"
                            }
                        }
                    },
                    "500": {
                        "description": ""
                    }
                }
            },
            "post": {
                "description": "Create a new short url. Users may specify a slug and an expiration date. If a slug is not supplied, an 8 character slug will automatically be generated for the short url.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "shorturls"
                ],
                "summary": "Create a new short url",
                "parameters": [
                    {
                        "description": "New short URL parameters",
                        "name": "shorturl",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.ShortUrlCreateFields"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.ShortUrlReadFields"
                        }
                    },
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.ShortUrlReadFields"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/e.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": ""
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/e.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": ""
                    }
                }
            }
        },
        "/shorturls/{slug}": {
            "get": {
                "description": "Get information about an existing short URL",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "shorturls"
                ],
                "summary": "Get information about an existing short URL",
                "parameters": [
                    {
                        "type": "string",
                        "description": "slug of short URL to get information about",
                        "name": "slug",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.ShortUrlReadFields"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/e.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": ""
                    }
                }
            },
            "delete": {
                "description": "Delete an existing short URL by supplying the slug.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "shorturls"
                ],
                "summary": "Delete an existing short URL",
                "parameters": [
                    {
                        "type": "string",
                        "description": "slug of short URL to delete",
                        "name": "slug",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": ""
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/e.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": ""
                    }
                }
            }
        },
        "/shorturls/{slug}/clicks": {
            "get": {
                "description": "Get clicks (statistics) for a short URL. Time periods of all time, 24 hours, and 1 week are permitted.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "shorturls"
                ],
                "summary": "Get clicks for a short URL",
                "parameters": [
                    {
                        "type": "string",
                        "description": "slug of short URL to retrieve statistics for",
                        "name": "slug",
                        "in": "path",
                        "required": true
                    },
                    {
                        "enum": [
                            "24_HOURS",
                            "1_WEEK",
                            "ALL_TIME"
                        ],
                        "type": "string",
                        "description": "time period to retrieve statistics for",
                        "name": "time_period",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/clicks.GetShortUrlClicksResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/e.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": ""
                    }
                }
            }
        }
    },
    "definitions": {
        "clicks.GetShortUrlClicksResponse": {
            "type": "object",
            "properties": {
                "count": {
                    "type": "integer"
                },
                "time_period": {
                    "type": "string"
                }
            }
        },
        "e.ErrorResponse": {
            "type": "object",
            "properties": {
                "errors": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/e.ValidationError"
                    }
                }
            }
        },
        "e.ValidationError": {
            "type": "object",
            "properties": {
                "field": {
                    "type": "string"
                },
                "reason": {
                    "type": "string"
                }
            }
        },
        "models.ShortUrlCreateFields": {
            "type": "object",
            "required": [
                "long_url"
            ],
            "properties": {
                "expires_on": {
                    "type": "string",
                    "format": "dateTime",
                    "example": "2023-01-01T16:30:00Z"
                },
                "long_url": {
                    "type": "string",
                    "format": "url",
                    "example": "http://www.google.com"
                },
                "slug": {
                    "type": "string",
                    "example": "myslug"
                }
            }
        },
        "models.ShortUrlReadFields": {
            "type": "object",
            "required": [
                "long_url"
            ],
            "properties": {
                "created_at": {
                    "type": "string",
                    "format": "dateTime",
                    "example": "2022-05-11T11:30:00Z"
                },
                "expires_on": {
                    "type": "string",
                    "format": "dateTime",
                    "example": "2023-01-01T16:30:00Z"
                },
                "long_url": {
                    "type": "string",
                    "format": "url",
                    "example": "http://www.google.com"
                },
                "slug": {
                    "type": "string",
                    "example": "myslug"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "URL Shortener",
	Description:      "A basic URL shortener",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
