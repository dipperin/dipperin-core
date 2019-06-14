#pragma once
#include <dipc/dipc.hpp>
#include <string>

using namespace dipc;

char tmp[7] = "supply";
char bal[8] = "balance";
char na[5] = "dipc";
char allow[10] = "allowance";
char ow[4] = "own";
char sy[7] = "symbol";
class TestToken : public Contract {
public: 
    void init(char* tokenName, char* symbol, uint64_t supply);
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
        //当使用非指针类型对存储型变量赋值时，会报unreachable错误
        *name = _name;
    }
    void transfer(const char* to, uint64_t value);
    void transferFrom(const char* from, const char* to, uint64_t value);
    void approve(const char* spender, uint64_t value);
    void burn(uint64_t _value);
    uint64_t getBalance(const char* own);
    uint64_t getApproveBalance(const char* from, const char* approved);
private: 
    String<na> name ;
    String<sy> symbol;
    uint8_t decimals = 6;
    Map<bal, std::string, uint64_t >  balance;
    //Map<allow, std::string, Map<bal, std::string, uint64_t>> allowance;
    //         存放授权地址    存放授权金额
    Map<allow, std::string, uint64_t> allowance;
    Uint64<tmp> total_supply;
    bool stopped = false;
    // 当owner和name使用同一个char*标记na时，会导致数据出错
    String<ow> owner;
    
    inline void isOwner(){
        Address2 callerAddr = caller2();
        std::string callerStr = callerAddr.toString();
        prints("isOwner");
        prints(&callerStr[0]);
        prints(&owner.get()[0]);
        println(owner.get() == callerStr);
        println(owner.get().compare(callerStr));
        DipcAssert(owner.get() == callerStr);
    }
};
// 没有加这个宏  导致编译wasm的时候没通过也没报错  待优化
DIPC_ABI(TestToken, init);
DIPC_ABI(TestToken, transfer);
DIPC_ABI(TestToken, start);
DIPC_ABI(TestToken, stop);
DIPC_ABI(TestToken, setName);
DIPC_ABI(TestToken, transferFrom);
DIPC_ABI(TestToken, getBalance);
DIPC_ABI(TestToken, getApproveBalance);
DIPC_ABI(TestToken, approve);
DIPC_ABI(TestToken, burn);
DIPC_EVENT(Tranfer, const char*, const char*, uint64_t);
DIPC_EVENT(Approval, const char*, const char*, uint64_t);