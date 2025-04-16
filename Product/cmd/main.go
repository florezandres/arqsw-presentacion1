package main

import (
	product "Taller2/Product/gen/go/api"
	"Taller2/Product/internal/api"
	"Taller2/Product/internal/command"
	kafkaComm "Taller2/Product/internal/command/infrastructure/kafka"
	"Taller2/Product/internal/command/infrastructure/persistence"
	"Taller2/Product/internal/query"
	kafkaQue "Taller2/Product/internal/query/infrastructure/kafka"
	persistenceQuery "Taller2/Product/internal/query/infrastructure/persistence"
	config "Taller2/Product/internal/shared"
	"context"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {

	// Cargar archivo .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error cargando archivo .env:", err)
	}

	// Cargar configuración
	cfg := config.Load()

	log.Printf("Config loaded: %+v", cfg)

	// 1. Inicializar Command
	commandRepo, err := persistence.NewPostgresProductRepository(cfg.PostgreSQL.CommandDSN)
	if err != nil {
		log.Fatal("Failed to init Command DB:", err)
	}

	kafkaProducer := kafkaComm.NewProducer(cfg.Kafka.Brokers)
	commandApp := command.NewApplication(commandRepo, kafkaProducer)

	// 2. Inicializar Query
	queryRepo, err := persistenceQuery.NewPostgresProductReadRepository(cfg.PostgreSQL.QueryDSN)
	if err != nil {
		log.Fatal("Failed to init Query DB:", err)
	}

	kafkaConsumer := kafkaQue.NewConsumer(cfg.Kafka.Brokers, "product_events", queryRepo)
	go func() {
		kafkaConsumer.Start(context.Background())
	}()

	queryApp := query.NewApplication(queryRepo)

	// 3. Iniciar servidor GRPC
	grpcServer := grpc.NewServer()
	product.RegisterProductServiceServer(
		grpcServer,
		api.NewProductServer(commandApp, queryApp),
	)

	lis, err := net.Listen("tcp", ":"+cfg.GRPC.Port)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Servicio Producto GRPC escuchando en %s", cfg.GRPC.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
