package sqltypes

type stringType struct {
	typecommon[stringType, string]
}

func Char(name string, length int) *stringType {
	ins := &stringType{}
	ins.sqltype(name, "char", length)
	return ins
}

func VarChar(name string, length int) *stringType {
	ins := &stringType{}
	ins.sqltype(name, "varchar", length)
	return ins
}

func Text(name string) *stringType {
	ins := &stringType{}
	ins.sqltype(name, "text")
	return ins
}
