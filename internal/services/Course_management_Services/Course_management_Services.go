package CourseManagementServices

import (
	"context"
	"course_online_backend/internal/dto"
	"course_online_backend/internal/models"
	"course_online_backend/internal/services"
	"errors"
	"fmt"
	"math"
	"mime/multipart"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CourseService struct {
	db          *gorm.DB
	minioHelper *services.MinioHelper
}

func NewCourseService(db *gorm.DB) *CourseService {
	return &CourseService{
		db:          db,
		minioHelper: &services.MinioHelper{},
	}
}

func (s *CourseService) mapToCourseResponse(course *models.Course) dto.CourseResponse {
	resp := dto.CourseResponse{
		ID:                course.ID,
		Name:              course.Name,
		Description:       course.Description,
		Thumbnail:         course.Thumbnail,
		Price:             course.Price,
		IsProgressLimited: course.IsProgressLimited,
		CourseTypeID:      course.CourseTypeID,
		CreatedBy:         course.CreatedBy,
		CreatedAt:         course.CreatedAt,
		UpdatedAt:         course.UpdatedAt,
		ModulesCount:      len(course.Modules),
		EnrollmentsCount:  len(course.Enrollments),
	}

	if course.DeletedAt.Valid {
		deletedAt := course.DeletedAt.Time
		resp.DeletedAt = &deletedAt
	}

	if course.CourseType != nil {
		resp.CourseType = &dto.CourseTypeResponse{
			ID:          course.CourseType.ID,
			Name:        course.CourseType.Name,
			Description: course.CourseType.Description,
		}
	}

	if course.Creator != nil {
		resp.Creator = &dto.CreatorResponse{
			ID:       course.Creator.ID,
			FullName: course.Creator.Username,
			Email:    course.Creator.EmailAddress,
		}
	}

	lessonsCount := 0
	for _, module := range course.Modules {
		lessonsCount += len(module.Lessons)
	}
	resp.LessonsCount = lessonsCount

	return resp
}

func (s *CourseService) mapToCourseDetailResponse(course *models.Course) *dto.CourseDetailResponse {
	resp := &dto.CourseDetailResponse{
		ID:                course.ID,
		Name:              course.Name,
		Description:       course.Description,
		Thumbnail:         course.Thumbnail,
		Price:             course.Price,
		IsProgressLimited: course.IsProgressLimited,
		CourseTypeID:      course.CourseTypeID,
		CreatedBy:         course.CreatedBy,
		CreatedAt:         course.CreatedAt,
		UpdatedAt:         course.UpdatedAt,
		EnrollmentsCount:  len(course.Enrollments),
	}

	if course.DeletedAt.Valid {
		deletedAt := course.DeletedAt.Time
		resp.DeletedAt = &deletedAt
	}

	if course.CourseType != nil {
		resp.CourseType = &dto.CourseTypeResponse{
			ID:          course.CourseType.ID,
			Name:        course.CourseType.Name,
			Description: course.CourseType.Description,
		}
	}

	if course.Creator != nil {
		resp.Creator = &dto.CreatorResponse{
			ID:       course.Creator.ID,
			FullName: course.Creator.Username,
			Email:    course.Creator.EmailAddress,
		}
	}

	modules := make([]dto.ModuleResponse, len(course.Modules))
	for i, module := range course.Modules {
		lessons := make([]dto.LessonResponse, len(module.Lessons))
		for j, lesson := range module.Lessons {
			lessons[j] = dto.LessonResponse{
				ID:           lesson.ID,
				Name:         lesson.Name,
				VideoID:      lesson.VideoID,
				VideoDetails: lesson.VideoDetails,
				CourseOrder:  lesson.CourseOrder,
			}
		}

		modules[i] = dto.ModuleResponse{
			ID:           module.ID,
			Name:         module.Name,
			Number:       module.Number,
			LessonsCount: len(module.Lessons),
			Lessons:      lessons,
		}
	}
	resp.Modules = modules

	return resp
}

func (s *CourseService) CreateCourse(
	ctx context.Context,
	req *dto.CreateCourseRequest,
	courseTypeID uuid.UUID,
	firebaseUID string, 
	thumbnail *multipart.FileHeader,
) (*dto.CourseDetailResponse, error) {

	var user models.User
	if err := s.db.Where("firebase_uid = ?", firebaseUID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found in database")
		}
		return nil, err
	}

	var courseType models.CourseType
	if err := s.db.First(&courseType, courseTypeID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("course type not found")
		}
		return nil, err
	}

	course := &models.Course{
		Name:              req.Name,
		Description:       req.Description,
		Price:             req.Price,
		IsProgressLimited: req.IsProgressLimited,
		CourseTypeID:      courseTypeID,
		CreatedBy:         &user.ID,
	}

	if thumbnail != nil {
		file, err := thumbnail.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open thumbnail: %w", err)
		}
		defer file.Close()

		thumbnailURL, err := s.minioHelper.UploadCourseThumbnail(file, thumbnail)
		if err != nil {
			return nil, fmt.Errorf("failed to upload thumbnail: %w", err)
		}
		course.Thumbnail = thumbnailURL
	}

	if err := s.db.Create(course).Error; err != nil {
		if course.Thumbnail != "" {
			_ = s.minioHelper.DeleteCourseThumbnail(course.Thumbnail)
		}
		return nil, err
	}

	return s.GetByIDCourse(ctx, course.ID, false)
}


func (s *CourseService) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateCourseRequest, thumbnail *multipart.FileHeader) (*dto.CourseDetailResponse, error) {
	var course models.Course
	if err := s.db.First(&course, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("course not found")
		}
		return nil, err
	}

	if req.CourseTypeID != uuid.Nil && req.CourseTypeID != course.CourseTypeID {
		var courseType models.CourseType
		if err := s.db.First(&courseType, req.CourseTypeID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("course type not found")
			}
			return nil, err
		}
		course.CourseTypeID = req.CourseTypeID
	}

	if req.Name != "" {
		course.Name = req.Name
	}
	if req.Description != "" {
		course.Description = req.Description
	}
	if req.Price >= 0 {
		course.Price = req.Price
	}
	if req.IsProgressLimited != nil {
		course.IsProgressLimited = *req.IsProgressLimited
	}

	if thumbnail != nil {
		file, err := thumbnail.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open thumbnail: %w", err)
		}
		defer file.Close()

		if course.Thumbnail != "" {
			_ = s.minioHelper.DeleteCourseThumbnail(course.Thumbnail)
		}

		thumbnailURL, err := s.minioHelper.UploadCourseThumbnail(file, thumbnail)
		if err != nil {
			return nil, fmt.Errorf("failed to upload thumbnail: %w", err)
		}
		course.Thumbnail = thumbnailURL
	}

	if err := s.db.Save(&course).Error; err != nil {
		return nil, err
	}

	return s.GetByIDCourse(ctx, course.ID, false)
}

func (s *CourseService) SoftDeleteCourse(ctx context.Context, id uuid.UUID) error {
	var course models.Course
	if err := s.db.First(&course, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("course not found")
		}
		return err
	}

	return s.db.Delete(&course).Error
}

func (s *CourseService) RestoreCourse(ctx context.Context, id uuid.UUID) (*dto.CourseDetailResponse, error) {
	var course models.Course
	if err := s.db.Unscoped().First(&course, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("course not found")
		}
		return nil, err
	}

	if course.DeletedAt.Time.IsZero() {
		return nil, errors.New("course is not deleted")
	}

	if err := s.db.Unscoped().Model(&course).Update("deleted_at", nil).Error; err != nil {
		return nil, err
	}

	return s.GetByIDCourse(ctx, course.ID, false)
}

func (s *CourseService) PermanentDeleteCourse(ctx context.Context, id uuid.UUID) error {
	var course models.Course
	if err := s.db.Unscoped().First(&course, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("course not found")
		}
		return err
	}

	if course.Thumbnail != "" {
		if err := s.minioHelper.DeleteCourseThumbnail(course.Thumbnail); err != nil {
			fmt.Printf("Failed to delete thumbnail: %v\n", err)
		}
	}

	return s.db.Unscoped().Delete(&course).Error
}

//course browsing

func (s *CourseService) GetAllCourse(ctx context.Context, params *dto.CourseQueryParams) (*dto.PaginationResponse, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 10
	}
	if params.SortBy == "" {
		params.SortBy = "created_at"
	}
	if params.SortOrder == "" {
		params.SortOrder = "desc"
	}

	query := s.db.Model(&models.Course{})

	if params.IncludeDeleted {
		query = query.Unscoped()
	}

	if params.Search != "" {
		searchPattern := "%" + params.Search + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ?", searchPattern, searchPattern)
	}

	if params.CourseTypeID != uuid.Nil {
		query = query.Where("course_type_id = ?", params.CourseTypeID)
	}

	if params.CreatedBy != uuid.Nil {
		query = query.Where("created_by = ?", params.CreatedBy)
	}

	if params.MinPrice > 0 {
		query = query.Where("price >= ?", params.MinPrice)
	}

	if params.MaxPrice > 0 {
		query = query.Where("price <= ?", params.MaxPrice)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	orderClause := fmt.Sprintf("%s %s", params.SortBy, params.SortOrder)
	query = query.Order(orderClause)

	offset := (params.Page - 1) * params.Limit
	query = query.Offset(offset).Limit(params.Limit)

	query = query.Preload("CourseType").
		Preload("Creator").
		Preload("Modules").
		Preload("Enrollments")

	var courses []models.Course
	if err := query.Find(&courses).Error; err != nil {
		return nil, err
	}

	courseResponses := make([]dto.CourseResponse, len(courses))
	for i, course := range courses {
		courseResponses[i] = s.mapToCourseResponse(&course)
	}

	totalPages := int(math.Ceil(float64(total) / float64(params.Limit)))

	return &dto.PaginationResponse{
		Total:       total,
		Page:        params.Page,
		Limit:       params.Limit,
		TotalPages:  totalPages,
		HasNext:     params.Page < totalPages,
		HasPrevious: params.Page > 1,
		Data:        courseResponses,
	}, nil
}

func (s *CourseService) GetByIDCourse(ctx context.Context, id uuid.UUID, includeDeleted bool) (*dto.CourseDetailResponse, error) {
	query := s.db.Preload("CourseType").
		Preload("Creator").
		Preload("Modules.Lessons").
		Preload("Enrollments")

	if includeDeleted {
		query = query.Unscoped()
	}

	var course models.Course
	if err := query.First(&course, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("course not found")
		}
		return nil, err
	}

	return s.mapToCourseDetailResponse(&course), nil
}


func (s *CourseService) GetCoursesByCategory(ctx context.Context, categoryID uuid.UUID, params *dto.SimplePaginationParams) (*dto.CourseByCategoryResponse, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 10
	}

	var category models.CourseType
	if err := s.db.First(&category, categoryID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}

	query := s.db.Model(&models.Course{}).Where("course_type_id = ?", categoryID)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (params.Page - 1) * params.Limit
	query = query.Offset(offset).Limit(params.Limit)

	query = query.Preload("CourseType").
		Preload("Creator").
		Preload("Modules").
		Preload("Enrollments").
		Order("created_at DESC")

	var courses []models.Course
	if err := query.Find(&courses).Error; err != nil {
		return nil, err
	}

	courseResponses := make([]dto.CourseResponse, len(courses))
	for i, course := range courses {
		courseResponses[i] = s.mapToCourseResponse(&course)
	}

	return &dto.CourseByCategoryResponse{
		Category: dto.CourseTypeWithCountResponse{
			ID:           category.ID,
			Name:         category.Name,
			Description:  category.Description,
			CoursesCount: total,
			CreatedAt:    category.CreatedAt,
			UpdatedAt:    category.UpdatedAt,
		},
		Courses: courseResponses,
		Total:   total,
	}, nil
}

func (s *CourseService) GetPopularCourses(ctx context.Context, limit int) ([]dto.CourseResponse, error) {
	if limit < 1 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	type CourseWithEnrollmentCount struct {
		models.Course
		EnrollmentCount int64
	}

	var coursesWithCount []CourseWithEnrollmentCount
	err := s.db.Model(&models.Course{}).
		Select("courses.*, COUNT(enrollments.id) as enrollment_count").
		Joins("LEFT JOIN enrollments ON enrollments.course_id = courses.id").
		Group("courses.id").
		Order("enrollment_count DESC").
		Limit(limit).
		Preload("CourseType").
		Preload("Creator").
		Preload("Modules").
		Preload("Enrollments").
		Find(&coursesWithCount).Error

	if err != nil {
		return nil, err
	}

	courseResponses := make([]dto.CourseResponse, len(coursesWithCount))
	for i, item := range coursesWithCount {
		courseResponses[i] = s.mapToCourseResponse(&item.Course)
	}

	return courseResponses, nil
}

func (s *CourseService) GetLatestCourses(ctx context.Context, limit int) ([]dto.CourseResponse, error) {
	if limit < 1 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	var courses []models.Course
	err := s.db.Model(&models.Course{}).
		Preload("CourseType").
		Preload("Creator").
		Preload("Modules").
		Preload("Enrollments").
		Order("created_at DESC").
		Limit(limit).
		Find(&courses).Error

	if err != nil {
		return nil, err
	}

	courseResponses := make([]dto.CourseResponse, len(courses))
	for i, course := range courses {
		courseResponses[i] = s.mapToCourseResponse(&course)
	}

	return courseResponses, nil
}

func (s *CourseService) GetAllCourseTypes(ctx context.Context) ([]dto.CourseTypeWithCountResponse, error) {
	type CourseTypeWithCount struct {
		models.CourseType
		CoursesCount int64
	}

	var typesWithCount []CourseTypeWithCount
	err := s.db.Model(&models.CourseType{}).
		Select("course_types.*, COUNT(courses.id) as courses_count").
		Joins("LEFT JOIN courses ON courses.course_type_id = course_types.id AND courses.deleted_at IS NULL").
		Group("course_types.id").
		Order("course_types.name ASC").
		Find(&typesWithCount).Error

	if err != nil {
		return nil, err
	}

	responses := make([]dto.CourseTypeWithCountResponse, len(typesWithCount))
	for i, item := range typesWithCount {
		responses[i] = dto.CourseTypeWithCountResponse{
			ID:           item.ID,
			Name:         item.Name,
			Description:  item.Description,
			CoursesCount: item.CoursesCount,
			CreatedAt:    item.CreatedAt,
			UpdatedAt:    item.UpdatedAt,
		}
	}

	return responses, nil
}

func (s *CourseService) GetCoursesByInstructor(ctx context.Context, instructorID uuid.UUID, params *dto.SimplePaginationParams) (*dto.CoursesByInstructorResponse, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 10
	}

	var instructor models.User
	if err := s.db.First(&instructor, instructorID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("instructor not found")
		}
		return nil, err
	}

	query := s.db.Model(&models.Course{}).Where("created_by = ?", instructorID)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (params.Page - 1) * params.Limit
	query = query.Offset(offset).Limit(params.Limit)

	query = query.Preload("CourseType").
		Preload("Creator").
		Preload("Modules").
		Preload("Enrollments").
		Order("created_at DESC")

	var courses []models.Course
	if err := query.Find(&courses).Error; err != nil {
		return nil, err
	}

	courseResponses := make([]dto.CourseResponse, len(courses))
	for i, course := range courses {
		courseResponses[i] = s.mapToCourseResponse(&course)
	}

	return &dto.CoursesByInstructorResponse{
		Instructor: dto.InstructorResponse{
			ID:       instructor.ID,
			FullName: instructor.Username,
			Email:    instructor.EmailAddress,
		},
		Courses: courseResponses,
		Total:   total,
	}, nil
}

func (s *CourseService) GetRelatedCourses(ctx context.Context, courseID uuid.UUID, limit int) ([]dto.CourseResponse, error) {
	if limit < 1 {
		limit = 5
	}
	if limit > 20 {
		limit = 20
	}

	var currentCourse models.Course
	if err := s.db.First(&currentCourse, courseID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("course not found")
		}
		return nil, err
	}

	var courses []models.Course
	err := s.db.Model(&models.Course{}).
		Where("course_type_id = ? AND id != ?", currentCourse.CourseTypeID, courseID).
		Preload("CourseType").
		Preload("Creator").
		Preload("Modules").
		Preload("Enrollments").
		Order("created_at DESC").
		Limit(limit).
		Find(&courses).Error

	if err != nil {
		return nil, err
	}

	courseResponses := make([]dto.CourseResponse, len(courses))
	for i, course := range courses {
		courseResponses[i] = s.mapToCourseResponse(&course)
	}

	return courseResponses, nil
}



func (s *CourseService) GetCourseStats(ctx context.Context) (*dto.CourseStatsResponse, error) {
	var totalCourses int64
	if err := s.db.Model(&models.Course{}).Count(&totalCourses).Error; err != nil {
		return nil, err
	}

	var totalEnrollments int64
	if err := s.db.Model(&models.Enrollment{}).Count(&totalEnrollments).Error; err != nil {
		return nil, err
	}

	var avgPrice struct {
		Average float64
	}
	if err := s.db.Model(&models.Course{}).Select("AVG(price) as average").Scan(&avgPrice).Error; err != nil {
		return nil, err
	}

	type CategoryCount struct {
		CategoryID   uuid.UUID
		CategoryName string
		Count        int64
	}

	var categoryCounts []CategoryCount
	err := s.db.Model(&models.Course{}).
		Select("course_types.id as category_id, course_types.name as category_name, COUNT(courses.id) as count").
		Joins("JOIN course_types ON course_types.id = courses.course_type_id").
		Where("courses.deleted_at IS NULL").
		Group("course_types.id, course_types.name").
		Order("count DESC").
		Find(&categoryCounts).Error

	if err != nil {
		return nil, err
	}

	categoryResponses := make([]dto.CourseCountByCategoryResponse, len(categoryCounts))
	for i, cc := range categoryCounts {
		categoryResponses[i] = dto.CourseCountByCategoryResponse{
			CategoryID:   cc.CategoryID,
			CategoryName: cc.CategoryName,
			Count:        cc.Count,
		}
	}

	return &dto.CourseStatsResponse{
		TotalCourses:      totalCourses,
		TotalEnrollments:  totalEnrollments,
		AveragePrice:      avgPrice.Average,
		CoursesByCategory: categoryResponses,
	}, nil
}

func (s *CourseService) GetMyCourses(ctx context.Context, firebaseUID string, userRole string, params *dto.CourseQueryParams) (*dto.PaginationResponse, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 10
	}
	if params.SortBy == "" {
		params.SortBy = "created_at"
	}
	if params.SortOrder == "" {
		params.SortOrder = "desc"
	}

	var user models.User
	if err := s.db.Where("firebase_uid = ?", firebaseUID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	query := s.db.Model(&models.Course{})

	switch userRole {
	case "super_admin":
		
	case "admin":
		query = query.Where("created_by = ?", user.ID)
		
	case "instructor":
		query = query.Where("1 = 0") 
		
	case "student":
		query = query.Joins("INNER JOIN enrollments ON enrollments.course_id = courses.id").
			Where("enrollments.user_id = ?", user.ID)
		
	default:
		return nil, errors.New("invalid user role")
	}

	if params.IncludeDeleted {
		query = query.Unscoped()
	}

	if params.Search != "" {
		searchPattern := "%" + params.Search + "%"
		query = query.Where("courses.name ILIKE ? OR courses.description ILIKE ?", searchPattern, searchPattern)
	}

	if params.CourseTypeID != uuid.Nil {
		query = query.Where("courses.course_type_id = ?", params.CourseTypeID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	var orderClause string
	if userRole == "student" {
		orderClause = fmt.Sprintf("courses.%s %s", params.SortBy, params.SortOrder)
	} else {
		orderClause = fmt.Sprintf("%s %s", params.SortBy, params.SortOrder)
	}
	query = query.Order(orderClause)

	offset := (params.Page - 1) * params.Limit
	query = query.Offset(offset).Limit(params.Limit)

	query = query.Preload("CourseType").
		Preload("Creator").
		Preload("Modules", func(db *gorm.DB) *gorm.DB {
			return db.Order("number ASC")
		}).
		Preload("Enrollments")

	if userRole == "student" {
		query = query.Distinct("courses.id")
	}

	var courses []models.Course
	if err := query.Find(&courses).Error; err != nil {
		return nil, err
	}

	courseResponses := make([]dto.CourseResponse, len(courses))
	for i, course := range courses {
		courseResponses[i] = s.mapToCourseResponse(&course)
	}

	totalPages := int(math.Ceil(float64(total) / float64(params.Limit)))

	return &dto.PaginationResponse{
		Total:       total,
		Page:        params.Page,
		Limit:       params.Limit,
		TotalPages:  totalPages,
		HasNext:     params.Page < totalPages,
		HasPrevious: params.Page > 1,
		Data:        courseResponses,
	}, nil
}