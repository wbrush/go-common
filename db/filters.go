package db

import (
	"fmt"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/types"
	"github.com/go-pg/pg/v9/orm"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

const MaxLimit = 90000

const (
	TagNameJSON = "json"
	TagNameDB   = "pg"
)

const (
	DefaultFilterFirst   = "first"
	DefaultFilterAfter   = "after"
	DefaultFilterOrderBy = "orderBy"
)

/* --- moved from go-pg url.Values to use with relation filtering --- */

type DBFilter struct {
	values url.Values
}

func NewDBFilter(values url.Values) *DBFilter {
	return &DBFilter{
		values: values,
	}
}

func bytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func forEachValue(
	b []byte, field string, values []string, opValue string,
) []byte {
	for _, value := range values {
		b = types.AppendField(b, field, 1)
		b = append(b, opValue...)
		b = types.AppendString(b, value, 1)
	}
	return b
}

func forAllValues(
	b []byte, field string, values []string, singleOpValue, multiOpValue string,
) []byte {
	if len(values) <= 1 {
		return forEachValue(b, field, values, singleOpValue)
	}

	b = types.AppendField(b, field, 1)
	b = append(b, multiOpValue...)
	b = append(b, '(')
	b = pg.Strings(values).AppendValue(b, 1)
	b = append(b, ')')
	return b
}

func addOperator(b []byte, field, op string, values []string) []byte {
	switch op {
	case "", "include":
		b = forAllValues(b, field, values, " = ", " IN ")
	case "exclude", "neq":
		b = forAllValues(b, field, values, " != ", " NOT IN ")
	case "gt":
		b = forEachValue(b, field, values, " > ")
	case "gte":
		b = forEachValue(b, field, values, " >= ")
	case "lt":
		b = forEachValue(b, field, values, " < ")
	case "lte":
		b = forEachValue(b, field, values, " <= ")
	case "ieq":
		b = forEachValue(b, field, values, " ILIKE ")
	case "match":
		b = forEachValue(b, field, values, " SIMILAR TO ")
	}
	return b
}

func (f *DBFilter) Filters(q *orm.Query) (*orm.Query, error) {
	if f == nil {
		return q, nil
	}

	var b []byte
	for filter, values := range f.values {
		if strings.HasSuffix(filter, "[]") {
			filter = filter[:len(filter)-2]
		}

		var operation string
		if ind := strings.Index(filter, "__"); ind != -1 {
			filter, operation = filter[:ind], filter[ind+2:]
		}

		if !IsDefaultFilter(filter) {
			if b != nil {
				b = append(b, " AND "...)
			}
			b = addOperator(b, filter, operation, values)
		}
	}

	if len(b) > 0 {
		q = q.Where(bytesToString(b))
	}
	return q, nil
}

/* --- --- */

func IsDefaultFilter(name string) bool {
	switch name {
	case DefaultFilterFirst,
		DefaultFilterAfter,
		DefaultFilterOrderBy:
		return true
	}

	return false
}

func ApplyDefaultFilters(q *orm.Query, f url.Values, applyLimits ...bool) (*orm.Query, error) {
	q.Limit(MaxLimit) //set default max limit

	if f == nil {
		return q, nil //nothing to do
	}

	var first int64 = 0
	var err error

	if (f.Get(DefaultFilterFirst) != "") && ((len(applyLimits) == 0) || ((len(applyLimits) > 0) && (applyLimits[0]))) {
		first, err = strconv.ParseInt(f.Get(DefaultFilterFirst), 10, 64)
		if err != nil {
			return q, fmt.Errorf("can't parse first filter (%s): %s", f.Get(DefaultFilterFirst), err.Error())
		}

		if first > 0 && first < MaxLimit {
			q.Limit(int(first))
		}
	}

	if (f.Get(DefaultFilterAfter) != "") && ((len(applyLimits) == 0) || ((len(applyLimits) > 0) && (applyLimits[0]))) {
		after, err := DecodeCursorToId(f.Get(DefaultFilterAfter))
		if err != nil {
			return q, err
		}
		//q.Where("?TableAlias.id > ?", after)
		q.Offset(int(after))
	}

	if f.Get(DefaultFilterOrderBy) != "" {
		list := strings.Split(f.Get(DefaultFilterOrderBy), ",")

		for _, o := range list {
			ind := strings.LastIndex(o, "_")
			if ind < 0 {
				return q, fmt.Errorf("bad orderBy format. Use field_DESC or field_ASC")
			}

			r := []rune(o)
			r[ind] = ' '
			order := string(r)

			q.Order(order)
		}
	} else {
		q.Order("id ASC")
	}

	return q, nil
}

//this type is used to collect errors on filters
type PrepareFiltersByModelResult struct {
	Prepared url.Values

	//map is used to set field as a key for message
	Errors map[string]error
}

func PrepareFiltersByModel(filters url.Values, model interface{}, ignoreUnknowns ...bool) (res PrepareFiltersByModelResult, err error) {
	if model == nil {
		return res, fmt.Errorf("provided model is nil")
	}

	var ignoreUnknown = false
	if len(ignoreUnknowns) > 0 {
		ignoreUnknown = ignoreUnknowns[0]
	}

	res.Prepared = make(url.Values)
	res.Errors = make(map[string]error)

	modelValues := reflect.ValueOf(model)
	modelType := modelValues.Type()
	if modelValues.Kind() == reflect.Ptr {
		modelValues = modelValues.Elem()
	}

	//prepare and fill tags map
	tags := make(map[string]string)
	var jsonTagName string
	for i := 0; i < modelValues.NumField(); i++ {
		jsonTagName = ""
		field := modelType.Field(i)
		if tagv := field.Tag.Get(TagNameJSON); tagv != "" {
			tag := strings.Split(tagv, ",")
			if len(tag) > 0 && tag[0] != "" && tag[0] != "-" {
				jsonTagName = tag[0]
			}
		}

		//if json tag is not filled,there is no way to fill values
		if jsonTagName == "" {
			continue
		}

		if tagv := field.Tag.Get(TagNameDB); tagv != "" {
			tag := strings.Split(tagv, ",")
			if len(tag) > 0 && tag[0] != "" && tag[0] != "-" {
				tags[jsonTagName] = tag[0]
			}
		}
	}

	var (
		dbName         string
		preparedFilter string
		ok             bool
	)
	for filterWithAct := range filters {
		f := strings.SplitN(filterWithAct, "__", 2)
		if len(f) > 0 {
			switch f[0] {
			case DefaultFilterFirst:
				res.Prepared[filterWithAct] = filters[filterWithAct]
				res.Prepared["limit"] = filters[filterWithAct] //add first aw limit
			case DefaultFilterAfter:
				res.Prepared[filterWithAct] = filters[filterWithAct]
			case DefaultFilterOrderBy:
				if len(filters[filterWithAct]) > 0 {
					res.Prepared["orderBy"] = []string{}
					res.Prepared["orderByUnknown"] = []string{}
					for i := range filters[filterWithAct] {
						fields := strings.SplitN(filters[filterWithAct][i], ",", 2) //support for comma-separated filters
						for f := range fields {
							o := strings.SplitN(strings.TrimSpace(fields[f]), "_", 2) //get the field name from a value

							dbName, ok = tags[o[0]]
							if !ok {
								if !ignoreUnknown {
									res.Errors[filterWithAct] = fmt.Errorf("unknown %s field for orderBy is provided. Please, check field name and %s datamodel's tags", filterWithAct, reflect.TypeOf(model))
								} else {
									//leave as is
									preparedFilter = o[0]
									if len(o) > 1 { //add direction if provided
										preparedFilter += "_" + o[1]
									}
									res.Prepared["orderByUnknown"] = append(res.Prepared["orderByUnknown"], preparedFilter)
								}
								continue
							}

							preparedFilter = dbName
							if len(o) > 1 { //add direction if provided
								preparedFilter += "_" + o[1]
							}
							res.Prepared["orderBy"] = append(res.Prepared["orderBy"], preparedFilter)
						}
					}
				}
			default: //all other names means filters
				dbName, ok = tags[f[0]]
				if !ok {
					if !ignoreUnknown {
						res.Errors[filterWithAct] = fmt.Errorf("unknown %s filter provided. Please, check filter name and %s datamodel's tags", filterWithAct, reflect.TypeOf(model))
					} else {
						res.Prepared[filterWithAct] = filters[filterWithAct]
					}
					continue
				}

				preparedFilter = dbName
				if len(f) > 1 { //add action if provided
					preparedFilter += "__" + strings.ToLower(f[1])
				}

				res.Prepared[preparedFilter] = filters[filterWithAct]
			}

		}
	}

	//calculate page for urlstruct.Pager
	var limit int64 = MaxLimit
	if limitStr, ok := res.Prepared["limit"]; ok && len(limitStr) > 0 && limitStr[0] != "" {
		limit, err = strconv.ParseInt(limitStr[0], 10, 64)
		if err != nil {
			res.Errors[DefaultFilterFirst] = fmt.Errorf("can't parse %s = %s as proper int value: %s", DefaultFilterFirst, limitStr[0], err.Error())
		}
		if limit == 0 {
			res.Errors[DefaultFilterFirst] = fmt.Errorf("%s can't be zero", DefaultFilterFirst)
		}
	}
	var offset int64 = 0
	if offsetStr, ok := res.Prepared[DefaultFilterAfter]; ok && len(offsetStr) > 0 && offsetStr[0] != "" {
		offset, err = DecodeCursorToId(offsetStr[0])
		if err != nil {
			res.Errors[DefaultFilterAfter] = fmt.Errorf("can't parse %s = %s as int value: %s", DefaultFilterAfter, offsetStr[0], err.Error())
		}
	}
	res.Prepared["page"] = []string{strconv.FormatInt(offset/limit+1, 10)} //pages starts from 1

	return res, nil
}

func ProcessFilters(filters url.Values, value string) (list []string, ok bool) {
	f, ok := filters[value]
	if ok && len(f) > 0 {
		delete(filters, value)

		//  NOTE: it's hard to pass an array through graphql-svc for the query so we pass a string (ie "2", "2,3", etc)
		for i := range f {
			temp := strings.Split(f[i], ",")
			if len(temp) > 0 {
				for j := range temp {
					list = append(list, temp[j])
				}
			} else {
				list = append(list, f[i])
			}
		}
	}

	return
}

type Ordering struct {
	FieldName string
	Distinct  string
}

func ProcessOrdering(filters url.Values, orderFieldsList []string) (ordering []Ordering, isOrderByFound bool) {
	var orderByAfterProcessingList []string

	oFilter, isOrderByFound := filters["orderBy"]
	if !isOrderByFound || len(oFilter) == 0 {
		return ordering, false
	}

	findRecordInOrderFieldsListFN := func(fieldName string) bool {
		for i := range orderFieldsList {
			if strings.EqualFold(orderFieldsList[i], strings.TrimSpace(fieldName)) {
				return true
			}
		}
		return false
	}

	oList := strings.SplitN(oFilter[0], ",", -1) //fond all orderBy clauses

	for i := range oList {
		o := strings.SplitN(oList[i], "_", 2) //slit on field name and distinct
		if len(o) == 0 {
			continue //nothing to do
		}

		distinct := ""
		if len(o) > 1 { //if direction was provided
			distinct = " " + strings.ToUpper(o[1])
		}

		if findRecordInOrderFieldsListFN(o[0]) {
			ordering = append(ordering, Ordering{
				FieldName: strings.TrimSpace(o[0]),
				Distinct:  distinct,
			})
		} else {
			orderByAfterProcessingList = append(orderByAfterProcessingList, oList[i]) //to leave "as is"
		}
	}

	if len(orderByAfterProcessingList) > 0 { //if we have some unprocessed filters
		filters["orderBy"] = []string{strings.Join(orderByAfterProcessingList, ",")}
	} else {
		delete(filters, "orderBy")
	}

	return ordering, isOrderByFound
}

func MakeEmptyFilterFunc() func(q *orm.Query) (*orm.Query, error) {
	return func(q *orm.Query) (*orm.Query, error) {
		return q, nil
	}
}
