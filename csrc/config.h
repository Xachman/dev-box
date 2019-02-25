#ifndef DEV_BOX_CONFIG
#define DEV_BOX_CONFIG


typedef struct Config {
    char namespacestr[32];
    char filepath[0];
} Config;


struct Config *NewConfig(struct Config *c, char namespacestr[], char filepath[]);

#endif