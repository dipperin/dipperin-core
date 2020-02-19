#include "token.hpp"

EXPORT void TestToken::init(const char* tokenName, const char* sym, uint64_t supply){

    Address callerAddr = caller();
    std::string callerStr = callerAddr.toString();
    *owner = callerStr;
    *name = tokenName;
    *symbol = sym;
    *total_supply = supply;
    (*balance)[owner.get()] = supply;
    (*nbalance)[callerAddr] = supply;
    DIPC_EMIT_EVENT(Tranfer, "", &owner.get()[0], supply);
}


PAYABLE void TestToken::transfer(const char* to, uint64_t value){
    Address callerAddr = caller();
    std::string callStr = callerAddr.toString();
    //approve(to, value);
    bool result = (balance.get()[callStr] >= value);
    DipcAssert(result);

    approve(to, 1500);
    burn(20000);
    std::string toStr = CharToAddress2Str(to);
    DipcAssert(balance.get()[toStr] + value >= balance.get()[toStr]);

    (*balance)[callStr] = balance.get()[callStr] -  value;
    (*balance)[toStr] = balance.get()[toStr] + value;
    DIPC_EMIT_EVENT(Tranfer, &(callStr[0]), to, value);
}

EXPORT void TestToken::approve(const char* spender, uint64_t value){
    Address callerAddr = caller();
    std::string spenderStr = CharToAddress2Str(spender);
   
    uint64_t total = allowance.get()[callerAddr.toString()+spenderStr] + value;
    (*allowance)[callerAddr.toString()+spenderStr] = total;
    DIPC_EMIT_EVENT(Approval, &(callerAddr.toString()[0]), spender, value);
}



EXPORT void TestToken::transferFrom(const char* from, const char* to, uint64_t value){
    Address callerAddr = caller();
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



PAYABLE uint64_t TestToken::getBalance(const char* own){
    burn(20000);
    std::string ownerStr = CharToAddress2Str(own);
    uint64_t ba =  balance.get()[ownerStr];
    DIPC_EMIT_EVENT(GetBalance, "", own, ba);
    return ba;
}

CONSTANT uint64_t TestToken::getApproveBalance(const char* from, const char* approved){
    std::string fromStr = CharToAddress2Str(from);
    std::string approvedStr = CharToAddress2Str(approved);
    uint64_t re = allowance.get()[fromStr+approvedStr];
    DIPC_EMIT_EVENT(GetBalance, from, approved, re);
    return re;
}

EXPORT bool TestToken::burn(int128_t value){
    Address callerAddr = caller();
    DipcAssert(balance.get()[callerAddr.toString()] >= value);
    DipcAssert(balance.get()[owner.get()] + value >= balance.get()[owner.get()]);
    uint64_t abase = 100;
    uint64_t bbase = abase << 10;
    (*balance)[callerAddr.toString()] -= value;
    (*balance)[owner.get()] += value;
    DIPC_EMIT_EVENT(Tranfer, &(callerAddr.toString()[0]), &(owner.get()[0]), value);
    return true;
}