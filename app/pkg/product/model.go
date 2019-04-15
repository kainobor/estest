package product

import (
	"fmt"
	"html"
	"strconv"
	"strings"
)

type (
	Product struct {
		Title string
		Brand string
		Price float64
		Stock int64
	}

	Filter interface {
		Name() string
		MarshalJSON() ([]byte, error)
	}

	StringFilter struct {
		name  string
		Match string
	}

	FloatFilter struct {
		name string
		Eq   float64
	}

	IntFilter struct {
		name string
		Eq   int64
	}

	Page struct {
		Num  int64
		Size int64
	}

	Sort struct {
		Field string
		Desc  bool
	}
)

const (
	Sep = ":"

	nameBrandFilter = "brand"
	namePriceFilter = "price"
	nameStockFilter = "stock"
)

func ParseFilter(query string) (Filter, error) {
	filterAttrs := strings.Split(query, Sep)
	if len(filterAttrs) != 2 {
		return nil, fmt.Errorf("wrong format of filter. Required format: filter_name:filter_val")
	}

	switch filterName := filterAttrs[0]; filterName {
	case nameBrandFilter:
		escapedString := html.EscapeString(filterAttrs[1])
		return &StringFilter{name: filterName, Match: escapedString}, nil
	case namePriceFilter:
		price, err := strconv.ParseFloat(filterAttrs[1], 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse `%s` as float", filterAttrs[1])
		}
		return &FloatFilter{name: filterName, Eq: price}, nil
	case nameStockFilter:
		price, err := strconv.ParseInt(filterAttrs[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse `%s` as integer", filterAttrs[1])
		}
		return &IntFilter{name: filterName, Eq: price}, nil
	}

	return nil, fmt.Errorf("wrong filter name")
}

func ParseSort(query string) (*Sort, error) {
	sortAttrs := strings.Split(query, Sep)
	switch len(sortAttrs) {
	case 1:
		return &Sort{Field: html.EscapeString(sortAttrs[0])}, nil
	case 2:
		switch dir := sortAttrs[1]; dir {
		case "desc":
			return &Sort{Field: html.EscapeString(sortAttrs[0]), Desc: true}, nil
		case "asc":
			return &Sort{Field: html.EscapeString(sortAttrs[0])}, nil
		default:
			return nil, fmt.Errorf("wrong sort direstion. Available: desc, asc")
		}
	default:
		return nil, fmt.Errorf("wrong sort format. Required format: field_name:direction")
	}
}

func (f *StringFilter) Name() string {
	return f.name
}

func (f *StringFilter) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"%s": "%s"}`, f.name, f.Match)), nil
}

func (f *FloatFilter) Name() string {
	return f.name
}

func (f *FloatFilter) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"%s": %.2f}`, f.name, f.Eq)), nil
}

func (f *IntFilter) Name() string {
	return f.name
}

func (f *IntFilter) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"%s": %d}`, f.name, f.Eq)), nil
}
