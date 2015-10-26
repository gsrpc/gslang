package gslang.annotations;

using gslang.Flag;

// attribute target flag
@Flag
enum Target{
    Module(1),Script(2),Table(4),Method(8),Param(16),Enum(32)
}

// attribute Usage attribute
@Usage(Target.Table)
table Usage{
    Target Target;
}
