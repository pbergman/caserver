package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/pbergman/caserver/ca"
	"github.com/pbergman/caserver/config"
	"github.com/pbergman/caserver/controller"
	"github.com/pbergman/caserver/router"
	"github.com/pbergman/caserver/util"
	"github.com/pbergman/logger"
	"github.com/pbergman/logger/handlers"
	"github.com/spf13/pflag"
)

func main() {
	var file string
	var debug bool
	pflag.StringVarP(&file, "config", "c", "/etc/caserver.cnf", "The application config file.")
	pflag.BoolVarP(&debug, "debug", "d", false, "This will print debug levels and add debug routing (see: 'net/http/pprof').")
	pflag.Parse()
	conf, err := getConfig(file)
	log := getLogger(debug)
	if err != nil {
		log.Error(err)
		return
	}
	manager, err := ca.NewManager(conf, nil)
	if err != nil {
		log.Error(err)
		return
	}
	log.Debug(fmt.Sprintf("Starting server '%s'", conf.Address))
	if err := http.ListenAndServe(conf.Address, router.NewRouter(log, getControllers(manager, debug)...)); err != nil {
		log.Error(err)
	}
}

func getControllers(manager *ca.Manager, debug bool) []router.ControllerInterface {
	controllers := []router.ControllerInterface{
		controller.NewApiCa(manager),
		controller.NewApiCertSign(manager),
		controller.NewApiCertCreate(manager),
		controller.NewApiCertDelete(manager),
		controller.NewApiCertGet(manager),
		controller.NewApiList(manager),
		controller.NewDebug(),
	}
	if debug {
		return controllers
	} else {
		return controllers[:len(controllers)-1]
	}
}

func getLogger(debug bool) *logger.Logger {
	var handler logger.HandlerInterface = handlers.NewWriterHandler(os.Stdout, logger.DEBUG)
	if !debug {
		handler = handlers.NewThresholdLevelHandler(
			handler,
			logger.ERROR,
			5,
		)
	}
	return logger.NewLogger("main", handler)
}

func getConfig(file string) (*config.Config, error) {
	cnf := new(config.Config)
	util.SetDefaults(cnf)
	return cnf, cnf.Read(file)
}
