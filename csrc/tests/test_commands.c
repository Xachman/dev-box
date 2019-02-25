#include "acutest.h"
#include "../config.h"

void test_command_strings(void)
{
    struct Config *config;

    config = NewConfig(config, "namespace", "/path/to/volume");

    printf("namespace: %s \n", config->namespacestr);
    printf("filepath: %s \n", config->filepath);

    TEST_CHECK_(0, "Expected %d, got %d", 3, 5);
}


TEST_LIST = {
    { "command-strings", test_command_strings },
    {NULL,NULL}
};