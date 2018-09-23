package util

type OrdererType string

const (
	Kafka OrdererType	= "kafka"
	Solo OrdererType	= "solo"
)