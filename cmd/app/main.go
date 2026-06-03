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

	departmentService := service.NewDepartmentService(departmentRepo, employeeRepo)
	departmentHandler := handler.NewDepartmentHandler(departmentService)

	employeeService := service.NewEmployeeService(employeeRepo, departmentRepo)
	employeeHandler := handler.NewEmployeeHandler(employeeService)

	http.HandleFunc("POST /departments", departmentHandler.CreateDepartment)
	http.HandleFunc("POST /departments/{id}/employees", employeeHandler.CreateEmployee)

	http.HandleFunc("GET /departments/{id}", departmentHandler.GetDepartment)
	http.HandleFunc("GET /departments/{id}/employees", employeeHandler.GetDepartmentEmployees)

	http.HandleFunc("PATCH /departments/{id}", departmentHandler.UpdateDepartment)

	http.HandleFunc("DELETE /departments/{id}", departmentHandler.DeleteDepartment)

	loggedMux := middleware.Logger(http.DefaultServeMux)

	log.Println("server started in :8080")
	if err := http.ListenAndServe(":"+cfg.AppPort, loggedMux); err != nil {
		log.Fatalf("Server failed to start: %s\n", err)
	}

}
