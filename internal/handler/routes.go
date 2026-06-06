package handler

import "net/http"

func RegisterRoute(departmentHandler *DepartmentHandler, employeeHandler *EmployeeHandler) {
	http.HandleFunc("POST /departments", departmentHandler.CreateDepartment)
	http.HandleFunc("POST /departments/{id}/employees", employeeHandler.CreateEmployee)

	http.HandleFunc("GET /departments/{id}", departmentHandler.GetDepartment)
	http.HandleFunc("GET /departments/{id}/employees", employeeHandler.GetDepartmentEmployees)

	http.HandleFunc("PATCH /departments/{id}", departmentHandler.UpdateDepartment)

	http.HandleFunc("DELETE /departments/{id}", departmentHandler.DeleteDepartment)

}
