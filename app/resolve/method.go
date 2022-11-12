package resolve

import (
	"fmt"
	"strings"

	"github.com/dop251/goja"
	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/app/script"
	"rxdrag.com/entify/model"
	"rxdrag.com/entify/model/meta"
	"rxdrag.com/entify/utils"
)

func argsString(method *meta.MethodMeta) string {
	names := []string{}
	for _, arg := range method.Args {
		names = append(names, arg.Name)
	}
	return strings.Join(names, ", ")
}

func MethodResolveFn(method *meta.MethodMeta, model *model.Model) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		defer utils.PrintErrorStack()
		vm := goja.New()
		script.Enable(vm)
		vm.Set("args", p.Args)
		script.Enable(vm)
		funcStr := fmt.Sprintf(
			`
			%s

			function doMethod() {
				const {%s} = args;
			%s
			}`,
			script.GetCodes(model),
			argsString(method),
			method.Script,
		)

		_, err := vm.RunString(funcStr)
		if err != nil {
			panic(err)
		}
		var doMethod func() interface{}
		err = vm.ExportTo(vm.Get("doMethod"), &doMethod)
		if err != nil {
			panic(err)
		}

		result := doMethod()
		return result, nil
	}
}
