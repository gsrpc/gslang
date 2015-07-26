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
table RemoteException {}

// HttpREST API
contract HttpREST {
    // invoke http post method
    void Post(byte[] content) throws (RemoteException,CodeException);
    // get invoke http get method
    byte[] Get() throws (RemoteException);
}

// remote exception
@Exception
enum CodeException {

}
