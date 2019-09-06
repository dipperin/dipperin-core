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
