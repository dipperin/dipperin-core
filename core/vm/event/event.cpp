#include "event.hpp"

void envEvent::init() {}

void envEvent::hello(const char* name, int64_t num) {
  // println("hello", name);
  PLATON_EMIT_EVENT(logName, name, num);
}