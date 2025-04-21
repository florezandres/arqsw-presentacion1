package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	product "Taller2/Product/gen/go/api" // Ajusta esta importación según tu estructura de carpetas
	"Taller2/Product/internal/api"
	"Taller2/Product/internal/command"
	kafkaComm "Taller2/Product/internal/command/infrastructure/kafka"
	"Taller2/Product/internal/command/infrastructure/persistence"
	"Taller2/Product/internal/query"
	kafkaQue "Taller2/Product/internal/query/infrastructure/kafka"
	persistenceQuery "Taller2/Product/internal/query/infrastructure/persistence"
	config "Taller2/Product/internal/shared"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4321")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Soporte para preflight
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func main() {
	// Cargar archivo .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error cargando archivo .env:", err)
	}

	// Cargar configuración
	cfg := config.Load()

	log.Printf("Config loaded: %+v", cfg)

	// Inicializar Command y Query (como ya lo tienes configurado)
	commandRepo, err := persistence.NewPostgresProductRepository(cfg.PostgreSQL.CommandDSN)
	if err != nil {
		log.Fatal("Failed to init Command DB:", err)
	}

	kafkaProducer := kafkaComm.NewProducer(cfg.Kafka.Brokers)
	commandApp := command.NewApplication(commandRepo, kafkaProducer)

	queryRepo := persistenceQuery.NewMemoryProductReadRepository()

	kafkaConsumer := kafkaQue.NewConsumer(cfg.Kafka.Brokers, "product_events", queryRepo)
	go func() {
		kafkaConsumer.Start(context.Background())
	}()

	queryApp := query.NewApplication(queryRepo)

	// Iniciar servidor GRPC
	grpcServer := grpc.NewServer()
	product.RegisterProductServiceServer(
		grpcServer,
		api.NewProductServer(commandApp, queryApp),
	)

	// Iniciar servidor gRPC en paralelo
	go func() {
		lis, err := net.Listen("tcp", ":"+cfg.GRPC.Port)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Servicio Producto GRPC escuchando en %s", cfg.GRPC.Port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	// Iniciar gRPC-Gateway (HTTP REST)
	go func() {
		ctx := context.Background()
		mux := runtime.NewServeMux()

		// Conectar al servidor gRPC (especifica la dirección gRPC)
		opts := []grpc.DialOption{grpc.WithInsecure()}
		err := product.RegisterProductServiceHandlerFromEndpoint(
			ctx,
			mux,
			"localhost:"+cfg.GRPC.Port, // Conéctate al servidor gRPC
			opts,
		)
		if err != nil {
			log.Fatalf("Error iniciando gRPC-Gateway: %v", err)
		}

		corsHandler := allowCORS(mux)

		// Servir el tráfico HTTP en el puerto configurado
		log.Printf("gRPC-Gateway REST escuchando en %s", cfg.HTTP.Port)
		if err := http.ListenAndServe(":"+cfg.HTTP.Port, corsHandler); err != nil {
			log.Fatalf("Fallo servidor HTTP: %v", err)
		}
	}()

	// Mantener el servidor activo
	select {}
}
