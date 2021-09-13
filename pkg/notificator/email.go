package notificator

type Email struct {
	Subject   string
	Recipient string
	Sender    string
	Data      map[string]interface{}
}
