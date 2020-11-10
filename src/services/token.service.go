package services

import (
	"errors"
	"os"
)

var (
	TokenService ApiTokenServiceInterface = &apiTokenservice{}
)

type  apiTokenservice struct{}

type ApiTokenServiceInterface interface {
	GetApiToken()( token string, err error)
	ValidateToken(string)( token bool, err error)
}


func (t apiTokenservice) GetApiToken()( token string, err error){
     requiredToken := os.Getenv("API_TOKEN")
	//token exist
	if requiredToken != "" {

		return requiredToken,nil
	}

	return "" ,errors.New("Environment  Api key Token is not find.")



}

func  (t apiTokenservice) ValidateToken(tokenHeader string)( validation bool, err error){
	// token valide
	tokenEnvironment , err := t.GetApiToken()
	if err != nil{
		return false, err
	}
		return tokenHeader == tokenEnvironment , nil

}