package littlepipe

type Plugin interface {
	Name() string
	Init(config interface{}) error
	Run() error
	Stop()
}

type InputPlugin interface {
	Plugin
	InChan() chan *Pack
}
