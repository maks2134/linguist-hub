package containers

import (
	h "github.com/makskozlov/linguist-hub/internal/handlers"
	r "github.com/makskozlov/linguist-hub/internal/repository"
	s "github.com/makskozlov/linguist-hub/internal/services"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type Repositories struct {
	userRepo    *r.UserRepository
	profileRepo *r.ProfileRepository
	planRepo    *r.PlanRepository
	chatRepo    *r.ChatRepository
}

func newRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		userRepo:    r.NewUserRepository(db),
		profileRepo: r.NewProfileRepository(db),
		planRepo:    r.NewPlanRepository(db),
		chatRepo:    r.NewChatRepository(db),
	}
}

type Services struct {
	authService    *s.AuthService
	profileService *s.ProfileService
	aiService      *s.AIService
}

func (c *Container) newServices(jwtSecret string) *Services {
	return &Services{
		authService:    s.NewAuthService(c.repos.userRepo, jwtSecret),
		profileService: s.NewProfileService(c.repos.profileRepo),
		aiService:      s.NewAIService(c.geminiService, c.repos.chatRepo, c.repos.profileRepo, c.repos.planRepo),
	}
}

type Handlers struct {
	authHandler      *h.AuthHandler
	profileHandler   *h.ProfileHandler
	aiHandler        *h.AIHandler
	webSocketHandler *h.WebSocketHandler
}

func (c *Container) newHandlers() *Handlers {
	return &Handlers{
		authHandler:      h.NewAuthHandler(c.services.authService),
		profileHandler:   h.NewProfileHandler(c.services.profileService),
		aiHandler:        h.NewAIHandler(c.services.aiService),
		webSocketHandler: h.NewWebSocketHandler(c.services.aiService),
	}
}

type Container struct {
	repos         *Repositories
	services      *Services
	handlers      *Handlers
	geminiService *s.GeminiService
}

func NewContainer(db *gorm.DB, geminiService *s.GeminiService, jwtSecret string) *Container {
	c := &Container{geminiService: geminiService}
	c.repos = newRepositories(db)
	c.services = c.newServices(jwtSecret)
	c.handlers = c.newHandlers()
	return c
}

func (c *Container) RegisterRoutes(api fiber.Router, authMiddleware fiber.Handler) {
	c.handlers.webSocketHandler.RegisterRoutes(api, authMiddleware)
	c.handlers.authHandler.RegisterRoutes(api, authMiddleware)
	c.handlers.profileHandler.RegisterRoutes(api, authMiddleware)
	c.handlers.aiHandler.RegisterRoutes(api, authMiddleware)
}
