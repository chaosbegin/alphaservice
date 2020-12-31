package models

type HelpToc struct {
	Id       int        `json:"id"`
	Name     string     `json:"name"`
	Key      string     `json:"key"`
	Children []*HelpToc `json:"children"`
}
