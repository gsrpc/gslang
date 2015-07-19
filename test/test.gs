using gslang.AttributeUsage;
using gslang.AttributeTarget;
using gslang.Exception;


@AttributeUsage(AttributeTarget.Package|AttributeTarget.Script)
table Description {
    String Text; // Description text
}


@Exception
table RemoteException {}

// HttpREST API
contract HttpREST {
    // invoke http post method
    void Post(byte[] content) throws (RemoteException);
    // get invoke http get method
    byte[] Get() throws (RemoteException);
}
