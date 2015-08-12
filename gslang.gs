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


enum TimeUnit{
    Hour,Minute,Second,Millisecond,Micoseconds,Naosecond
}

table Duration{
    int32 Value;
    TimeUnit Unit;
}
