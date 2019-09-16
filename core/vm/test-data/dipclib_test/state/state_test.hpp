#pragma once
#include <dipc/dipc.hpp>

using namespace dipc;

class statetest
{
public:
    void stateTests();
};

bytes generateHashSrcData(uint64_t datalen);
void sha3Test();