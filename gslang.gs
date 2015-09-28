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

@Usage(Target.Module)
table Package{
    string Lang; // language name
    string Name;
    string Redirect; //define language package name
}
