package common

var Assert, Xk = assert, xk

func assert(b bool, msg interface{}) {
	if !b {
		panic(msg)
	}
}

func xk(err error) {
	if err != nil {
		panic(err)
	}
}
