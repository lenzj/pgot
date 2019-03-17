package pgot

type frontMatter struct {
	fd string                 // source directory
	fm map[string]interface{} // key,value pairs
	ns string                 // namespace
}

type fmStack struct {
	fma   []*frontMatter
	count int
}

func newfmStack() *fmStack {
	return &fmStack{
		fma:   make([]*frontMatter, 0),
		count: 0,
	}
}

func (s *fmStack) push(fm *frontMatter) {
	s.fma = append(s.fma[:s.count], fm)
	s.count++
}

func (s *fmStack) pop() *frontMatter {
	if s.count == 0 {
		return nil
	}
	s.count--
	return s.fma[s.count]
}

func (s *fmStack) last() *frontMatter {
	return s.fma[s.count-1]
}

func (s *fmStack) first() *frontMatter {
	return s.fma[0]
}
