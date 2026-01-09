package utils // Or your util package

import (
	"fmt"
	"reflect"
)

// buildSets takes a pointer-based patch struct and dynamically generates
// the SET clauses and values for an UPDATE query.
// [] string : fields
// []any : actual values of the fields
func BuildSets(patch any) ([]string, []any, int) {
	v := reflect.ValueOf(patch)
	// Dereference the pointer to get the struct value
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, nil, 0
	}

	structType := v.Type()
	var sets []string
	var values []any
	index := 1 // Start counter for $1, $2, etc.

	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)
		fieldType := structType.Field(i)
		columnName := fieldType.Tag.Get("db")

		// Skip if no valid db tag or if the field is not a pointer
		if columnName == "" || columnName == "-" || fieldValue.Kind() != reflect.Pointer {
			continue
		}

		// CRITICAL CHECK: If the pointer is NOT nil, it means we have an update.
		if !fieldValue.IsNil() {
			// Append the SET clause (e.g., "title = $1")
			sets = append(sets, fmt.Sprintf("%s = $%d", columnName, index))
			// Append the dereferenced value (the actual string, bool, etc.)
			values = append(values, fieldValue.Elem().Interface())
			// Increment the index for the next parameter
			index++
		}
	}
	return sets, values, index
}

// Return Type,Example Value,Role in Final Query
// []string (sets),"[""title = $1"", ""content = $2""]",These are the dynamic assignments.
// []any (values),"[""New Title"", ""Updated Content""]",These are the actual values passed securely to the database.
// int (nextIndex),3,This is the next available placeholder number. (Count of dynamic fields: 2 + 1 = 3)

// StructToMap converts a struct to a map[string]interface{} based on a tag name.
// It skips fields that are nil pointers.
func StructToMap(data interface{}, tagName string) map[string]interface{} {
	result := make(map[string]interface{})
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return result
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get(tagName)
		if tag == "" || tag == "-" {
			continue
		}

		val := v.Field(i)
		// Handle pointers: if nil, skip. If not nil, dereference.
		if val.Kind() == reflect.Ptr {
			if val.IsNil() {
				continue
			}
			result[tag] = val.Elem().Interface()
		} else {
			result[tag] = val.Interface()
		}
	}
	return result
}
