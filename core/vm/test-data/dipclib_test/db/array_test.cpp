#include "array_test.hpp"

char arrayStrName[] = "arraystr";
char arrayIntName[] = "arrayint";
char arraySetName[] = "arrayset";

typedef dipc::db::Array<arrayIntName, int, 20> ArrayInt;
typedef dipc::db::Array<arrayStrName, std::string, 2> ArrayStr;
typedef dipc::db::Array<arraySetName, std::string, 10> ArraySet;

void arraytest::arrayTest() {
    print("\r\n dipc lib test arrayTest start\r\n");
    arrayIntTest();
    arrayStrTest();
    arraySetTest();
    print("\r\n dipc lib test arrayTest end\r\n");
}

void arrayIntTest() {
    {
        ArrayInt arrayInt;
        for (size_t i = 0; i < 20; i++) {
            arrayInt[i] = i;
        }
        for (size_t i = 0; i < 10; i++) {
            arrayInt[i] = 0;
        }
    }

    {
        ArrayInt arrayInt;
        for (size_t i = 0; i < 10; i++) {
            DipcAssert(arrayInt[i] == 0, "array[", i, "]", arrayInt[i]);
        }

        for (size_t i = 10; i < 20; i++) {
            DipcAssert(arrayInt[i] == i, "array[", i, "]", arrayInt[i]);
        }
    }

    {
        DEBUG("test iterator");
        ArrayInt arrayInt;
        ArrayInt::Iterator iter = arrayInt.begin();
        for (size_t i = 0; i < 10 && iter != arrayInt.end(); i++, iter++) {
            println("i:", i, "iter:", *iter);
            DipcAssert(*iter == 0);
        }

        for (size_t i = 10; i < 20 && iter != arrayInt.end(); i++, iter++) {
            println("i:", i, "iter:", *iter);
            DipcAssert(*iter == i, "iter:", *iter, "i:", i);
        }

        ArrayInt::ConstIterator citer = arrayInt.cbegin();
        for (size_t i = 0; i < 10 && citer != arrayInt.cend(); i++, citer++) {
            println("i:", i, "iter:", *citer);
            DipcAssert(*citer == 0);
        }

        for (size_t i = 10; i < 20 && citer != arrayInt.cend(); i++, citer++) {
            println("i:", i, "iter:", *citer);
            DipcAssert(*citer == i, "iter:", *citer, "i:", i);
        }
    }

    {
        DEBUG("test  reserve iterator");
        ArrayInt arrayInt;

        ArrayInt::ConstIterator citer = arrayInt.cbegin();
        for (size_t i = 0; i < 10 && citer != arrayInt.cend(); i++, citer++) {
            DipcAssert(*citer == 0);
        }

        for (size_t i = 10; i < 20 && citer != arrayInt.cend(); i++, citer++) {
            DipcAssert(*citer == i, "iter:", *citer, "i:", i);
        }

        ArrayInt::ReverseIterator iter = arrayInt.rbegin();
        for (size_t i = 19; i >= 10 && iter != arrayInt.rend(); i--, iter++) {
            DipcAssert(*iter == i, "iter:", *iter, "i:", i);
        }

        for (size_t i = 9; i >= 0 && iter != arrayInt.rend(); i--, iter++) {
            DipcAssert(*iter == 0, "iter:", *iter);
        }

        ArrayInt::ConstReverseIterator criter = arrayInt.crbegin();
        for (size_t i = 19; i >= 10 && criter != arrayInt.crend(); i--, ++criter) {
            DipcAssert(*criter == i, "criter:", *criter, "i:", i);
        }

        for (size_t i = 9; i >= 0 && criter != arrayInt.crend(); i--, criter++) {
            DipcAssert(*criter == 0, "criter:", *criter, "i:", i);
        }
    }
}

void arrayStrTest() {
    {
        ArrayStr arrayStr;
        arrayStr[0] = "hello";
        arrayStr[1] = "world";
    }

    {
        DEBUG("test reopen");
        ArrayStr arrayStr;

        DipcAssert(arrayStr[0] == "hello", "arrayStr[0]:", arrayStr[0]);
        DipcAssert(arrayStr[1] == "world", "arrayStr[1]:", arrayStr[1]);
    }

    {
        DEBUG("test del");
        ArrayStr arrayStr;
        arrayStr[0] = "";
    }

    {
        DEBUG("test reopen deled");
        ArrayStr arrayStr;
        DipcAssert(arrayStr[0] == "", "arrayStr[0]:", arrayStr[0]);
    }
}

void arraySetTest() {
    ArraySet array;
    array[0] = "hello";
    array.setConst(0, "helloworld");
    DipcAssert(array[0] == "helloworld");
    array.setConst(1, "world");
    DipcAssert(array[1] == "world");
}

