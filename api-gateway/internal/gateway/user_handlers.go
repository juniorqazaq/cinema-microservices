package gateway

import (
	"net/http"

	userpb "github.com/cinema-booking-system/user-service/gen/go/user"
	"github.com/gin-gonic/gin"
)

func (s *Server) listCities(c *gin.Context) {
	ok(c, []string{"Astana", "Almaty"})
}

func (s *Server) getProfile(c *gin.Context) {
	user, okUser := currentUser(c)
	if !okUser {
		fail(c, http.StatusUnauthorized, "missing authenticated user")
		return
	}
	ctx, cancel := s.requestContext(c)
	defer cancel()
	resp, err := s.clients.User.GetProfile(ctx, &userpb.GetProfileRequest{UserId: user.ID})
	if err != nil {
		failGRPC(c, err)
		return
	}
	ok(c, userFromProto(resp.GetUser()))
}

func (s *Server) topUpWallet(c *gin.Context) {
	user, okUser := currentUser(c)
	if !okUser {
		fail(c, http.StatusUnauthorized, "missing authenticated user")
		return
	}
	var req struct {
		Amount float64 `json:"amount"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Amount <= 0 {
		fail(c, http.StatusBadRequest, "amount must be positive")
		return
	}
	ctx, cancel := s.requestContext(c)
	defer cancel()
	resp, err := s.clients.User.TopUpBalance(ctx, &userpb.TopUpBalanceRequest{
		UserId: user.ID,
		Amount: req.Amount,
	})
	if err != nil {
		failGRPC(c, err)
		return
	}
	ok(c, userFromProto(resp.GetUser()))
}

func (s *Server) adminUpdateUserRole(c *gin.Context) {
	admin, okUser := currentUser(c)
	if !okUser {
		fail(c, http.StatusUnauthorized, "missing authenticated user")
		return
	}
	var req struct {
		Role string `json:"role"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, "invalid request body")
		return
	}
	var role userpb.Role
	switch req.Role {
	case "admin":
		role = userpb.Role_ROLE_ADMIN
	default:
		role = userpb.Role_ROLE_USER
	}
	ctx, cancel := s.requestContext(c)
	defer cancel()
	resp, err := s.clients.User.UpdateUserRole(ctx, &userpb.UpdateUserRoleRequest{
		AdminId: admin.ID,
		UserId:  c.Param("id"),
		Role:    role,
	})
	if err != nil {
		failGRPC(c, err)
		return
	}
	ok(c, userFromProto(resp.GetUser()))
}
