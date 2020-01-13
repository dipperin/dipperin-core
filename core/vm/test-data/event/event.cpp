#include "event.hpp"

CONSTANT void envEvent::init() {}

CONSTANT char* envEvent::returnString(char* name) {
    DIPC_EMIT_EVENT(topic, name);
    return name;
}

CONSTANT int64_t envEvent::returnInt(char* name) {
    DIPC_EMIT_EVENT(topic, name);
    int64_t num = 50;
    return num;
}

CONSTANT uint64_t envEvent::returnUint(char* name) {
    DIPC_EMIT_EVENT(topic, name);
    uint64_t num = 50;
    return num;
}