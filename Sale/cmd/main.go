package main

import (
	"context"
	"database/sql"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	_ "modernc.org/sqlite"

	sale "Taller2/Sale/gen/go/api"
	"Taller2/Sale/internal/api"
	"Taller2/Sale/internal/command"
	kafkaComm "Taller2/Sale/internal/command/infrastructure/kafka"
	"Taller2/Sale/internal/command/infrastructure/persistence"
	"Taller2/Sale/internal/query"
	persistenceQuery "Taller2/Sale/internal/query/infrastructure/persistence"
)

type Config struct {
	GRPC struct {
		Port string
	}
	HTTP struct {
		Port string
	}
	SQLite struct {
		CommandDBPath string
		QueryDBPath   string
	}
	Kafka struct {
		Brokers []string
	}
}

func LoadConfig() *Config {
	cfg := &Config{}

	// Configuración gRPC
	cfg.GRPC.Port = os.Getenv("GRPC_PORT")
	if cfg.GRPC.Port == "" {
		cfg.GRPC.Port = "50052"
	}

	// Configuración HTTP
	cfg.HTTP.Port = os.Getenv("HTTP_PORT")
	if cfg.HTTP.Port == "" {
		cfg.HTTP.Port = "8082"
	}

	// Configuración SQLite
	cfg.SQLite.CommandDBPath = os.Getenv("SQLITE_COMMAND_DB_PATH")
	if cfg.SQLite.CommandDBPath == "" {
		cfg.SQLite.CommandDBPath = "./sale_command.db"
	}

	cfg.SQLite.QueryDBPath = os.Getenv("SQLITE_QUERY_DB_PATH")
	if cfg.SQLite.QueryDBPath == "" {
		cfg.SQLite.QueryDBPath = "./sale_query.db"
	}

	// Configuración Kafka
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		kafkaBrokers = "localhost:9092"
	}
	cfg.Kafka.Brokers = []string{kafkaBrokers}

	return cfg
}

func main() {
	// Cargar configuración
	cfg := LoadConfig()
	log.Printf("Configuración del servicio Sale cargada: %+v", cfg)

	// 1. Inicialización de bases de datos
	commandDB, err := sql.Open("sqlite", cfg.SQLite.CommandDBPath)
	if err != nil {
		log.Fatalf("Error al abrir la base de datos SQLite de Command: %v", err)
	}
	defer commandDB.Close()

	queryDB, err := sql.Open("sqlite", cfg.SQLite.QueryDBPath)
	if err != nil {
		log.Fatalf("Error al abrir la base de datos SQLite de Query: %v", err)
	}
	defer queryDB.Close()

	// Migraciones (omitiendo por brevedad)

	// 2. Configuración de aplicaciones
	commandRepo := persistence.NewSQLiteSaleRepository(commandDB)
	kafkaProducer := kafkaComm.NewProducer(cfg.Kafka.Brokers)
	commandApp := command.NewApplication(commandRepo, kafkaProducer)

	queryRepo := persistenceQuery.NewSQLiteSaleReadRepository(queryDB)
	queryApp := query.NewApplication(queryRepo)

	// 3. Configuración del consumidor Kafka (omitiendo por brevedad)

	// 4. Servidor gRPC
	grpcServer := grpc.NewServer()
	sale.RegisterSalesServiceServer(grpcServer, api.NewSaleServer(commandApp, queryApp))

	go func() {
		lis, err := net.Listen("tcp", ":"+cfg.GRPC.Port)
		if err != nil {
			log.Fatalf("Error al escuchar gRPC: %v", err)
		}
		log.Printf("Servicio Sale gRPC escuchando en %s", cfg.GRPC.Port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Error en servidor gRPC: %v", err)
		}
	}()

	// 5. Configuración del gRPC-Gateway (Versión actualizada)
	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(
			runtime.MIMEWildcard,
			&runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{
					UseProtoNames:   true, // Equivalente a OrigName: true
					EmitUnpopulated: true, // Equivalente a EmitDefaults: true
				},
				UnmarshalOptions: protojson.UnmarshalOptions{
					DiscardUnknown: true,
				},
			},
		),
	)

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := sale.RegisterSalesServiceHandlerFromEndpoint(
		context.Background(),
		mux,
		"localhost:"+cfg.GRPC.Port,
		opts,
	); err != nil {
		log.Fatalf("Error al registrar el handler HTTP: %v", err)
	}

	// 6. Servidor HTTP
	log.Printf("Servicio Sale HTTP escuchando en %s", cfg.HTTP.Port)
	if err := http.ListenAndServe(":"+cfg.HTTP.Port, mux); err != nil {
		log.Fatalf("Error en servidor HTTP: %v", err)
	}
}
