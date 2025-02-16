package dbx

func QueryOne[T any](filter map[string]any) (target *T, err error) {
	query := DB.Model(new(T))
	query = setQueryFilter(query, filter)
	err = query.First(target).Error
	return
}
func QueryOneTarget[T any](target *T, filter map[string]any) error {
	query := DB.Model(new(T))
	query = setQueryFilter(query, filter)
	return query.First(target).Error
}

func QueryOneByStruct(params any, target any) (err error) {
	return DB.Where(params).First(&target).Error
}

func QueryOneByPk[T any](pk any) (result *T, err error) {
	err = DB.Model(new(T)).Where(map[string]any{Option.ModelPrimaryKey: pk}).First(&result).Error
	return result, err
}

func QueryOneByMap[T any](filter map[string]any) (result *T, err error) {
	err = DB.Model(new(T)).Where(filter).First(&result).Error
	return result, err
}

func QueryOneFieldsMap[T any](fields []string, filter map[string]any) (*map[string]any, error) {
	t := &map[string]any{}
	query := DB.Model(new(T)).Select(fields)
	query = setQueryFilter(query, filter)
	err := query.First(t).Error
	return t, err
}

func QueryOneWithPreload[T any](filter map[string]any, preloadMap map[string][]any) (result *T, err error) {
	query := DB.Model(new(T))
	query = setQueryFilter(query, filter)
	query = setQueryPreload(query, preloadMap)
	err = query.First(&result).Error
	return
}

func QueryOneByPkWithPreload[T any](pk any, preloadMap map[string][]any) (*T, error) {
	return QueryOneWithPreload[T](map[string]any{Option.ModelPrimaryKey: pk}, preloadMap)
}
