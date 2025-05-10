package v1

import (
        "github.com/gin-gonic/gin"

        "dnd-combat/config"
        "dnd-combat/internal/auth"
        "dnd-combat/internal/character"
        "dnd-combat/internal/combat"
        "dnd-combat/internal/game"
        "dnd-combat/pkg/database"
        "dnd-combat/pkg/dnd5e"
        "dnd-combat/pkg/websocket"
)

// SetupRoutes configures all API routes
func SetupRoutes(r *gin.Engine, db *database.DB, srdClient *dnd5e.SRDClient, wsHub *websocket.Hub, cfg *config.Config) {
        // Create dice roller and combat rules
        diceRoller := dnd5e.NewDiceRoller()
        combatRules := dnd5e.NewCombatRules(diceRoller)

        // Auth setup
        authRepo := auth.NewRepository(db)
        authService := auth.NewService(authRepo, cfg)
        authHandler := auth.NewHandler(authService)
        authMiddleware := auth.NewMiddleware(authService)

        // Character setup
        characterRepo := character.NewRepository(db)
        characterService := character.NewService(characterRepo)
        characterHandler := character.NewHandler(characterService)

        // Game setup
        gameRepo := game.NewRepository(db)
        gameService := game.NewService(gameRepo)
        gameHandler := game.NewHandler(gameService)

        // Combat setup
        combatRepo := combat.NewRepository(db)
        combatService := combat.NewService(combatRepo, diceRoller, combatRules)
        srdClientAdapter := dnd5e.NewSRDClientAdapter(srdClient)
        combatHandler := combat.NewHandler(combatService, characterService, srdClientAdapter, wsHub)

        // Public routes (no auth required)
        publicRoutes := r.Group("/api/v1")
        {
                // Auth routes
                authGroup := publicRoutes.Group("/auth")
                {
                        authGroup.POST("/register", authHandler.Register)
                        authGroup.POST("/login", authHandler.Login)
                }
                
                // Websocket routes (auth is checked inside the handler)
                wsGroup := publicRoutes.Group("/ws")
                {
                        wsGroup.GET("/combat/:id", authMiddleware.RequireAuth(), combatHandler.WebSocketHandler)
                }
        }

        // Protected routes (auth required)
        protectedRoutes := r.Group("/api/v1")
        protectedRoutes.Use(authMiddleware.RequireAuth())
        {
                // Character routes
                characterGroup := protectedRoutes.Group("/characters")
                {
                        characterGroup.POST("", characterHandler.Create)
                        characterGroup.GET("/:id", characterHandler.Get)
                        characterGroup.GET("", characterHandler.List)
                }

                // Game routes
                gameGroup := protectedRoutes.Group("/games")
                {
                        gameGroup.POST("", gameHandler.Create)
                        gameGroup.GET("/:id", gameHandler.Get)
                        gameGroup.GET("", gameHandler.List)
                        gameGroup.PUT("/:id", gameHandler.Update)
                }

                // Combat routes
                combatGroup := protectedRoutes.Group("/combat")
                {
                        combatGroup.POST("", combatHandler.InitiateCombat)
                        combatGroup.GET("/:id", combatHandler.GetCombat)
                        combatGroup.POST("/:id/action", combatHandler.PerformAction)
                        combatGroup.POST("/:id/end-turn", combatHandler.EndTurn)
                }
        }

        // Debug routes (only available in development)
        if cfg.Environment == "development" {
                debugRoutes := r.Group("/api/v1/debug")
                {
                        debugRoutes.GET("/ping", func(c *gin.Context) {
                                c.JSON(200, gin.H{
                                        "message": "pong",
                                })
                        })
                }
        }
}
