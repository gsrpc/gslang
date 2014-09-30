//test gslang's import instruction
import (
    "github.com/gsdocker/gslang/testing" //import testing
)

[gslang.AttrUsage(gslang.AttrTarget.Script)]
table test {status bool;}

//hello
[test()] //world
//import test

[test] //world
[test(false)] //world
