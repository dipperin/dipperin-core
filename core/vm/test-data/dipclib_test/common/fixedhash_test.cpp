#include "fixedhash_test.hpp"

void fixedhashtest::fixedHashTest() {
    print("\r\n dipc lib test fixedHashTest start\r\n");
    compareTest();
    xorTest();
    orTest();
    andTest();
    insertTest();
    containTest();
    print("\r\n dipc lib test fixedHashTest success\r\n");
}

void compareTest() {
    DEBUG("compareTest");
    FixedHash<4> h1(::dipc::sha3("abcd"));
    FixedHash<4> h2(::dipc::sha3("abcd"));
    FixedHash<4> h3(::dipc::sha3("aadd"));
    FixedHash<4> h4(0xBAADF00D);
    FixedHash<4> h5(0xAAAAAAAA);
    FixedHash<4> h6(0xBAADF00D);

    DipcAssertEQ(h1, h2);
    DipcAssertNE(h2, h3);

    DipcAssert(h4 > h5);
    DipcAssert(h5 < h4);
    DipcAssert(h6 <= h4);
    DipcAssert(h6 >= h4);
}

void xorTest() {
    FixedHash<2> h1("0xAAAA");
    FixedHash<2> h2("0xBBBB");

    DipcAssertEQ((h1 ^ h2), FixedHash<2>("0x1111"));
    h1 ^= h2;
    DipcAssertEQ(h1, FixedHash<2>("0x1111"));
}

void orTest() {
    FixedHash<4> h1("0xD3ADB33F");
    FixedHash<4> h2("0xBAADF00D");
    FixedHash<4> res("0xFBADF33F");

    DipcAssertEQ((h1 | h2), res);
    h1 |= h2;
    DipcAssertEQ(h1, res);
}

void andTest() {
    FixedHash<4> h1("0xD3ADB33F");
    FixedHash<4> h2("0xBAADF00D");
    FixedHash<4> h3("0x92aDB00D");

    DipcAssertEQ((h1 & h2), h3);
    h1 &= h2;
    DipcAssertEQ(h1, h3);
}

void insertTest() {
    FixedHash<4> h1("0xD3ADB33F");
    FixedHash<4> h2("0x2C524CC0");

    DipcAssertEQ(~h1, h2);
}

void containTest() {
    FixedHash<4> h1("0xD3ADB331");
    FixedHash<4> h2("0x0000B331");
    FixedHash<4> h3("0x0000000C");

    DipcAssert(h1.contains(h2));
    DipcAssert(!h1.contains(h3));
}

