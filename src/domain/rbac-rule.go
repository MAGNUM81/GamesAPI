package domain

import (
	"context"
	"errors"
	"strings"
)

type Rule struct {
	Key		 string `yaml:"key"`
	Operator string `yaml:"operator"`
	Value 	 string `yaml:"value"`
}

//Strongly inspired from https://dev.to/bastianrob/rbac-in-rest-api-using-go-5gg0
func (rule Rule) FromContext(ctx context.Context) (interface{}, error) {
	if !strings.HasPrefix(rule.Value, "ctx") {
		return rule.Value, nil
	}
	paths := strings.Split(rule.Value, ".")
	var ctxval interface{}
	var err error = nil
	// starts from 1, as we exclude the ctx part
	for i := 1; i < len(paths); i++ {
		ctxkey := paths[i]

		//Get current context index
		if i == 1 {
			ctxval = ctx.Value(ContextKey(ctxkey))
		} else {
			// if rule.Value is nested more than 1 level, we assume the context value is of type map[string]interface{}

			kvp, ok := ctxval.(map[string]interface{})
			if !ok {
				err = errors.New("nested rule value is not a map or cannot be found")
				ctxval = nil
				break
			}
			ctxval, ok = kvp[ctxkey]
			if !ok || ctxval == nil {
				ctxval = nil
				err = errors.New("value does not exist")
				break
			}
		}
	}
	return ctxval, err
}
