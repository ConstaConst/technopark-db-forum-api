// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"github.com/ConstaConst/technopark-db-forum-api/postgres"
	"log"
	"net/http"

	"github.com/ConstaConst/technopark-db-forum-api/restapi/operations"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
)

//go:generate swagger generate server --target .. --name  --spec ../swagger.yml

func configureFlags(api *operations.ForumAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.ForumAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.BinConsumer = runtime.ByteStreamConsumer()

	api.JSONProducer = runtime.JSONProducer()

	dbConn, err := postgres.MakeDBConn()
	if err != nil {
		dbConn.Close()
		log.Fatalln("Can't create db connection: ", err)
	}

	api.ClearHandler = operations.ClearHandlerFunc(dbConn.ClearService)
	api.ForumCreateHandler = operations.ForumCreateHandlerFunc(dbConn.CreateForum)
	api.ForumGetOneHandler = operations.ForumGetOneHandlerFunc(dbConn.GetOneForum)
	api.ForumGetThreadsHandler = operations.ForumGetThreadsHandlerFunc(dbConn.GetForumThreads)
	api.ForumGetUsersHandler = operations.ForumGetUsersHandlerFunc(dbConn.GetForumUsers)
	api.PostGetOneHandler = operations.PostGetOneHandlerFunc(dbConn.GetOnePost)
	api.PostUpdateHandler = operations.PostUpdateHandlerFunc(dbConn.UpdatePost)
	api.PostsCreateHandler = operations.PostsCreateHandlerFunc(dbConn.CreatePosts)
	api.StatusHandler = operations.StatusHandlerFunc(dbConn.ShowStatus)
	api.ThreadCreateHandler = operations.ThreadCreateHandlerFunc(dbConn.CreateThread)
	api.ThreadGetOneHandler = operations.ThreadGetOneHandlerFunc(dbConn.GetOneThread)
	api.ThreadGetPostsHandler = operations.ThreadGetPostsHandlerFunc(dbConn.GetPosts)
	api.ThreadUpdateHandler = operations.ThreadUpdateHandlerFunc(dbConn.UpdateThread)
	api.ThreadVoteHandler = operations.ThreadVoteHandlerFunc(dbConn.ThreadVote)
	api.UserCreateHandler = operations.UserCreateHandlerFunc(dbConn.CreateUser)
	api.UserGetOneHandler = operations.UserGetOneHandlerFunc(dbConn.GetOneUser)
	api.UserUpdateHandler = operations.UserUpdateHandlerFunc(dbConn.UpdateUser)

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
