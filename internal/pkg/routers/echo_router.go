package routers

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"reflect"
	"runtime"
	"server-template/internal/models/constants"
)

// IEchoRouter интерфейс, которому удовлетворяют echo.Echo и echo.Group
type IEchoRouter interface {
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	Any(path string, handler echo.HandlerFunc, middleware ...echo.MiddlewareFunc) []*echo.Route

	Group(prefix string, middleware ...echo.MiddlewareFunc) (sg *echo.Group)
	Use(middleware ...echo.MiddlewareFunc)
}

type EchoRouter struct {
	DefaultWrapper func(handler echo.HandlerFunc, opts ...RouteOption) echo.HandlerFunc
	Router         IEchoRouter
}

/* Примеры допустимых обработчиков */
type (
	inputType  struct{}
	outputType struct{}

	_ func(ctx context.Context) (*outputType, error)
	_ func(ctx context.Context, req *inputType) (*outputType, error)
	_ func(ctx context.Context, params Params) (*outputType, error)
	_ func(ctx context.Context, params Params, req *inputType) (*outputType, error)
)

func (r *EchoRouter) wrap(handler interface{}, opts ...RouteOption) echo.HandlerFunc {
	var internalHandler echo.HandlerFunc

	// TODO Сделать обработку опций.
	//  Пока опций нет и обработки нет

	switch h := handler.(type) {
	case func(c echo.Context) error:
		internalHandler = h
	case echo.HandlerFunc:
		internalHandler = h
	case func(w http.ResponseWriter, r *http.Request):
		internalHandler = echo.WrapHandler(http.HandlerFunc(h))
	case http.Handler:
		internalHandler = echo.WrapHandler(h)
	default:
		internalHandler = echoRequest(handler)
	}

	if r.DefaultWrapper == nil {
		return internalHandler
	}
	// Обработать опции
	return r.DefaultWrapper(internalHandler, opts...)
}

func echoRequest(handler interface{}) echo.HandlerFunc {
	fVal := reflect.ValueOf(handler)
	fType := fVal.Type()
	fName := runtime.FuncForPC(fVal.Pointer()).Name()

	paramsType := reflect.TypeOf(Params{})
	contextType := reflect.TypeOf((*context.Context)(nil)).Elem()
	errorType := reflect.TypeOf((*error)(nil)).Elem()

	// 1-3 параметра входящие
	//  1 - ctx - контекст
	//  2 - params - параметры запроса типа Params (необязательный)
	//  2 или 3 - req - структура запроса (необязательный)

	// 2 параметра результата
	//  1 - res - структура ответа
	//  2 - error - ошибка

	if fType.NumIn() == 0 ||
		fType.NumIn() > 3 ||
		fType.NumOut() == 0 ||
		fType.NumOut() > 2 {
		panic(fmt.Errorf("ошибка в обработчике %s %s: некорректное число параметров", fName, fVal.String()))
	}
	if fType.In(0) != contextType {
		panic(fmt.Errorf("ошибка в обработчике %s %s: первый аргумент должен быть context.Context", fName, fVal.String()))
	}

	if fType.Out(fType.NumOut()-1) != errorType {
		panic(fmt.Errorf("ошибка в обработчике %s %s: последний аргумент должен быть error", fName, fVal.String()))
	}

	hasParams, hasReq, hasResp := false, false, false
	var reqType reflect.Type

	if fType.NumIn() == 3 {
		hasParams = true
		hasReq = true
		reqType = fType.In(2)
	} else if fType.NumIn() == 2 {
		if fType.In(1) == paramsType {
			hasParams = true
		} else {
			hasReq = true
			reqType = fType.In(1)
		}
	}

	if fType.NumOut() == 2 {
		hasResp = true
	}

	var newParams func() reflect.Value
	var newReq func() reflect.Value
	var isReqPtr bool

	if hasParams {
		newParams = func() reflect.Value {
			return reflect.MakeMapWithSize(paramsType, 0)
		}
	}

	if hasReq {
		t := reqType
		if reqType.Kind() == reflect.Ptr {
			t = reqType.Elem()
			isReqPtr = true
		}

		newReq = func() reflect.Value {
			return reflect.New(t)
		}
	}

	invokeHandler := func(c echo.Context) error {
		// Добавляем контекст первым параметром
		in := []reflect.Value{reflect.ValueOf(c.Request().Context())}

		if hasParams {
			// paramsMap - reflect.Value с видом Map
			paramsMap := newParams()

			// paramsPointer - reflect.Value с видом указатель на Map
			paramsPointer := reflect.New(paramsType)
			paramsPointer.Elem().Set(paramsMap)

			if err := (&echo.DefaultBinder{}).BindPathParams(c, paramsPointer.Interface()); err != nil {
				return fmt.Errorf("не удалось преобразовать параметры запроса: %w", err)
			}

			in = append(in, paramsMap)
		}

		if hasReq {
			req := newReq()

			if err := c.Bind(req.Interface()); err != nil {
				return fmt.Errorf("не удалось преобразовать тело запроса: %w", err)
			}

			inReq := req
			if !isReqPtr {
				inReq = req.Elem()
			}

			in = append(in, inReq)
		}

		out := fVal.Call(in)

		var err error
		var resp interface{}
		resp = struct {
			Message string
		}{
			Message: constants.ResponseSuccess,
		}

		if hasResp {
			resp = out[0].Interface()
			err, _ = out[1].Interface().(error)
		} else {
			err, _ = out[0].Interface().(error)
		}

		// Если есть ошибка, то возвращаем ее, никак не анализируя resp
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, resp)
	}

	return invokeHandler
}

func (r *EchoRouter) All(path string, handler interface{}, opts ...RouteOption) Router {
	r.Router.Any(path, r.wrap(handler, opts...))
	return r
}

func (r *EchoRouter) Group(prefix string, middlewares ...interface{}) Router {
	echoMiddlewares := make([]echo.MiddlewareFunc, 0, len(middlewares))
	for _, middleware := range middlewares {
		echoMiddleware, ok := middleware.(echo.MiddlewareFunc)
		if !ok {
			panic("указанные в Group middlewares не удовлетворяют типу echo.MiddlewareFunc")
		}
		echoMiddlewares = append(echoMiddlewares, echoMiddleware)
	}

	return &EchoRouter{
		DefaultWrapper: r.DefaultWrapper,
		Router:         r.Router.Group(prefix, echoMiddlewares...),
	}
}

func (r *EchoRouter) Use(middlewares ...interface{}) {
	echoMiddlewares := make([]echo.MiddlewareFunc, 0, len(middlewares))
	for _, middleware := range middlewares {
		echoMiddleware, ok := middleware.(echo.MiddlewareFunc)
		if !ok {
			panic("указанные в Group middlewares не удовлетворяют типу echo.MiddlewareFunc")
		}
		echoMiddlewares = append(echoMiddlewares, echoMiddleware)
	}

	r.Router.Use(echoMiddlewares...)
}

func (r *EchoRouter) Post(path string, handler interface{}, opts ...RouteOption) Router {
	r.Router.POST(path, r.wrap(handler, opts...))
	return r
}

func (r *EchoRouter) Get(path string, handler interface{}, opts ...RouteOption) Router {
	r.Router.GET(path, r.wrap(handler, opts...))
	return r
}

func (r *EchoRouter) Put(path string, handler interface{}, opts ...RouteOption) Router {
	r.Router.PUT(path, r.wrap(handler, opts...))
	return r
}

func (r *EchoRouter) Patch(path string, handler interface{}, opts ...RouteOption) Router {
	r.Router.PATCH(path, r.wrap(handler, opts...))
	return r
}

func (r *EchoRouter) Delete(path string, handler interface{}, opts ...RouteOption) Router {
	r.Router.DELETE(path, r.wrap(handler, opts...))
	return r
}

func (r *EchoRouter) Static(pathPrefix string, fsRoot string) Router {
	// Тут нужен type switch, т.к. echo.Echo и echo.Group имеют разные методы Static
	switch routerWithStatic := r.Router.(type) {
	case interface {
		Static(pathPrefix string, fsRoot string)
	}:
		routerWithStatic.Static(pathPrefix, fsRoot)
	case interface {
		Static(pathPrefix string, fsRoot string) *echo.Route
	}:
		routerWithStatic.Static(pathPrefix, fsRoot)
	default:
		panic("указан неверный тип роутера")
	}

	return r
}
