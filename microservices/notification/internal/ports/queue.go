package ports

type Queue interface {
	Consume(queueName string, callback func(args []byte) error) error
}
