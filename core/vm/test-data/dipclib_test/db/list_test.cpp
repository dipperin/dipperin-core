#include "list_test.hpp"

char listStrName[] = "liststr";
char listIntName[] = "listint";
char listPushName[] = "listPush";
char listInsertName[] = "listInsert";

typedef dipc::db::List<listIntName, int> ListInt;
typedef dipc::db::List<listStrName, std::string> ListStr;
typedef dipc::db::List<listPushName, std::string> ListPush;
typedef dipc::db::List<listInsertName, std::string> ListInsert;

void listtest::listTest() {
    print("\r\n dipc lib test listTest start\r\n");
    listPushTest();
    listIntTest();
    listStrTest();
    listInsertTest();
    print("\r\n dipc lib test listTest end\r\n");
}

void listPushTest() {
    {
        ListPush pl;
        for (int i = 0; i < 500; i++) {
            pl.push("hello");
        }

        for (int i = 0; i < 400; i++) {
            pl.del(0);
        }
    }

    {
        ListPush pl;
        for (int i = 0; i < 100; i++) {
            DipcAssert(pl[i] == "hello");
        }
    }
}

void listIntTest() {
    {
        ListInt listInt;
        for (size_t i = 0; i < 20; i++) {
            listInt.push(i);
        }
        for (size_t i = 0; i < 10; i++) {
            listInt.del(i);
        }
    }

    {
        ListInt listInt;
        for (size_t i = 0; i < 10; i++) {
            DipcAssert(listInt[i] == i + 10);
        }
    }
}

void listStrTest() {
    {
        ListStr listStr;
        listStr.push("hello");
        listStr.push("world");
    }

    {
        DEBUG("test reopen");
        ListStr listStr;
        DipcAssert(listStr[0] == "hello");
        DipcAssert(listStr[1] == "world");
        println("listStr size:", listStr.size());
    }

    {
        DEBUG("test del");
        ListStr listStr;
        println("listStr size:", listStr.size());
        listStr.del(0);
        println("listStr size:", listStr.size());
        DipcAssert(listStr.size() == 1);
    }

    {
        DEBUG("test reopen deled");
        ListStr listStr;
        DipcAssert(listStr.size() == 1);
        DipcAssert(listStr[0] == "world", "listStr[0] ", listStr[0]);
    }
    {
        DEBUG("test reopen iterator");
        ListStr listStr;
        listStr.push("world");
        DipcAssert(listStr.size() == 2);

        size_t size = listStr.size();
        println("list size:", size);
        size_t count = 0;
        for (ListStr::Iterator iter = listStr.begin(); iter != listStr.end();
                iter++) {
            println("++++++++++");
            DipcAssert(*iter == "world");
            count++;
            println("--------");
        }
        DipcAssert(count == size);
        count = 0;
        for (ListStr::ReverseIterator iter = listStr.rbegin();
                iter != listStr.rend(); iter++) {
            DipcAssert(*iter == "world");
            count++;
        }
        DipcAssert(count == size);

        count = 0;
        for (ListStr::ConstIterator iter = listStr.cbegin(); iter != listStr.cend();
                iter++) {
            DipcAssert(*iter == "world");
            count++;
        }

        DipcAssert(count == size);
        count = 0;
        for (ListStr::ConstReverseIterator iter = listStr.crbegin();
                iter != listStr.crend(); iter++) {
            DipcAssert(*iter == "world");
            count++;
        }
        DipcAssert(count == size);
    }
}

void listInsertTest() {
    ListInsert listInsert;
    listInsert.push("hello");
    listInsert.setConst(0, "helloworld");
    DipcAssert(listInsert[0] == "helloworld");
}