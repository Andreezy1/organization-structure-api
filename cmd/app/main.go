package main

import (
	"log"
	"net/http"
	"org_struct_api/internal/config"
	"org_struct_api/internal/database"
	"org_struct_api/internal/handler"
	"org_struct_api/internal/middleware"
	"org_struct_api/internal/repository"
	"org_struct_api/internal/service"
)

func main() {
	cfg := config.LoadConfig()

	dbconn := database.NewPostgresDB(cfg)

	departmentRepo := repository.NewDepartmentRepository(dbconn)
	employeeRepo := repository.NewEmployeeRepository(dbconn)
	transactionManager := repository.NewTransactionManager(dbconn)

	departmentService := service.NewDepartmentService(departmentRepo, employeeRepo, transactionManager)
	departmentHandler := handler.NewDepartmentHandler(departmentService)

	employeeService := service.NewEmployeeService(employeeRepo, departmentRepo)
	employeeHandler := handler.NewEmployeeHandler(employeeService)

	handler.RegisterRoute(departmentHandler, employeeHandler)
	loggedMux := middleware.Logger(http.DefaultServeMux)

	log.Printf("server started in :%s", cfg.AppPort)
	if err := http.ListenAndServe(":"+cfg.AppPort, loggedMux); err != nil {
		log.Fatalf("Server failed to start: %s\n", err)
	}

}
