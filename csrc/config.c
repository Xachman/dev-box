#include <string.h>
#include <stdlib.h>
#include "config.h"


struct Config *NewConfig(struct Config *c, char namespacestr[32], char filepath[]) {
    c = malloc( sizeof(*c) + sizeof(char) * strlen(filepath) );

    strcpy(c->namespacestr, namespacestr);
    strcpy(c->filepath, filepath);

    return c;
}