package clause

func (s *ForUpdate) TestForUpdate() {
	fu := ForUpdate{}

	stmt := fu.Parse()
	expected := "FOR UPDATE"
	if stmt != expected {
		panic("TestForUpdate failed: expected " + expected + ", got " + stmt)
	}

	fu.NoWaitMode()
	stmt = fu.Parse()
	expected = "FOR UPDATE NOWAIT"
	if stmt != expected {
		panic("TestForUpdate with NoWait failed: expected " + expected + ", got " + stmt)
	}

	fu.SkipLockedMode()
	stmt = fu.Parse()
	expected = "FOR UPDATE NOWAIT SKIP LOCKED"
	if stmt != expected {
		panic("TestForUpdate with SkipLocked failed: expected " + expected + ", got " + stmt)
	}
}
