typedef i32 a(k="k");


enum b {
	a b = 1
}(k="v")

senum c{
	"k"
	"e";
}

struct A xsd_all {
	1: optional i32& A = 1 xsd_optional
	xsd_nillable xsd_attrs {
		i32 A
	} (a="b")
} (k="v")

exception B {
	1: optional i32& A = 1 xsd_optional
	xsd_nillable xsd_attrs {
		i32 A
	} (a="b")
} (k="v")

service D{}

service C extends D {
	/*
	* normal
	 */void foo(
1:i32 a
)throws (1:B b) (k="v")
// oneway
oneway i32 bar() (a="b")
}(k="v")