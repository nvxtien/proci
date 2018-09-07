package proci

type OrdererType string

const (
	Kafka OrdererType	= "kafka"
	Solo OrdererType	= "solo"
)