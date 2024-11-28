package pipeline

type Source interface {
	Read() (*Message, error)
}

type Sink interface {
	Write(data *Message) error
}

type Stage interface {
	Process(data *Message) (*Message, error)
}
