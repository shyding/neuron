#ifndef NEURON_MODULES_DEFAULT_CONNECT_H
#define NEURON_MODULES_DEFAULT_CONNECT_H

#include <stdint.h>
#include <neuron.h>
#include <neuron/plugin.h>
#include <neuron/utils/log.h>

static inline int default_send_msg(neu_plugin_t *plugin, uint16_t len, uint8_t *buf)
{
    (void) plugin;
    (void) len;
    (void) buf;
    return 0;
}

#endif
