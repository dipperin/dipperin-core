#pragma once
#include <dipc/dipc.hpp>
#include "common/buildins_test.hpp"
#include "common/convert_test.hpp"
#include "common/fixedhash_test.hpp"
#include "common/rlp_test.hpp"
#include "db/array_test.hpp"
#include "db/list_test.hpp"
#include "db/map_test.hpp"
#include "state/state_test.hpp"
#include "state/test_token.hpp"
#include "storage/storagetest.hpp"

using namespace dipc;

class dipcLibTest:
public buildins,
public converttest,
public fixedhashtest,
public rlptest,
public arraytest,
public listtest,
public maptest,
public statetest,
public testtoken,
public storagetest
{
public:
    void init();
    void libTest();
};

DIPC_ABI(dipcLibTest, init);
DIPC_ABI(dipcLibTest, libTest);
