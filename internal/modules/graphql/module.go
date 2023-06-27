package graphql

import (
	"go.uber.org/fx"
	"server-template/internal/modules/graphql/resolvers"
)

var Module = fx.Module("graphql",
	fx.Provide(resolvers.NewResolver),
	fx.Invoke(RegisterGraphQL))
