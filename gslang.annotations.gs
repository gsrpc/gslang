package gslang.annotations;


// attribute target flag
enum Target{
    Package,Script,Table,Method,Param,Enum
}

// attribute Usage attribute
@Usage(Target.Table)
table Usage{
    Target Target;
}
