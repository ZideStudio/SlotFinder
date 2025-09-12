package server

import (
	"app/config"
)

func Init() {

	c := config.GetConfig()

	r := NewRouter()
	r.Run(c.Host + ":" + c.Port)

}
