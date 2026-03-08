package clause

type ForUpdate struct {
	NoWait     bool
	SkipLocked bool
}

func (f *ForUpdate) NoWaitMode() *ForUpdate {
	f.NoWait = true
	return f
}

func (f *ForUpdate) SkipLockedMode() *ForUpdate {
	f.SkipLocked = true
	return f
}

func (f ForUpdate) Parse() string {
	clause := "FOR UPDATE"

	if f.NoWait {
		clause += " NOWAIT"
	}

	if f.SkipLocked {
		clause += " SKIP LOCKED"
	}

	return clause
}
