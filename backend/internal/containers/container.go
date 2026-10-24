package containers

import (
	h "github.com/makskozlov/linguist-hub/internal/handlers"
	r "github.com/makskozlov/linguist-hub/internal/repository"
	s "github.com/makskozlov/linguist-hub/internal/services"
	"github.com/makskozlov/linguist-hub/internal/services/speech"
	"github.com/makskozlov/linguist-hub/pkg/middleware"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type Repositories struct {
	userRepo        *r.UserRepository
	profileRepo     *r.ProfileRepository
	planRepo        *r.PlanRepository
	chatRepo        *r.ChatRepository
	voiceLessonRepo *r.VoiceLessonRepository
}

func newRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		userRepo:        r.NewUserRepository(db),
		profileRepo:     r.NewProfileRepository(db),
		planRepo:        r.NewPlanRepository(db),
		chatRepo:        r.NewChatRepository(db),
		voiceLessonRepo: r.NewVoiceLessonRepository(db),
	}
}

type Services struct {
	authService    *s.AuthService
	profileService *s.ProfileService
	aiService      *s.AIService
	voiceService   *s.VoiceService
}

func (c *Container) newServices(jwtSecret string, serverSTT bool) *Services {
	ai := s.NewAIService(c.llmService, c.repos.chatRepo, c.repos.profileRepo, c.repos.planRepo)
	var recognizer speech.Recognizer = speech.NewPassthroughRecognizer()
	if serverSTT {
		recognizer = speech.NewLLMRecognizer(c.llmService)
	}
	return &Services{
		authService:    s.NewAuthService(c.repos.userRepo, jwtSecret),
		profileService: s.NewProfileService(c.repos.profileRepo),
		aiService:      ai,
		voiceService: s.NewVoiceService(
			ai, c.llmService, c.repos.voiceLessonRepo,
			c.repos.chatRepo, c.repos.profileRepo, c.repos.planRepo,
			recognizer,
		),
	}
}

type Handlers struct {
	authHandler      *h.AuthHandler
	profileHandler   *h.ProfileHandler
	aiHandler        *h.AIHandler
	voiceHandler     *h.VoiceHandler
	webSocketHandler *h.WebSocketHandler
}

func (c *Container) newHandlers(serverSTT bool) *Handlers {
	vh := h.NewVoiceHandler(c.services.voiceService)
	if serverSTT {
		var recognizer speech.Recognizer = speech.NewLLMRecognizer(c.llmService)
		vh = h.NewVoiceHandlerWithSTT(c.services.voiceService, recognizer)
	}
	return &Handlers{
		authHandler:      h.NewAuthHandler(c.services.authService),
		profileHandler:   h.NewProfileHandler(c.services.profileService),
		aiHandler:        h.NewAIHandler(c.services.aiService),
		voiceHandler:     vh,
		webSocketHandler: h.NewWebSocketHandler(
			c.services.aiService,
			c.services.voiceService,
			middleware.DefaultAIMinuteLimit,
			middleware.DefaultAIDailyQuota,
		),
	}
}

type Container struct {
	repos         *Repositories
	services      *Services
	handlers      *Handlers
	llmService *s.LLMService
}

func NewContainer(db *gorm.DB, llmService *s.LLMService, jwtSecret string, serverSTT bool) *Container {
	c := &Container{llmService: llmService}
	c.repos = newRepositories(db)
	c.services = c.newServices(jwtSecret, serverSTT)
	c.handlers = c.newHandlers(serverSTT)
	return c
}

func (c *Container) RegisterRoutes(api fiber.Router, authMiddleware fiber.Handler) {
	registerRL := middleware.IPRateLimit(
		middleware.DefaultRegisterLimit,
		"REGISTER_RATE_LIMIT",
		"Too many registrations from this IP, try later",
	)
	loginRL := middleware.IPRateLimit(
		middleware.DefaultLoginLimit,
		"LOGIN_RATE_LIMIT",
		"Too many login attempts, try later",
	)
	aiRL := middleware.AIProtect(
		middleware.DefaultAIMinuteLimit,
		middleware.DefaultAIDailyQuota,
	)
	voiceRL := middleware.VoiceRateLimit(middleware.DefaultVoiceRateLimitConfig())
	voiceDaily := middleware.AIDailyQuota(middleware.DefaultAIDailyQuota)

	c.handlers.webSocketHandler.RegisterRoutes(api, authMiddleware)
	c.handlers.authHandler.RegisterRoutes(api, authMiddleware, registerRL, loginRL)
	c.handlers.profileHandler.RegisterRoutes(api, authMiddleware)
	c.handlers.aiHandler.RegisterRoutes(api, authMiddleware, aiRL)
	c.handlers.voiceHandler.RegisterRoutes(api, authMiddleware, voiceDaily, voiceRL)
}
