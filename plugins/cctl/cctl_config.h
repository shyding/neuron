#ifndef _CCTL_CONFIG_H_
#define _CCTL_CONFIG_H_

#include <neuron.h>
#include <neuron/plugin.h>

#define CCTL_DEFAULT_HOST "127.0.0.1"
#define CCTL_DEFAULT_PORT 18884
#define CCTL_DEFAULT_DEVICE_ID "pump-03"
#define CCTL_DEFAULT_KEEPALIVE 30

typedef struct {
    char host[256];
    int  port;
    char device_id[128];
    int  keepalive;
} cctl_config_t;

int  cctl_config_parse(neu_plugin_t *plugin, const char *setting, cctl_config_t *config);
void cctl_config_fini(cctl_config_t *config);

#endif
