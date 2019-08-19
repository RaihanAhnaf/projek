package models

import (
	"github.com/eaciit/dbox"
)

// KendoDatasourceQuery is the model of KendoJS datasource query
type KendoDatasourceQuery struct {
	Take     int                        `json:"take"`
	Skip     int                        `json:"skip"`
	Page     int                        `json:"page"`
	PageSize int                        `json:"pageSize"`
	Filter   KendoDatasourceQueryFilter `json:"filter"`
}

// KendoDatasourceQueryFilter is the model of KendoJS filter
type KendoDatasourceQueryFilter struct {
	Logic   string                                 `json:"logic"`
	Filters []KendoDatasourceQueryFilterExpression `json:"filters"`
}

// ToDboxFilter converts kendo data source query filter to dbox.Filter
func (f *KendoDatasourceQueryFilter) ToDboxFilter() *dbox.Filter {
	var filters = []*dbox.Filter{}
	for _, filter := range f.Filters {
		op := filter.Operator
		switch op {
		case "contains":
			if filter.Value != nil {
				fval, ok := filter.Value.(string)
				if ok {
					filters = append(filters, dbox.Contains(filter.Field, fval))
				}
			}
		case "equal":
			filters = append(filters, dbox.Eq(filter.Field, filter.Value))
		}
	}

	if f.Logic == "and" {
		return dbox.And(filters...)
	}
	return dbox.Or(filters...)
}

// KendoDatasourceQueryFilterExpression is the model of KendoJS filter expression
type KendoDatasourceQueryFilterExpression struct {
	Value      interface{} `json:"value"`
	Operator   string      `json:"operator"`
	Field      string      `json:"field"`
	IgnoreCase bool        `json:"ignoreCase"`
}

// KendoDatasourceResult is the model of KendoJS data source result
type KendoDatasourceResult struct {
	Data  interface{} `json:"data"`
	Count int         `json:"count"`
}
