package model

import "github.com/proci"

type Configuration struct {
	MumberOfOrg		int `json:"numberOfOrg"`
	OrdererType		proci.OrdererType `json:"ordererType"`
	NumberOfKafka	int		`json:"numberOfKafka"`
	Profile 		string 	`json:"profile"`
}