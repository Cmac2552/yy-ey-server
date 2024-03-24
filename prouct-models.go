package main

type productTypeBody struct {
	ProductTypeName string `json:"productTypeName"`
}

type productValue struct {
	AttributeName string `json:"attributeName"`
	Value         string `json:"value"`
}

type productTypeAndAttributes struct {
	ProductTypeName string   `json:"productTypeName"`
	Attributes      []string `json:"attributes"`
}

type productTypeAndAttributeValues struct {
	ProductTypeName string            `json:"productTypeName"`
	AttrValues      map[string]string `json:"attrValues"`
}

type productTypeAndNumberWithAttributeValues struct {
	ProductTypeName string            `json:"productTypeName"`
	ProductNumber   int               `json:"productNumber"`
	AttrValues      map[string]string `json:"attrValues"`
}
