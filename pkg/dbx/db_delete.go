package dbx

func Delete(target any) error {
	return DB.Delete(target).Error
}

func DeleteMany(targets []any) error {
	return DB.Delete(targets).Error
}

func DeleteById[T any](id any) error {
	return DeleteByMap[T](map[string]any{"id": id})
}

func DeleteByUid[T any](uid any) error {
	return DeleteByMap[T](map[string]any{"uid": uid})
}

func DeleteByMap[T any](filter map[string]any) error {
	return DB.Where(filter).Delete(new(T)).Error
}
