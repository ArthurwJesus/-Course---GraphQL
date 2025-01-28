package main

import (
	"database/sql"
	"github.com/ArthurwJesus/graphql/internal/database"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/ArthurwJesus/graphql/graph"
	"github.com/vektah/gqlparser/v2/ast"

	//Adicionando pacote do SQLITE
	_ "github.com/mattn/go-sqlite3"
)

const defaultPort = "8080"

func main() {

	//Criando conexão com o banco de dados
	db, err := sql.Open("sqlite3", "./data.db")

	if err != nil {
		log.Fatal("Falhar ao conectar ao database: %v", err)
	}

	defer db.Close()

	//Criando a conexão para a criação de categoria
	categoryDB := database.NewCategory(db)
	courseDB := database.NewCourse(db) // ADD para courses

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		//resolver injetado
		CategoryDB: categoryDB,
		CourseDB:   courseDB, //ADD para courses
	}}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
