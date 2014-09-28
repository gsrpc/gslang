//The attribute's target enum
enum AttrTarget {
    Package(1),
    Script(2),
    Table(4),
    Struct(8),
    Contract(16),
    Method(32),
    Return(64),
    Param(128)
}

[AttrUsage(AttrUsage.Table)]
table AttrUsage {
    Target AttrTarget;
}

//indicate AST node table is a struct
[AttrUsage(AttrUsage.Table)]
table Struct {}
