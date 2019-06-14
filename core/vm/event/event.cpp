#include "event.hpp"

void envEvent::init(char *tokenName, char *sym, uint64_t supply) {

}

void envEvent::hello(const char* name, int64_t num) {
  // println("hello", name);
  PLATON_EMIT_EVENT(logName, name, num);
}