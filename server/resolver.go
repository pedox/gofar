package main

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/graphql-go/graphql"
)

func resolve(fieldName string, param graphql.ResolveParams, fields map[string]interface{}, parentFieldName *string, parentFields *map[string]interface{}) (map[string]interface{}, error) {
	fieldSet := map[string]bool{}
	fieldKey := []string{}

	id, _ := param.Args["ID"].(string)

	if _, ok := fields["ID"]; !ok {
		fields["ID"] = "string"
	}

	//Pre defined foreign Key
	for name, typeData := range fields {
		if reflect.TypeOf(typeData).String()[0:4] == "map[" {
			foreignRow := name + "_ID"
			fields[foreignRow] = "string"
			if _, ok := fieldSet[foreignRow]; !ok {
				fieldSet[foreignRow] = true
				fieldKey = append(fieldKey, foreignRow)
			}
		} else {
			if _, ok := fieldSet[name]; !ok {
				fieldSet[name] = true
				fieldKey = append(fieldKey, name)
			}
		}
	}

	//-------- EXECUTE DATABASE HERE -------------//
	fmt.Printf("SELECT %s FROM %s WHERE id = \"%s\"\n\n", strings.Join(fieldKey, ", "), fieldName, id)
	//------- END ---------------//

	for name, typeData := range fields {
		if reflect.TypeOf(typeData).String()[0:4] == "map[" {
			fieldValue, _ := resolve(name, param, typeData.(map[string]interface{}), &fieldName, &fields)
			fields[name] = fieldValue
		} else {
			// FILTER OPT
			_, opt := getFieldTypeData(fields[name].(string))
			if _, ok := opt["hide"]; ok {
				fields[name] = nil
			}
		}
	}

	return fields, nil
}
