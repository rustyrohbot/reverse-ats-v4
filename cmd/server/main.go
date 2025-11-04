package main

import (
	"log"
	"net/http"
	"os"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"

	"reverse-ats/internal/handlers"
	_ "reverse-ats/pb_migrations"
)

func main() {
	// Initialize PocketBase
	app := pocketbase.New()

	// Get port from environment
	port := os.Getenv("REVERSE_ATS_PORT")
	if port == "" {
		port = "5627"
	}

	// Hook into the serve event to add custom routes
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		// Create handlers with PocketBase app
		companiesHandler := handlers.NewCompaniesHandler(app)
		rolesHandler := handlers.NewRolesHandler(app)
		contactsHandler := handlers.NewContactsHandler(app)
		interviewsHandler := handlers.NewInterviewsHandler(app)
		statsHandler := handlers.NewStatsHandler(app)
		exportHandler := handlers.NewExportHandler(app)
		importHandler := handlers.NewImportHandler(app)

		// Static files - serve from ./static directory
		se.Router.GET("/static/{path...}", func(e *core.RequestEvent) error {
			// Get the file path from the URL
			path := e.Request.URL.Path
			// Serve the static file
			http.ServeFile(e.Response, e.Request, "."+path)
			return nil
		})

		// Page routes
		se.Router.GET("/", func(e *core.RequestEvent) error {
			return e.Redirect(http.StatusFound, "/companies")
		})

		// Companies routes
		se.Router.GET("/companies", func(e *core.RequestEvent) error {
			return companiesHandler.List(e.Response, e.Request)
		})
		se.Router.POST("/companies", func(e *core.RequestEvent) error {
			return companiesHandler.Create(e.Response, e.Request)
		})
		se.Router.GET("/companies/new", func(e *core.RequestEvent) error {
			return companiesHandler.New(e.Response, e.Request)
		})
		se.Router.GET("/companies/{id}/edit", func(e *core.RequestEvent) error {
			return companiesHandler.Edit(e.Response, e.Request)
		})
		se.Router.POST("/companies/{id}", func(e *core.RequestEvent) error {
			return companiesHandler.Update(e.Response, e.Request)
		})
		se.Router.PUT("/companies/{id}", func(e *core.RequestEvent) error {
			return companiesHandler.Update(e.Response, e.Request)
		})
		se.Router.DELETE("/companies/{id}", func(e *core.RequestEvent) error {
			return companiesHandler.Delete(e.Response, e.Request)
		})

		// Roles routes
		se.Router.GET("/roles", func(e *core.RequestEvent) error {
			return rolesHandler.List(e.Response, e.Request)
		})
		se.Router.POST("/roles", func(e *core.RequestEvent) error {
			return rolesHandler.Create(e.Response, e.Request)
		})
		se.Router.GET("/roles/new", func(e *core.RequestEvent) error {
			return rolesHandler.New(e.Response, e.Request)
		})
		se.Router.GET("/roles/{id}/edit", func(e *core.RequestEvent) error {
			return rolesHandler.Edit(e.Response, e.Request)
		})
		se.Router.POST("/roles/{id}", func(e *core.RequestEvent) error {
			return rolesHandler.Update(e.Response, e.Request)
		})
		se.Router.PUT("/roles/{id}", func(e *core.RequestEvent) error {
			return rolesHandler.Update(e.Response, e.Request)
		})
		se.Router.DELETE("/roles/{id}", func(e *core.RequestEvent) error {
			return rolesHandler.Delete(e.Response, e.Request)
		})

		// Contacts routes
		se.Router.GET("/contacts", func(e *core.RequestEvent) error {
			return contactsHandler.List(e.Response, e.Request)
		})
		se.Router.POST("/contacts", func(e *core.RequestEvent) error {
			return contactsHandler.Create(e.Response, e.Request)
		})
		se.Router.GET("/contacts/new", func(e *core.RequestEvent) error {
			return contactsHandler.New(e.Response, e.Request)
		})
		se.Router.GET("/contacts/{id}/edit", func(e *core.RequestEvent) error {
			return contactsHandler.Edit(e.Response, e.Request)
		})
		se.Router.POST("/contacts/{id}", func(e *core.RequestEvent) error {
			return contactsHandler.Update(e.Response, e.Request)
		})
		se.Router.PUT("/contacts/{id}", func(e *core.RequestEvent) error {
			return contactsHandler.Update(e.Response, e.Request)
		})
		se.Router.DELETE("/contacts/{id}", func(e *core.RequestEvent) error {
			return contactsHandler.Delete(e.Response, e.Request)
		})

		// Interviews routes
		se.Router.GET("/interviews", func(e *core.RequestEvent) error {
			return interviewsHandler.List(e.Response, e.Request)
		})
		se.Router.POST("/interviews", func(e *core.RequestEvent) error {
			return interviewsHandler.Create(e.Response, e.Request)
		})
		se.Router.GET("/interviews/new", func(e *core.RequestEvent) error {
			return interviewsHandler.New(e.Response, e.Request)
		})
		se.Router.GET("/interviews/{id}/edit", func(e *core.RequestEvent) error {
			return interviewsHandler.Edit(e.Response, e.Request)
		})
		se.Router.POST("/interviews/{id}", func(e *core.RequestEvent) error {
			return interviewsHandler.Update(e.Response, e.Request)
		})
		se.Router.PUT("/interviews/{id}", func(e *core.RequestEvent) error {
			return interviewsHandler.Update(e.Response, e.Request)
		})
		se.Router.DELETE("/interviews/{id}", func(e *core.RequestEvent) error {
			return interviewsHandler.Delete(e.Response, e.Request)
		})

		// API route for cascading dropdowns
		se.Router.GET("/api/roles-by-company", func(e *core.RequestEvent) error {
			return interviewsHandler.GetRolesByCompany(e.Response, e.Request)
		})

		// Stats route
		se.Router.GET("/stats", func(e *core.RequestEvent) error {
			return statsHandler.Show(e.Response, e.Request)
		})

		// Export route
		se.Router.GET("/export", func(e *core.RequestEvent) error {
			return exportHandler.Export(e.Response, e.Request)
		})

		// Import route
		se.Router.POST("/import", func(e *core.RequestEvent) error {
			return importHandler.Import(e.Response, e.Request)
		})

		return se.Next()
	})

	// Start the server
	log.Printf("Starting server on http://localhost:%s", port)
	log.Printf("PocketBase admin UI available at http://localhost:%s/_/", port)

	// Set command line args to force serve mode
	// PocketBase expects: program serve --http=host:port
	os.Args = []string{os.Args[0], "serve", "--http=0.0.0.0:" + port}

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
