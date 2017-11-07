package main

import (
	"github.com/zenazn/goji"
	"github.com/sourcecode/ImageScrapGolang/system"
	"github.com/sourcecode/ImageScrapGolang/route"
	"flag"
)

func main()  {
	var application = &system.Application{}
	route.PrepareRoutes(application)
	flag.Set("bind","172.20.10.4:8000")
	goji.Serve()

}

