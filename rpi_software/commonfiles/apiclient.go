package commonfiles

import "net/http"

type api_client struct {
	api_key    string
	httpclient *http.Client
}
