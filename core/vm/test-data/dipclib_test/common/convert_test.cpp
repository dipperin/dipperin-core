#include "convert_test.hpp"

void converttest::convertTest() {
    print("\r\n dipc lib test convertTest start111\r\n");
    convertStr();
    operatorTest();
    printTest();
    print("\r\n dipc lib test convertTest success\r\n");
}

void convertStr() {
    u64 i = 1234567890;
    DipcAssert(i.convert_to<std::string>() == "1234567890");

    u128 i128 = 1234567890;
    DipcAssert(i128.convert_to<std::string>() == "1234567890");

    u160 i160 = 1234567890;
    DipcAssert(i160.convert_to<std::string>() == "1234567890");

    u256 i256 = 1234567890;
    DipcAssert(i256.convert_to<std::string>() == "1234567890");

    u512 i512 = 1234567890;
    DipcAssert(i512.convert_to<std::string>() == "1234567890");

    std::string hex = "0x00ff";
    bytes bs = fromHex(hex);
    u64 ii = fromBigEndian<u64, bytes>(bs);
    DipcAssert(ii == 255);

    bs.resize(8);
    toBigEndian(ii, bs);
    hex = toHex(bs.begin(), bs.end(), "0x");
    DipcAssert(hex == "0x00000000000000ff");
}

void operatorTest() {
    u256 i = 100;
    DipcAssert(i - 50 == 50);
    DipcAssert(i + 10 == 110);
    DipcAssert(i / 2 == 50);
    DipcAssert(i * 2 == 200);
}

void printTest() {
    print("printTest called");
    int128_t large_int(-87654323456);
    printi128(&large_int);

    uint128_t large_uint(87654323456);
    printui128(&large_uint);
}