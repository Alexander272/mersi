package utils

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/Alexander272/mersi/backend/internal/constants"
	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/gin-gonic/gin"
)

func GetFilterParams(c *gin.Context) *models.GetSiDTO {
	// u, exists := c.Get(constants.CtxUser)
	// if !exists {
	// 	return nil
	// }
	// user := u.(models.User)
	permissions, exists := c.Get(constants.IdentityCookie)
	if !exists {
		return nil
	}

	params := &models.GetSiDTO{
		SectionId: "",
		Page:      parsePage(c),
		Sort:      parseSort(c.Query("sort_by")),
		Filters:   make([]*models.Filter, 0),
	}

	// 1. Проверка прав (all)
	//TODO это не работает т.к. у user нет permissions в контексте
	// targetPermission := fmt.Sprintf("%s:%s", constants.SI, constants.Write)
	// hasWritePermission := slices.Contains(user.Permissions, targetPermission)
	// logger.Debug("User has write permission", logger.BoolAttr("hasWritePermission", hasWritePermission), logger.AnyAttr("perms", user.Permissions))
	//TODO можно добавить cookie с разрешениями и получать их оттуда
	// есть еще такой вариант: когда я сделаю доступ по подразделениям в том списке закрепить резерв и настраивать к нему доступ также как и к подразделениям
	// хотя резерв это статус, может это и не очень хорошая идея
	// all := c.Query("all")
	// hasWritePermission := all == "true"
	//* пока реализовал это через cookie
	targetPermission := fmt.Sprintf("%s:%s", constants.SI, constants.Write)
	hasWritePermission := slices.Contains(strings.Split(permissions.(string), ","), targetPermission)

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

			filterValues = append(filterValues, &models.FilterValue{
				CompareType: compareType,
				Value:       val,
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

// Выносим маппинг и реплейсер на уровень пакета, чтобы не пересоздавать их при каждом вызове
var (
	placeSuffixes = map[string]string{
		"_reserve": "reserve",
		"_moved":   "moved",
	}
	// Реплейсер подготовит строку, удалив все суффиксы и лишние запятые за один проход
	placeCleaner = strings.NewReplacer("_reserve", "", "_moved", "", ",", "")
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
	// strings.Trim удалит возможные "ошметки" запятых по краям после замены
	cleanedValue := strings.Trim(placeCleaner.Replace(*value), " ")

	*value = cleanedValue
	*field = "department"

	// var statuses []string
	// if strings.Contains(*value, "_reserve") {
	// 	statuses = append(statuses, "reserve")
	// }
	// if strings.Contains(*value, "_moved") {
	// 	statuses = append(statuses, "moved")
	// }

	// if len(statuses) > 0 {
	// 	params.Filters = append(params.Filters, &models.Filter{
	// 		Field:     "status",
	// 		FieldType: "list",
	// 		Values:    []*models.FilterValue{{CompareType: "in", Value: strings.Join(statuses, ",")}},
	// 	})
	// }

	// *value = strings.NewReplacer("_reserve", "", "_moved", "", ",", "").Replace(*value)
	// *field = "department"
}

func applyUserRestrictions(params *models.GetSiDTO) {
	params.Filters = append(params.Filters,
		&models.Filter{
			Field:  "status",
			Values: []*models.FilterValue{{CompareType: "nlike", Value: "reserve"}},
		},
		&models.Filter{
			Field:     "department",
			FieldType: "list",
			Values:    []*models.FilterValue{{CompareType: "nin", Value: "cc718041-f3da-4490-b647-380297bd3344"}},
		},
	)
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

// func GetFilterParams(c *gin.Context) *models.GetSiDTO {
// 	params := &models.GetSiDTO{
// 		SectionId: "",
// 		Page:      &models.Page{},
// 		Sort:      []*models.Sort{},
// 		Filters:   []*models.Filter{},
// 	}

// 	page := c.Query("page")
// 	size := c.Query("size")

// 	// all := c.Query("all")
// 	all := false
// 	u, exists := c.Get(constants.CtxUser)
// 	if !exists {
// 		return nil
// 	}
// 	user := u.(models.User)
// 	for _, p := range user.Permissions {
// 		if p == fmt.Sprintf("%s:%s", constants.SI, constants.Write) {
// 			all = true
// 		}
// 	}

// 	sortLine := c.Query("sort_by")
// 	filters := c.QueryMap("filters")

// 	limit, err := strconv.Atoi(size)
// 	if err != nil {
// 		params.Page.Limit = 15
// 	} else {
// 		params.Page.Limit = limit
// 	}

// 	p, err := strconv.Atoi(page)
// 	if err != nil {
// 		params.Page.Offset = 0
// 	} else {
// 		params.Page.Offset = (p - 1) * params.Page.Limit
// 	}

// 	if sortLine != "" {
// 		sort := strings.Split(sortLine, ",")
// 		for _, v := range sort {
// 			field, found := strings.CutPrefix(v, "-")
// 			t := "ASC"
// 			if found {
// 				t = "DESC"
// 			}

// 			params.Sort = append(params.Sort, &models.Sort{
// 				Field: field,
// 				Type:  t,
// 			})
// 		}
// 	}

// 	// можно сделать массив с именами полей, а потом передавать для каждого поля значение фильтра, например
// 	// filter[0]=nextVerificationDate&nextVerificationDate[lte]=somevalue&nextVerificationDate[qte]=somevalue&filter[1]=name&name[eq]=somevalue
// 	// qte - >=; lte - <=
// 	// нужен еще тип как-то передать
// 	// как вариант можно передать filter[nextVerificationDate]=date, filter[name]=string
// 	// только надо проверить как это все будет читаться на сервере и записываться на клиенте

// 	// можно сделать следующие варианты compareType (это избавит от необходимости знать тип поля)
// 	// number or date: eq, qte, lte
// 	// string: like, con, start, end
// 	// list: in

// 	for k, v := range filters {
// 		valueMap := c.QueryMap(k)
// 		values := []*models.FilterValue{}
// 		for key, value := range valueMap {
// 			if k == "place" {
// 				statusFilter := &models.Filter{Field: "status", FieldType: "list", Values: []*models.FilterValue{{CompareType: "in"}}}
// 				tmp := []string{}
// 				if strings.Contains(value, "_reserve") {
// 					tmp = append(tmp, "reserve")
// 				}
// 				if strings.Contains(value, "_moved") {
// 					tmp = append(tmp, "moved")
// 				}
// 				statusFilter.Values[0].Value = strings.Join(tmp, ",")
// 				params.Filters = append(params.Filters, statusFilter)

// 				value = strings.Replace(value, "_reserve", "", -1)
// 				value = strings.Replace(value, "_moved", "", -1)
// 				value = strings.Trim(value, ",")
// 				k = "department"

// 				if value != "" {
// 					// params.Filters = append(params.Filters, &models.SIFilter{Field: "last_place", Values: []*models.SIFilterValue{}})
// 					values = append(values, &models.FilterValue{CompareType: key, Value: value})
// 				}
// 			}

// 			values = append(values, &models.FilterValue{
// 				CompareType: key,
// 				Value:       value,
// 			})
// 		}
// 		if values[0].Value == "" {
// 			continue
// 		}

// 		f := &models.Filter{
// 			Field:     k,
// 			FieldType: v,
// 			Values:    values,
// 		}

// 		params.Filters = append(params.Filters, f)
// 	}

// 	// if all != "true" {
// 	if !all {
// 		params.Filters = append(params.Filters, &models.Filter{
// 			Field:     "status",
// 			FieldType: "",
// 			Values:    []*models.FilterValue{{CompareType: "nlike", Value: "reserve"}},
// 		}, &models.Filter{ //TODO задавать id подразделения не очень хорошая идея
// 			Field:     "department",
// 			FieldType: "list",
// 			//TODO если сделаю фильтрацию по подразделениям, то надо будет убрать этот фильтр
// 			Values: []*models.FilterValue{{CompareType: "nin", Value: "cc718041-f3da-4490-b647-380297bd3344"}},
// 		})

// 	}

// 	search := c.QueryMap("search")
// 	for key, value := range search {
// 		params.Search = &models.Search{
// 			Value:  value,
// 			Fields: strings.Split(key, ","),
// 		}
// 	}

// 	status := c.Query("status")
// 	if status == "" {
// 		params.Status = models.InstrumentStatusWork
// 	} else {
// 		params.Status = models.InstrumentStatus(status)
// 	}

// 	return params
// }
