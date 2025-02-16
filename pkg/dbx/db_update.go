package dbx

func Update(obj any) error {
	return DB.Save(obj).Error
}

func UpdateTarget(target any, val any) error {
	return DB.Model(target).Updates(val).Error
}

func UpdateTargetFileds(target any, val any, fields []string) error {
	return DB.Model(target).Select(fields).Updates(val).Error
}

func UpdateMap[T any](where map[string]any, mapValue map[string]any) error {
	return DB.Model(new(T)).Where(where).Updates(mapValue).Error
}

func UpdateByMap[T any](where map[string]any, val T) error {
	return DB.Model(new(T)).Where(where).Updates(val).Error
}

func UpdateByMapWithFields[T any](where map[string]any, val T, fields []string) error {
	return DB.Model(new(T)).Where(where).Select(fields).Updates(val).Error
}
