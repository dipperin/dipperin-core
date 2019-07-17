#include "storagetest.hpp"

char testInt[] = "testInt";
char testString[] = "testString";
char testVector[] = "testVector";
char testSet[] = "testSet";
char testMap[] = "testMap";

char testArray[] = "testArray";

char testTuple[] = "testTuple";
char testDeque[] = "testDeque";


void storagetest::StorageTests()
{
    //print_f("\r\n dipc lib test StorageTests start\r\n");
    testIntStorage();
    testStringStorage();
    testVectorStorage();
    testSetStorage();
    testMapStorage();
    testArrayStorage();
    testTupleStorage();
    testDequeStorage();
    //print_f("\r\n dipc lib test StorageTests success\r\n");
}

void storagetest::testIntStorage()
{
    intValue = 12;
    DipcAssertEQ((intValue + 1), 13);
    DipcAssertEQ((intValue - 1), 11);
    DipcAssertEQ((intValue * 2), 24);
    DipcAssertEQ((intValue / 2), 6);
    DipcAssertEQ((intValue % 2), 0);
    DipcAssertEQ((intValue++), 13);
    DipcAssertEQ((intValue >> 1), 6);
    DipcAssertEQ((intValue << 1), 24);
    DipcAssert((intValue < 13));
    DipcAssert((intValue > 11));
    return;
}

void storagetest::testStringStorage()
{
    stringValue = "hello";
    DipcAssert((stringValue < "i"));
    DipcAssert((stringValue > "a"));

    std::string tmpValue;
    tmpValue = *stringValue;
    DipcAssertEQ(tmpValue.length(), 5);
    tmpValue.insert(0, "hello");
    DipcAssertEQ(stringValue, "hellohello");
    int index;
    index = tmpValue.find("he", 0);
    DipcAssertEQ(index, 0);
}

void storagetest::testVectorStorage()
{
    std::vector<int> tmpVector{1,2,3,4,5,6,7,8,9,10};
    vectorValue = tmpVector;
    DipcAssertEQ(vectorValue,tmpVector);
    
    std::vector<int> tmpData;
    tmpData = *vectorValue;
    DipcAssertEQ(tmpData,tmpVector);
    
    tmpData.insert(tmpData.begin()+1,6);
    std::vector<int> changeVector{1,6,2,3,4,5,6,7,8,9,10};
    DipcAssertEQ(tmpVector,changeVector);

    tmpData.clear();
    tmpData = *vectorValue;
    DipcAssert(tmpData.empty());
}

void storagetest::testSetStorage()
{
    std::set<int> tmpSet{6,2,3,1,4,5,};
    setValue = tmpSet;
    DipcAssertEQ(setValue,tmpSet);
    std::set<int> tmpValue;
    tmpValue = *setValue;
    DipcAssertEQ(tmpValue,tmpSet);
    
    tmpValue.insert(1);
    DipcAssertEQ(setValue,tmpSet);
}

void storagetest::testMapStorage()
{
    std:: map<std::string,int> tmpMap;
    std:: pair<std::string,int> p1("key0",0);
    std:: pair<std::string,int> p2("key1",1);
    std:: pair<std::string,int> p3("key2",2);
    tmpMap.insert(p1);
    tmpMap.insert(p2);
    tmpMap.insert(p3);
    mapValue = tmpMap;
    DipcAssertEQ(mapValue,tmpMap);

    std:: map<std::string,int> tmpValue;
    tmpValue = *mapValue;
    tmpValue["key1"] = 3;

    tmpMap["key1"] = 3;
    DipcAssertEQ(mapValue,tmpMap);
}

void storagetest::testArrayStorage()
{
    std::array<int,10> tmpArray{1,2,3,4,5,6,7,8,9,10};
    arrayValue = tmpArray;
    DipcAssertEQ(arrayValue,tmpArray);
    
    std::array<int,10> tmpValue;
    tmpValue = *arrayValue;
    tmpValue[0] = 11;
    std::array<int,10> changeArray{11,2,3,4,5,6,7,8,9,10};
    DipcAssertEQ(arrayValue,changeArray);
}

void storagetest::testTupleStorage()
{
    std::tuple<int,int,std::string> tmpTuple{0,1,"testTuple"};
    tupleValue = tmpTuple;
    DipcAssertEQ(tupleValue,tmpTuple);

    std::tuple<int,int,std::string> tmpValue;
    tmpValue = *tupleValue;
    
    std::string p3 = std::get<2>(tmpValue);
    DipcAssertEQ(p3,"testTuple");
    std::get<0>(tmpValue) = 3;

    return;
}

void storagetest::testDequeStorage()
{
    std::deque<int> tmpDeque {1,2,3,4,5,6};
    dequeValue = tmpDeque;
    DipcAssertEQ(dequeValue,tmpDeque);

    std::deque<int> tmpValue;
    tmpValue = *dequeValue;

    tmpValue.insert(tmpValue.begin(),0);
    std::deque<int> changeDeque {0,1,2,3,4,5,6};
    DipcAssertEQ(dequeValue,changeDeque);

    return;
}