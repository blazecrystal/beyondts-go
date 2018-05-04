package web

import (
    "github.com/blazecrystal/beyondts-go/properties"
    "github.com/blazecrystal/beyondts-go/logq2"
    "os"
    "net/http"
    "flag"
)

const (
    LOGGER = "github.com/blazecrystal/beyondts-go/web"
    WEB_CONTEXT_ROOT = "web.context.root"
    WEB_LISTENER_ADDR = "web.listener.addr"
)

func configAll(configFile string) {

}

func CreateServlet(configFile string) {
    /*config, err := properties.LoadPropertiesFromFile(configFile)
    if err != nil {
        logq2.GetLogger(LOGGER).Fatal(err, "can't parse web config file : ", configFile)
        os.Exit(1)
    }
    contextRoot := config.Get(WEB_CONTEXT_ROOT)
    listenerAddr := config.Get(WEB_LISTENER_ADDR)
    flag.String()*/
}

type WebFilter interface {
    Filter(rw http.ResponseWriter, r *http.Request)
}

