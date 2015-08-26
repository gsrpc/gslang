package gslang.test;

using gslang.annotations.Usage; // tail comment
/*same line comment*/using gslang.annotations.Target;
// header line comment
using gslang.Exception;
using gslang.Flag;

enum TimeUnit{
    Second
}

table Duration {
    int32 Value;
    TimeUnit Unit;
}


// Description define new Attribute
@Usage(Target.Module|Target.Script)
table Description {
    string Text; // Description text
    //long texts
    string LongText;
}

@Usage(Target.Method)
table Async {
}

@Usage(Target.Param)
table Out {
}

@Usage(Target.Method)
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

@Exception
table CodeException {
}
