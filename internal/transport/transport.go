package transport

type Channels struct {
	NavigationCh chan int
	InputCh      chan string
}

type Core interface {
	GetPasswordHidden() (string, error)
	StartInputScanner() error
	SendMessageToUser(message string)
	GetChannels() *Channels
	SwitchFocus(b bool)
}
