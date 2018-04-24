package notification

//Notifier a notifier is able to send a message
type Notifier interface {
	//Send send a message to something
	Send(message interface{}) error
}
