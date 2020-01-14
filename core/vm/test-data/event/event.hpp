#pragma once
#include <dipc/dipc.hpp>
using namespace dipc;

class envEvent : public Contract {
 public: 
  CONSTANT void init();
  CONSTANT char* returnString(char* name);
  CONSTANT int64_t returnInt(char* name);
  CONSTANT uint64_t returnUint(char* name);
};

DIPC_EVENT(topic, const char*);