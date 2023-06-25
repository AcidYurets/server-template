package routers

import "server-template/internal/pkg/option"

type RouteOption interface {
	option.Interface
}

type Router interface {
	Group(path string, handlers ...interface{}) Router

	Post(path string, f interface{}, opts ...RouteOption) Router
	All(path string, f interface{}, opts ...RouteOption) Router
	Get(path string, f interface{}, opts ...RouteOption) Router
	Put(path string, f interface{}, opts ...RouteOption) Router
	Patch(path string, f interface{}, opts ...RouteOption) Router
	Delete(path string, f interface{}, opts ...RouteOption) Router
}

type PermissionIdent struct{}
type PermissionReIdent struct{}
