#pragma once
#include <dipc/dipc.hpp>
#include <string>

using namespace dipc;

char tmp[7] = "supply";
char bal[8] = "balance";
char na[5] = "dipc";
char allow[10] = "allowance";
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
        name = _name;
    }
    void transfer(const char* to, uint64_t value);
    void transferFrom(const char* from, const char* to, uint64_t value);
    void approve(const char* spender, uint64_t value);
    void burn(uint64_t _value);
private: 
    String<na> name ;
    String<na> symbol;
    uint8_t decimals = 6;
    Map<bal, std::string, uint64_t >  balance;
    //Map<allow, std::string, Map<bal, std::string, uint64_t>> allowance;
    //         存放授权地址    存放授权地址与被授权地址的拼接值
    Map<allow, std::string, std::string> allowance;
    Uint64<tmp> total_supply;
    bool stopped = false;
    //std::string addr = "0x000062be10f46b5d01Ecd9b502c4bA3d6131f6fc2e41";
    bytes result = fromHex("0x000062be10f46b5d01Ecd9b502c4bA3d6131f6fc2e41");
    Address2 owner = Address2(&result[0], 22);
    
    inline void isOwner(){
        DipcAssertEQ(owner, address2());
    }
};
// 没有加这个宏  导致编译wasm的时候没通过也没报错  待优化
DIPC_ABI(TestToken, init);
DIPC_ABI(TestToken, transfer);
DIPC_ABI(TestToken, start)
DIPC_ABI(TestToken, stop)
DIPC_ABI(TestToken, setName)
DIPC_ABI(TestToken, transferFrom)
DIPC_ABI(TestToken, approve)
DIPC_ABI(TestToken, burn)
DIPC_EVENT(Tranfer, const char*, const char*, uint64_t)
DIPC_EVENT(Approval, const char*, const char*, uint64_t)