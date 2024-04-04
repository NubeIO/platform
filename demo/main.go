package main

import (
	"github.com/NubeIO/platform"
	"github.com/gin-gonic/gin"
	"log"
	"sync"
)

func main() {
	r := gin.Default()

	im := &platform.InstanceManager{
		Instances: make(map[string]*platform.Instance),
		Lock:      sync.Mutex{},
	}

	err := im.LoadFromFile("./db.yaml")
	if err != nil {
		log.Fatal(err)
	}

	platform.NewInstanceManagerHandler(im, r)

	r.Run(":8080")
}
