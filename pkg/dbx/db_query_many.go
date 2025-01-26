package dbx

// 查询分页
func QueryPage[T any](options ...*QueryOption) (PageBody[T], error) {
	query, opt := NewQuery[T](options...)
	result := PageBody[T]{Page: opt.Page, PageSize: opt.PageSize}
	if err := query.Count(&result.Total).Error; err != nil {
		return result, err
	}
	query = setQueryPage(query, result.Page, result.PageSize)
	if err := query.Find(&result.List).Error; err != nil {
		return result, err
	}
	return result, nil
}

// 查询所有
func QueryAll[T any](options ...*QueryOption) ([]T, error) {
	query, _ := NewQuery[T](options...)
	result := []T{}
	if err := query.Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

// 查询数量
func QueryCount[T any](options ...*QueryOption) (int64, error) {
	query, _ := NewQuery[T](options...)
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return count, err
	}
	return count, nil
}
