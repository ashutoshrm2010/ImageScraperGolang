package route

import (
	"github.com/zenazn/goji"
	"github.com/sourcecode/ImageScrapGolang/api"
	"github.com/sourcecode/ImageScrapGolang/system"
)

func PrepareRoutes(application *system.Application) {
	goji.Post("/services/images/search", application.Route(&api.Controller{}, "ImageScrapfromGoogle"))
	goji.Post("/services/list/searched/keys", application.Route(&api.Controller{}, "ListUserSearchInputs"))
	goji.Post("/services/list/image/urls/by/searchedkey", application.Route(&api.Controller{}, "GetSearchedImageUrlsFromDB"))

}

