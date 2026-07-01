package rexMiddleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/rootexit/rexLib/rexCtx"
	"github.com/rootexit/rexLib/rexDatabase"
	"github.com/rootexit/rexLib/rexHeaders"
	"github.com/zeromicro/go-zero/core/logc"
)

type BotHttpInterceptorMiddleware struct {
	debug bool
}

func NewBotHttpInterceptorMiddleware(isDebug bool) *BotHttpInterceptorMiddleware {
	return &BotHttpInterceptorMiddleware{
		debug: isDebug,
	}
}

func (m *BotHttpInterceptorMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		ctx := r.Context()

		clientInfo := rexDatabase.Client{}

		if ctx.Value(rexCtx.CtxClientInfo{}) == nil {
			clientInfo = rexDatabase.Client{}
			ctx = context.WithValue(ctx, rexCtx.CtxClientInfo{}, clientInfo)
		} else {
			clientInfo = ctx.Value(rexCtx.CtxClientInfo{}).(rexDatabase.Client)
		}
		if m.debug {
			logc.Infof(ctx, "BotHttpInterceptorMiddleware clientInfo: %s", clientInfo)
		}

		userAgent := ""
		if ctx.Value(rexCtx.CtxUserAgent{}) == nil {
			userAgent = r.Header.Get(rexHeaders.HeaderUserAgent)
			ctx = context.WithValue(ctx, rexCtx.CtxUserAgent{}, userAgent)
		} else {
			userAgent = ctx.Value(rexCtx.CtxUserAgent{}).(string)
		}
		if m.debug {
			logc.Infof(ctx, "BotHttpInterceptorMiddleware userAgent: %s", userAgent)
		}

		isBot, botName, botCategory, detectionMethod := m.DetectBotByUA(userAgent)
		if isBot {
			if m.debug {
				logc.Infof(ctx,
					"BotHttpInterceptorMiddleware detected bot: ua=%s, name=%s, category=%s, method=%s",
					userAgent, botName, botCategory, detectionMethod,
				)
			}

			clientInfo.ClientBot = rexDatabase.ClientBot{
				IsBot:       true,
				BotName:     botName,
				BotCategory: botCategory,
				BotReason:   detectionMethod,
			}
		} else {
			clientInfo.ClientBot.IsBot = false
		}
		ctx = context.WithValue(ctx, rexCtx.CtxClientInfo{}, clientInfo)

		endTime := time.Now()
		if m.debug {
			logc.Infof(ctx, "BotHttpInterceptorMiddleware time consumption: %s", endTime.Sub(startTime).String())
		}

		r = r.WithContext(ctx)
		next(w, r)
	}
}

func (m *BotHttpInterceptorMiddleware) DetectBotByUA(ua string) (bool, string, string, string) {
	u := strings.ToLower(ua)
	switch {
	case strings.Contains(u, "googlebot"):
		return true, "Googlebot", "search_engine", "ua_rule"
	case strings.Contains(u, "bingbot"):
		return true, "Bingbot", "search_engine", "ua_rule"
	case strings.Contains(u, "baiduspider"):
		return true, "Baiduspider", "search_engine", "ua_rule"
	case strings.Contains(u, "telegrambot"):
		return true, "TelegramBot", "social_preview", "ua_rule"
	case strings.Contains(u, "twitterbot"):
		return true, "TwitterBot", "social_preview", "ua_rule"
	case strings.Contains(u, "facebookexternalhit"):
		return true, "FacebookExternalHit", "social_preview", "ua_rule"
	case strings.Contains(u, "curl"):
		return true, "curl", "cli", "ua_rule"
	case strings.Contains(u, "wget"):
		return true, "wget", "cli", "ua_rule"
	case strings.Contains(u, "python-requests"):
		return true, "python-requests", "script", "ua_rule"
	case strings.Contains(u, "headlesschrome"):
		return true, "HeadlessChrome", "automation", "ua_rule"
	case strings.Contains(u, "playwright"):
		return true, "Playwright", "automation", "ua_rule"
	case strings.Contains(u, "selenium"):
		return true, "Selenium", "automation", "ua_rule"
	default:
		return false, "", "", ""
	}
}
