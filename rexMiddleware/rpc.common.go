package rexMiddleware

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/rootexit/rexLib/rexCodes"
	"github.com/rootexit/rexLib/rexCtx"
	"github.com/rootexit/rexLib/rexHeaders"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// note: 基于grpc的中间件，实现读取metadata中的信息映射到context中

func StreamHeaderParseInterceptor() grpc.StreamServerInterceptor {
	return func(svr any, stream grpc.ServerStream, info *grpc.StreamServerInfo,
		handler grpc.StreamHandler) error {
		ctx := stream.Context()

		result := &Resp{
			Code: rexCodes.EngineStatusOK,
			Msg:  rexCodes.StatusText(rexCodes.EngineStatusOK),
			Path: info.FullMethod,
		}

		ctx = context.WithValue(ctx, rexCtx.CtxFullMethod{}, info.FullMethod)

		// note: metadata中尝试获取requestId, 如果不存在就生成一个
		tempMD, isExist := metadata.FromIncomingContext(ctx)
		if !isExist {
			result.Code = rexCodes.EngineStatusNotFoundMetadata
			result.Msg = rexCodes.StatusText(rexCodes.EngineStatusNotFoundMetadata)
			return errors.New(rexCodes.StatusText(rexCodes.EngineStatusNotFoundMetadata))
		}

		requestId := tempMD.Get(rexHeaders.HeaderXRequestIDFor)
		if len(requestId) > 0 {
			ctx = context.WithValue(ctx, rexCtx.CtxRequestId{}, requestId[0])
			result.RequestID = requestId[0]
		} else {
			tempRequestId := uuid.NewString()
			ctx = context.WithValue(ctx, rexCtx.CtxRequestId{}, tempRequestId)
			result.RequestID = tempRequestId
		}

		//note: 读取metadata中的信息
		xTenantIDFor := tempMD.Get(rexHeaders.HeaderXTenantIDFor)
		if len(requestId) > 0 {
			ctx = context.WithValue(ctx, rexCtx.CtxTenantId{}, xTenantIDFor[0])
		}

		xDomainIdFor := tempMD.Get(rexHeaders.HeaderXDomainIDFor)
		if len(requestId) > 0 {
			ctx = context.WithValue(ctx, rexCtx.CtxDomainId{}, xDomainIdFor[0])
		}

		return handler(svr, stream)
	}
}

func UnaryHeaderParseInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (any, error) {

		result := &Resp{
			Code: rexCodes.EngineStatusOK,
			Msg:  rexCodes.StatusText(rexCodes.EngineStatusOK),
			Path: info.FullMethod,
		}

		ctx = context.WithValue(ctx, rexCtx.CtxFullMethod{}, info.FullMethod)

		// note: metadata中尝试获取requestId, 如果不存在就生成一个
		tempMD, isExist := metadata.FromIncomingContext(ctx)
		if !isExist {
			result.Code = rexCodes.EngineStatusNotFoundMetadata
			result.Msg = rexCodes.StatusText(rexCodes.EngineStatusNotFoundMetadata)
			return result, nil
		}

		requestId := tempMD.Get(rexHeaders.HeaderXRequestIDFor)
		if len(requestId) > 0 {
			ctx = context.WithValue(ctx, rexCtx.CtxRequestId{}, requestId[0])
			result.RequestID = requestId[0]
		} else {
			tempRequestId := uuid.NewString()
			ctx = context.WithValue(ctx, rexCtx.CtxRequestId{}, tempRequestId)
			result.RequestID = tempRequestId
		}

		xTenantIDFor := tempMD.Get(rexHeaders.HeaderXTenantIDFor)
		if len(requestId) > 0 {
			ctx = context.WithValue(ctx, rexCtx.CtxTenantId{}, xTenantIDFor[0])
		}

		xDomainIdFor := tempMD.Get(rexHeaders.HeaderXDomainIDFor)
		if len(requestId) > 0 {
			ctx = context.WithValue(ctx, rexCtx.CtxDomainId{}, xDomainIdFor[0])
		}

		return handler(ctx, req)
	}
}
