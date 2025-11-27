package app

import (
	"net/http"
	"os"
	"synk/gateway/app/controller"
	"synk/gateway/app/util"
)

func Router(service *Service) {
	aboutController := controller.NewAbout(service.DB)
	userController := controller.NewUsers(service.DB)

	http.HandleFunc("GET /about", aboutController.HandleAbout)
	http.HandleFunc("GET /users", userController.HandleShow)
	http.HandleFunc("POST /users/register", userController.HandleRegister)
	http.HandleFunc("POST /users/login", userController.HandleLogin)
	http.HandleFunc("GET /users/refresh", userController.HandleRefresh)
	http.HandleFunc("GET /users/check", userController.HandleCheck)
	http.HandleFunc("GET /users/logout", userController.HandleLogout)

	port := os.Getenv("PORT")
	util.Log("app running on port " + port)

	var err error
	env := os.Getenv("ENV")

	if env == "production" {
		util.Log("Running in PRODUCTION mode (HTTP)")
		err = http.ListenAndServe(":"+port, controller.Cors(http.DefaultServeMux))
	} else {
		util.Log("Running in DEV mode (HTTPS)")
		err = http.ListenAndServeTLS(
			":"+port,
			"/cert/cert.pem",
			"/cert/key.pem",
			controller.Cors(http.DefaultServeMux),
		)
	}

	if err != nil {
		util.Log("app failed on running on port " + port + ": " + err.Error())
	}
}
