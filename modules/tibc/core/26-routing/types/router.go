package types

import (
	"fmt"
)

// Route is a TIBCModule and a Port
type Route struct {
	Module TIBCModule
	Port   string
}

// IsManyPerContainerType implements the depinject.ManyPerContainerType interface.
func (Route) IsManyPerContainerType() {}

// Router is a map from module name to the TIBCModule
// which contains all the module-defined callbacks required by TICS-26
type Router struct {
	routes map[string]TIBCModule
	sealed bool
}

// NewRouter creates a new TIBC Router
func NewRouter() *Router {
	return &Router{
		routes: make(map[string]TIBCModule),
	}
}

// Seal prevents the Router from any subsequent route handlers to be registered.
// Seal will panic if called more than once.
func (rtr *Router) Seal() {
	if rtr.sealed {
		panic("router already sealed")
	}
	rtr.sealed = true
}

// Sealed returns a boolean signifying if the Router is sealed or not.
func (rtr Router) Sealed() bool {
	return rtr.sealed
}

// AddRoute adds TIBCModule for a given module name. It returns the Router
// so AddRoute calls can be linked. It will panic if the Router is sealed.
func (rtr *Router) AddRoute(port Port, cbs TIBCModule) *Router {
	if rtr.sealed {
		panic(fmt.Sprintf("router sealed; cannot register %s route callbacks", port))
	}
	if rtr.HasRoute(port) {
		panic(fmt.Sprintf("route %s has already been registered", port))
	}

	rtr.routes[string(port)] = cbs
	return rtr
}

// Add adds TIBCModule for a given module name. It returns the Router
// so Add calls can be linked. It will panic if the Router is sealed.
func (rtr *Router) Add(r Route) *Router {
	if rtr.sealed {
		panic(fmt.Sprintf("router sealed; cannot register %s route callbacks", r.Port))
	}
	if rtr.HasRoute(Port(r.Port)) {
		panic(fmt.Sprintf("route %s has already been registered", r.Port))
	}

	rtr.routes[r.Port] = r.Module
	return rtr
}

// HasRoute returns true if the Router has a module registered or false otherwise.
func (rtr *Router) HasRoute(port Port) bool {
	_, ok := rtr.routes[string(port)]
	return ok
}

// GetRoute returns a TIBCModule for a given module.
func (rtr *Router) GetRoute(port Port) (TIBCModule, bool) {
	if !rtr.HasRoute(port) {
		return nil, false
	}
	return rtr.routes[string(port)], true
}
