package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"org_struct_api/internal/models"
	"testing"
)

type mockDepartmentService struct{}

func (m *mockDepartmentService) GetDepartment(id uint, depth int, includeEmployees bool) (*models.Department, error) {
	return &models.Department{
		ID:   1,
		Name: "IT",
	}, nil
}

func (m *mockDepartmentService) CreateDepartment(department *models.Department) (*models.Department, error) {
	return nil, nil
}

func (m *mockDepartmentService) UpdateDepartment(id uint, name string, parentID *uint) (*models.Department, error) {
	return nil, nil
}

func (m *mockDepartmentService) DeleteDepartment(departmentID uint, mode string, reassignTo *uint) error {
	return nil
}

func (m *mockDepartmentService) GetDepartmentTree() ([]models.Department, error) {
	return nil, nil
}

func TestGetDepartment_InvalidID(t *testing.T) {
	handler := NewDepartmentHandler(&mockDepartmentService{})

	req := httptest.NewRequest(http.MethodGet, "/departments/abc", nil)

	req.SetPathValue("id", "abc")

	rr := httptest.NewRecorder()

	handler.GetDepartment(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d",
			http.StatusBadRequest,
			rr.Code)
	}
}

func TestGetDepartment_WrongMethod(t *testing.T) {
	h := NewDepartmentHandler(&mockDepartmentService{})

	req := httptest.NewRequest(http.MethodPost, "/departments/1", nil)

	req.SetPathValue("id", "1")

	rr := httptest.NewRecorder()

	h.GetDepartment(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected %d got %d",
			http.StatusMethodNotAllowed,
			rr.Code,
		)
	}
}

func TestGetDepartment_Success(t *testing.T) {
	h := NewDepartmentHandler(&mockDepartmentService{})

	req := httptest.NewRequest(http.MethodGet, "/departments/1", nil)

	req.SetPathValue("id", "1")

	rr := httptest.NewRecorder()

	h.GetDepartment(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf(
			"expected %d got %d",
			http.StatusOK,
			rr.Code,
		)
	}

	body := rr.Body.String()

	if !strings.Contains(body, `"name":"IT"`) {
		t.Fatalf("unexpected response body: %s", body)
	}
}
