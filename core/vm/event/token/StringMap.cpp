#include "StringMap.hpp"

void StringMap::init() {}

void StringMap::setBalance( char* key, int value){
	std::string strKey;
	strKey = key;
	(*balance)[strKey] = value;
	(*balance)[strKey+"2"] = value;
};
CONSTANT int StringMap::getBalance(char* key){
	std::string strKey;
	strKey = key;
	return balance.get()[strKey];
};
