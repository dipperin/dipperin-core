#pragma once
#include <dipc/dipc.hpp>

using namespace dipc;

extern char tmp[7];
extern char bal[8];
extern char na[5];
extern char allow[10];
extern char ow[4];
extern char sy[7];

class testtoken
{
public:
    void TestTokens();

private:
    String<na> name;
    String<sy> symbol;
    uint8_t decimals = 6;
    Map<bal, std::string, uint64_t> balance;
    //         存放授权地址    存放授权金额
    Map<allow, std::string, uint64_t> allowance;
    Uint64<tmp> total_supply;
    bool stopped = false;
    // 当owner和name使用同一个char*标记na时，会导致数据出错
    String<ow> owner;

    inline void isOwner();
    void start();
    void stop();
    void setName(const char *_name);
    void init(char *tokenName, char *sym, uint64_t supply);
    void transfer(const char *to, uint64_t value);
    void transferFrom(const char *from, const char *to, uint64_t value);
    void approve(const char *spender, uint64_t value);
    void burn(uint64_t value);
    uint64_t getBalance(const char *own);
    uint64_t getApproveBalance(const char *from, const char *approved);
};

