#ifndef NEURON_MODULES_DEFAULT_PLUGIN_H
#define NEURON_MODULES_DEFAULT_PLUGIN_H

#include <stdlib.h>
#include <neuron.h>
#include <neuron/plugin.h>
#include <neuron/utils/log.h>

#ifndef NEU_PLUGIN_STRUCT_DEFINED
#define NEU_PLUGIN_STRUCT_DEFINED
struct neu_plugin {
    neu_plugin_common_t common;
};
#endif

static inline neu_plugin_t *default_driver_open(void)
{
    neu_plugin_t *plugin = calloc(1, sizeof(neu_plugin_t));
    if (plugin) {
        neu_plugin_common_init(&plugin->common);
    }
    return plugin;
}

static inline int default_driver_close(neu_plugin_t *plugin)
{
    if (plugin) {
        free(plugin);
    }
    return 0;
}

static inline int default_driver_load_tags(neu_plugin_t *plugin, const char *group,
                                           neu_datatag_t *tags, int n_tags)
{
    (void) plugin;
    (void) group;
    (void) tags;
    (void) n_tags;
    return 0;
}

static inline int default_driver_add_tags(neu_plugin_t *plugin, const char *group,
                                          neu_datatag_t *tags, int n_tags)
{
    (void) plugin;
    (void) group;
    (void) tags;
    (void) n_tags;
    return 0;
}

static inline int default_driver_del_tags(neu_plugin_t *plugin, int n_tags)
{
    (void) plugin;
    (void) n_tags;
    return 0;
}

static inline int default_driver_request(neu_plugin_t *plugin,
                                         neu_reqresp_head_t *head, void *data)
{
    (void) head;
    (void) data;
    zlog_warn(plugin->common.log, "driver request not implemented");
    return 0;
}

#endif
