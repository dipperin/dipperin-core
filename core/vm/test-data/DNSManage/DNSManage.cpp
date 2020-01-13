#include "DNSManage.hpp"

/**
 * @brief init method used to setup own addr
 */ 
PAYABLE void PasswordManage::init(){
    std::string ownerStr = caller().toString();
    *owner = ownerStr;
}
/**
 * @brief registerPassword method used to register a password
 * @param passwordName: password name
 */ 
EXPORT void PasswordManage::registerPassword(char* passwordName){
    std::string toAddr = caller().toString();
    if (passwordStore.get()[passwordName] != "") {
        DIPC_EMIT_EVENT(ErrEvent, "registerPassword  password exist", passwordName);
        return;
    }
    (*passwordStore)[passwordName] = toAddr;
    print("passwordName ---");
    print(passwordName);
    (*addrStore)[&toAddr[0]] = passwordName;
    DIPC_EMIT_EVENT(registerPasswordEvent, &toAddr[0], passwordName);
}

/**
 * @brief queryPasswordByAddr method used to query password of a addr
 * @param addr: account address
 */ 
CONSTANT char* PasswordManage::queryPasswordByAddr(char* _addr){
    // print("queryPasswordByAddr");
    // print(addrStore.get()[_addr]);
    // char* result = &addrStore.get()[_addr][0];
    // print("queryPasswordByAddr result");
    // print(result);
    std::map<std::string,std::string>::iterator iter;
    auto mapStore = passwordStore.get();
    iter = mapStore.begin();
    while(iter != mapStore.end()){
        if(iter->second == _addr){
            auto result = iter->first.c_str();
            char *buf = new char[strlen(result)+1];
　　         strcpy(buf, result);
            return buf;
        }
        iter++;
    }
   return "";
}

/**
 * @brief queryAddrByPassword method used to query addr of a password
 * @param _passwordName: password name of some account address
 */ 
CONSTANT char* PasswordManage::queryAddrByPassword(char* passwordName){
   return &(passwordStore.get()[passwordName][0]);
}