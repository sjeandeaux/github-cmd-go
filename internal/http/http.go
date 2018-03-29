package http

import gohttp "net/http"

//Close the body response
func Close(resp *gohttp.Response) {
	if resp != nil && resp.Body != nil {
		resp.Body.Close()
	}
}
