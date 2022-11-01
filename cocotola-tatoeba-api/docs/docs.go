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
        "/v1/admin/link/import": {
            "post": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "description": "import links",
                "tags": [
                    "tatoeba"
                ],
                "summary": "import links",
                "parameters": [
                    {
                        "type": "file",
                        "description": "links.csv",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    },
                    "400": {
                        "description": ""
                    },
                    "401": {
                        "description": ""
                    },
                    "500": {
                        "description": ""
                    }
                }
            }
        },
        "/v1/admin/sentence/import": {
            "post": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "description": "import sentences",
                "tags": [
                    "tatoeba"
                ],
                "summary": "import sentences",
                "parameters": [
                    {
                        "type": "file",
                        "description": "***_sentences_detailed.tsv",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    },
                    "400": {
                        "description": ""
                    },
                    "401": {
                        "description": ""
                    },
                    "500": {
                        "description": ""
                    }
                }
            }
        },
        "/v1/user/sentence/{sentenceNumber}": {
            "get": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "description": "import links",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tatoeba"
                ],
                "summary": "import links",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Sentence number",
                        "name": "sentenceNumber",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.TatoebaSentenceResponse"
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
        "/v1/user/sentence_pair/find": {
            "post": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "description": "import links",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tatoeba"
                ],
                "summary": "import links",
                "parameters": [
                    {
                        "description": "parameter to find sentences",
                        "name": "param",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entity.TatoebaSentenceFindParameter"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.TatoebaSentencePairFindResponse"
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
        "entity.TatoebaSentenceFindParameter": {
            "type": "object",
            "required": [
                "pageNo",
                "pageSize"
            ],
            "properties": {
                "keyword": {
                    "type": "string"
                },
                "pageNo": {
                    "type": "integer",
                    "minimum": 1
                },
                "pageSize": {
                    "type": "integer",
                    "minimum": 1
                },
                "random": {
                    "type": "boolean"
                }
            }
        },
        "entity.TatoebaSentencePair": {
            "type": "object",
            "properties": {
                "dst": {
                    "$ref": "#/definitions/entity.TatoebaSentenceResponse"
                },
                "src": {
                    "$ref": "#/definitions/entity.TatoebaSentenceResponse"
                }
            }
        },
        "entity.TatoebaSentencePairFindResponse": {
            "type": "object",
            "properties": {
                "results": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entity.TatoebaSentencePair"
                    }
                },
                "totalCount": {
                    "type": "integer"
                }
            }
        },
        "entity.TatoebaSentenceResponse": {
            "type": "object",
            "properties": {
                "author": {
                    "type": "string"
                },
                "lang2": {
                    "type": "string",
                    "enum": [
                        "ja",
                        "en"
                    ]
                },
                "sentenceNumber": {
                    "type": "integer"
                },
                "text": {
                    "type": "string"
                },
                "updatedAt": {
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
