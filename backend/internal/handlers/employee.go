package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"tatapps/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type EmployeeHandler struct {
	db *gorm.DB
}

func NewEmployeeHandler(db *gorm.DB) *EmployeeHandler {
	return &EmployeeHandler{db: db}
}

type employeePayload struct {
	EmployeeCode    string  `json:"employeeCode"`
	NIK             string  `json:"nik"`
	FullName        string  `json:"fullName" binding:"required"`
	BirthPlace      string  `json:"birthPlace"`
	BirthDate       *string `json:"birthDate"`
	Gender          string  `json:"gender"`
	BloodType       string  `json:"bloodType"`
	MaritalStatus   string  `json:"maritalStatus"`
	Religion        string  `json:"religion"`
	IdentityType    string  `json:"identityType"`
	IdentityNumber  string  `json:"identityNumber"`
	AddressKtp      string  `json:"addressKtp"`
	AddressDomicile string  `json:"addressDomicile"`
	Address         string  `json:"address"`
	Phone           string  `json:"phone"`
	Email           string  `json:"email"`
	Timezone        string  `json:"timezone"`
	DivisionID      *uint   `json:"divisionId"`
	PositionID      *uint   `json:"positionId"`
	EmploymentType  string  `json:"employmentType"`
	Status          string  `json:"status"`
	JoinDate        *string `json:"joinDate"`
	Photo           string  `json:"photo"`
}

type employeeResponse struct {
	ID              uint      `json:"id"`
	EmployeeCode    string    `json:"employeeCode"`
	NIK             string    `json:"nik"`
	FullName        string    `json:"fullName"`
	BirthPlace      string    `json:"birthPlace"`
	BirthDate       *string   `json:"birthDate"`
	Gender          string    `json:"gender"`
	BloodType       string    `json:"bloodType"`
	MaritalStatus   string    `json:"maritalStatus"`
	Religion        string    `json:"religion"`
	IdentityType    string    `json:"identityType"`
	IdentityNumber  string    `json:"identityNumber"`
	AddressKtp      string    `json:"addressKtp"`
	AddressDomicile string    `json:"addressDomicile"`
	Address         string    `json:"address"`
	Phone           string    `json:"phone"`
	Email           string    `json:"email"`
	Timezone        string    `json:"timezone"`
	DivisionID      *uint     `json:"divisionId"`
	DivisionName    string    `json:"divisionName"`
	PositionID      *uint     `json:"positionId"`
	PositionTitle   string    `json:"positionTitle"`
	EmploymentType  string    `json:"employmentType"`
	Status          string    `json:"status"`
	JoinDate        *string   `json:"joinDate"`
	Photo           string    `json:"photo"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type divisionPayload struct {
	Name              string `json:"name" binding:"required"`
	Description       string `json:"description"`
	RecruitmentStatus string `json:"recruitmentStatus"`
	HeadEmployeeID    *uint  `json:"headEmployeeId"`
	HeadPositionID    *uint  `json:"headPositionId"`
	Head              string `json:"head"`
	HeadTitle         string `json:"headTitle"`
}

type divisionResponse struct {
	ID                uint      `json:"id"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	RecruitmentStatus string    `json:"recruitmentStatus"`
	HeadEmployeeID    *uint     `json:"headEmployeeId"`
	HeadPositionID    *uint     `json:"headPositionId"`
	Head              string    `json:"head"`
	HeadTitle         string    `json:"headTitle"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

type positionPayload struct {
	Title       string `json:"title" binding:"required"`
	Code        string `json:"code"`
	DivisionID  *uint  `json:"divisionId"`
	ParentID    *uint  `json:"parentId"`
	Notes       string `json:"notes"`
	Grade       string `json:"grade"`
	SalaryRange string `json:"salaryRange"`
}

type positionResponse struct {
	ID           uint      `json:"id"`
	Title        string    `json:"title"`
	Code         string    `json:"code"`
	DivisionID   *uint     `json:"divisionId"`
	DivisionName string    `json:"divisionName"`
	ParentID     *uint     `json:"parentId"`
	Notes        string    `json:"notes"`
	Grade        string    `json:"grade"`
	SalaryRange  string    `json:"salaryRange"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type batchDeleteRequest struct {
	IDs []uint `json:"ids"`
}

func (h *EmployeeHandler) ListEmployees(c *gin.Context) {
	var employees []models.Employee

	query := h.db.
		Preload("Division").
		Preload("Position")

	search := strings.TrimSpace(c.Query("search"))
	if search != "" {
		q := "%" + strings.ToLower(search) + "%"
		query = query.Where(
			"LOWER(full_name) LIKE ? OR LOWER(employee_code) LIKE ? OR LOWER(email) LIKE ? OR LOWER(phone) LIKE ? OR LOWER(nik) LIKE ?",
			q, q, q, q, q,
		)
	}

	if divisionID := strings.TrimSpace(c.Query("division_id")); divisionID != "" {
		query = query.Where("division_id = ?", divisionID)
	}

	if positionID := strings.TrimSpace(c.Query("position_id")); positionID != "" {
		query = query.Where("position_id = ?", positionID)
	}

	if status := strings.TrimSpace(c.Query("status")); status != "" {
		query = query.Where("status = ?", status)
	}

	if employment := strings.TrimSpace(c.Query("employment_type")); employment != "" {
		query = query.Where("employment_type = ?", employment)
	}

	if err := query.Order("full_name ASC").Find(&employees).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch employees"})
		return
	}

	responses := make([]employeeResponse, 0, len(employees))
	for i := range employees {
		responses = append(responses, buildEmployeeResponse(&employees[i]))
	}

	c.JSON(http.StatusOK, gin.H{"data": responses})
}

func (h *EmployeeHandler) GetEmployee(c *gin.Context) {
	id := c.Param("id")
	var employee models.Employee
	if err := h.db.Preload("Division").Preload("Position").First(&employee, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch employee"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": buildEmployeeResponse(&employee)})
}

func (h *EmployeeHandler) CreateEmployee(c *gin.Context) {
	var payload employeePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	employee, err := h.buildEmployeeModel(payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Create(&employee).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create employee"})
		return
	}

	if err := h.db.Preload("Division").Preload("Position").First(&employee, employee.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load employee"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": buildEmployeeResponse(&employee)})
}

func (h *EmployeeHandler) UpdateEmployee(c *gin.Context) {
	id := c.Param("id")

	var payload employeePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var employee models.Employee
	if err := h.db.First(&employee, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch employee"})
		return
	}

	updated, err := h.buildEmployeeModel(payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	employee.EmployeeCode = updated.EmployeeCode
	employee.NIK = updated.NIK
	employee.FullName = updated.FullName
	employee.BirthPlace = updated.BirthPlace
	employee.BirthDate = updated.BirthDate
	employee.Gender = updated.Gender
	employee.BloodType = updated.BloodType
	employee.MaritalStatus = updated.MaritalStatus
	employee.Religion = updated.Religion
	employee.IdentityType = updated.IdentityType
	employee.IdentityNumber = updated.IdentityNumber
	employee.AddressKTP = updated.AddressKTP
	employee.AddressDomicile = updated.AddressDomicile
	employee.Address = updated.Address
	employee.Phone = updated.Phone
	employee.Email = updated.Email
	employee.Timezone = updated.Timezone
	employee.DivisionID = updated.DivisionID
	employee.PositionID = updated.PositionID
	employee.EmploymentType = updated.EmploymentType
	employee.Status = updated.Status
	employee.JoinDate = updated.JoinDate
	employee.Photo = updated.Photo

	if err := h.db.Save(&employee).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update employee"})
		return
	}

	if err := h.db.Preload("Division").Preload("Position").First(&employee, employee.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load employee"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": buildEmployeeResponse(&employee)})
}

func (h *EmployeeHandler) DeleteEmployee(c *gin.Context) {
	idParam := c.Param("id")
	parsed, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee id"})
		return
	}
	id := uint(parsed)

	if err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := h.cleanupAfterEmployeeDelete(tx, []uint{id}); err != nil {
			return err
		}
		return tx.Delete(&models.Employee{}, id).Error
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete employee"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *EmployeeHandler) DeleteEmployeesBatch(c *gin.Context) {
	var payload batchDeleteRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(payload.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No employee ids provided"})
		return
	}

	if err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := h.cleanupAfterEmployeeDelete(tx, payload.IDs); err != nil {
			return err
		}
		return tx.Where("id IN ?", payload.IDs).Delete(&models.Employee{}).Error
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete employees"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *EmployeeHandler) ListDivisions(c *gin.Context) {
	var divisions []models.EmployeeDivision
	if err := h.db.Order("name ASC").Find(&divisions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch divisions"})
		return
	}

	responses := make([]divisionResponse, 0, len(divisions))
	for i := range divisions {
		responses = append(responses, buildDivisionResponse(&divisions[i]))
	}

	c.JSON(http.StatusOK, gin.H{"data": responses})
}

func (h *EmployeeHandler) CreateDivision(c *gin.Context) {
	var payload divisionPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	division := models.EmployeeDivision{
		Name:              strings.TrimSpace(payload.Name),
		Description:       strings.TrimSpace(payload.Description),
		RecruitmentStatus: defaultString(payload.RecruitmentStatus, "Stabil"),
		HeadEmployeeID:    payload.HeadEmployeeID,
		HeadPositionID:    payload.HeadPositionID,
		Head:              strings.TrimSpace(payload.Head),
		HeadTitle:         strings.TrimSpace(payload.HeadTitle),
	}

	if err := h.db.Create(&division).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create division"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": buildDivisionResponse(&division)})
}

func (h *EmployeeHandler) UpdateDivision(c *gin.Context) {
	id := c.Param("id")

	var payload divisionPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var division models.EmployeeDivision
	if err := h.db.First(&division, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Division not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch division"})
		return
	}

	division.Name = strings.TrimSpace(payload.Name)
	division.Description = strings.TrimSpace(payload.Description)
	division.RecruitmentStatus = defaultString(payload.RecruitmentStatus, "Stabil")
	division.HeadEmployeeID = payload.HeadEmployeeID
	division.HeadPositionID = payload.HeadPositionID
	division.Head = strings.TrimSpace(payload.Head)
	division.HeadTitle = strings.TrimSpace(payload.HeadTitle)

	if err := h.db.Save(&division).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update division"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": buildDivisionResponse(&division)})
}

func (h *EmployeeHandler) DeleteDivision(c *gin.Context) {
	id := c.Param("id")

	if err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Employee{}).Where("division_id = ?", id).Update("division_id", nil).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.EmployeePosition{}).Where("division_id = ?", id).Update("division_id", nil).Error; err != nil {
			return err
		}
		return tx.Delete(&models.EmployeeDivision{}, id).Error
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete division"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *EmployeeHandler) ListPositions(c *gin.Context) {
	var positions []models.EmployeePosition

	query := h.db.Order("title ASC")
	if divisionID := strings.TrimSpace(c.Query("division_id")); divisionID != "" {
		query = query.Where("division_id = ?", divisionID)
	}

	if err := query.Preload("Division").Find(&positions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch positions"})
		return
	}

	responses := make([]positionResponse, 0, len(positions))
	for i := range positions {
		responses = append(responses, buildPositionResponse(&positions[i]))
	}

	c.JSON(http.StatusOK, gin.H{"data": responses})
}

func (h *EmployeeHandler) CreatePosition(c *gin.Context) {
	var payload positionPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	position := models.EmployeePosition{
		Title:       strings.TrimSpace(payload.Title),
		Code:        strings.TrimSpace(payload.Code),
		DivisionID:  payload.DivisionID,
		ParentID:    payload.ParentID,
		Notes:       strings.TrimSpace(payload.Notes),
		Grade:       strings.TrimSpace(payload.Grade),
		SalaryRange: strings.TrimSpace(payload.SalaryRange),
	}

	if err := h.db.Create(&position).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create position"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": buildPositionResponse(&position)})
}

func (h *EmployeeHandler) UpdatePosition(c *gin.Context) {
	id := c.Param("id")

	var payload positionPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if payload.ParentID != nil {
		if parsed, err := strconv.ParseUint(id, 10, 64); err == nil && *payload.ParentID == uint(parsed) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Position cannot be its own parent"})
			return
		}
	}

	var position models.EmployeePosition
	if err := h.db.First(&position, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Position not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch position"})
		return
	}

	position.Title = strings.TrimSpace(payload.Title)
	position.Code = strings.TrimSpace(payload.Code)
	position.DivisionID = payload.DivisionID
	position.ParentID = payload.ParentID
	position.Notes = strings.TrimSpace(payload.Notes)
	position.Grade = strings.TrimSpace(payload.Grade)
	position.SalaryRange = strings.TrimSpace(payload.SalaryRange)

	if err := h.db.Save(&position).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update position"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": buildPositionResponse(&position)})
}

func (h *EmployeeHandler) DeletePosition(c *gin.Context) {
	id := c.Param("id")

	if err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Employee{}).Where("position_id = ?", id).Update("position_id", nil).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.EmployeePosition{}).Where("parent_id = ?", id).Update("parent_id", nil).Error; err != nil {
			return err
		}
		return tx.Delete(&models.EmployeePosition{}, id).Error
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete position"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *EmployeeHandler) cleanupAfterEmployeeDelete(tx *gorm.DB, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}

	if err := tx.Model(&models.EmployeeDivision{}).
		Where("head_employee_id IN ?", ids).
		Updates(map[string]interface{}{
			"head_employee_id": nil,
			"head_position_id": nil,
			"head":             "",
			"head_title":       "",
		}).Error; err != nil {
		return err
	}

	if err := tx.Model(&models.Employee{}).Where("manager_id IN ?", ids).Update("manager_id", nil).Error; err != nil {
		return err
	}

	if err := tx.Model(&models.Warehouse{}).Where("manager_id IN ?", ids).Update("manager_id", nil).Error; err != nil {
		return err
	}

	return nil
}

func (h *EmployeeHandler) buildEmployeeModel(payload employeePayload) (models.Employee, error) {
	now := time.Now()
	employeeCode := strings.TrimSpace(payload.EmployeeCode)
	nik := strings.TrimSpace(payload.NIK)
	if employeeCode == "" && nik != "" {
		employeeCode = nik
	}
	if nik == "" {
		nik = employeeCode
	}

	joinDate := now
	if payload.JoinDate != nil && strings.TrimSpace(*payload.JoinDate) != "" {
		parsed, err := parseISOTime(*payload.JoinDate)
		if err != nil {
			return models.Employee{}, err
		}
		joinDate = *parsed
	}

	var birthDate *time.Time
	if payload.BirthDate != nil && strings.TrimSpace(*payload.BirthDate) != "" {
		parsed, err := parseISOTime(*payload.BirthDate)
		if err != nil {
			return models.Employee{}, err
		}
		birthDate = parsed
	}

	return models.Employee{
		EmployeeCode:    employeeCode,
		NIK:             nik,
		FullName:        strings.TrimSpace(payload.FullName),
		BirthPlace:      strings.TrimSpace(payload.BirthPlace),
		BirthDate:       birthDate,
		Gender:          strings.TrimSpace(payload.Gender),
		BloodType:       strings.TrimSpace(payload.BloodType),
		MaritalStatus:   strings.TrimSpace(payload.MaritalStatus),
		Religion:        strings.TrimSpace(payload.Religion),
		IdentityType:    strings.TrimSpace(payload.IdentityType),
		IdentityNumber:  strings.TrimSpace(payload.IdentityNumber),
		AddressKTP:      strings.TrimSpace(payload.AddressKtp),
		AddressDomicile: strings.TrimSpace(payload.AddressDomicile),
		Address:         strings.TrimSpace(payload.Address),
		Phone:           strings.TrimSpace(payload.Phone),
		Email:           strings.TrimSpace(payload.Email),
		Timezone:        defaultString(payload.Timezone, "WIB (+7)"),
		DivisionID:      payload.DivisionID,
		PositionID:      payload.PositionID,
		EmploymentType:  strings.TrimSpace(payload.EmploymentType),
		Status:          defaultString(payload.Status, "Aktif"),
		JoinDate:        joinDate,
		Photo:           payload.Photo,
	}, nil
}

func buildEmployeeResponse(employee *models.Employee) employeeResponse {
	var birthDate *string
	if employee.BirthDate != nil {
		formatted := employee.BirthDate.UTC().Format(time.RFC3339)
		birthDate = &formatted
	}

	var joinDate *string
	if !employee.JoinDate.IsZero() {
		formatted := employee.JoinDate.UTC().Format(time.RFC3339)
		joinDate = &formatted
	}

	divisionName := ""
	if employee.Division != nil {
		divisionName = employee.Division.Name
	}

	positionTitle := ""
	if employee.Position != nil {
		positionTitle = employee.Position.Title
	}

	return employeeResponse{
		ID:              employee.ID,
		EmployeeCode:    employee.EmployeeCode,
		NIK:             employee.NIK,
		FullName:        employee.FullName,
		BirthPlace:      employee.BirthPlace,
		BirthDate:       birthDate,
		Gender:          employee.Gender,
		BloodType:       employee.BloodType,
		MaritalStatus:   employee.MaritalStatus,
		Religion:        employee.Religion,
		IdentityType:    employee.IdentityType,
		IdentityNumber:  employee.IdentityNumber,
		AddressKtp:      employee.AddressKTP,
		AddressDomicile: employee.AddressDomicile,
		Address:         employee.Address,
		Phone:           employee.Phone,
		Email:           employee.Email,
		Timezone:        employee.Timezone,
		DivisionID:      employee.DivisionID,
		DivisionName:    divisionName,
		PositionID:      employee.PositionID,
		PositionTitle:   positionTitle,
		EmploymentType:  employee.EmploymentType,
		Status:          employee.Status,
		JoinDate:        joinDate,
		Photo:           employee.Photo,
		CreatedAt:       employee.CreatedAt,
		UpdatedAt:       employee.UpdatedAt,
	}
}

func buildDivisionResponse(division *models.EmployeeDivision) divisionResponse {
	return divisionResponse{
		ID:                division.ID,
		Name:              division.Name,
		Description:       division.Description,
		RecruitmentStatus: division.RecruitmentStatus,
		HeadEmployeeID:    division.HeadEmployeeID,
		HeadPositionID:    division.HeadPositionID,
		Head:              division.Head,
		HeadTitle:         division.HeadTitle,
		CreatedAt:         division.CreatedAt,
		UpdatedAt:         division.UpdatedAt,
	}
}

func buildPositionResponse(position *models.EmployeePosition) positionResponse {
	divisionName := ""
	if position.Division != nil {
		divisionName = position.Division.Name
	}

	return positionResponse{
		ID:           position.ID,
		Title:        position.Title,
		Code:         position.Code,
		DivisionID:   position.DivisionID,
		DivisionName: divisionName,
		ParentID:     position.ParentID,
		Notes:        position.Notes,
		Grade:        position.Grade,
		SalaryRange:  position.SalaryRange,
		CreatedAt:    position.CreatedAt,
		UpdatedAt:    position.UpdatedAt,
	}
}

func parseISOTime(value string) (*time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}
	if t, err := time.Parse(time.RFC3339, value); err == nil {
		return &t, nil
	}
	if t, err := time.Parse("2006-01-02", value); err == nil {
		return &t, nil
	}
	return nil, errors.New("invalid date format, expected RFC3339")
}

func defaultString(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return strings.TrimSpace(value)
}
