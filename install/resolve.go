package install

import (
	"log"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/mitchellh/mapstructure"
	"rxdrag.com/entify/app"
	"rxdrag.com/entify/consts"
	"rxdrag.com/entify/logs"
	"rxdrag.com/entify/model/data"
	"rxdrag.com/entify/model/meta"
	"rxdrag.com/entify/orm"
	"rxdrag.com/entify/service"
	"rxdrag.com/entify/utils"
)

type InstallArg struct {
	Admin    string     `json:"admin"`
	Password string     `json:"password"`
	WithDemo bool       `json:"withDemo"`
	Meta     utils.JSON `json:"meta"`
}

const INPUT = "input"

func InstallResolve(p graphql.ResolveParams) (interface{}, error) {
	defer utils.PrintErrorStack()

	systemAppData := meta.SystemAppData
	input := InstallArg{}
	mapstructure.Decode(p.Args[INPUT], &input)

	if input.Meta != nil {
		systemAppData = input.Meta
	}

	nextMeta := systemAppData["meta"].(meta.MetaContent)
	app.PublishMeta(&meta.MetaContent{}, &nextMeta)

	systemApp := app.GetSystemApp()

	instance := data.NewInstance(
		systemAppData,
		systemApp.GetEntityByName(meta.APP_ENTITY_NAME),
	)

	_, err := service.InsertOne(instance)

	if err != nil {
		log.Panic(err)
	}

	systemApp, err = app.Get(1)
	if err != nil {
		log.Panic(err)
	}

	if input.Admin != "" {
		instance = data.NewInstance(
			adminInstance(input.Admin, input.Password),
			systemApp.GetEntityByName(consts.META_USER),
		)
		_, err = service.SaveOne(instance)
		if err != nil {
			logs.WriteBusinessLog(systemApp.Model, p, logs.INSTALL, logs.FAILURE, err.Error())
			return nil, err
		}
		if input.WithDemo {
			instance = data.NewInstance(
				demoInstance(),
				systemApp.GetEntityByName(consts.META_USER),
			)
			_, err = service.SaveOne(instance)
			if err != nil {
				logs.WriteBusinessLog(systemApp.Model, p, logs.INSTALL, logs.FAILURE, err.Error())
				return nil, err
			}
		}
	}
	isExist := orm.IsEntityExists(meta.APP_ENTITY_NAME)
	logs.WriteBusinessLog(systemApp.Model, p, logs.INSTALL, logs.SUCCESS, "")
	app.Installed = true
	return isExist, nil
}

func adminInstance(name string, password string) map[string]interface{} {
	return map[string]interface{}{
		consts.NAME:           "Admin",
		consts.LOGIN_NAME:     name,
		consts.PASSWORD:       utils.BcryptEncode(password),
		consts.IS_SUPPER:      true,
		consts.META_CREATEDAT: time.Now(),
		consts.META_UPDATEDAT: time.Now(),
	}
}

func demoInstance() map[string]interface{} {
	return map[string]interface{}{
		consts.NAME:           "Demo",
		consts.LOGIN_NAME:     "demo",
		consts.PASSWORD:       utils.BcryptEncode("demo"),
		consts.IS_DEMO:        true,
		consts.META_CREATEDAT: time.Now(),
		consts.META_UPDATEDAT: time.Now(),
	}
}
