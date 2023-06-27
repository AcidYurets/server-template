package graphql

import (
	"context"
	"errors"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"go.uber.org/zap"
	"net/http"
	"server-template/internal/modules/graphql/generated"

	"server-template/internal/models/err_const"
	"server-template/internal/modules/graphql/resolvers"
	api "server-template/internal/modules/http"
	"server-template/internal/modules/logger"
	"server-template/internal/pkg/http/error_handler/http_errors"
)

//go:generate go run -mod=mod github.com/99designs/gqlgen@v0.17.34 generate

func RegisterGraphQL(router api.ApiRouter, resolver *resolvers.Resolver) {
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: resolver,
	}))

	// Обработка паник
	srv.SetRecoverFunc(func(ctx context.Context, panicInstance interface{}) error {
		panicErr := err_const.FromPanic(panicInstance)

		logger.GetFromCtx(ctx).Error("Паника graphql", zap.Any("panic", panicInstance))

		return graphql.DefaultErrorPresenter(ctx, panicErr)
	})

	// Преобразования ошибок к необходимому виду
	srv.AroundResponses(handleResponse)

	router.All("/graphql", srv)
}

func handleResponse(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
	// Получаем сформированный ответ
	resp := next(ctx)
	resp.Extensions = make(map[string]interface{})

	commonStatusCode := http.StatusOK
	// Карта полученных кодов
	statusCodes := make(map[int]struct{})

	for _, err := range resp.Errors {
		// Фактическая ошибка
		actualError := errors.Unwrap(err)
		err.Extensions = make(map[string]interface{})

		// Если фактической ошибки нет -> ошибка в самом запросе (BadRequest)
		if actualError == nil {
			statusCodes[http.StatusBadRequest] = struct{}{}
			continue
		}

		apiErr := http_errors.NewErrorResponse(actualError)
		err.Extensions["error_response"] = apiErr
		err.Message = apiErr.Message
		statusCodes[apiErr.Code] = struct{}{}
	}

	// Если присутствуют ошибки -> по умолчанию ставим код InternalServerError
	if len(resp.Errors) != 0 {
		commonStatusCode = http.StatusInternalServerError
	}

	// Обрабатываем коды в порядке их приоритета
	if _, ok := statusCodes[http.StatusNotFound]; ok {
		commonStatusCode = http.StatusNotFound
	}
	if _, ok := statusCodes[http.StatusForbidden]; ok {
		commonStatusCode = http.StatusForbidden
	}
	if _, ok := statusCodes[http.StatusBadRequest]; ok {
		commonStatusCode = http.StatusBadRequest
	}
	if _, ok := statusCodes[http.StatusInternalServerError]; ok {
		commonStatusCode = http.StatusInternalServerError
	}

	// Устанавливаем код статуса в Extensions
	resp.Extensions["status_code"] = commonStatusCode

	return resp
}
