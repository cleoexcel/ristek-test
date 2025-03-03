package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/cleoexcel/ristek-test/app/answer"
	"github.com/cleoexcel/ristek-test/app/auth"
	"github.com/cleoexcel/ristek-test/app/question"
	"github.com/cleoexcel/ristek-test/app/tryout"
	"github.com/cleoexcel/ristek-test/config"
	"github.com/cleoexcel/ristek-test/middleware"
)

var db *gorm.DB

func main() {
	db = config.InitDatabase()
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	repo := auth.NewRepository(db)
	service := auth.NewAuthService(repo)
	handler := auth.NewAuthHandler(service)

	r.GET("/auth/get-all-user", handler.GetAllUsers)
	r.POST("/auth/register", handler.Register)
	r.POST("/auth/login", handler.Login)

	r.Use(middleware.AuthMiddleware())

	tryoutrepo := tryout.NewRepository(db)
	tryoutservice := tryout.NewTryoutService(tryoutrepo)
	tryouthandler := tryout.NewTryoutHandler(tryoutservice)

	r.POST("/tryout/create-tryout", tryouthandler.CreateTryout)
	r.GET("/tryout/get-all-tryout", tryouthandler.GetAllTryout)
	r.GET("/tryout/get-detail-tryout/:id", tryouthandler.GetDetailTryout)
	r.PATCH("/tryout/edit-tryout/:id", tryouthandler.EditTryout)
	r.DELETE("/tryout/delete-tryout/:id", tryouthandler.DeleteTryoutById)

	questionrepo := question.NewQuestionRepository(db)
	answerrepo := answer.NewAnswerRepository(db)
	answerservice := answer.NewAnswerService(answerrepo)
	questionservice := question.NewQuestionService(questionrepo, answerservice)
	questionhandler := question.NewQuestionHandler(questionservice)

	r.POST("/question/create-question", questionhandler.CreateQuestion)
	r.GET("/question/get-all-question/:id", questionhandler.GetAllQuestions)
	r.PATCH("/question/edit-question/:id", questionhandler.EditQuestion)
	r.DELETE("/question/delete-question/:id", questionhandler.DeleteQuestion)

	r.Run(":8080")
}
