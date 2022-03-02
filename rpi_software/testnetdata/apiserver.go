package testnetdata

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	comfiles "github.com/jax-b/ModulerKTNE/rpi_software/commonfiles"
)

type api_server struct {
	config_api_keys []string
	server          *http.Server
	cstatusptr      *comfiles.Status
}

var loadedmod [][4]rune = {
	{'b', 't', 't', 'n'},
	{'s', 'm', 'o', 'n'},
	{' ', ' ', ' ', ' '},
	{' ', ' ', ' ', ' '},
	{' ', ' ', ' ', ' '},
	{' ', ' ', ' ', ' '},
	{' ', ' ', ' ', ' '},
	{' ', ' ', ' ', ' '},
	{' ', ' ', ' ', ' '},
	{' ', ' ', ' ', ' '}
}

var modsolvedstat [10]bool = {0,1,0,0,0,0,0,0,0,0}

func NewAPIServer(config *Config, Status *comfiles.Status) *apiapi_server {
	apiSRV := &api_server{
		config_api_keys: config.Network.api_keys,
		cstatusptr:      Status,
	}
	apiSRV.server = &http.Server{
		Addr:    ":" + config.Network.api_port,
		Handler: apiSRV.Handler,
	}
	return apiSRV
}

func (sapi *api_server) Handler(w http.ResponseWriter, r *http.Request) {
	apiurl := strings.Split(r.URL.Path, "/")
	switch r.Method {
	case "GET":
		switch apiurl[1] {
		case "status":
			outstat := &cstatusptr
			outstat.StrTime = fmt.Sprintf("%s", cstatusptr.Time)
			io.WriteString(w, json.Marshal(outstat))
		case "module":
			selectedModule := int(apiurl[2])
			switch apiurl[3] {
			case "loaded":
				outString := ""
				for _, j:= range loadedmod[selectedModule] {
					outString += string(j)
				}
				io.WriteString(w, outString)
			case "solved":
				if modsolvedstat[selectedModule] {
					io.WriteString(w, "True")
				} else {
					io.WriteString(w, "False")
				}
			}
		}
	case "POST":
		switch apiurl[1] {
		case "module":
			selectedModule := int(apiurl[2])
			switch apiurl[3] {
			case "solved":
				if apiurl[4] == "True" {
					modsolvedstat[selectedModule] = true
				} else {
					modsolvedstat[selectedModule] = false
				}
		}
		case "status":
			var newStatus comfiles.Status
			iinfo := r.Body
			err := json.Unmarshal(iinfo, &newStatus) // In the main program each field should be checked then accessor methods used
			if err != nil {
				io.WriteString(w, "Error: "+err.Error())
			}
			if iinfo contains "NumStrike" {// fix this
				cstatusptr.NumStrike = newStatus.NumStrike
			}
			io.WriteString(w, "Success")
	}	

	http.Handle("/api", fooHandler)
}
