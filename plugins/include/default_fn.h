#ifndef NEURON_MODULES_DEFAULT_FN_H
#define NEURON_MODULES_DEFAULT_FN_H

#include <stdint.h>
#include <neuron.h>
#include <neuron/plugin.h>

typedef enum {
    PLUGIN_BIT_DEFAULT = 0,
} plugin_bit_e;

static inline int default_enforce_tag_count(neu_plugin_t *plugin, neu_plugin_group_t *group)
{
    (void) plugin;
    (void) group;
    return 0;
}

static inline int default_timer_license_check(neu_plugin_t *plugin, plugin_bit_e bit, neu_plugin_group_t *group)
{
    (void) plugin;
    (void) bit;
    (void) group;
    return 0;
}

static inline int default_update_error(neu_plugin_t *plugin, const char *group, int err)
{
    (void) plugin;
    (void) group;
    (void) err;
    return 0;
}

#endif
