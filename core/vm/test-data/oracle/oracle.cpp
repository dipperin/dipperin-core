#include "oracle.hpp"

EXPORT void oracle::init() {}

EXPORT void oracle::setHeader(uint64_t key, char* value)
{
  (*chain)[key] = value;
};

CONSTANT char* oracle::getHeader(uint64_t key)
{
  std::string header = chain.get()[key];
  return &header[0];
};
