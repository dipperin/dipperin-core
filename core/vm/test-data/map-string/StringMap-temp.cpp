#include "StringMap.hpp"

void StringMap::init() {}

void StringMap::setBalance( char* key, int value){
	std::string strKey;
	strKey = key;
	(*balance)[strKey] = value;
};
CONSTANT int StringMap::getBalance(char* key){
	return balance.get()[key];
};
//dipc autogen begin
extern "C" { 
void init() {
StringMap StringMap_dipc;
StringMap_dipc.init();
}
void setBalance(char * key ,int data) {
StringMap StringMap_dipc;
StringMap_dipc.setBalance(key, data);
}
int getBalance(char * key) {
StringMap StringMap_dipc;
return StringMap_dipc.getBalance(key);
}

}
//dipc autogen end