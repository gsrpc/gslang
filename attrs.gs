//The attribute's target enum
enum AttrTarget(uint16) {
    Package(1),
    Script(2),
    Table(4),
    Struct(8),
    Enum(16),
    EnumVal(32),
    Field(64),
    Contract(128),
    Method(256),
    Return(512),
    Param(1024)
}

@AttrUsage(AttrTarget.Table)
table AttrUsage {
    Target AttrTarget;
}

//indicate AST node table is a struct
@AttrUsage(AttrTarget.Struct)
table Struct {}


// the enum is error code
@AttrUsage(AttrTarget.Enum)
table Error {
    UUID string;
}
