package product

import (
	"fmt"
	"github.com/kainobor/estest/app/lib/logger"
	"github.com/kainobor/estest/app/pkg/product"
	"github.com/kainobor/estest/app/srv/ctrl/response"
	"net/http"
	"strconv"
)

type Product struct {
	prodSrv product.Service
}

const (
	attrPageNum  = "page_num"
	attrPageSize = "page_size"
	attrSort     = "sort"
)

func New(prodSrv product.Service) *Product {
	return &Product{prodSrv: prodSrv}
}

func (ctrl *Product) GetProduct(w http.ResponseWriter, r *http.Request) {
	log := logger.New(r.Context())
	var err error

	query := r.URL.Query().Get("q")
	if query == "" {
		response.WriteError(r.Context(), w, http.StatusBadRequest, "Attr `q` is required")
		return
	}

	var filter product.Filter
	filterString := r.URL.Query().Get("filter")
	if filterString != "" {
		filter, err = product.ParseFilter(filterString)
		if err != nil {
			response.WriteError(r.Context(), w, http.StatusBadRequest, "Error: "+err.Error())
			return
		}
	}

	var page *product.Page
	pageNumStr := r.URL.Query().Get(attrPageNum)
	pageSizeStr := r.URL.Query().Get(attrPageSize)
	if pageNumStr != "" && pageSizeStr != "" {
		pageNum, numErr := strconv.ParseInt(pageNumStr, 10, 64)
		pageSize, sizeErr := strconv.ParseInt(pageSizeStr, 10, 64)
		if numErr != nil || sizeErr != nil {
			response.WriteError(r.Context(), w, http.StatusBadRequest, fmt.Sprintf("Wrong format of paging: `%s` and `%s` need to be int", attrPageNum, attrPageSize))
			return
		}
		page = &product.Page{Num: pageNum, Size: pageSize}
	}

	var sort *product.Sort
	sortStr := r.URL.Query().Get(attrSort)
	if sortStr != "" {
		sort, err = product.ParseSort(sortStr)
		if err != nil {
			response.WriteError(r.Context(), w, http.StatusBadRequest, "Error: "+err.Error())
			return
		}
	}

	products, err := ctrl.prodSrv.List(r.Context(), query, filter, sort, page)
	if err != nil {
		log.Errorw("failed to list products", "err", err, "query", query, "filter", filter, "page", page)
		response.WriteError(r.Context(), w, http.StatusInternalServerError, "Failed to load products")
		return
	}

	response.WriteJSON(r.Context(), w, products)
	return
}
