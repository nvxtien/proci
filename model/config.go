package model

import "github.com/proci/util"

type Configuration struct {
	MumberOfOrg   int              `json:"numberOfOrg"`
	OrdererType   util.OrdererType `json:"ordererType"`
	NumberOfKafka int              `json:"numberOfKafka"`
	Profile       string           `json:"profile"`
}