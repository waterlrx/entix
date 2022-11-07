package service

import (
	"context"
	"log"

	"rxdrag.com/entify/common/auth"
	"rxdrag.com/entify/common/contexts"
	"rxdrag.com/entify/consts"
	"rxdrag.com/entify/model/graph"
	"rxdrag.com/entify/model/meta"
	"rxdrag.com/entify/orm"
)

type Service struct {
	isSystem bool
	ctx      context.Context
	roleIds  []uint64
	model    *graph.Model
}

func New(ctx context.Context, model *graph.Model) *Service {

	return &Service{
		isSystem: false,
		ctx:      ctx,
		model:    model,
		roleIds:  QueryRoleIds(ctx, model),
	}
}

func NewSystem() *Service {
	return &Service{
		isSystem: true,
	}
}

func (s *Service) me() *auth.User {
	return contexts.Values(s.ctx).Me
}

func (s *Service) appId() uint64 {
	return contexts.Values(s.ctx).AppId
}

func (s *Service) canReadEntity(entity *graph.Entity) (bool, graph.QueryArg) {
	whereArgs := map[string]interface{}{}
	return false, whereArgs
}

func QueryRoleIds(ctx context.Context, model *graph.Model) []uint64 {
	ids := []uint64{
		consts.GUEST_ROLE_ID,
	}

	me := contexts.Values(ctx).Me

	if me == nil {
		return ids
	}

	session, err := orm.Open()
	if err != nil {
		log.Panic(err.Error())
	}

	result := session.QueryEntity(model.GetEntityByName(meta.ROLE_ENTITY_NAME), map[string]interface{}{
		"users": map[string]interface{}{
			"id": map[string]interface{}{
				consts.ARG_EQ: me.Id,
			},
		},
	})

	for _, role := range result.Nodes {
		ids = append(ids, role[consts.ID].(uint64))
	}

	return ids
}
