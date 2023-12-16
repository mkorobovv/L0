package app

import (
	"fmt"
	"log"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mkorobovv/L0/config"
	"github.com/mkorobovv/L0/internal/cache"
	"github.com/mkorobovv/L0/internal/controller"
	"github.com/mkorobovv/L0/internal/database"
	"github.com/mkorobovv/L0/internal/generate"
	"github.com/mkorobovv/L0/internal/nats"
)

type Server struct {
	echo *echo.Echo
	port string
}

func NewServer() *Server {
	return &Server{
		echo: echo.New(),
		port: ":8000",
	}
}

func (s *Server) ListenAndServe(orderHand echo.HandlerFunc, allOrdersHand echo.HandlerFunc) error {
	s.echo.GET("/:order", orderHand)
	s.echo.GET("/", allOrdersHand)
	return s.echo.Start(s.port)
}

func ServerRun(cfg *config.Config) {
	ns := nats.NewNats(&cfg.Nats)
	fmt.Println("Nats connecrion succes")

	pgConn, err := database.Connect(&cfg.Postgres)
	if err != nil {
		fmt.Printf("Postgres connection error: %v\n", err)
	}
	defer pgConn.Close()

	db := database.NewDB(pgConn)
	dbErr := db.CreateTable()

	if dbErr != nil {
		fmt.Printf("Error creating table: %v\n", dbErr)
	}

	cache := cache.NewCache(db)
	cache.Preload()

	go func() {
		for {
			order := generate.Generate()
			fmt.Println("Successfully generated")
			err := ns.Publish(*order)
			if err != nil {
				fmt.Printf("Error publish: %v\n", err)
			}
			time.Sleep(30 * time.Second)
		}
	}()

	go func() {
		for {
			msg, err := ns.Subscribe()

			if err != nil {
				fmt.Printf("Error at receiving %v\n", err)
			}
			fmt.Println("Order received successfully")
			cache.SetOrderInfo(*msg)
			time.Sleep(60 * time.Second)
		}
	}()

	httpServer := NewServer()
	controller := controller.NewOrderController(cache)

	serverErr := httpServer.ListenAndServe(controller.GetOrderController, controller.GetAllOrders)

	if serverErr != nil {
		log.Fatalf("Error at server starting: %v", serverErr)
	}

}
