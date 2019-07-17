#pragma once
#include <dipc/dipc.hpp>

using namespace dipc;

extern char testInt[];
extern char testString[];
extern char testVector[];
extern char testSet[];
extern char testMap[];
extern char testArray[];
extern char testTuple[];
extern char testDeque[];

class storagetest
{
public:
    void StorageTests();

private:
    Int<testInt> intValue;
    String<testString> stringValue;
    Vector<testVector, int> vectorValue;
    Set<testSet, int> setValue;
    Map<testMap, std::string, int> mapValue;
    Array<testArray, int, 10> arrayValue;
    Tuple<testTuple, int, int,std::string> tupleValue;
    Deque<testDeque, int> dequeValue;

    void testIntStorage();
    void testStringStorage();
    void testVectorStorage();
    void testSetStorage();
    void testMapStorage();
    void testArrayStorage();
    void testTupleStorage();
    void testDequeStorage();
};
