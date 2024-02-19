package main

type productTypeBody struct {
	ProductTypeName string `json:"productTypeName"`
}

type productValue struct {
	AttributeName string `json:"attributeName"`
	Value         string `json:"value"`
}
