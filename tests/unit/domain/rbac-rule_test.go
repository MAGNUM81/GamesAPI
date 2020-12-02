package domain

import (
	"GamesAPI/src/domain"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type rbacRuleTestSuite struct {
	suite.Suite
}

func TestRbacRuleTestSuite(t *testing.T) {
	suite.Run(t, new(rbacRuleTestSuite))
}

//strongly inspired from https://dev.to/bastianrob/rbac-in-rest-api-using-go-5gg0
func (s *rbacRuleTestSuite) TestRule_Comply() {
	type args struct {
		expected interface{}
		actual   interface{}
	}
	tests := []struct {
		given string
		then  string
		rule  domain.Rule
		args  args
		want  bool
	}{{
		given: "With rule: actual must be = expected", then: "query complies with our rule",
		rule: domain.Rule{
			Operator: "=",
		},
		args: args{
			expected: "something",
			actual:   "something",
		},
		want: true,
	}, {
		given: "With rule: actual must be != expected", then: "query complies with our rule",
		rule: domain.Rule{
			Operator: "!=",
		},
		args: args{
			expected: "something",
			actual:   "another",
		},
		want: true,
	}, {
		given: "With rule operator not known", then: "query does not comply",
		rule: domain.Rule{
			Operator: "unknown",
		},
		want: false,
	}}
	for _, tt := range tests {
		s.T().Run(tt.given, func(t *testing.T) {
			got := tt.rule.Comply(tt.args.expected, tt.args.actual)
			assert.Equal(t, tt.want, got, tt.then)
		})
	}
}

func (s *rbacRuleTestSuite) TestRule_FromContext() {
	tests := []struct {
		given     string
		then      string
		rule      domain.Rule
		ctx       func() context.Context
		want      interface{}
		shouldErr bool
	}{{
		given: "Non ctx rule.Value", then: "return value should be rule.Value as is",
		rule: domain.Rule{Value: "something"},
		ctx:  func() context.Context { return context.Background() },
		want: "something",
	}, {
		given: "rule.Value with ctx", then: "return value should be taken from ctx",
		rule: domain.Rule{Value: "ctx.email"},
		ctx: func() context.Context {
			return context.WithValue(context.Background(), domain.ContextKey("email"), "someone@email.com")
		},
		want: "someone@email.com",
	}, {
		given: "rule.Value with deep nested ctx", then: "return value should be taken from ctx",
		rule: domain.Rule{Value: "ctx.access.id"},
		ctx: func() context.Context {
			return context.WithValue(context.Background(), domain.ContextKey("access"), map[string]interface{}{
				"id": "IDX-0001",
			})
		},
		want: "IDX-0001",
	}, {
		given: "rule.Value with deep nested ctx, but at 4th level its not a map", then: "code should throw err",
		rule: domain.Rule{Value: "ctx.access.id.name"},
		ctx: func() context.Context {
			return context.WithValue(context.Background(), domain.ContextKey("access"), map[string]interface{}{
				"id": "IDX-0001",
			})
		},
		want:      errors.New("nested rule value is not a map or cannot be found"),
		shouldErr: true,
	}, {
		given: "rule.Value with deep nested ctx, but does not exists", then: "code should throw err",
		rule: domain.Rule{Value: "ctx.something.not.exists"},
		ctx: func() context.Context {
			return context.Background()
		},
		want:      errors.New("nested rule value is not a map or cannot be found"),
		shouldErr: true,
	}}
	for _, tt := range tests {
		s.T().Run(tt.given, func(t *testing.T) {
			if !tt.shouldErr {
				got, _ := tt.rule.FromContext(tt.ctx())
				assert.Equal(t, tt.want, got, tt.then)
			} else {
				_, got := tt.rule.FromContext(tt.ctx())
				assert.Equal(t, tt.want, got, tt.then)
			}
		})
	}
}
