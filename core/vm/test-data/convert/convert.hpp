#pragma once
#include <dipc/dipc.hpp>
using namespace dipc;

class convert : public Contract {
 public: 
  void init();
  void toString();
  void getBlockInfo();
  void transferTest(char *to, uint64_t amount);
  void rlpTest();
  void printTest();
};

// You must define ABI here.
DIPC_ABI(convert, init);
DIPC_ABI(convert, toString);
DIPC_ABI(convert, getBlockInfo);
DIPC_ABI(convert, transferTest);
DIPC_ABI(convert, rlpTest);
DIPC_ABI(convert, printTest);
DIPC_EVENT(int64, const char*, int64_t);
DIPC_EVENT(uint64, const char*, uint64_t);
DIPC_EVENT(string, const char*, const char*);

bytes generateHashSrcData(uint64_t datalen);
void sha3Test();