#include <string.h>
#include <stdlib.h>
#include <jansson.h>

#include "cctl_config.h"
#include "cctl_plugin.h"
#include <neuron/utils/log.h>

int cctl_config_parse(neu_plugin_t *plugin, const char *setting, cctl_config_t *config)
{
    if (!config) {
        return -1;
    }

    memset(config, 0, sizeof(cctl_config_t));
    strncpy(config->host, CCTL_DEFAULT_HOST, sizeof(config->host) - 1);
    config->port = CCTL_DEFAULT_PORT;
    strncpy(config->device_id, CCTL_DEFAULT_DEVICE_ID, sizeof(config->device_id) - 1);
    config->keepalive = CCTL_DEFAULT_KEEPALIVE;

    if (!setting || strlen(setting) == 0) {
        return 0;
    }

    if (plugin) {
        plog_notice(plugin, "raw setting json: %s", setting);
    }

    json_error_t err;
    json_t *root = json_loads(setting, 0, &err);
    if (!root) {
        if (plugin) {
            plog_error(plugin, "json_loads fail: %s", err.text);
        }
        return -1;
    }

    json_t *p_obj = json_object_get(root, "params");
    json_t *params = (p_obj && json_is_object(p_obj)) ? p_obj : root;

    json_t *j_host = json_object_get(params, "host");
    if (j_host && json_is_string(j_host)) {
        strncpy(config->host, json_string_value(j_host), sizeof(config->host) - 1);
    }

    json_t *j_port = json_object_get(params, "port");
    if (j_port && json_is_integer(j_port)) {
        config->port = (int) json_integer_value(j_port);
    }

    json_t *j_dev = json_object_get(params, "device_id");
    if (j_dev && json_is_string(j_dev)) {
        strncpy(config->device_id, json_string_value(j_dev), sizeof(config->device_id) - 1);
    }

    json_t *j_ka = json_object_get(params, "keepalive");
    if (j_ka && json_is_integer(j_ka)) {
        config->keepalive = (int) json_integer_value(j_ka);
    }

    json_decref(root);
    return 0;
}

void cctl_config_fini(cctl_config_t *config)
{
    if (config) {
        memset(config, 0, sizeof(cctl_config_t));
    }
}
