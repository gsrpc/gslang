package gslang;

using gslang.annotations.Target;
using gslang.annotations.Usage;

// gslang exception attribute
@Usage(Target.Table)
table Exception{
}


@Usage(Target.Enum)
table Flag{
}

@Usage(Target.Script)
table Lang{
    string Name; // language name
    string Package; //define language package name
}
