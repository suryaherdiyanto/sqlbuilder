package clause

func (s *ForShare) TestForShare() {
	fs := ForShare{}

	stmt := fs.Parse()
	expected := "FOR SHARE"
	if stmt != expected {
		panic("TestForShare failed: expected " + expected + ", got " + stmt)
	}

	fs.NoWaitMode()
	stmt = fs.Parse()
	expected = "FOR SHARE NOWAIT"
	if stmt != expected {
		panic("TestForShare with NoWait failed: expected " + expected + ", got " + stmt)
	}

	fs.SkipLockedMode()
	stmt = fs.Parse()
	expected = "FOR SHARE NOWAIT SKIP LOCKED"
	if stmt != expected {
		panic("TestForShare with SkipLocked failed: expected " + expected + ", got " + stmt)
	}
}
