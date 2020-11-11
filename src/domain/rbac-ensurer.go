package domain

import (
	"GamesAPI/src/utils/errorUtils"
	"context"
	"fmt"
	"net/url"
)

type Ensurer struct {
	Query  []Rule `yaml:"query"`
	Header []Rule `yaml:"header"`
	Path   []Rule `yaml:"path"`
}

type Enforcer Ensurer

// QueryComplies checks whether query request complies with rules
func (ens Ensurer) QueryComplies(ctx context.Context, url *url.URL) error {
	if ens.Query == nil || len(ens.Query) <= 0 {
		return nil
	}

	for _, rule := range ens.Query {
		actual := url.Query().Get(rule.Key)
		expected, err := rule.FromContext(ctx)
		if err != nil {
			return err
		}
		if !rule.Comply(expected, actual) {
			return fmt.Errorf("query rule violation: ensure '%s' %s '%v', instead got: '%s'",
				rule.Key, rule.Operator, expected, actual)
		}
	}
	return nil
}

// QueryComplies enforces query request from rule
func (enf Enforcer) QueryComplies(ctx context.Context, url *url.URL) error {
	q := url.Query()
	for _, rule := range enf.Query {
		expected, err := rule.FromContext(ctx)
		if err != nil {
			return err
		}
		valueStr, isString := expected.(string)
		if !isString {
			return errorUtils.ErrNotString
		}

		q.Set(rule.Key, valueStr)
	}
	url.RawQuery = q.Encode()
	// whole enforced with rules
	return nil
}
