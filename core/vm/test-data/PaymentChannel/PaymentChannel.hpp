#include <dipc/dipc.hpp>
using namespace dipc;

char senderc[] = "sender";
char recipientc[] = "recipient";
char expirationc[] = "expiration";
char balancec[] = "bal";
char closec[] = "close";

class PaymentChannel : public Contract {
private: 
   //String<senderc> sender;
   AddressStore<senderc> sender;
   String<recipientc> recipient;
   Uint64<balancec> balance; 
   Uint64<expirationc> expiration;
   Bool<closec> closed;
public:
    PAYABLE void init(char* _recipient, uint64_t duration, uint64_t _balance);
    //PAYABLE void init(std::string _recipient, uint64_t duration, uint64_t _balance);
    
    EXPORT void close(uint64_t amount, char* signature);
    EXPORT void extend(uint64_t newExpiration);
    EXPORT void claimTimeout();
};