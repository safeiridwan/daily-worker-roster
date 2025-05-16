package utils

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"slices"

	"gorm.io/gorm"
)

var operator = []string{"or", "eq", "ne", "lt", "lte", "gt", "gte", "isnull", "like", "dlt", "dlte", "dgt", "dgte"}

// getOperator returns the SQL operator string corresponding to the given operator code.
//
// Parameters:
//   - op: string - The operator code.
//
// Returns:
//   - string - The corresponding SQL operator.
//
// Example:
//
//	getOperator("eq") returns "="
//	getOperator("lt") returns "<"
//	getOperator("gte") returns ">="
func getOperator(op string) string {
	switch op {
	case "eq":
		return "="
	case "ne":
		return "!="
	case "lt":
		return "<"
	case "lte":
		return "<="
	case "gt":
		return ">"
	case "gte":
		return ">="
	default:
		return "="
	}
}

// ParseQueryParams parses the query parameters from an HTTP request and returns a map of filters.
// It handles various filter types including 'or' conditions, custom operators, and equality checks.
//
// Parameters:
//   - r: *http.Request - The HTTP request containing query parameters.
//
// Returns:
//   - map[string]any - A map of parsed filters.
//
// Example:
//
//	Given the following URL query parameters:
//	  ?name=John&age=gte:25&city=New York,London&offset=10&limit=20&order_by=-created_at&order_by=+updated_at
//
//	The function will return:
//	  {
//	    "name": {"eq": ["John"]},
//	    "age": {"gte": "25"},
//	    "city": {"or": ["New York", "London"]},
//	    "offset": ["10"],
//	    "limit": ["20"],
//	    "order_by": ["-created_at", "+updated_at"]
//	  }
func ParseQueryParams(r *http.Request) map[string]any {
	filters := make(map[string]any)
	for key, values := range r.URL.Query() {
		if len(values) == 0 {
			continue
		}

		if _, ok := filters[key]; !ok {
			filters[key] = make(map[string]any)
		}

		for _, value := range values {
			if value != "" {
				var next bool
				if strings.Contains(value, ":") {
					parts := strings.SplitN(value, ":", 2)
					op, val := parts[0], parts[1]

					if slices.Contains(operator, op) {
						if op == "or" {
							filters[key].(map[string]any)[op] = strings.Split(val, ",")
						} else {
							filters[key].(map[string]any)[op] = val
						}
						next = true
					}
				}

				if !next {
					if strings.Contains(value, ",") {
						filters[key].(map[string]any)["or"] = strings.Split(value, ",")
					} else {
						filters[key].(map[string]any)["eq"] = value
					}
				}
			}
		}
	}

	for _, param := range []string{"offset", "limit", "order_by"} {
		if values := r.URL.Query()[param]; len(values) > 0 {
			filters[param] = values
			if param == "order_by" {
				filters[param] = values
			} else {
				intVal, _ := strconv.Atoi(values[0])
				filters[param] = intVal
			}
		}
	}

	return filters
}

// BuildWhere constructs a GORM query based on the provided filters.
// It applies various conditions, including 'or' clauses, equality checks, and custom operators.
// It also handles offset, limit, and order_by parameters.
//
// Parameters:
//   - query: *gorm.DB - The initial GORM query to be modified.
//   - filters: map[string]any - A map of filters to apply to the query.
//
// Returns:
//   - *gorm.DB - The modified GORM query with applied filters.
//
// Example:
//
//	Given the following filters:
//	  {
//	    "name": {"eq": ["John"]},
//	    "age": {"gte": "25"},
//	    "city": {"or": ["New York", "London"]},
//	    "offset": ["10"],
//	    "limit": ["20"],
//	    "order_by": ["-created_at", "+updated_at"]
//	  }
//
//	The function will modify the GORM query as follows:
//	  query = query.Where("name = ?", "John")
//	  query = query.Where("age >= ?", "25")
//	  query = query.Where("city = ? OR city = ?", "New York", "London")
//	  query = query.Offset(10)
//	  query = query.Limit(20)
//	  query = query.Order("created_at DESC")
//	  query = query.Order("updated_at ASC")
func BuildWhere(query *gorm.DB, filters map[string]any) *gorm.DB {
	whereFilters := make(map[string]any)
	for key, value := range filters {
		switch key {
		case "offset", "limit", "order_by":
			continue
		default:
			whereFilters[key] = value
		}
	}

	for key, value := range whereFilters {
		// if key == "created_at" {
		// 	key = key + "::date"
		// }
		// if len(table) > 0 {
		// 	key = table[0] + "." + key
		// }
		switch v := value.(type) {
		case map[string]any:
			for op, val := range v {
				switch op {
				case "or":
					var conditions []string
					var args []any
					for _, v := range val.([]string) {
						conditions = append(conditions, fmt.Sprintf("%s = ?", key))
						args = append(args, v)
					}
					query = query.Where(strings.Join(conditions, " OR "), args...)
				case "isnull":
					query = query.Where(key + " IS NULL")
				case "like":
					v := val.(string)
					query = query.Where(fmt.Sprintf("UPPER(%s) LIKE UPPER(?)", key), "%"+v+"%")
				case "dlt":
					query = query.Where(key+" < TO_TIMESTAMP(?, 'YYYY-MM-DD HH24:MI:SS')", val)
				case "dlte":
					query = query.Where(key+" <= TO_TIMESTAMP(?, 'YYYY-MM-DD HH24:MI:SS')", val)
				case "dgt":
					query = query.Where(key+" > TO_TIMESTAMP(?, 'YYYY-MM-DD HH24:MI:SS')", val)
				case "dgte":
					query = query.Where(key+" >= TO_TIMESTAMP(?, 'YYYY-MM-DD HH24:MI:SS')", val)
				default:
					query = query.Where(key+" "+getOperator(op)+" ?", val)
				}
			}
		default:
			query = query.Where(key+" = ?", value)
		}
	}

	if offsetValue, ok := filters["offset"]; ok {
		if offsetSlice, ok := offsetValue.([]string); ok && len(offsetSlice) > 0 {
			if offset, err := strconv.Atoi(offsetSlice[0]); err == nil {
				query = query.Offset(offset)
			}
		}
	}

	if limitValue, ok := filters["limit"]; ok {
		if limitSlice, ok := limitValue.([]string); ok && len(limitSlice) > 0 {
			if limit, err := strconv.Atoi(limitSlice[0]); err == nil {
				query = query.Limit(limit)
			}
		}
	}

	if orderByValue, ok := filters["order_by"]; ok {
		if orderBySlice, ok := orderByValue.([]string); ok {
			for _, orderByStr := range orderBySlice {
				direction := "ASC"
				switch {
				case strings.HasPrefix(orderByStr, "-"):
					direction = "DESC"
					orderByStr = strings.TrimPrefix(orderByStr, "-")
				default:
					orderByStr = strings.TrimPrefix(orderByStr, "+")
				}
				query = query.Order(orderByStr + " " + direction)
			}
		}
	}
	return query
}

// Omit Non Valid fields in order to prevent GORM inserting default value like
// string = ""
// int64 = 0
// bool = false
// time = 0001-01-01 00:00:00 +0000 UTC
func FilterValidFields(entityStrut any) []string {
	result := []string{}
	reflectValue := reflect.ValueOf(entityStrut)
	reflectType := reflect.TypeOf(entityStrut)

	for i := range reflectValue.NumField() {
		field := reflectValue.Field(i)
		fieldType := reflectType.Field(i)

		kind := reflect.ValueOf(field).Kind()
		if kind == reflect.String || kind == reflect.Int || kind == reflect.Float64 || kind == reflect.Bool {
			// Handle primitive types directly
			result = append(result, reflectValue.Type().Field(i).Name)
			continue
		} else if field.Kind() == reflect.Struct && field.Type().PkgPath() == "database/sql" {
			valid := field.FieldByName("Valid").Bool()
			if !valid {
				var columnName string
				gormTag := fieldType.Tag.Get("gorm")
				tagParts := strings.Split(gormTag, ";")
				for _, tag := range tagParts {
					if len(tag) > 7 && tag[:7] == "column:" {
						columnName = tag[7:]
						break
					}
				}

				result = append(result, columnName)
			}
		}
	}
	return result
}

func FilterFields(input any) []string {
	var validFields []string
	val := reflect.ValueOf(input)
	typeOfStruct := val.Type()

	for i := range val.NumField() {
		field := val.Field(i)
		fieldType := typeOfStruct.Field(i)

		// Check if the field is a gorm-declared column
		gormTag := fieldType.Tag.Get("gorm")
		if gormTag == "-" {
			continue // Skip fields explicitly ignored by GORM
		}

		// Check if field is a sql.Null* type and has a valid value
		validField := false
		if field.Kind() == reflect.Struct {
			validField = field.FieldByName("Valid").Bool()
		} else {
			validField = !field.IsZero()
		}

		if validField {
			validFields = append(validFields, fieldType.Name)
		}
	}

	return validFields
}
