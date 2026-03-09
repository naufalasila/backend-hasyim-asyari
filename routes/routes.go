// routes/routes.go
package routes

import (
	"backend/controllers"
	"backend/middleware"
	"database/sql"
	"net/http"
)

func Setup(db *sql.DB) http.Handler {
	mux := http.NewServeMux()

	// Static file server
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	// === Auth Routes ===
	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		controllers.Register(db)(w, r)
	})

	mux.HandleFunc("/login", controllers.Login(db))

	mux.HandleFunc("/forgot-password", func(w http.ResponseWriter, r *http.Request) {
		controllers.ForgotPassword(db)(w, r)
	})

	mux.HandleFunc("/reset-password", func(w http.ResponseWriter, r *http.Request) {
		controllers.ResetPassword(db)(w, r)
	})

	mux.HandleFunc("/verify", func(w http.ResponseWriter, r *http.Request) {
		controllers.VerifyEmail(db)(w, r)
	})

	// === Protected Routes - User Role ===
	mux.Handle("/profile", middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		controllers.GetProfile(db)(w, r)
	}))

	mux.Handle("/profile/update", middleware.Auth(middleware.Role("user")(func(w http.ResponseWriter, r *http.Request) {
		controllers.UpdateProfile(db)(w, r)
	})))

	// === Pendaftaran Routes ===
	// Public route for submitting registration
	mux.HandleFunc("/api/pendaftar", controllers.SubmitPendaftaran(db))

	// Admin protected routes
	mux.Handle("/api/admin/pendaftar", middleware.Auth(middleware.Role("admin")(controllers.GetAllPendaftar(db))))
	mux.Handle("/api/admin/pendaftar/status", middleware.Auth(middleware.Role("admin")(controllers.UpdateStatusPendaftar(db))))

	// === Berita Routes ===
	// Public route for getting all berita
	mux.HandleFunc("/api/berita", controllers.GetAllBerita(db))

	// Admin protected routes for berita
	mux.Handle("/api/admin/berita", middleware.Auth(middleware.Role("admin")(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controllers.CreateBerita(db)(w, r)
		case http.MethodDelete:
			controllers.DeleteBerita(db)(w, r)
		default:
			controllers.GetAllBerita(db)(w, r)
		}
	})))

	// === Guru (Tim Pengajar) Routes ===
	// Public route for getting guru by jenjang
	mux.HandleFunc("/api/guru", controllers.GetGuru(db))

	// Admin protected routes for guru
	mux.Handle("/api/admin/guru", middleware.Auth(middleware.Role("admin")(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controllers.CreateGuru(db)(w, r)
		case http.MethodPut:
			controllers.UpdateGuru(db)(w, r)
		case http.MethodDelete:
			controllers.DeleteGuru(db)(w, r)
		default:
			controllers.GetGuru(db)(w, r)
		}
	})))

	return mux
}
