package mxf

type Status int

const (
	StatusOK Status = iota
	StatusTransaction
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
	iid := x.GetIid()
	if iid >= len(m.p) {
		m.Grow(iid * 2)
	}
	m.p[iid] = x
	if m.s == StatusTransaction {
		m.cacheSet(iid)
	}
	return m
}

func (m *Imap) Get(iid int) Identifier {
	if iid < len(m.p) {
		return m.p[iid]
	}
	return nil
}

func (m *Imap) BulkSet(l []Identifier) {
	if len(l) == 0 {
		return
	}

	m.Begin()

	_ = l[len(l)-1]
	for len(l) > 8 {
		m.Set(l[0]).Set(l[1]).Set(l[2]).Set(l[3]).Set(l[4]).Set(l[5]).Set(l[6]).Set(l[7])
		l = l[8:]
	}
	for i := range l {
		m.Set(l[i])
	}

	m.Commit()
}

func (m *Imap) Begin() bool {
	if m.s == StatusTransaction {
		return false
	}
	m.s = StatusTransaction
	for k := range m.cid {
		delete(m.cid, k)
	}
	return true
}

func (m *Imap) Commit() bool {
	if m.s != StatusTransaction {
		return false
	}
	_ = m.p[len(m.p)-1]
	for i := 0; i < len(m.p); i += 8 {
		m.clear(i).clear(i + 1).clear(i + 2).clear(i + 3).clear(i + 4).clear(i + 5).clear(i + 6).clear(i + 7)
	}
	m.s = StatusOK
	return true
}

func (m *Imap) clear(iid int) *Imap {
	if iid >= len(m.p) || m.p[iid] == nil {
		return m
	}
	if _, ok := m.cid[iid]; !ok {
		m.p[iid] = nil
	}
	return m
}

func (m *Imap) cacheSet(iid int) {
	if m.cid == nil {
		m.cid = make(map[int]bool, len(m.p))
	}
	m.cid[iid] = true
}
