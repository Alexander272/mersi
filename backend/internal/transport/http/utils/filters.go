package utils

import (
	"strconv"
	"strings"
	"time"

	"github.com/Alexander272/mersi/backend/internal/constants"
	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/gin-gonic/gin"
)

func GetFilterParams(c *gin.Context) *models.GetSiDTO {
	u, exists := c.Get(constants.CtxUser)
	if !exists {
		return nil
	}
	user, ok := u.(models.User)
	if !ok {
		return nil
	}

	params := &models.GetSiDTO{
		SectionId: "",
		Page:      parsePage(c),
		Sort:      parseSort(c.Query("sort_by")),
		Filters:   make([]*models.Filter, 0),
		UserID:    user.ID,
	}

	hasWritePermission := false
	if wp, exists := c.Get(constants.CtxHasWritePermission); exists {
		if wp, ok := wp.(bool); ok {
			hasWritePermission = wp
		}
	}

	// 2. Обработка фильтров
	filtersMap := c.QueryMap("filters")
	for field, fieldType := range filtersMap {
		rawValues := c.QueryMap(field)
		if len(rawValues) == 0 {
			continue
		}

		filterValues := make([]*models.FilterValue, 0, len(rawValues))

		for compareType, val := range rawValues {
			if val == "" {
				continue
			}

			// Специфичная логика для "place"
			if field == "place" {
				handlePlaceFilter(params, &field, &val)
			}

			if val == "" {
				continue
			}

			filterValues = append(filterValues, &models.FilterValue{
				CompareType: compareType,
				Value:       NormalizeFilterValue(field, val),
			})
		}

		if len(filterValues) > 0 {
			params.Filters = append(params.Filters, &models.Filter{
				Field:     field,
				FieldType: fieldType,
				Values:    filterValues,
			})
		}
	}

	// 3. Ограничения для обычных пользователей
	if !hasWritePermission {
		if deptIDs, exists := c.Get(constants.CtxDepartmentAccess); exists {
			if deptIDs, ok := deptIDs.([]string); ok {
				params.DepartmentAccess = deptIDs
			}
		}
		applyUserRestrictions(params)
	}
	// 4. Поиск, статус и прочее
	applySearchAndStatus(c, params)

	return params
}

// --- Вспомогательные функции для чистоты кода ---

func parsePage(c *gin.Context) *models.Page {
	limit, _ := strconv.Atoi(c.DefaultQuery("size", "15"))
	p, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if limit <= 0 {
		limit = 15
	}
	if p <= 0 {
		p = 1
	}

	return &models.Page{
		Limit:  limit,
		Offset: (p - 1) * limit,
	}
}

func parseSort(sortLine string) []*models.Sort {
	if sortLine == "" {
		return []*models.Sort{}
	}

	rawSorts := strings.Split(sortLine, ",")
	sorts := make([]*models.Sort, 0, len(rawSorts))
	for _, v := range rawSorts {
		field, found := strings.CutPrefix(v, "-")
		order := "ASC"
		if found {
			order = "DESC"
		}
		sorts = append(sorts, &models.Sort{Field: field, Type: order})
	}
	return sorts
}

// dateFilterFields — поля, значения которых сравниваются со столбцами типа date/time в БД.
// Для них нужно приводить legacy-значения (unix-время в секундах) к строке, понятной PostgreSQL.
var dateFilterFields = map[string]struct{}{
	"nextVerificationDate": {},
	"verificationDate":     {},
	"dateOfReceipt":        {},
	"repairInfo":           {},
}

// NormalizeFilterValue приводит значение date-фильтра к строке, которую PostgreSQL сможет
// привести к типу столбца. Старые фильтры, сохранённые прежней версией фронтенда, содержат
// unix-время в секундах (например 1774983600), а колонки БД уже являются datetime.
// ISO-строки ("2026-04-01T00:00:00.000Z") и NullDate оставляем без изменений.
func NormalizeFilterValue(field, value string) string {
	if _, ok := dateFilterFields[field]; !ok {
		return value
	}
	sec, err := strconv.ParseInt(value, 10, 64)
	if err != nil || sec <= 0 {
		return value
	}
	return time.Unix(sec, 0).In(time.Local).Format("2006-01-02 15:04:05")
}

// Выносим маппинг и реплейсер на уровень пакета, чтобы не пересоздавать их при каждом вызове
var (
	placeSuffixes = map[string]string{
		"_reserve": "reserve",
		"_moved":   "moved",
	}
	// Реплейсер подготовит строку, удалив все суффиксы
	placeCleaner = strings.NewReplacer("_reserve", "", "_moved", "")
)

func handlePlaceFilter(params *models.GetSiDTO, field *string, value *string) {
	if *value == "" {
		return
	}

	// 1. Собираем статусы. Используем make с небольшой емкостью, чтобы избежать расширения слайса
	statuses := make([]string, 0, len(placeSuffixes))
	for suffix, status := range placeSuffixes {
		if strings.Contains(*value, suffix) {
			statuses = append(statuses, status)
		}
	}

	// 2. Если найдены статусы, добавляем фильтр
	if len(statuses) > 0 {
		params.Filters = append(params.Filters, &models.Filter{
			Field:     "status",
			FieldType: "list",
			Values:    []*models.FilterValue{{CompareType: "in", Value: strings.Join(statuses, ",")}},
		})
	}

	// 3. Очищаем значение и меняем имя поля
	cleanedValue := strings.Trim(placeCleaner.Replace(*value), " ,")

	// Если после удаления суффиксов ничего не осталось (только "Резерв"/"Перемещение"),
	// то department filter не добавляем, он не нужен
	if cleanedValue == "" {
		*value = ""
		return
	}

	*value = cleanedValue
	*field = "department"
}

func applyUserRestrictions(params *models.GetSiDTO) {
	params.Filters = append(params.Filters, &models.Filter{
		Field: "status", Values: []*models.FilterValue{{CompareType: "nlike", Value: "reserve"}},
	})
	if len(params.DepartmentAccess) > 0 {
		filterVal := strings.Join(params.DepartmentAccess, ",")
		params.Filters = append(params.Filters, &models.Filter{
			Field: "department", FieldType: "list",
			Values: []*models.FilterValue{{
				CompareType: "in", Value: filterVal,
			}},
		})
	}
}

func applySearchAndStatus(c *gin.Context, params *models.GetSiDTO) {
	for fields, val := range c.QueryMap("search") {
		params.Search = &models.Search{
			Value:  val,
			Fields: strings.Split(fields, ","),
		}
		break
	}

	status := c.DefaultQuery("status", string(models.InstrumentStatusWork))
	params.Status = models.InstrumentStatus(status)
}
