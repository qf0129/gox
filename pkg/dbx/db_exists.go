package dbx

func Exist[T any](filter map[string]any) (bool, error) {
	count, err := QueryCount[T](&QueryOption{Filter: filter})
	return count > 0, err
}

func ExistByPk[T any](pk any) (bool, error) {
	return Exist[T](map[string]any{Option.ModelPrimaryKey: pk})
}
