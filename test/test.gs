package gslang.test;

using gslang.AttrUsage; // tail comment
/*same line comment*/using gslang.AttrTarget;
// header line comment
using gslang.Exception;


// Description define new Attribute
@AttrUsage(AttrTarget.Package|AttrTarget.Script)
table Description {
    string Text; // Description text
    //long texts
    string LongText;
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
