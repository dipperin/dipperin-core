#include "dipc/dipc.hpp"
#include "rapidjson/stringbuffer.h"
#include "rapidjson/writer.h"
#include "rapidjson/document.h"
using namespace dipc;


char passwordStorec[] = "passwordStore";
char addrStorec[] = "addrStore";
char ownerc[] = "owner";
class PasswordManage : public Contract {
private:
    Map<passwordStorec, std::string, std::string> passwordStore;
    Map<addrStorec, std::string, std::string> addrStore; 
    String<ownerc> owner;

    void isOwner(){
        std::string callerStr = caller().toString();
        std::string ownerStr = owner.get();
        DipcAssert(callerStr == owner.get());
    }

public:
   PAYABLE void init();
   EXPORT void registerPassword(char* passwordName);
   CONSTANT char* queryPasswordByAddr(char* _addr);
   CONSTANT char* queryAddrByPassword(char* _password);
};


DIPC_EVENT(depositingValue,uint64_t);
DIPC_EVENT(ErrEvent, const char*, const char*);
DIPC_EVENT(registerPasswordEvent, const char*, const char*);



