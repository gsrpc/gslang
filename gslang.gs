package gslang;

using gslang.annotations.Target;
using gslang.annotations.Usage;

// gslang exception attribute
@Usage(Target.Table)
table Exception{}

@Usage(Target.Enum)
table Flag{}

// POD indicate serializer this table has compact struct
@Usage(Target.Table)
table POD{}

@Usage(Target.Module)
table Package{
    string Lang; // language name
    string Name;
    string Redirect; //define language package name
}


// indicate this method don't expect call response
@Usage(Target.Method)
table Async {}
