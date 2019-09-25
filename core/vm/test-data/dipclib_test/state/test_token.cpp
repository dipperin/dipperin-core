#include "test_token.hpp"

char tmp[7] = "supply";
char bal[8] = "balance";
char na[5] = "dipc";
char allow[10] = "allowance";
char ow[4] = "own";
char sy[7] = "symbol";

DIPC_EVENT(Tranfer, const char*, const char*, uint64_t);
DIPC_EVENT(Approval, const char*, const char*, uint64_t);

void testtoken::TestTokens()
{
    print("\r\n dipc lib test TestTokens start\r\n");

    //init testToken
    std::string t = "testToken";
    char* tokenName = (char*)t.data();
    std::string e = "erc20";
    char* symbol = (char*)e.data();
    uint64_t supply = 2000;
    init(tokenName, symbol, supply);

    //isowner
    isOwner();

    //get balance
    uint64_t balance;
    balance = getBalance(owner.get().data());
    DipcAssertEQ(supply, balance);

    //transfer
    char toAddr[] = "0x0000A1B7a7B7BA883E1BEF8B9D312a430e857Ee20B17";
    transfer(toAddr, 2);
    balance = getBalance(toAddr);
    DipcAssertEQ(balance, 2);
    uint64_t ownerBalance;
    ownerBalance = getBalance(owner.get().data());
    DipcAssertEQ(ownerBalance, (supply - 2));

    //approve
    char approveAddr[] = "0x000002430bcfEE64A8DEeDD44C1Cb447662c383e8520";
    approve(approveAddr, 100);
    uint64_t approveValue;

    approveValue = getApproveBalance(caller2().toString().data(), approveAddr);
    DipcAssertEQ(approveValue, 100);

    //burn
    burn(100);

    print("\r\n dipc lib test TestTokens success\r\n");
    return;
}
void testtoken::init(char *tokenName, char *sym, uint64_t supply)
{
    Address2 callerAddr = caller2();
    std::string callerStr = callerAddr.toString();
    *owner = callerStr;
    *name = tokenName;
    *symbol = sym;
    *total_supply = supply;
    (*balance)[owner.get()] = supply;
    //prints(&owner.get()[0]);
    DIPC_EMIT_EVENT(Tranfer, "", &owner.get()[0], supply);
}

void testtoken::transfer(const char *to, uint64_t value)
{
    Address2 callerAddr = caller2();
    std::string callStr = callerAddr.toString();
    //uint64_t originValue = balance.get()[callStr];
    //prints(&callStr[0]);
    //printui(originValue);
    //printui(value);
    bool result = (balance.get()[callStr] >= value);
    DipcAssert(result);

    std::string toStr = CharToAddress2Str(to);
    DipcAssert(balance.get()[toStr] + value >= balance.get()[toStr]);

    (*balance)[callStr] = balance.get()[callStr] - value;
    (*balance)[toStr] = balance.get()[toStr] + value;
    DIPC_EMIT_EVENT(Tranfer, &(callStr[0]), to, value);
}

void testtoken::transferFrom(const char *from, const char *to, uint64_t value)
{
    Address2 callerAddr = caller2();
    std::string fromStr = CharToAddress2Str(from);
    std::string toStr = CharToAddress2Str(to);

    DipcAssert(balance.get()[fromStr] >= value);
    DipcAssert(balance.get()[toStr] + value >= balance.get()[toStr]);
    DipcAssert(allowance.get()[fromStr + callerAddr.toString()] >= value);

    (*balance)[toStr] = balance.get()[toStr] + value;
    (*balance)[fromStr] = balance.get()[fromStr] - value;
    (*allowance)[fromStr + callerAddr.toString()] = allowance.get()[fromStr + callerAddr.toString()] - value;
    DIPC_EMIT_EVENT(Tranfer, from, to, value);
}
void testtoken::approve(const char *spender, uint64_t value)
{
    Address2 callerAddr = caller2();
    std::string spenderStr = CharToAddress2Str(spender);

    uint64_t total = allowance.get()[callerAddr.toString() + spenderStr] + value;
    (*allowance)[callerAddr.toString() + spenderStr] = total;
    DIPC_EMIT_EVENT(Approval, &(callerAddr.toString()[0]), spender, value);
}
void testtoken::burn(uint64_t value)
{
    Address2 callerAddr = caller2();
    DipcAssert(balance.get()[callerAddr.toString()] >= value);
    DipcAssert(balance.get()[owner.get()] + value >= balance.get()[owner.get()]);
    //prints("burn======");
    //prints(&owner.get()[0]);
    //printui(balance.get()[owner.get()]);
    (*balance)[callerAddr.toString()] -= value;
    (*balance)[owner.get()] += value;
    //printui(balance.get()[owner.get()]);
    DIPC_EMIT_EVENT(Tranfer, &(callerAddr.toString()[0]), &(owner.get()[0]), value);
}

uint64_t testtoken::getBalance(const char *own)
{
    //prints(own);
    std::string ownerStr = CharToAddress2Str(own);
    uint64_t ba = balance.get()[ownerStr];
    //printui(ba);
    DIPC_EMIT_EVENT(Tranfer, "", own, ba);
    return ba;
}

uint64_t testtoken::getApproveBalance(const char *from, const char *approved)
{
    //prints("getApproveBalance");
    //prints(from);
    // prints(approved);
    std::string fromStr = CharToAddress2Str(from);
    std::string approvedStr = CharToAddress2Str(approved);
    uint64_t re = allowance.get()[fromStr + approvedStr];
    //printui(re);
    DIPC_EMIT_EVENT(Tranfer, from, approved, re);
    return re;
}

inline void testtoken::isOwner()
{
    Address2 callerAddr = caller2();
    std::string callerStr = callerAddr.toString();
    DipcAssert(owner.get() == callerStr);
}

void testtoken::stop()
{
    isOwner();
    stopped = true;
}
void testtoken::start()
{
    isOwner();
    stopped = false;
}
void testtoken::setName(const char *_name)
{
    isOwner();
    //当使用非指针类型对存储型变量赋值时，会报unreachable错误
    *name = _name;
}