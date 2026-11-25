/**
 * NEURON Industrial Protocol Plugin Implementation
 * Module: Inovance Modbus TCP
 * Schema: inovance-modbus-tcp
 **/

#include <stdlib.h>
#include <string.h>
#include <neuron.h>
#include <neuron/plugin.h>
#include <neuron/errcodes.h>

#include "default_plugin.h"
#include "default_fn.h"
#include "default_connect.h"

static neu_plugin_t *driver_open(void)
{
    return default_driver_open();
}

static int driver_close(neu_plugin_t *plugin)
{
    return default_driver_close(plugin);
}

static int driver_init(neu_plugin_t *plugin, bool load)
{
    (void) plugin;
    (void) load;
    return 0;
}

static int driver_uninit(neu_plugin_t *plugin)
{
    (void) plugin;
    return 0;
}

static int driver_start(neu_plugin_t *plugin)
{
    (void) plugin;
    return 0;
}

static int driver_stop(neu_plugin_t *plugin)
{
    (void) plugin;
    return 0;
}

static int driver_config(neu_plugin_t *plugin, const char *config)
{
    (void) plugin;
    (void) config;
    return 0;
}

static int driver_request(neu_plugin_t *plugin, neu_reqresp_head_t *head,
                          void *data)
{
    return default_driver_request(plugin, head, data);
}

static int driver_connect(neu_plugin_t *plugin)
{
    (void) plugin;
    return 0;
}

static int driver_validate_tag(neu_plugin_t *plugin, neu_datatag_t *tag)
{
    (void) plugin;
    if (!tag || !tag->name || !tag->address) {
        return NEU_ERR_NODE_NAME_EMPTY;
    }
    return 0;
}

static int driver_group_timer(neu_plugin_t *plugin, neu_plugin_group_t *group)
{
    if (!plugin || !group) {
        return 0;
    }
    default_enforce_tag_count(plugin, group);
    return 0;
}

static int driver_write(neu_plugin_t *plugin, void *req, neu_datatag_t *tag,
                        neu_value_u value)
{
    (void) plugin;
    (void) req;
    (void) tag;
    (void) value;
    return 0;
}

static int driver_write_tags(neu_plugin_t *plugin, void *req, UT_array *tags)
{
    (void) plugin;
    (void) req;
    (void) tags;
    return 0;
}

static int driver_test_read_tag(neu_plugin_t *plugin, void *req,
                                neu_datatag_t tag)
{
    (void) plugin;
    (void) req;
    (void) tag;
    return 0;
}

static const neu_plugin_intf_funs_t plugin_intf_funs = {
    .open        = driver_open,
    .close       = driver_close,
    .init        = driver_init,
    .uninit      = driver_uninit,
    .start       = driver_start,
    .stop        = driver_stop,
    .setting     = driver_config,
    .request     = driver_request,
    .try_connect = driver_connect,

    .driver.validate_tag  = driver_validate_tag,
    .driver.group_timer   = driver_group_timer,
    .driver.group_sync    = driver_group_timer,
    .driver.write_tag     = driver_write,
    .driver.tag_validator = NULL,
    .driver.write_tags    = driver_write_tags,
    .driver.test_read_tag = driver_test_read_tag,
    .driver.add_tags      = default_driver_add_tags,
    .driver.load_tags     = default_driver_load_tags,
    .driver.del_tags      = default_driver_del_tags,
};

const neu_plugin_module_t neu_plugin_module = {
    .version         = NEURON_PLUGIN_VER_1_0,
    .schema          = "inovance-modbus-tcp",
    .module_name     = "Inovance Modbus TCP",
    .module_descr    = "This plugin is used to connect Inovance PLCs using the Modbus TCP protocol. ",
    .module_descr_zh = "该插件用于连接使用 modbus tcp 协议的汇川 plc。",
    .intf_funs       = &plugin_intf_funs,
    .kind            = NEU_PLUGIN_KIND_SYSTEM,
    .type            = NEU_NA_TYPE_DRIVER,
    .display         = true,
    .single          = false,
};
