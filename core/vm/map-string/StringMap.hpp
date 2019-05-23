#pragma once
#include <dipc/dipc.hpp>
using namespace dipc;

char temp[8] = "balance";
char tempkey[3] = "cc";
class StringMap : public Contract {
 public: 
  void init();
  Map<temp, std::string, int> balance;

  void setBalance(char *, int );
  CONSTANT int getBalance(char* );
};

// You must define ABI here.
DIPC_ABI(StringMap, init);
DIPC_ABI(StringMap, setBalance);
DIPC_ABI(StringMap, getBalance);
