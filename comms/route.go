package comms

import r "github.com/rubikorg/rubik"

var msgpRouter = r.Create("/_msgp")

var messageRoute = r.Route{
	Method: r.POST,
	Path:   "/message",
}

var newPunchInRoute = r.Route{
	Path: "/new/service",
}

var listServicesRoute = r.Route{
	Path: "/list",
}

func addRoutes() {
	msgpRouter.Add(messageRoute)
	msgpRouter.Add(newPunchInRoute)
	msgpRouter.Add(listServicesRoute)
}
