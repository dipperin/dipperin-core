#include "map_test.hpp"

char mapStrName[] = "mapstr";
typedef dipc::db::Map<mapStrName, std::string, std::string> MapStr;

void maptest::mapTest() {
    print("\r\n dipc lib test mapTest start\r\n");
    {
        DEBUG("test insertConst");
        MapStr map;
        map.insertConst("hello", "helloworld");
        DipcAssert(map["hello"] == "helloworld");

        map.insert("hello", "world");
        DipcAssert(map["hello"] == "world", map["hello"].c_str());
    }

    {
        DEBUG("test reopen insertConst");
        MapStr map;
        DipcAssert(map["hello"] == "world", map["hello"].c_str());
    }

    {
        DEBUG("test insert");
        MapStr map;
        map["hello1"] = "world";
        DipcAssert(map.size() == 2);
    }

    {
        DEBUG("test del");
        MapStr map;
        map.del("hello");
        DipcAssert(map.size() == 1);
    }

    {
        DEBUG("test reopen del");
        MapStr map;
        DipcAssert(map.size() == 1);
    }

    {
        DEBUG("test iterator");
        MapStr map;
        size_t size = map.size();
        println("map size:", size);
        size_t count = 0;
        for (MapStr::Iterator iter = map.begin(); iter != map.end(); iter++) {
            println("iter second:", iter->second());
            DipcAssert(iter->second() == "world");
            count++;
        }
        println("count:", count, "size:", size);
        DipcAssert(count == size);
        count = 0;
        for (MapStr::ReverseIterator iter = map.rbegin(); iter != map.rend();
                ++iter) {
            println("iter....");
            println("iter second:", iter->second());
            DipcAssert(iter->second() == "world", "reverse iterator error");
            println("iter second:", iter->second());
            count++;
        }
        println("count:", count, "size:", size);
        DipcAssert(count == size);

        count = 0;
        for (MapStr::ConstIterator citer = map.cbegin(); citer != map.cend();
                citer++) {
            println("iter second:", citer->second());
            DipcAssert(citer->second() == "world");
            count++;
        }
        println("count :", count, "size:", size);
        DipcAssert(count == size);
        count = 0;
        for (MapStr::ConstReverseIterator criter = map.crbegin();
                criter != map.crend(); ++criter) {
            println("iter....");
            println("iter second:", criter->second());
            DipcAssert(criter->second() == "world", "reverse iterator error");
            println("iter second:", criter->second());
            count++;
        }
        println("count:", count, "size:", size);
        DipcAssert(count == size);
    }
    print("\r\n dipc lib test mapTest success\r\n");
}