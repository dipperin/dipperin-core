#include "token.hpp"

void TestToken::init(char* tokenName, char* sym, uint64_t supply){
    //isOwner();
    name = tokenName;
    symbol = sym;
    total_supply = supply;
    (*balance)[owner.toString()] = supply;
    DIPC_EMIT_EVENT(Tranfer, "", &(owner.toString()[0]), supply);
}


void TestToken::transfer(const char* to, uint64_t value){
    Address2 callerAddr = caller2();
    DipcAssert(balance.get()[callerAddr.toString()] >= value);
    DipcAssert(balance.get()[to] + value >= balance.get()[to]);
    (*balance)[callerAddr.toString()] = balance.get()[callerAddr.toString()] -  value;
    (*balance)[to] = balance.get()[to] + value;
    DIPC_EMIT_EVENT(Tranfer, &(callerAddr.toString()[0]), to, value);
}
void TestToken::transferFrom(const char* from, const char* to, uint64_t value){
    Address2 callerAddr = caller2();
    DipcAssert(balance.get()[from] >= value);
    DipcAssert(balance.get()[to] + value >= balance.get()[to]);
    //DipcAssert(allowance.get()[from].get()[callerAddr.toString()] >= value);
    //std::string fromStr = from;
    DipcAssert(balance.get()[from+callerAddr.toString()] >= value);
    (*balance)[to] = balance.get()[to] + value;
    (*balance)[from] = balance.get()[from] - value;
    //allowance.get()[from].get()[callerAddr.toString()] -= value;
    (*balance)[from+callerAddr.toString()] = balance.get()[from+callerAddr.toString()] - value; 
    DIPC_EMIT_EVENT(Tranfer, from, to, value);
}
void TestToken::approve(const char* spender, uint64_t value){
    Address2 callerAddr = caller2();
   // std::string spenderStr = spender;
    //DipcAssert(value == 0 || allowance.get()[callerAddr.toString()].get()[spender] == 0);
    //(*(*allowance)[callerAddr.toString()])[spender] = value;
    DipcAssert(value == 0 || balance.get()[callerAddr.toString()+spender] == 0);
    (*balance)[callerAddr.toString()+spender] = value;
    DIPC_EMIT_EVENT(Approval, &(callerAddr.toString()[0]), spender, value);
}
void TestToken::burn(uint64_t value){
    Address2 callerAddr = caller2();
    DipcAssert(balance.get()[callerAddr.toString()] >= value);
    (*balance)[callerAddr.toString()] -= value;
    (*balance)[owner.toString()] += value;
    DIPC_EMIT_EVENT(Tranfer, &(callerAddr.toString()[0]), &(owner.toString()[0]), value);
}