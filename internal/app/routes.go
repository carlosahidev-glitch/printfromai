package app

import (
	"net/http"
	// We will uncomment this import in our next step once we create our handlers

	"github.com/printfromai/core/internal/handlers"
)

// NewRouter initializes the main HTTP multiplexer and registers all routes.
// This keeps our main.go file clean and gives us a single map of our entire application.
func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	// 1. Static Asset Server (CSS, JS, Images)
	// This serves files from our web/static directory so the browser can load our Tailwind CSS and images.
	staticDir := http.Dir("web/static")
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(staticDir)))

	// 2. Public Core Pages (Full Page Loads)
	mux.HandleFunc("GET /", handlers.Home)
	// mux.HandleFunc("GET /faq", handlers.FAQ)
	// mux.HandleFunc("GET /how-it-works", handlers.HowItWorks)
	// mux.HandleFunc("GET /about-us", handlers.AboutUs)

	// 3. Product & Print Flow
	// Using Go 1.22 path variables. E.g., /product/business-cards
	mux.HandleFunc("GET /product/{slug}", handlers.ProductPage)
	// E-Commerce & Cart Routes
	mux.HandleFunc("GET /cart", handlers.GetCart)
	mux.HandleFunc("POST /cart/add", handlers.AddToCart)

	// 4. HTMX Dynamic Endpoints (The Magic)
	// HTMX on the frontend will call these via AJAX.
	// Instead of returning JSON or a full page, these return small HTML partials to swap into the UI.
	// mux.HandleFunc("POST /htmx/calculate-price", handlers.HTMXCalculatePrice)
	// mux.HandleFunc("POST /htmx/toggle-back-design", handlers.HTMXToggleBackDesign)
	// mux.HandleFunc("POST /htmx/upload-art", handlers.HTMXUploadArt)

	// 5. Authentication (Firebase integration)
	// mux.HandleFunc("GET /login", handlers.LoginView)
	// mux.HandleFunc("POST /login", handlers.LoginAction) // Email/Password
	// mux.HandleFunc("GET /auth/google", handlers.GoogleAuthAction)

	// 6. Protected Admin Routes & SEO
	// mux.HandleFunc("GET /admin", handlers.AdminDashboard)
	// mux.HandleFunc("GET /llms.txt", handlers.LLMsTxt) // For Generative Engine Optimization

	return mux
}
