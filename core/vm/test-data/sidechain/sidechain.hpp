#pragma once
#include <dipc/dipc.hpp>
using namespace dipc;

char ow[6] = "owner";
char fr[5] = "from";
char am[7] = "amount";
char ch[8] = "chainid";
char he[7] = "height";

class sidechain : public Contract
{
public:
    PAYABLE void init(char* addr, uint64_t chainid);
    EXPORT void transfer(char* proof, char* oracleAddr);
    EXPORT void withdraw();

private:
    String<ow> owner;
    String<fr> from;
    Uint64<am> amount;
    Uint64<ch> chainid;
    Uint64<he> height;

    inline void isZeroAddress(std::string const &_s)
    {
        DipcAssert(_s != "0x0");
    }
};
