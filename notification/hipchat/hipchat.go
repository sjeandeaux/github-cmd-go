package hipchat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	internalhttp "github.com/sjeandeaux/toolators/internal/http"
	"github.com/sjeandeaux/toolators/notification"
)

//Notifier send a notificiation
type Notifier struct {
	token      string
	url        string
	httpClient *http.Client
}

// Payload https://www.hipchat.com/docs/apiv2/method/send_room_notification
type Payload struct {
	// From see the documentation hipchat
	From string `json:"from"`
	// Notify see the documentation hipchat
	Notify bool `json:"notify"`
	// Message see the documentation hipchat
	Message string `json:"message"`
	// Color see the documentation hipchat
	Color string `json:"color"`
}

var _ notification.Notifier = &Notifier{}

//NewNotifier create a notifier hipchat
func NewNotifier(url, token string) notification.Notifier {
	return &Notifier{
		httpClient: http.DefaultClient,
		url:        url,
		token:      token,
	}
}

//Send send a msg in room hipchat.
func (n *Notifier) Send(message interface{}) error {
	jsonValue, _ := json.Marshal(message)
	println(string(jsonValue))
	request, _ := http.NewRequest(http.MethodPost, n.url, bytes.NewBuffer(jsonValue))
	request.Header.Add(internalhttp.ContentType, internalhttp.ApplicationJSON)
	request.Header.Add(internalhttp.Authorization, internalhttp.Bearer+n.token)
	resp, err := n.httpClient.Do(request)
	defer internalhttp.Close(resp)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNoContent {
		b, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("failed %d %s", resp.StatusCode, b)
	}
	return nil
}
