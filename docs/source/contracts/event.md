### Events
Events provides an abstraction of the logging capabilities of WAVM. Applications can subscribe to and listen to these events through the client's RPC interface. The event is declared by the keyword DIPC_EVENT. The event only needs to declare the event name and parameters, and no return value. The event parameter type is consistent with the parameter type restrictions of the externally accessible function.

```c++
//Declaration
DIPC_EVENT(event_name,int32, string);

//Call
int32 val1;
string val2;
DIPC_EMIT_EVENT(event_name,val1,val2）;
```