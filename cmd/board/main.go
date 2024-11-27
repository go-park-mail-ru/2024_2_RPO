package main

import (
	AuthGRPC "RPO_back/internal/pkg/auth/delivery/grpc/gen"
	BoardDelivery "RPO_back/internal/pkg/board/delivery"
	BoardRepository "RPO_back/internal/pkg/board/repository"
	BoardUsecase "RPO_back/internal/pkg/board/usecase"
	"RPO_back/internal/pkg/config"
	"RPO_back/internal/pkg/middleware/cors"
	"RPO_back/internal/pkg/middleware/csrf"
	"RPO_back/internal/pkg/middleware/logging_middleware"
	"RPO_back/internal/pkg/middleware/no_panic"
	"RPO_back/internal/pkg/middleware/performance"
	"RPO_back/internal/pkg/middleware/session"
	"RPO_back/internal/pkg/utils/logging"
	"RPO_back/internal/pkg/utils/misc"
	"net"
	"time"

	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Костыль
	log.Info("Sleeping 10 seconds waiting Postgres to start...")
	time.Sleep(10 * time.Second)

	// Формирование конфига
	err := config.LoadConfig()
	if err != nil {
		log.Fatalf("environment configuration is invalid: %s", err.Error())
		return
	}

	// Настройка движка логов
	logsFile, err := os.OpenFile(config.CurrentConfig.Board.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Error while opening log file %s: %s\n", config.CurrentConfig.Board.LogFile, err.Error())
		return
	}
	defer logsFile.Close()
	logging.SetupLogger(logsFile)
	log.Info("Config: ", fmt.Sprintf("%#v", config.CurrentConfig))

	// Подключение к PostgreSQL
	postgresDB, err := misc.ConnectToPgx(config.CurrentConfig.Board.PostgresPoolSize)
	if err != nil {
		log.Fatal("error connecting to PostgreSQL: ", err)
		return
	}
	defer postgresDB.Close()

	// Подключение к GRPC сервису авторизаци
	dialer := func(ctx context.Context, addr string) (net.Conn, error) {
		d := net.Dialer{
			Timeout: 5 * time.Second, // Установите подходящее время ожидания
		}
		return d.DialContext(ctx, "tcp4", addr) // Используем "tcp4" для IPv4
	}
	grpcAddr := config.CurrentConfig.AuthURL
	conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(dialer))
	if err != nil {
		log.Fatal("error connecting to GRPC: ", err)
	}
	authGRPC := AuthGRPC.NewAuthClient(conn)

	// Проверка подключения к GRPC
	sess := &AuthGRPC.CheckSessionRequest{SessionID: "12345678"}
	_, err = authGRPC.CheckSession(context.Background(), sess)
	if err != nil {
		log.Fatal("error while pinging GRPC: ", err)
	}

	//Board
	boardRepository := BoardRepository.CreateBoardRepository(postgresDB)
	boardUsecase := BoardUsecase.CreateBoardUsecase(boardRepository)
	boardDelivery := BoardDelivery.CreateBoardDelivery(boardUsecase)

	// Создаём новый маршрутизатор
	router := mux.NewRouter()

	// Применяем middleware
	router.Use(no_panic.PanicMiddleware)
	pm, err := performance.CreateHTTPPerformanceMiddleware("board")
	if err != nil {
		log.Fatal("create HTTP middleware: ", err)
	}
	router.Use(pm.Middleware)
	router.Use(logging_middleware.LoggingMiddleware)
	router.Use(cors.CorsMiddleware)
	router.Use(csrf.CSRFMiddleware)
	sm := session.CreateSessionMiddleware(authGRPC)
	router.Use(sm.Middleware)

	// Регистрируем обработчики
	router.HandleFunc("/prometheus/metrics", promhttp.Handler().ServeHTTP)
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
	router.HandleFunc("/cards/{cardID}", boardDelivery.UpdateCard).Methods("PATCH", "OPTIONS")
	router.HandleFunc("/cards/{cardID}", boardDelivery.DeleteCard).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/cardDetails/{cardID}", boardDelivery.GetCardDetails).Methods("GET", "OPTIONS")
	router.HandleFunc("/columns/{boardID}", boardDelivery.CreateColumn).Methods("POST", "OPTIONS")
	router.HandleFunc("/columns/{columnID}", boardDelivery.UpdateColumn).Methods("PUT", "OPTIONS")
	router.HandleFunc("/columns/{columnID}", boardDelivery.DeleteColumn).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/assignedUser/{cardID}", boardDelivery.AssignUser).Methods("PUT", "OPTIONS")
	router.HandleFunc("/assignedUser/{cardID}/{userID}", boardDelivery.DeassignUser).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/comments/{cardID}", boardDelivery.AddComment).Methods("POST", "OPTIONS")
	router.HandleFunc("/comments/{commentID}", boardDelivery.UpdateComment).Methods("PUT", "OPTIONS")
	router.HandleFunc("/comments/{commentID}", boardDelivery.DeleteComment).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/checkList/{cardID}", boardDelivery.AddCheckListField).Methods("POST", "OPTIONS")
	router.HandleFunc("/checkList/{fieldID}", boardDelivery.UpdateCheckListField).Methods("PATCH", "OPTIONS")
	router.HandleFunc("/checkList/{fieldID}", boardDelivery.DeleteCheckListField).Methods("DELETE", "OPTIONS")
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

	// Регистрируем обработчик Prometheus
	metricsRouter := mux.NewRouter()
	metricsRouter.Handle("/prometheus/metrics", promhttp.Handler())

	// Объявляем серверы
	mainAddr := fmt.Sprintf(":%s", os.Getenv("SERVER_PORT"))
	promAddr := ":8087"
	mainServer := &http.Server{
		Addr:    mainAddr,
		Handler: router,
	}
	promServer := &http.Server{
		Addr:    promAddr,
		Handler: metricsRouter,
	}

	misc.StartServers(mainServer, promServer)
}
