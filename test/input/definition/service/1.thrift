service D{} /* 123*/service C extends D {
	/*
	* normal
	 */void foo(
1:i32 a
)throws (1:B b) (k="v")
// oneway
oneway i32 bar() (a="b")
}(k="v")