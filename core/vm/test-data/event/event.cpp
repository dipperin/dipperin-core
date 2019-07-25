#include "event.hpp"

void envEvent::init() {}

CONSTANT const char* envEvent::returnString(const char* name) {
    DIPC_EMIT_EVENT(topic, name);
    return name;
}

CONSTANT int64_t envEvent::returnInt(const char* name) {
    DIPC_EMIT_EVENT(topic, name);
    int64_t num = 50;
    return num;
}

CONSTANT uint64_t envEvent::returnUint(const char* name) {
    DIPC_EMIT_EVENT(topic, name);
    uint64_t num = 50;
    return num;
}