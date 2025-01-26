package setx

import mapset "github.com/deckarep/golang-set/v2"

type Set struct {
	MapSet mapset.Set[any]
}

func New() *Set {
	return &Set{
		MapSet: mapset.NewSet[any](),
	}
}

func (set *Set) Add(val any) bool {
	return set.MapSet.Add(val)
}
func (set *Set) Append(val ...any) int {
	return set.MapSet.Append(val)
}
func (set *Set) Clear() {
	set.MapSet.Clear()
}
func (set *Set) Contains(val ...any) bool {
	return set.MapSet.Contains(val)
}
func (set *Set) IsEmpty() bool {
	return set.MapSet.IsEmpty()
}
func (set *Set) Remove(i any) {
	set.MapSet.Remove(i)
}
func (set *Set) String() string {
	return set.MapSet.String()
}
