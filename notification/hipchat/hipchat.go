package hipchat

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	internalhttp "github.com/sjeandeaux/toolators/internal/http"
	"github.com/sjeandeaux/toolators/notification"
)

//URLRoom we need to replace hostname and room number.
const URLRoom = "https://%s/v2/room/%s/notification"

//Notifier send a notificiation
type Notifier struct {
	token      string
	url        string
	httpClient *http.Client
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
// https://www.hipchat.com/docs/apiv2/method/send_room_notification
func (n *Notifier) Send(message io.Reader) error {
	request, _ := http.NewRequest(http.MethodPost, n.url, message)
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
