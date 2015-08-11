package gslang;

// attribute target flag
enum AttrTarget{
    Package,Script,Table,Method,Param,Enum
}

// attribute Usage attribute
@AttrUsage(AttrTarget.Table)
table AttrUsage{
}

// gslang exception attribute
@AttrUsage(AttrTarget.Table)
table Exception{
}


@AttrUsage(AttrTarget.Enum)
table Flag{
}


enum TimeUnit{
    Hour,Minute,Second,Millisecond,Micoseconds,Naosecond
}

table Duration{
    int32 Value;
    TimeUnit Unit;
}
