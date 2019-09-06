#pragma once
#include <dipc/dipc.hpp>
using namespace dipc;

class envEvent : public Contract {
 public: 
  void init();
  CONSTANT const char* returnString(const char* name);
  CONSTANT int64_t returnInt(const char* name);
  CONSTANT uint64_t returnUint(const char* name);
};

// You must define ABI here.
DIPC_ABI(envEvent, init);
DIPC_ABI(envEvent, returnString);
DIPC_ABI(envEvent, returnInt);
DIPC_ABI(envEvent, returnUint);
DIPC_EVENT(topic, const char*);