package clause

type ForShare struct {
	NoWait     bool
	SkipLocked bool
}

func (f *ForShare) NoWaitMode() *ForShare {
	f.NoWait = true
	return f
}

func (f *ForShare) SkipLockedMode() *ForShare {
	f.SkipLocked = true
	return f
}

func (f ForShare) Parse() string {
	clause := "FOR SHARE"

	if f.NoWait {
		clause += " NOWAIT"
	}

	if f.SkipLocked {
		clause += " SKIP LOCKED"
	}

	return clause
}
