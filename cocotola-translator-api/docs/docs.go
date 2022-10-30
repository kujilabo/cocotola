// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/swaggo/swag"
)

var doc = `{
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
        "/v1/admin/find": {
            "post": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "description": "find translations with first letter",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "translator"
                ],
                "summary": "find translations with first letter",
                "parameters": [
                    {
                        "description": "parameter to find translations",
                        "name": "param",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entity.TranslationFindParameterHTTPEntity"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.TranslationFindResponseHTTPEntity"
                        }
                    },
                    "400": {
                        "description": ""
                    },
                    "401": {
                        "description": ""
                    }
                }
            }
        },
        "/v1/admin/text/{text}": {
            "get": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "description": "find translations with text",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "translator"
                ],
                "summary": "find translations with text",
                "parameters": [
                    {
                        "type": "string",
                        "description": "text",
                        "name": "text",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.TranslationFindResponseHTTPEntity"
                        }
                    },
                    "400": {
                        "description": ""
                    },
                    "401": {
                        "description": ""
                    }
                }
            }
        },
        "/v1/admin/text/{text}/pos/{pos}": {
            "get": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "description": "find translations with text and pos",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "translator"
                ],
                "summary": "find translations with text and pos",
                "parameters": [
                    {
                        "type": "string",
                        "description": "text",
                        "name": "text",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "pos",
                        "name": "pos",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.TranslationHTTPEntity"
                        }
                    },
                    "400": {
                        "description": ""
                    },
                    "401": {
                        "description": ""
                    }
                }
            }
        },
        "/v1/user/dictionary/lookup": {
            "get": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "description": "dictionary lookup",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "translator"
                ],
                "summary": "dictionary lookup",
                "parameters": [
                    {
                        "type": "string",
                        "description": "text",
                        "name": "text",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "pos",
                        "name": "pos",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.TranslationFindResponseHTTPEntity"
                        }
                    },
                    "400": {
                        "description": ""
                    },
                    "401": {
                        "description": ""
                    }
                }
            }
        }
    },
    "definitions": {
        "entity.TranslationFindParameterHTTPEntity": {
            "type": "object",
            "properties": {
                "letter": {
                    "type": "string"
                }
            }
        },
        "entity.TranslationFindResponseHTTPEntity": {
            "type": "object",
            "properties": {
                "results": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entity.TranslationHTTPEntity"
                    }
                }
            }
        },
        "entity.TranslationHTTPEntity": {
            "type": "object",
            "properties": {
                "lang2": {
                    "type": "string"
                },
                "pos": {
                    "type": "integer"
                },
                "provider": {
                    "type": "string"
                },
                "text": {
                    "type": "string"
                },
                "translated": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "BasicAuth": {
            "type": "basic"
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
	Version:     "",
	Host:        "",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "",
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
		"escape": func(v interface{}) string {
			// escape tabs
			str := strings.Replace(v.(string), "\t", "\\t", -1)
			// replace " with \", and if that results in \\", replace that with \\\"
			str = strings.Replace(str, "\"", "\\\"", -1)
			return strings.Replace(str, "\\\\\"", "\\\\\\\"", -1)
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
	swag.Register("swagger", &s{})
}