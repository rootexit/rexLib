package rexPgPool

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/zeromicro/go-zero/core/logc"
)

type QueryTracer struct{}

func (t *QueryTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	logc.Infof(ctx, "SQLC START: %s args=%v", data.SQL, data.Args)
	return ctx
}

func (t *QueryTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	logc.Infof(ctx, "SQLC END: err=%v", data.Err)
}
