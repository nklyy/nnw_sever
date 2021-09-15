package notificator

type Email struct {
	Subject   string
	Recipient string
	Sender    string
	Template  string
	Data      map[string]interface{}
}
