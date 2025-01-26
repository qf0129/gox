package dbx

func Update(obj any) error {
	return DB.Save(obj).Error
}

func UpdateOneFileds(target any, fields []string) error {
	return DB.Model(target).Select(fields).Updates(target).Error
}

func UpdateOneWithMap(target any, mapValue map[string]any) error {
	return DB.Model(target).Updates(mapValue).Error
}

func UpdateOneByPk[T any](pk any, structValue T) error {
	return DB.Model(new(T)).Where(map[string]any{Option.ModelPrimaryKey: pk}).Updates(structValue).Error
}

func UpdateOneFiledsByPk[T any](pk any, structValue T, fields []string) error {
	return DB.Model(new(T)).Where(map[string]any{Option.ModelPrimaryKey: pk}).Select(fields).Updates(structValue).Error
}

func UpdateMany[T any](where map[string]any, structValue T) error {
	return DB.Model(new(T)).Where(where).Updates(structValue).Error
}

func UpdateManyWithMap[T any](where map[string]any, mapValue map[string]any) error {
	return DB.Model(new(T)).Where(where).Updates(mapValue).Error
}
