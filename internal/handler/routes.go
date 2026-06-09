package handler

import "net/http"

func RegisterRoute(departmentHandler *DepartmentHandler, employeeHandler *EmployeeHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /departments", departmentHandler.CreateDepartment)
	mux.HandleFunc("POST /departments/", departmentHandler.CreateDepartment)

	mux.HandleFunc("POST /departments/{id}/employees", employeeHandler.CreateEmployee)
	mux.HandleFunc("POST /departments/{id}/employees/", employeeHandler.CreateEmployee)

	mux.HandleFunc("GET /departments/{id}", departmentHandler.GetDepartment)
	mux.HandleFunc("GET /departments/{id}/", departmentHandler.GetDepartment)

	mux.HandleFunc("PATCH /departments/{id}", departmentHandler.UpdateDepartment)
	mux.HandleFunc("PATCH /departments/{id}/", departmentHandler.UpdateDepartment)

	mux.HandleFunc("DELETE /departments/{id}", departmentHandler.DeleteDepartment)
	mux.HandleFunc("DELETE /departments/{id}/", departmentHandler.DeleteDepartment)

	return mux
}
