#pragma once
#include <dipc/dipc.hpp>
using namespace dipc;

char temp[7] = "header";
class oracle : public Contract
{
public:
    EXPORT void init();
    EXPORT void setHeader(uint64_t key, char* value);
    CONSTANT char* getHeader(uint64_t key);

    Map<temp, uint64_t, std::string> chain;
};
