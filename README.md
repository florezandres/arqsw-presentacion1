
## 🚀 Project Structure

Inside of your Astro project, you'll see the following folders and files:

```text
├── Venta (Go))
│   
├── .env
│   
├── docker-compose.yml
│   
├── Taller2.iml
│   
├── Product/
	│
	├── buf.gen.yaml
	│
	├── Dockerfile
	│
	├── go.mod
	│   ├── go.sum
	│
	├── cmd/
	│   ├── main.go	
	│   
	├── gen/                 # Capa GRPC
	│   ├── go/
	│       ├── internal/ 
	│           ├── api/ 
	│               ├──product_service.pb.go 
	│               ├──product_service_grpc.pb.go
	│   	
	├── internal/
	    │
	    ├── shared/
	    │   ├── init_db/
	    │   │   ├── command/
	    │   │   │    ├── init.sql
	    │   │   │
	    │   │   ├── query/
	    │   │        ├── init.sql
	    │   │   
	    │   ├── config.go
	    │
	    ├── api/
	    │   ├── grpc_server.go
	    │   ├── product_service.proto
	    │
	    ├── command/
	    │   ├── application/          # Casos de uso (comandos)
	    │   │   ├── commands/        # Definiciones de comandos
	    │   │   │    ├── create_product_command.go
	    │   │   └── handlers/       # Manejadores de comandos
	    │   │        ├── create_product_handler.go
	    │   │
	    │   ├── domain/
	    │   │   ├── entities/        # Entidades de negocio (Command)
	    │   │   │    ├── product.go
	    │   │   ├── events/         # Eventos de dominio
	    │   │   │    ├── product_events.go
	    │   │   └── repositories/  # Interfaces de repositorio
	    │   │        ├── product_repository.go
	    │   │
	    │   └── infrastructure/     # Implementaciones (Kafka, PostgreSQL)
	    │   │    ├── kafka/ 
	    │   │    │    ├── producer.go
	    │   │    ├── persistence/
	    │   │         ├── postgres_repository.go
	    │   │
	    │   ├── doc.go
	    │
	    ├── query/	
	        ├── application/         # Consultas (queries)
	        │   ├── queries/       # Definiciones de queries
	        │   │    ├── get_product_query.go
	        │   └── handlers/     # Manejadores de queries
	        │        ├── get_product_handler.go
	        │        
	        ├── domain/
	        │   ├── models/       # Modelos de lectura (Query)
	        │   │    ├── product_read.go
	        │   └── repositories/ # Interfaces de repositorio
	        │        ├── product_repository.go
	        │        
	        └── infrastructure/   # Implementaciones (PostgreSQL, Kafka Consumer)
	        │   ├── kafka/ 
	        │   │    ├── consumer.go
	        │   ├── persistence/
	        │        ├── postgres_repository.go
	        │        
	        ├── doc.go

```

To learn more about the folder structure of an Astro project, refer to [our guide on project structure](https://docs.astro.build/en/basics/project-structure/).

## 🧞 Commands
All commands are run from the root of the project, from a terminal:

| Command                              | Action                                                 |
|:-------------------------------------|:-------------------------------------------------------|
| `docker-compose up -d`               | Runs zookeeper and kafka on docker                     |
| `go run Product/cmd/main.go`         | Starts microservice listening in port `50051` |                      |

## 👀 Want to learn more?

Feel free to check [our documentation](https://docs.astro.build) or jump into our [Discord server](https://astro.build/chat).
