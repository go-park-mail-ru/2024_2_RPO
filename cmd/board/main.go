package main

import (
	BoardDelivery "RPO_back/internal/pkg/board/delivery"
	BoardRepository "RPO_back/internal/pkg/board/repository"
	BoardUsecase "RPO_back/internal/pkg/board/usecase"
	"RPO_back/internal/pkg/config"
	"RPO_back/internal/pkg/middleware/cors"
	"RPO_back/internal/pkg/middleware/csrf"
	"RPO_back/internal/pkg/middleware/logging_middleware"
	"RPO_back/internal/pkg/middleware/no_panic"
	"RPO_back/internal/pkg/utils/logging"

	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Настройка движка логов
	if _, exists := os.LookupEnv("LOGS_FILE"); exists == false {
		fmt.Printf("You should provide log file env variable: LOGS_FILE\n")
		return
	}
	logsFile, err := os.OpenFile(os.Getenv("LOGS_FILE"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Error while opening log file %s: %s\n", os.Getenv("LOGS_FILE"), err.Error())
		return
	}
	defer logsFile.Close()
	logging.SetupLogger(logsFile)

	// Загрузка переменных окружения
	err = godotenv.Load(".env")
	if err != nil {
		log.Warn("warning: no .env file loaded: ", err.Error())
		fmt.Print()
	} else {
		log.Info(".env file loaded")
	}

	// Формирование конфига
	err = config.LoadConfig()
	if err != nil {
		log.Fatalf("environment configuration is invalid: %w", err)
		return
	}

	// Подключение к PostgreSQL
	postgresDb, err := pgxpool.New(context.Background(), config.CurrentConfig.PostgresDSN)
	if err != nil {
		log.Error("error connecting to PostgreSQL: ", err)
		return
	}
	defer postgresDb.Close()

	// Проверка подключения к PostgreSQL
	if err = postgresDb.Ping(context.Background()); err != nil {
		log.Fatal("error while pinging PostgreSQL: ", err)
	}

	//Board
	boardRepository := BoardRepository.CreateBoardRepository(postgresDb)
	boardUsecase := BoardUsecase.CreateBoardUsecase(boardRepository)
	boardDelivery := BoardDelivery.CreateBoardDelivery(boardUsecase)

	// Создаём новый маршрутизатор
	router := mux.NewRouter()

	// Применяем middleware
	router.Use(no_panic.PanicMiddleware)
	router.Use(logging_middleware.LoggingMiddleware)
	router.Use(cors.CorsMiddleware)
	router.Use(csrf.CSRFMiddleware)

	// Регистрируем обработчики
	router.HandleFunc("/boards", boardDelivery.CreateNewBoard).Methods("POST", "OPTIONS")
	router.HandleFunc("/boards/{boardID}", boardDelivery.DeleteBoard).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/boards/{boardID}", boardDelivery.UpdateBoard).Methods("PUT", "OPTIONS")
	router.HandleFunc("/boards/{boardID}/backgroundImage", boardDelivery.SetBoardBackground).Methods("PUT", "OPTIONS")
	router.HandleFunc("/boards/my", boardDelivery.GetMyBoards).Methods("GET", "OPTIONS")
	router.HandleFunc("/userPermissions/{boardID}", boardDelivery.GetMembersPermissions).Methods("GET", "OPTIONS")
	router.HandleFunc("/userPermissions/{boardID}", boardDelivery.AddMember).Methods("POST", "OPTIONS")
	router.HandleFunc("/userPermissions/{boardID}/{userID}", boardDelivery.UpdateMemberRole).Methods("PUT", "OPTIONS")
	router.HandleFunc("/userPermissions/{boardID}/{userID}", boardDelivery.RemoveMember).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/cards/{boardID}/allContent", boardDelivery.GetBoardContent).Methods("GET", "OPTIONS")
	router.HandleFunc("/cards/{boardID}", boardDelivery.CreateNewCard).Methods("POST", "OPTIONS")
	router.HandleFunc("/cards/{cardID}", boardDelivery.UpdateCard).Methods("PUT", "OPTIONS")
	router.HandleFunc("/cards/{cardID}", boardDelivery.DeleteCard).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/columns/{boardID}", boardDelivery.CreateColumn).Methods("POST", "OPTIONS")
	router.HandleFunc("/columns/{columnID}", boardDelivery.UpdateColumn).Methods("PUT", "OPTIONS")
	router.HandleFunc("/columns/{columnID}", boardDelivery.DeleteColumn).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/assignedUser/{cardID}/{userID}", boardDelivery.AssignUser).Methods("PUT", "OPTIONS")
	router.HandleFunc("/assignedUser/{cardID}/{userID}", boardDelivery.DeassignUser).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/comments/{cardID}", boardDelivery.AddComment).Methods("POST", "OPTIONS")
	router.HandleFunc("/comments/{commentID}", boardDelivery.UpdateComment).Methods("PUT", "OPTIONS")
	router.HandleFunc("/comments/{commentID}", boardDelivery.DeleteComment).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/checklist/{cardID}", boardDelivery.AddCheckListField).Methods("POST", "OPTIONS")
	router.HandleFunc("/checklist/{fieldID}", boardDelivery.UpdateCheckListField).Methods("PATCH", "OPTIONS")
	router.HandleFunc("/checklist/{fieldID}", boardDelivery.DeleteCheckListField).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/cardCover/{cardID}", boardDelivery.SetCardCover).Methods("PUT", "OPTIONS")
	router.HandleFunc("/cardCover/{cardID}", boardDelivery.DeleteCardCover).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/attachments/{cardID}", boardDelivery.AddAttachment).Methods("PUT", "OPTIONS")
	router.HandleFunc("/attachments/{attachmentID}", boardDelivery.DeleteAttachment).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/cardOrder/{cardID}", boardDelivery.MoveCard).Methods("PUT", "OPTIONS")
	router.HandleFunc("/columnOrder/{columnID}", boardDelivery.MoveColumn).Methods("PUT", "OPTIONS")
	router.HandleFunc("/sharedCard/{cardUUID}", boardDelivery.GetSharedCard).Methods("GET", "OPTIONS")
	router.HandleFunc("/inviteLink/{boardID}", boardDelivery.RaiseInviteLink).Methods("PUT", "OPTIONS")
	router.HandleFunc("/inviteLink/{boardID}", boardDelivery.DeleteInviteLink).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/joinBoard/{inviteUUID}", boardDelivery.FetchInvite).Methods("GET", "OPTIONS")
	router.HandleFunc("/joinBoard/{inviteUUID}", boardDelivery.AcceptInvite).Methods("POST", "OPTIONS")

	// Запускаем сервер
	addr := fmt.Sprintf(":%s", os.Getenv("SERVER_PORT"))
	log.Infof("server started at http://0.0.0.0%s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("error while starting server: %w", err)
	}
}
