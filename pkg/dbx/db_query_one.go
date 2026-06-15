package dbx

func QueryOneByStruct(params any, target any) (err error) {
	return DB.Where(params).First(&target).Error
}

func QueryOneByMap[T any](filter map[string]any) (result T, err error) {
	err = DB.Model(new(T)).Where(filter).First(&result).Error
	return result, err
}

func QueryOneByPk[T any](pk any) (T, error) {
	return QueryOneByMap[T](map[string]any{Option.ModelPrimaryKey: pk})
}

func QueryOneById[T any](id any) (T, error) {
	return QueryOneByMap[T](map[string]any{"id": id})
}

func QueryOneByUid[T any](uid any) (T, error) {
	return QueryOneByMap[T](map[string]any{"uid": uid})
}

func QueryOneWithPreload[T any](filter map[string]any, preloadMap map[string][]any) (result T, err error) {
	query := DB.Model(new(T))
	query = setQueryFilter(query, filter)
	query = setQueryPreload(query, preloadMap)
	err = query.First(&result).Error
	return
}

func QueryOneByIdWithPreload[T any](id any, preloadMap map[string][]any) (T, error) {
	return QueryOneWithPreload[T](map[string]any{"id": id}, preloadMap)
}

func QueryOneByUidWithPreload[T any](uid any, preloadMap map[string][]any) (T, error) {
	return QueryOneWithPreload[T](map[string]any{"uid": uid}, preloadMap)
}

func QueryMaxId[T any]() (int64, error) {
	var maxId int64
	err := DB.Model(new(T)).Select("MAX(id)").Scan(&maxId).Error
	return maxId, err
}

func QueryTargetByUid(uid any, tgt any) error {
	return DB.Model(tgt).Where(map[string]any{"uid": uid}).First(tgt).Error
}

func QueryPluck[T any, R any](cloumn string, filter map[string]any) (result []R, err error) {
	err = DB.Model(new(T)).Where(filter).Distinct().Pluck(cloumn, &result).Error
	return
}
