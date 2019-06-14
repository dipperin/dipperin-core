#include "token.hpp"

void TestToken::init(char* tokenName, char* sym, uint64_t supply){
    Address2 callerAddr = caller2();
    std::string callerStr = callerAddr.toString();
    *owner = callerStr;
    *name = tokenName;
    *symbol = sym;
    *total_supply = supply;
    (*balance)[owner.get()] = supply;
    prints(&owner.get()[0]);
    DIPC_EMIT_EVENT(Tranfer, "", &owner.get()[0], supply);
}

void TestToken::transfer(const char* to, uint64_t value){
    Address2 callerAddr = caller2();
    std::string callStr = callerAddr.toString();
    uint64_t originValue = balance.get()[callStr];

    prints(&callStr[0]);
    printui(originValue);
    printui(value);

    bool result = (originValue >= value);
    DipcAssert(result);

    std::string toStr = CharToAddress2Str(to);
    DipcAssert(balance.get()[toStr] + value >= balance.get()[toStr]);

    (*balance)[callStr] = balance.get()[callStr] -  value;
    (*balance)[toStr] = balance.get()[toStr] + value;
    DIPC_EMIT_EVENT(Tranfer, &(callStr[0]), to, value);
}
void TestToken::transferFrom(const char* from, const char* to, uint64_t value){
    Address2 callerAddr = caller2();
    std::string fromStr = CharToAddress2Str(from);
    std::string toStr = CharToAddress2Str(to);

    DipcAssert(balance.get()[fromStr] >= value);
    DipcAssert(balance.get()[toStr] + value >= balance.get()[toStr]);
    DipcAssert(allowance.get()[fromStr+callerAddr.toString()] >= value);

    (*balance)[toStr] = balance.get()[toStr] + value;
    (*balance)[fromStr] = balance.get()[fromStr] - value;
    (*allowance)[fromStr +callerAddr.toString()] = allowance.get()[fromStr+callerAddr.toString()] - value; 
    DIPC_EMIT_EVENT(Tranfer, from, to, value);
}
void TestToken::approve(const char* spender, uint64_t value){
    Address2 callerAddr = caller2();
    std::string spenderStr = CharToAddress2Str(spender);
   
    uint64_t total = allowance.get()[callerAddr.toString()+spenderStr] + value;
    (*allowance)[callerAddr.toString()+spenderStr] = total;
    DIPC_EMIT_EVENT(Approval, &(callerAddr.toString()[0]), spender, value);
}
void TestToken::burn(uint64_t value){
    Address2 callerAddr = caller2();
    DipcAssert(balance.get()[callerAddr.toString()] >= value);
    DipcAssert(balance.get()[owner.get()] + value >= balance.get()[owner.get()]);
    
    (*balance)[callerAddr.toString()] -= value;
    (*balance)[owner.get()] += value;
    DIPC_EMIT_EVENT(Tranfer, &(callerAddr.toString()[0]), &(owner.get()[0]), value);
}

uint64_t TestToken::getBalance(const char* own){
    prints(own);
    std::string ownerStr = CharToAddress2Str(own);
    uint64_t ba =  balance.get()[ownerStr];
    printui(ba);
    return ba;
}