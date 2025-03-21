package dbx

func Create(items any) error {
	return DB.Create(items).Error
}

func CreateWithMap[T any](value map[string]any) error {
	return DB.Model(new(T)).Create(value).Error
}

func CreateWithMaps[T any](values []map[string]any) error {
	return DB.Model(new(T)).Create(values).Error
}
