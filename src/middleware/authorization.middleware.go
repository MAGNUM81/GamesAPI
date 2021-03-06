package middleware

import (
	"GamesAPI/src/domain"
	"GamesAPI/src/services"
	"GamesAPI/src/utils/errorUtils"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
)

func InitAuthorization(g *gin.RouterGroup) {
	services.AuthorizationService = services.NewAuthorizationService(os.Getenv("RBAC_FILEPATH"))
	g.Use(AuthorizationHandler)
}

func handleAuthError(c *gin.Context, status int, err error) {
	if err != nil {
		c.AbortWithStatusJSON(status, gin.H{"Error": err.Error()})
	}
}

func AuthorizationHandler(c *gin.Context) {
	//1. Determine which role the user has in our system using ctx.userId
	//	 Assuming ctx has userId forces us to set userId in the Authentification layer's context
	ctx := c.Request.Context()
	userId := ctx.Value(domain.RbacUserId())
	if userId == nil {
		handleAuthError(c, http.StatusInternalServerError, errors.New("no user id could be found in context"))
		return
	}

	roles, err := services.UserRoleService.GetRolesByUserID(userId.(uint64))
	if err != nil {
		handleAuthError(c, err.Status(), err)
		return
	}
	if len(roles) < 1 {
		handleAuthError(c, 500, errorUtils.ErrNoRole)
		return
	}
	//Typically, each user has only one role, so we'll take the first one we get
	roleName := roles[0].Name

	//2. Determine which resource we're trying to access
	url := c.Request.URL
	resource, resourceErr := extractResource(url.Path)
	if resourceErr != nil {
		handleAuthError(c, 400, resourceErr)
		return
	}

	//3. Determine which endpoint we're trying to access
	httpMethod := c.Request.Method
	endpoint, endpointErr := extractEndpoint(httpMethod)
	if endpointErr != nil {
		handleAuthError(c, 400, endpointErr)
		return
	}
	roleName = strings.ToLower(roleName)
	//4. Authorize the request using all the info provided
	authErr := services.AuthorizationService.Authorize(ctx, url, roleName, resource, endpoint)
	if authErr != nil {
		handleAuthError(c, 403, authErr)
		return
	}

	c.Next()
}

func extractResource(urlPath string) (string, error) {
	if strings.Contains(urlPath, "/games") {
		return "game", nil
	}

	if strings.Contains(urlPath, "/users") {
		return "user", nil
	}

	if strings.Contains(urlPath, "/LinkSteamUser") {
		return "link_steam_user", nil
	}

	if strings.Contains(urlPath, "/SyncGames") {
		return "sync_games", nil
	}

	return "", errors.New("resource does not exist")
}

func extractEndpoint(httpMethod string) (string, error) {
	var ret = ""
	var err error = nil
	switch httpMethod {
	case "GET":
		ret = "read"
	case "POST":
		ret = "create"
	case "PATCH":
		ret = "update"
	case "PUT":
		ret = "update"
	case "DELETE":
		ret = "delete"
	default:
		err = errors.New("endpoint does not exist")
	}
	return ret, err
}
