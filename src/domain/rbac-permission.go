package domain

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Permission struct {
	Allow bool 		 `yaml:"allow"`
	Ensure Ensurer	 `yaml:"ensure,omitempty"`
	Enforce Enforcer `yaml:"enforce,omitempty"`
}

type Endpoint map[string]Permission

type Resource map[string]Endpoint

type RBAC map[string]Resource

func RbacFromFile(path string) (RBAC, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	rbac := RBAC{}
	err = yaml.Unmarshal(f, rbac)
	if err != nil {
		return nil, err
	}

	return rbac, err
}