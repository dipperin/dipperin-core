### Data Types

#### Storage Types

Dipc provides template types to provide data persistence to the dipperin chain.

- Uint8
- Int8
- Uint16
- Int16
- Uint
- Int
- Uint64
- Int64
- String
- Vector
- Set
- Map
- Array
- Tuple
- Deque

template types to provide data persistence to the dipperin chain.

Storage types usage example:

```c++
// Example one Map uses:
// Define the storage field name
char bal[] = "balance";
// Storage field name   key type    value type
Map<     bal,           std::string, uint64_t >  balance;

// Example two String uses:
// Define the storage field name
char name[] = "contract_name";
//  Storage field name   
String<name> contract_name;

```
In the contract, the field defined by the storage type is used and its value is automatically stored on the dipperin chain when the contract is created.

#### Fundamental Types

Dipc supports all basic types of C++, standard library types and their arithmetic operations
And the types defined in the dipclib package：

- Big integer types defined using the boost library
  - bigint
  - u64
  - u128
  - u256
  - u160
  - u512
- Integer and unsigned integers encoded using VLQ 
  - unsigned_int
  - signed_int
- Custom types for efficient use of memory
  - map
  - array
  - list
- Types defined by the custom FixedHash class
  - h256  //32bytes
  - h160
  - h128
  - h64
  - Address
