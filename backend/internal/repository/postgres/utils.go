package postgres

import (
	"fmt"
)

func getFilterLine(compare string, fieldName string, count int) string {
	switch compare {
	case "con":
		return fmt.Sprintf("LOWER(%s) LIKE LOWER('%%'||$%d||'%%')", fieldName, count)
	case "start":
		return fmt.Sprintf("LOWER(%s) LIKE LOWER($%d||'%%')", fieldName, count)
	case "end":
		return fmt.Sprintf("LOWER(%s) LIKE LOWER('%%'||$%d)", fieldName, count)
	case "like":
		return fmt.Sprintf("LOWER(%s) = LOWER($%d)", fieldName, count)
	case "nlike":
		return fmt.Sprintf("LOWER(%s) != LOWER($%d)", fieldName, count)

	case "in":
		// LOWER(place) ~* 'test|Отдел технического сервиса'
		// LOWER(place) ILIKE ANY (ARRAY['test %','Отдел технического сервиса %'])
		return fmt.Sprintf("LOWER(%s::text) ~* $%d", fieldName, count)
	case "nin":
		return fmt.Sprintf("LOWER(%s::text) !~* $%d", fieldName, count)

	case "eq":
		return fmt.Sprintf("%s = $%d", fieldName, count)
	case "neq":
		return fmt.Sprintf("%s != $%d", fieldName, count)
	case "gte":
		return fmt.Sprintf("%s >= $%d", fieldName, count)
	case "lte":
		return fmt.Sprintf("%s <= $%d", fieldName, count)

	case "null":
		return fmt.Sprintf("%s IS NULL", fieldName)
	}

	return ""
}
