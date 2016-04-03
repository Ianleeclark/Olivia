#include <stdlib.h>

uint32_t jenkins_hash(const char *key, const size_t len) {

    uint32_t i = 0;
    uint32_t hash = 0;

    for(hash = i = 0; i < len; ++i) {
        hash += key[i];
        hash += (hash << 10);
        hash ^= (hash >> 6);
    }

    hash += (hash << 3);
    hash ^= (hash >> 11);
    hash += (hash << 15);

    return hash;
}
