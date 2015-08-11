package gslang.test;

using gslang.AttrUsage; // tail comment
/*same line comment*/using gslang.AttrTarget;
// header line comment
using gslang.Exception;
using gslang.Flag;
using gslang.Duration;
using gslang.TimeUnit;


// Description define new Attribute
@AttrUsage(AttrTarget.Package|AttrTarget.Script)
table Description {
    string Text; // Description text
    //long texts
    string LongText;
}


@AttrUsage(AttrTarget.Method)
table Async {
}

@AttrUsage(AttrTarget.Param)
table Out {
}

@AttrUsage(AttrTarget.Method)
table Timeout {
    Duration Duration;
}

// remote exception
@Exception
table RemoteException {
    Description Description;
}

table KV {
    string Key;
    string Value;
}



// HttpREST API
contract HttpREST {
    @Async
    // invoke http post method
    @Timeout(Duration(-100,TimeUnit.Second))
    void Post(@Out byte[] content) throws (RemoteException,CodeException);
    // get invoke http get method
    byte[] Get(KV[] properties) throws (RemoteException);
}

// remote exception
@Flag
@Exception
enum CodeException {
    /*test*/
    Success, //test 2
    Unknown(2)
}
