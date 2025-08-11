package utils

import (
	"strconv"
	"strings"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/gin-gonic/gin"
)

func GetFilterParams(c *gin.Context) *models.GetSiDTO {
	params := &models.GetSiDTO{
		SectionId: "",
		Page:      &models.Page{},
		Sort:      []*models.Sort{},
		Filters:   []*models.Filter{},
	}

	page := c.Query("page")
	size := c.Query("size")

	all := c.Query("all")

	sortLine := c.Query("sort_by")
	filters := c.QueryMap("filters")

	limit, err := strconv.Atoi(size)
	if err != nil {
		params.Page.Limit = 15
	} else {
		params.Page.Limit = limit
	}

	p, err := strconv.Atoi(page)
	if err != nil {
		params.Page.Offset = 0
	} else {
		params.Page.Offset = (p - 1) * params.Page.Limit
	}

	if sortLine != "" {
		sort := strings.Split(sortLine, ",")
		for _, v := range sort {
			field, found := strings.CutPrefix(v, "-")
			t := "ASC"
			if found {
				t = "DESC"
			}

			params.Sort = append(params.Sort, &models.Sort{
				Field: field,
				Type:  t,
			})
		}
	}

	// можно сделать массив с именами полей, а потом передавать для каждого поля значение фильтра, например
	// filter[0]=nextVerificationDate&nextVerificationDate[lte]=somevalue&nextVerificationDate[qte]=somevalue&filter[1]=name&name[eq]=somevalue
	// qte - >=; lte - <=
	// нужен еще тип как-то передать
	// как вариант можно передать filter[nextVerificationDate]=date, filter[name]=string
	// только надо проверить как это все будет читаться на сервере и записываться на клиенте

	// можно сделать следующие варианты compareType (это избавит от необходимости знать тип поля)
	// number or date: eq, qte, lte
	// string: like, con, start, end
	// list: in

	for k, v := range filters {
		valueMap := c.QueryMap(k)
		values := []*models.FilterValue{}
		for key, value := range valueMap {
			if k == "place" {
				statusFilter := &models.Filter{Field: "status", FieldType: "list", Values: []*models.FilterValue{{CompareType: "in"}}}
				tmp := []string{}
				if strings.Contains(value, "_reserve") {
					tmp = append(tmp, "reserve")
				}
				if strings.Contains(value, "_moved") {
					tmp = append(tmp, "moved")
				}
				statusFilter.Values[0].Value = strings.Join(tmp, ",")
				params.Filters = append(params.Filters, statusFilter)

				value = strings.Replace(value, "_reserve", "", -1)
				value = strings.Replace(value, "_moved", "", -1)
				value = strings.Trim(value, ",")
				k = "department"

				if value != "" {
					// params.Filters = append(params.Filters, &models.SIFilter{Field: "last_place", Values: []*models.SIFilterValue{}})
					values = append(values, &models.FilterValue{CompareType: key, Value: value})
				}
			}

			values = append(values, &models.FilterValue{
				CompareType: key,
				Value:       value,
			})
		}
		if values[0].Value == "" {
			continue
		}

		f := &models.Filter{
			Field:     k,
			FieldType: v,
			Values:    values,
		}

		params.Filters = append(params.Filters, f)
	}

	if all != "true" {
		params.Filters = append(params.Filters, &models.Filter{
			Field:     "status",
			FieldType: "",
			Values:    []*models.FilterValue{{CompareType: "nlike", Value: "reserve"}},
		}, &models.Filter{ //TODO задавать id подразделения не очень хорошая идея
			Field:     "department",
			FieldType: "list",
			Values:    []*models.FilterValue{{CompareType: "nin", Value: "cc718041-f3da-4490-b647-380297bd3344"}},
		})

	}

	search := c.QueryMap("search")
	for key, value := range search {
		params.Search = &models.Search{
			Value:  value,
			Fields: strings.Split(key, ","),
		}
	}

	status := c.Query("status")
	if status == "" {
		params.Status = models.InstrumentStatusWork
	} else {
		params.Status = models.InstrumentStatus(status)
	}

	return params
}
