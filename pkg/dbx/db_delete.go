package dbx

func DeleteTarget[T any](target T) error {
	return DB.Delete(target).Error
}

func DeleteTargets[T any](targets []T) error {
	return DB.Delete(targets).Error
}

func DeleteByPk[T any](pk any) error {
	return DeleteByFilter[T](map[string]any{Option.ModelPrimaryKey: pk})
}

func DeleteByPks[T any](pks []any) error {
	return DeleteByFilter[T](map[string]any{Option.ModelPrimaryKey: pks})
}

func DeleteByFilter[T any](filter map[string]any) error {
	return DB.Where(filter).Delete(new(T)).Error
}
