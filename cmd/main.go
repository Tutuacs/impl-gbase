package main

import (
	"fmt"

	"github.com/Tutuacs/pkg/config"
	"github.com/Tutuacs/pkg/logs"

	"github.com/Tutuacs/cmd/api"
)

func main() {

	conf_API := config.GetAPI()

	server := api.NewApiServer(conf_API.Port)
	if err := server.Run(); err != nil {
		logs.ErrorLog(fmt.Sprintf("Error starting server: %s", err))
	}

}
