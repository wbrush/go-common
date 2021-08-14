package datamodels

import "time"

var Now = time.Now //this is used for mocks

type Edge struct {
	Node   interface{} `json:"node"`
	Cursor string      `json:"cursor"`
}

type PageInfo struct {
	EndCursor   string `json:"endCursor"`
	HasNextPage bool   `json:"hasNextPage"`
}

type List struct {
	TotalCount int      `json:"totalCount"`
	Edges      []Edge   `json:"edges"`
	PageInfo   PageInfo `json:"pageInfo"`
}
