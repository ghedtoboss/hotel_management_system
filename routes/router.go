package routes

import (
	"hotel_management_system/controllers"
	"hotel_management_system/middleware"
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func InitRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/register", controllers.RegisterHandler).Methods("POST")
	r.HandleFunc("/login", controllers.LoginHandler).Methods("POST")
	r.Handle("/customers", middleware.JWTAuth(middleware.Authorize("admin", "receptionist")(http.HandlerFunc(controllers.GetCustomers)))).Methods("GET")
	r.Handle("/users/{user_id}", middleware.JWTAuth(middleware.Authorize("admin", "receptionist")(http.HandlerFunc(controllers.GetUser)))).Methods("GET")
	r.Handle("/users/{user_id}", middleware.JWTAuth(middleware.Authorize("admin", "receptionist")(http.HandlerFunc(controllers.UpdateUser)))).Methods("PUT")
	r.Handle("/users/{user_id}", middleware.JWTAuth(middleware.Authorize("admin", "receptionist")(http.HandlerFunc(controllers.DeleteUser)))).Methods("DELETE")
	r.Handle("/users", middleware.JWTAuth(middleware.Authorize("admin")(http.HandlerFunc(controllers.GetAllUsers)))).Methods("GET")
	r.Handle("/profile", middleware.JWTAuth(http.HandlerFunc(controllers.GetProfile))).Methods("GET")
	r.Handle("/profile", middleware.JWTAuth(http.HandlerFunc(controllers.UpdateProfile))).Methods("PUT")
	r.Handle("/profile/password", middleware.JWTAuth(http.HandlerFunc(controllers.UpdatePassword))).Methods("PUT")

	r.Handle("/rooms", middleware.JWTAuth(middleware.Authorize("admin", "receptionist")(http.HandlerFunc(controllers.CreateRoom)))).Methods("POST")
	r.Handle("/rooms/{room_id}", middleware.JWTAuth(middleware.Authorize("admin", "receptionist")(http.HandlerFunc(controllers.UpdateRoom)))).Methods("PUT")
	r.Handle("/rooms/{room_id}", middleware.JWTAuth(middleware.Authorize("admin", "receptionist")(http.HandlerFunc(controllers.DeleteRoom)))).Methods("DELETE")
	r.Handle("/rooms", middleware.JWTAuth(middleware.Authorize("admin", "receptionist")(http.HandlerFunc(controllers.GetRooms)))).Methods("GET")
	r.Handle("/rooms/{room_id}", middleware.JWTAuth(middleware.Authorize("admin", "receptionist")(http.HandlerFunc(controllers.GetRoomDetails)))).Methods("GET")

	r.Handle("/reservations", middleware.JWTAuth(middleware.Authorize("admin", "receptionist")(http.HandlerFunc(controllers.CreateReservation)))).Methods("POST")
	r.Handle("/reservations/{reservation_id}", middleware.JWTAuth(middleware.Authorize("admin", "receptionist")(http.HandlerFunc(controllers.UpdateReservation)))).Methods("PUT")
	r.Handle("/reservations/{reservation_id}", middleware.JWTAuth(middleware.Authorize("admin", "receptionist")(http.HandlerFunc(controllers.DeleteReservation)))).Methods("DELETE")
	r.Handle("/reservations", middleware.JWTAuth(middleware.Authorize("admin", "receptionist")(http.HandlerFunc(controllers.GetReservations)))).Methods("GET")
	r.Handle("/reservations/{reservation_id}", middleware.JWTAuth(middleware.Authorize("admin", "receptionist")(http.HandlerFunc(controllers.GetReservationDetails)))).Methods("GET")
	r.Handle("/reservations/status/{reservation_id}", middleware.JWTAuth(middleware.Authorize("admin", "receptionist")(http.HandlerFunc(controllers.UpdateReservationStatus)))).Methods("PUT")

	r.Handle("/occupancy", middleware.JWTAuth(middleware.Authorize("admin", "receptionist")(http.HandlerFunc(controllers.Occupancy)))).Methods("POST")
	r.Handle("/revenue", middleware.JWTAuth(middleware.Authorize("admin", "receptionist")(http.HandlerFunc(controllers.GetTotalRevenue)))).Methods("POST")
	r.Handle("/revenue/daily", middleware.JWTAuth(middleware.Authorize("admin", "receptionist")(http.HandlerFunc(controllers.GetDailyRevenue)))).Methods("POST")
	r.Handle("/revenue/monthly", middleware.JWTAuth(middleware.Authorize("admin", "receptionist")(http.HandlerFunc(controllers.GetMonthlyRevenue)))).Methods("POST")

	// Swagger endpoint
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/docs/swagger.json"), // The url pointing to API definition
	))

	// Static files endpoint for serving the swagger docs
	fs := http.FileServer(http.Dir("./docs"))
	r.PathPrefix("/docs/").Handler(http.StripPrefix("/docs/", fs))

	return r
}
