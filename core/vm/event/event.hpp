#pragma once
#include <platon/platon.hpp>
using namespace platon;

class envEvent : public Contract {
 public: 
  void init(char *tokenName, char *sym, uint64_t supply);

  void hello(const char* name, int64_t num);
};

// You must define ABI here.
PLATON_ABI(envEvent, init);
PLATON_ABI(envEvent, hello);
PLATON_EVENT(logName, const char*, int64_t);