
table Any {
    hello [1]int32;
}

[Any]

//Who interface
contract Who {
    //WhoAreYou get user name
    WhoAreYou() -> (string /*user name*/,int32)
}
