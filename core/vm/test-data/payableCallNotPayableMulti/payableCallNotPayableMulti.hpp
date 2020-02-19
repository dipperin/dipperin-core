#pragma once
#include <dipc/dipc.hpp>
#include <string>

using namespace dipc;

char tmp[] = "supply";
char bal[] = "balance";
char na[] = "dipc";
char allow[] = "allowance";
char ow[] = "own";
char sy[] = "symbol";
char newbl[] = "nbalance";
class TestToken : public Contract {
public: 
    EXPORT void init(const char* tokenName,const char* symbol, uint64_t supply);
    PAYABLE void transfer(const char* to, uint64_t value);
    EXPORT void approve(const char* spender, uint64_t value);
    EXPORT void transferFrom(const char* from, const char* to, uint64_t value);
    PAYABLE uint64_t getBalance(const char* own);
    CONSTANT uint64_t getApproveBalance(const char* from, const char* approved);
    EXPORT bool burn(int128_t _value);

    void stop(){
        isOwner();
        stopped = true;    
    }
    void start() {
        isOwner();
        stopped = false;
    }
    void setName(const char* _name){
        isOwner();
        *name = _name;
    }

private: 
    String<na> name ;
    String<sy> symbol;
    uint8_t decimals = 6;
    Map<bal, std::string, uint64_t >  balance;
    Map<newbl, Address, uint64_t> nbalance;
    Map<allow, std::string, uint64_t> allowance;
    Uint64<tmp> total_supply;
    bool stopped = false;
    String<ow> owner;
    
    inline void isOwner(){
        Address callerAddr = caller();
        std::string callerStr = callerAddr.toString();
        DipcAssert(owner.get() == callerStr);
    }
};
DIPC_EVENT(Tranfer, const char*, const char*, uint64_t);
DIPC_EVENT(Approval, const char*, const char*, uint64_t);
DIPC_EVENT(GetBalance, const char*, const char*, uint64_t);
