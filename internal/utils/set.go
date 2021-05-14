package utils

type SetStr map[string]struct{}

func NewSetStr(l int) SetStr {
	return make(map[string]struct{}, l)
}
func FromList(l []string) SetStr {
	s := make(map[string]struct{}, 0)
	for _, elem := range l {
		s[elem] = struct{}{}
	}
	return s
}

func (set SetStr) Put(s string) {
	set[s] = struct{}{}
}

func (set SetStr) Contains(s string) bool {
	_, ok := set[s]
	return ok
}
func (set SetStr) ToList() []string {
	l := make([]string, 0, len(set))
	for key := range set {
		l = append(l, key)
	}
	return l
}
