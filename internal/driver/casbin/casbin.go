package casbin

import (
	"context"

	"github.com/aidapedia/gdk/log"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/util"
	"go.uber.org/zap"
)

// NewCasbin is a function to create a new casbin enforcer
func NewCasbin(ctx context.Context) *casbin.Enforcer {
	m := model.NewModel()
	m.AddDef("r", "r", "sub, obj, act")
	m.AddDef("p", "p", "sub, obj, act")
	m.AddDef("e", "e", "some(where (p.eft == allow))")
	m.AddDef("m", "m", `r.sub == "superadmin" || (r.sub == p.sub && regexMatch(r.obj, p.obj) && keyMatch2(r.act, p.act))`)

	authEnforcer, err := casbin.NewEnforcer(m)
	if err != nil {
		log.FatalCtx(ctx, "Error when creating casbin enforcer", zap.Error(err))
	}
	authEnforcer.AddNamedMatchingFunc("g", "KeyMatch2", util.KeyMatch2)

	return authEnforcer
}
