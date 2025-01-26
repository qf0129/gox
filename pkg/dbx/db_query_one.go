package dbx

func QueryOneTarget[T any](target *T, filter map[string]any) error {
	query := DB.Model(new(T))
	query = setQueryFilter(query, filter)
	return query.First(target).Error
}

func QueryOne[T any](filter map[string]any) (*T, error) {
	t := new(T)
	err := QueryOneTarget(t, filter)
	return t, err
}

func QueryOneByPk[T any](pk any) (result *T, err error) {
	err = DB.Model(new(T)).Where(map[string]any{Option.ModelPrimaryKey: pk}).First(&result).Error
	return result, err
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
