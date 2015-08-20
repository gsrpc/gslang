package gslang.annotations;

using gslang.Flag;

// attribute target flag
@Flag
enum Target{
    Package,Script,Table,Method,Param,Enum
}

// attribute Usage attribute
@Usage(Target.Table)
table Usage{
    Target Target;
}
