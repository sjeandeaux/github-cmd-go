package http

import gohttp "net/http"

const (
	//ContentType Content-type Header
	ContentType = "Content-Type"
	//Authorization Authorization Header
	Authorization = "Authorization"
	//Bearer before token the space is already here.
	Bearer = "Bearer "
	//ApplicationJSON application json value header
	ApplicationJSON = "application/json"
)

//Close the body response
func Close(resp *gohttp.Response) {
	if resp != nil && resp.Body != nil {
		resp.Body.Close()
	}
}
