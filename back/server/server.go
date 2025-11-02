package server

import (
	"app/config"
)

func Init() {

	c := config.GetConfig()

	r := NewRouter()
	err := r.Run(c.Host + ":" + c.Port)
	if err != nil {
		panic(err)
	}
}
