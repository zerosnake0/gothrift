exception B {
	1: optional i32& A = 1 xsd_optional
	xsd_nillable xsd_attrs {
		i32 A
	} (a="b")
} (k="v")