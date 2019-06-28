#pragma once
#include <dipc/dipc.hpp>
using namespace dipc;

class buildins : public Contract {
 public: 
  void init();
  void arithmeticTest();
};

// You must define ABI here.
DIPC_ABI(buildins, init);
DIPC_ABI(buildins, arithmeticTest);

