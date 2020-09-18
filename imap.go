package mxf

type Status int

const (
	StatusOK Status = iota
	StatusRestCollect
)

type Identifier interface {
	GetIid() int
}

type Imap struct {
	s   Status
	p   []Identifier
	cid map[int]bool
}

func NewImap(size int) *Imap {
	m := Imap{s: StatusOK}
	m.Grow(size)
	return &m
}

func (m *Imap) Len() int {
	return len(m.p)
}

func (m *Imap) Grow(size int) {
	if size < 0 {
		return
	}
	if m.p == nil {
		m.p = make([]Identifier, size)
	} else {
		grow := size - len(m.p)
		if grow <= 0 {
			return
		}
		m.p = append(m.p, make([]Identifier, grow)...)
	}
}

func (m *Imap) Set(x Identifier) *Imap {
	idx := x.GetIid()
	if idx >= len(m.p) {
		m.Grow(idx * 2)
	}
	m.p[idx] = x
	if m.s == StatusRestCollect {
		m.cid[idx] = true
	}
	return m
}

func (m *Imap) BulkSet(l []Identifier) {
	m.RestCollect()

	_ = l[len(l)]
	for len(l) > 8 {
		m.Set(l[0]).Set(l[1]).Set(l[2]).Set(l[3]).Set(l[4]).Set(l[5]).Set(l[6]).Set(l[7])
		l = l[8:]
	}
	for i := range l {
		m.Set(l[i])
	}

	m.RestClear()
}

func (m *Imap) RestCollect() {
	for k := range m.cid {
		delete(m.cid, k)
	}
	m.s = StatusRestCollect
}

func (m *Imap) RestClear() {
	for i := range m.cid {
		m.p[i] = nil
	}
	m.s = StatusOK
}

func (m *Imap) Get(idx int) Identifier {
	if idx < len(m.p) {
		return m.p[idx]
	}
	return nil
}
