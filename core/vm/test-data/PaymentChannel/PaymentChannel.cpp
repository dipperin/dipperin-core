#include "PaymentChannel.hpp"

PAYABLE void PaymentChannel::init(char* _recipient, uint64_t duration, uint64_t _balance){
//PAYABLE void PaymentChannel::init(std::string _recipient, uint64_t duration, uint64_t _balance){
    *sender = caller();
    *recipient = _recipient;
    print("_recipient");
    print(_recipient);
    print("recipient");
    prints_l(&recipient.get()[0], recipient.get().size());
    *expiration = duration;
    *balance = _balance;
    DipcAssert(dipc::callValue() == _balance);
}

EXPORT void PaymentChannel::close(uint64_t amount, char* signature){
    DipcAssert(!closed.get());
    Address callerAddr = caller();
    DipcAssert(callerAddr.toString() == recipient.get());
    std::string sign = signature;
    // joint contractAddr, amount, toAddress, use sha3()  encrypt
    Address contractAddr = dipc::address();
    std::string data = contractAddr.toString() + std::to_string(amount) + callerAddr.toString();
    h256 sha3Data = dipc::sha3(data);
    DipcAssert(dipc::getSignerAddress(sha3Data,sign).toString() == sender.get().toString());
    callTransfer(callerAddr, amount);
    DipcAssert(balance.get() -amount > 0);
    callTransfer(sender.get(), balance.get() - amount);
    *closed = true;
}

EXPORT void PaymentChannel::extend(uint64_t newExpiration) {
    DipcAssert(!closed.get());
    DipcAssert(caller().toString() == sender.get().toString());
    DipcAssert(newExpiration > expiration.get());
    *expiration = newExpiration;
}

EXPORT void PaymentChannel::claimTimeout(){
    DipcAssert(!closed.get());
    DipcAssert(dipc::timestamp() > expiration.get());
    callTransfer(sender.get(), balance.get());
    *closed = true;
}