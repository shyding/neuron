#include <stdlib.h>
#include <string.h>

#include <neuron.h>
#include <neuron/plugin.h>
#include <neuron/errcodes.h>
#include <neuron/utils/log.h>

#include "cctl_plugin.h"
#include "cctl_config.h"
#include "cctl_handle.h"
#include "cctl_tcp.h"

#define CCTL_MODULE_NAME "CCTL"
#define CCTL_DESCRIPTION "JetLinks CCTL IoT Central Control Northbound Plugin"
#define CCTL_DESCRIPTION_ZH "CCTL / JetLinks 物联网中控系统北向应用插件"

extern const neu_plugin_module_t neu_plugin_module;

static neu_plugin_t *plugin_open(void)
{
    neu_plugin_t *plugin = calloc(1, sizeof(neu_plugin_t));
    if (!plugin) {
        return NULL;
    }

    neu_plugin_common_init(&plugin->common);
    pthread_mutex_init(&plugin->mtx, NULL);
    plugin->started = false;

    // Set default config
    cctl_config_parse(plugin, NULL, &plugin->config);

    return plugin;
}

static int plugin_close(neu_plugin_t *plugin)
{
    if (!plugin) {
        return 0;
    }

    if (plugin->tcp) {
        cctl_tcp_free(plugin->tcp);
        plugin->tcp = NULL;
    }

    cctl_config_fini(&plugin->config);
    pthread_mutex_destroy(&plugin->mtx);
    free(plugin);

    return 0;
}

static int plugin_init(neu_plugin_t *plugin, bool load)
{
    (void) load;
    if (!plugin) {
        return NEU_ERR_EINTERNAL;
    }

    plugin->events = neu_event_new(plugin->common.name);
    if (!plugin->events) {
        plog_error(plugin, "neu_event_new fail");
        return NEU_ERR_EINTERNAL;
    }

    pthread_mutex_lock(&plugin->mtx);
    plugin->tcp = cctl_tcp_new(plugin, plugin->config.host, plugin->config.port);
    pthread_mutex_unlock(&plugin->mtx);

    plog_notice(plugin, "initialize plugin `%s` success", neu_plugin_module.module_name);
    return NEU_ERR_SUCCESS;
}

static int plugin_uninit(neu_plugin_t *plugin)
{
    if (!plugin) {
        return NEU_ERR_SUCCESS;
    }

    pthread_mutex_lock(&plugin->mtx);
    plugin->started = false;
    if (plugin->timer && plugin->events) {
        neu_event_del_timer(plugin->events, plugin->timer);
        plugin->timer = NULL;
    }
    if (plugin->tcp) {
        cctl_tcp_free(plugin->tcp);
        plugin->tcp = NULL;
    }
    pthread_mutex_unlock(&plugin->mtx);

    if (plugin->events) {
        neu_event_close(plugin->events);
        plugin->events = NULL;
    }

    plog_notice(plugin, "uninitialize plugin `%s` success", neu_plugin_module.module_name);
    return NEU_ERR_SUCCESS;
}

static int plugin_start(neu_plugin_t *plugin)
{
    if (!plugin) {
        return NEU_ERR_EINTERNAL;
    }

    pthread_mutex_lock(&plugin->mtx);
    plugin->started = true;

    // Connect to CCTL TCP gateway
    if (plugin->tcp) {
        cctl_tcp_connect(plugin->tcp);
    }

    // Register keepalive timer
    if (!plugin->timer && plugin->events) {
        int interval = plugin->config.keepalive > 0 ? plugin->config.keepalive : CCTL_DEFAULT_KEEPALIVE;
        neu_event_timer_param_t param = {
            .second      = interval,
            .millisecond = 0,
            .cb          = cctl_handle_timer,
            .usr_data    = plugin,
        };
        plugin->timer = neu_event_add_timer(plugin->events, param);
    }
    pthread_mutex_unlock(&plugin->mtx);

    plog_notice(plugin, "start plugin `%s` success, connecting to %s:%d",
                neu_plugin_module.module_name, plugin->config.host, plugin->config.port);
    return NEU_ERR_SUCCESS;
}

static int plugin_stop(neu_plugin_t *plugin)
{
    if (!plugin) {
        return NEU_ERR_SUCCESS;
    }

    pthread_mutex_lock(&plugin->mtx);
    plugin->started = false;

    if (plugin->timer && plugin->events) {
        neu_event_del_timer(plugin->events, plugin->timer);
        plugin->timer = NULL;
    }

    if (plugin->tcp) {
        cctl_tcp_close(plugin->tcp);
    }
    pthread_mutex_unlock(&plugin->mtx);

    plog_notice(plugin, "stop plugin `%s` success", neu_plugin_module.module_name);
    return NEU_ERR_SUCCESS;
}

static int plugin_config(neu_plugin_t *plugin, const char *setting)
{
    if (!plugin) {
        return NEU_ERR_EINTERNAL;
    }

    cctl_config_t new_config;
    if (cctl_config_parse(plugin, setting, &new_config) != 0) {
        plog_error(plugin, "config parse fail");
        return NEU_ERR_NODE_SETTING_INVALID;
    }

    pthread_mutex_lock(&plugin->mtx);
    plugin->config = new_config;

    if (plugin->tcp) {
        cctl_tcp_free(plugin->tcp);
        plugin->tcp = cctl_tcp_new(plugin, plugin->config.host, plugin->config.port);
        if (plugin->started) {
            cctl_tcp_connect(plugin->tcp);
        }
    }
    pthread_mutex_unlock(&plugin->mtx);

    plog_notice(plugin, "config plugin `%s` success: target %s:%d device_id: %s",
                neu_plugin_module.module_name, plugin->config.host, plugin->config.port, plugin->config.device_id);
    return NEU_ERR_SUCCESS;
}

static int plugin_request(neu_plugin_t *plugin, neu_reqresp_head_t *head, void *data)
{
    if (!plugin || !head) {
        return NEU_ERR_EINTERNAL;
    }

    int ret = 0;
    switch (head->type) {
    case NEU_REQRESP_TRANS_DATA:
        ret = cctl_handle_trans_data(plugin, (neu_reqresp_trans_data_t *) data);
        break;
    case NEU_REQ_SUBSCRIBE_GROUP:
    case NEU_REQ_UPDATE_SUBSCRIBE_GROUP:
        ret = cctl_handle_subscribe(plugin, (neu_req_subscribe_t *) data);
        break;
    case NEU_REQ_UNSUBSCRIBE_GROUP: {
        neu_req_unsubscribe_t *req = (neu_req_unsubscribe_t *) data;
        ret = cctl_handle_unsubscribe(plugin, req->driver, req->group);
        break;
    }
    case NEU_REQ_DEL_GROUP: {
        neu_req_del_group_t *req = (neu_req_del_group_t *) data;
        ret = cctl_handle_unsubscribe(plugin, req->driver, req->group);
        break;
    }
    case NEU_REQ_UPDATE_GROUP:
    case NEU_REQ_UPDATE_NODE:
    case NEU_REQRESP_NODE_DELETED:
        break;
    default:
        break;
    }

    return ret;
}

static const neu_plugin_intf_funs_t plugin_intf_funs = {
    .open        = plugin_open,
    .close       = plugin_close,
    .init        = plugin_init,
    .uninit      = plugin_uninit,
    .start       = plugin_start,
    .stop        = plugin_stop,
    .setting     = plugin_config,
    .request     = plugin_request,
    .try_connect = NULL,
};

const neu_plugin_module_t neu_plugin_module = {
    .version         = NEU_VERSION(2, 16, 0),
    .schema          = "cctl",
    .module_name     = CCTL_MODULE_NAME,
    .module_descr    = CCTL_DESCRIPTION,
    .module_descr_zh = CCTL_DESCRIPTION_ZH,
    .intf_funs       = &plugin_intf_funs,
    .kind            = NEU_PLUGIN_KIND_SYSTEM,
    .type            = NEU_NA_TYPE_APP,
    .display         = true,
    .single          = false,
};
