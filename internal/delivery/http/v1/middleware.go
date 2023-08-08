package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/begenov/region-llc-task/internal/domain"
	"github.com/begenov/region-llc-task/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (s *Server) userIdentity(c *gin.Context) {
	id, err := s.parseAuthHeader(c)
	if err != nil {
		newResponse(c, http.StatusUnauthorized, err.Error(), fmt.Sprintf("s.parseAuthHeader(): %v", err))
		return
	}

	c.Set(userCtx, id)
}

func (s *Server) parseAuthHeader(c *gin.Context) (string, error) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		return "", errors.New("empty auth header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", domain.ErrInvalidAuthHeader
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("token is empty")
	}

	res, err := s.tokenManager.Parse(headerParts[1])
	if err != nil {
		return "", domain.ErrInvalidAuthHeader
	}

	return res, nil
}

func getUserID(c *gin.Context, context string) (primitive.ObjectID, error) {
	idFromCtx, ok := c.Get(context)
	if !ok {
		return primitive.ObjectID{}, errors.New("userCtx not found")
	}

	idStr, ok := idFromCtx.(string)
	if !ok {
		return primitive.ObjectID{}, errors.New("userCtx is of invalid type")
	}

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return primitive.ObjectID{}, domain.ErrInternalServer
	}
	logger.Info(id)
	return id, nil
}
