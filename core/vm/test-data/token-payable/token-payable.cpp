#include "token-payable.hpp"

EXPORT void TestToken::init(char* tokenName, char* sym, uint64_t supply){
    std::string originStr = origin2().toString();
    *owner = originStr;
    *name = tokenName;
    *symbol = sym;
    *total_supply = supply;
    (*balance)[owner.get()] = supply;
    //prints(&owner.get()[0]);
    DIPC_EMIT_EVENT(Tranfer, "", &owner.get()[0], supply);
}

PAYABLE bool TestToken::transfer(const char* to, uint64_t value){
    isZeroAddress(to);
    Address2 callerAddr = caller2();
    std::string callStr = callerAddr.toString();
    std::string toStr = CharToAddress2Str(to);

    // check uint64 underflow and over flow
    DipcAssert(balance.get()[callStr] >= value);
    DipcAssert(balance.get()[toStr] + value > balance.get()[toStr]);

    (*balance)[callStr] -=  value;
    (*balance)[toStr] += value;
    DIPC_EMIT_EVENT(Tranfer, &(callStr[0]), to, value);
    return true;
}

EXPORT bool TestToken::withdraw(){
    isOwner();
    Address2 callerAddr = caller2();
    Address2 contractAddr = address2();

    // check contract balance and withdraw
    u256 originBalance = dipc::balance(contractAddr);
    DipcAssert(originBalance > 0);
    DipcAssert(callTransfer2(callerAddr,  originBalance) == 0);
    return true;
}

EXPORT bool TestToken::approve(const char* spender, uint64_t value){
    isZeroAddress(spender);
    std::string callerStr= caller2().toString();
    std::string spenderStr = CharToAddress2Str(spender);
   
    // check uint64 underflow and over flow
    DipcAssert(balance.get()[callerStr] >= value);
    DipcAssert(allowance.get()[callerStr+spenderStr] + value > allowance.get()[callerStr+spenderStr]);
    
    (*allowance)[callerStr+spenderStr] += value;
    DIPC_EMIT_EVENT(Approval, &(callerStr[0]), spender, value);
    return true;
}

EXPORT bool TestToken::transferFrom(const char* from, const char* to, uint64_t value){
    isZeroAddress(to);
    std::string callerStr= caller2().toString();
    std::string fromStr = CharToAddress2Str(from);
    std::string toStr = CharToAddress2Str(to);

    // check uint64 underflow and over flow
    DipcAssert(balance.get()[fromStr] >= value);
    DipcAssert(balance.get()[toStr] + value > balance.get()[toStr]);
    DipcAssert(allowance.get()[fromStr+callerStr] >= value);

    (*balance)[fromStr] -= value;
    (*balance)[toStr] += value;
    (*allowance)[fromStr +callerStr] -= value; 
    DIPC_EMIT_EVENT(Tranfer, from, to, value);
    return true;
}

EXPORT bool TestToken::burn(uint64_t value){
    std::string callerStr= caller2().toString();

    // check uint64 underflow and over flow
    DipcAssert(owner.get() != callerStr);
    DipcAssert(balance.get()[callerStr] >= value);
    DipcAssert(balance.get()[owner.get()] + value > balance.get()[owner.get()]);

    (*balance)[callerStr] -= value;
    (*balance)[owner.get()] += value;
    DIPC_EMIT_EVENT(Tranfer, &(callerStr[0]), &(owner.get()[0]), value);
    return true;
}

CONSTANT uint64_t TestToken::getBalance(const char* own){
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