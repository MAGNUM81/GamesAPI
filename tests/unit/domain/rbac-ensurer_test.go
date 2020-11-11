package domain

import (
	"GamesAPI/src/domain"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type ensurerEnforcerTestSuite struct {
	suite.Suite
}

func TestEnsurerTestSuite(t *testing.T) {
	suite.Run(t, new(ensurerEnforcerTestSuite))
}

func (s *ensurerEnforcerTestSuite) TestEnsurer_QueryComplies() {
	type args struct {
		method string
		url    string
	}
	tests := []struct {
		given   string
		then    string
		ensurer domain.Ensurer
		context func() context.Context
		args    args
		wantErr bool
	}{{
		given: "Query: id=0001&name=John and Rule: id=0001&name=ctx.name and ctx.name=John",
		then:  "QueryComplies must not return error",
		args: args{
			// query: {id: "0001", name: "John"}
			url: "http://api.example.com/resources?id=0001&name=John",
		},
		ensurer: domain.Ensurer{
			Query: []domain.Rule{
				// id IS 0001, and name EQUALS to value stored in context.name
				{Key: "id", Operator: "=", Value: "0001"},
				{Key: "name", Operator: "=", Value: "ctx.name"},
			},
		},
		context: func() context.Context {
			// we give the context.name = "John"
			return context.WithValue(context.Background(), domain.ContextKey("name"), "John")
		},
	}}
	for _, tt := range tests {
		s.T().Run(tt.given, func(t *testing.T) {
			r, _ := http.NewRequest(tt.args.method, tt.args.url, nil)
			r = r.WithContext(tt.context())

			err := tt.ensurer.QueryComplies(r.Context(), r.URL)
			if tt.wantErr {
				assert.Error(t, err, tt.given)
			} else {
				assert.NoError(t, err, tt.given)
			}
		})
	}
}

func (s *ensurerEnforcerTestSuite) TestEnforcer_QueryComplies() {
	type args struct {
		method string
		url    string
	}
	tests := []struct {
		given    string
		then     string
		enforcer domain.Enforcer
		context  func() context.Context
		args     args
		want     map[string]string
		wantErr  bool
	}{{
		given: "Query: id=nil&name=nil and Rule: id=0001&name=ctx.name and ctx.name=John",
		then:  "QueryComplies must not return error, and query must be re-written by enforcer",
		args: args{
			url: "http://api.example.com/resources?id=nil&name=nil",
		},
		enforcer: domain.Enforcer{
			Query: []domain.Rule{
				// query: {id: "0001", name: "John"}
				{Key: "id", Value: "0001"},
				{Key: "name", Value: "ctx.name"},
			},
		},
		context: func() context.Context {
			// we give the context.name = "John"
			return context.WithValue(context.Background(), domain.ContextKey("name"), "John")
		},
		want: map[string]string{
			"id":   "0001",
			"name": "John",
		},
	}}
	for _, tt := range tests {
		s.T().Run(tt.given, func(t *testing.T) {
			r, _ := http.NewRequest(tt.args.method, tt.args.url, nil)
			r = r.WithContext(tt.context())

			err := tt.enforcer.QueryComplies(r.Context(), r.URL)
			if tt.wantErr {
				assert.Error(t, err, tt.given)
			} else {
				assert.NoError(t, err, tt.given)
				assert.Equal(t, len(tt.want), len(r.URL.Query()))
				for key, val := range tt.want {
					assert.Equal(t, val, r.URL.Query().Get(key), tt.then)
				}
			}
		})
	}
}
