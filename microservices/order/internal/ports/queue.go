package ports

type Queue interface {
	Publish(message string) error
	Consume(queueName string, callback func(args []byte) error) error
}
