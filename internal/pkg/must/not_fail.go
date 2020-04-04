package must

func NotFail(err error) {
	if err != nil {
		panic(err)
	}
}
