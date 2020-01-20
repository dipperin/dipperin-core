#include "dipclib_test.hpp"

PAYABLE void dipcLibTest::init() {}

EXPORT void dipcLibTest::libTest()
{
    //test common function 
    arithmeticTest();
    convertTest();
    fixedHashTest();
    rlpTest();

    // // test db function
    arrayTest();
    listTest();
    mapTest();

    // // test state function 
    stateTests();  
    TestTokens();

    // // test storage
    StorageTests();
}