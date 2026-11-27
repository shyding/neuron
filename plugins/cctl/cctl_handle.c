#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <jansson.h>

#include <neuron.h>
#include <neuron/plugin.h>
#include <neuron/utils/log.h>
#include <neuron/utils/time.h>
#include <neuron/json/neu_json_rw.h>

#include "cctl_plugin.h"
#include "cctl_handle.h"
#include "cctl_tcp.h"

// Internal device state simulation for pump-03
static char g_run_state[32] = "running";
static double g_target_pressure = 2.0;

int cctl_send_event(neu_plugin_t *plugin, const char *event_id, json_t *event_data)
{
    if (!plugin || !plugin->tcp || !event_id) {
        if (event_data) json_decref(event_data);
        return -1;
    }

    uint64_t now_ms = neu_time_ms();
    json_t *root = json_object();
    json_object_set_new(root, "deviceId", json_string(plugin->config.device_id));
    json_object_set_new(root, "event", json_string(event_id));
    json_object_set_new(root, "data", event_data ? event_data : json_object());
    json_object_set_new(root, "timestamp", json_integer((json_int_t) now_ms));

    char msg_id[64];
    snprintf(msg_id, sizeof(msg_id), "evt_%llu", (unsigned long long) now_ms);
    json_object_set_new(root, "messageId", json_string(msg_id));

    char *dump = json_dumps(root, JSON_COMPACT);
    json_decref(root);
    if (!dump) {
        return -1;
    }

    size_t len = strlen(dump);
    char *buf = malloc(len + 2);
    if (!buf) {
        free(dump);
        return -1;
    }

    memcpy(buf, dump, len);
    buf[len] = '\n';
    buf[len + 1] = '\0';
    free(dump);

    plog_notice(plugin, "CCTL reported EVENT '%s' for device %s", event_id, plugin->config.device_id);
    int ret = cctl_tcp_send(plugin->tcp, buf, len + 1);
    free(buf);
    return ret;
}

int cctl_handle_downlink_message(neu_plugin_t *plugin, const char *msg_str)
{
    if (!plugin || !msg_str || strlen(msg_str) == 0) {
        return -1;
    }

    json_error_t err;
    json_t *root = json_loads(msg_str, 0, &err);
    if (!root) {
        plog_error(plugin, "Downlink json_loads error: %s (msg: %s)", err.text, msg_str);
        return -1;
    }

    json_t *j_type = json_object_get(root, "messageType");
    const char *type = j_type && json_is_string(j_type) ? json_string_value(j_type) : "";

    json_t *j_msg_id = json_object_get(root, "messageId");
    const char *msg_id = j_msg_id && json_is_string(j_msg_id) ? json_string_value(j_msg_id) : "";

    json_t *j_dev_id = json_object_get(root, "deviceId");
    const char *dev_id = j_dev_id && json_is_string(j_dev_id) ? json_string_value(j_dev_id) : plugin->config.device_id;

    plog_notice(plugin, "CCTL received downlink message type: '%s', messageId: '%s', deviceId: '%s'",
                type, msg_id, dev_id);

    if (strcmp(type, "INVOKE_FUNCTION") == 0) {
        json_t *j_func = json_object_get(root, "functionId");
        const char *func_id = j_func && json_is_string(j_func) ? json_string_value(j_func) : "";
        json_t *inputs = json_object_get(root, "inputs");

        plog_notice(plugin, "Executing function '%s' on device '%s'...", func_id, dev_id);

        json_t *output = json_object();
        bool success = true;

        if (strcmp(func_id, "startPump") == 0) {
            strncpy(g_run_state, "running", sizeof(g_run_state) - 1);
            json_object_set_new(output, "success", json_true());
            json_object_set_new(output, "runState", json_string("running"));
            json_object_set_new(output, "message", json_string("Water pump started successfully"));
        } else if (strcmp(func_id, "stopPump") == 0) {
            strncpy(g_run_state, "stopped", sizeof(g_run_state) - 1);
            json_object_set_new(output, "success", json_true());
            json_object_set_new(output, "runState", json_string("stopped"));
            json_object_set_new(output, "message", json_string("Water pump stopped successfully"));
        } else if (strcmp(func_id, "setPressureTarget") == 0) {
            json_t *j_target = inputs ? json_object_get(inputs, "targetPressure") : NULL;
            if (j_target && (json_is_number(j_target))) {
                g_target_pressure = json_number_value(j_target);
            }
            json_object_set_new(output, "success", json_true());
            json_object_set_new(output, "targetPressure", json_real(g_target_pressure));
            json_object_set_new(output, "message", json_string("Target pressure updated"));
        } else if (strcmp(func_id, "resetAlarm") == 0) {
            json_object_set_new(output, "success", json_true());
            json_object_set_new(output, "message", json_string("Device alarms cleared"));
        } else {
            json_object_set_new(output, "success", json_true());
            json_object_set_new(output, "message", json_string("Generic function executed"));
        }

        // Construct INVOKE_FUNCTION_REPLY
        json_t *reply = json_object();
        json_object_set_new(reply, "messageType", json_string("INVOKE_FUNCTION_REPLY"));
        json_object_set_new(reply, "deviceId", json_string(dev_id));
        json_object_set_new(reply, "messageId", json_string(msg_id));
        json_object_set_new(reply, "output", json_true());
        json_object_set_new(reply, "success", json_boolean(success));
        json_decref(output);

        char *dump = json_dumps(reply, JSON_COMPACT);
        json_decref(reply);

        if (dump) {
            size_t len = strlen(dump);
            char *buf = malloc(len + 2);
            if (buf) {
                memcpy(buf, dump, len);
                buf[len] = '\n';
                buf[len + 1] = '\0';
                cctl_tcp_send(plugin->tcp, buf, len + 1);
                free(buf);
            }
            free(dump);
        }
        plog_notice(plugin, "Function '%s' reply sent for messageId: %s", func_id, msg_id);
    } else if (strcmp(type, "WRITE_PROPERTY") == 0) {
        json_t *props = json_object_get(root, "properties");
        plog_notice(plugin, "Handling WRITE_PROPERTY on device '%s'...", dev_id);

        json_t *reply = json_object();
        json_object_set_new(reply, "messageType", json_string("WRITE_PROPERTY_REPLY"));
        json_object_set_new(reply, "deviceId", json_string(dev_id));
        json_object_set_new(reply, "messageId", json_string(msg_id));
        if (props) {
            json_incref(props);
            json_object_set_new(reply, "properties", props);
        } else {
            json_object_set_new(reply, "properties", json_object());
        }
        json_object_set_new(reply, "success", json_true());

        char *dump = json_dumps(reply, JSON_COMPACT);
        json_decref(reply);

        if (dump) {
            size_t len = strlen(dump);
            char *buf = malloc(len + 2);
            if (buf) {
                memcpy(buf, dump, len);
                buf[len] = '\n';
                buf[len + 1] = '\0';
                cctl_tcp_send(plugin->tcp, buf, len + 1);
                free(buf);
            }
            free(dump);
        }
    } else if (strcmp(type, "READ_PROPERTY") == 0) {
        json_t *reply = json_object();
        json_object_set_new(reply, "messageType", json_string("READ_PROPERTY_REPLY"));
        json_object_set_new(reply, "deviceId", json_string(dev_id));
        json_object_set_new(reply, "messageId", json_string(msg_id));
        json_t *props = json_object();
        json_object_set_new(props, "runState", json_string(g_run_state));
        json_object_set_new(props, "pressure", json_real(g_target_pressure));
        json_object_set_new(reply, "properties", props);
        json_object_set_new(reply, "success", json_true());

        char *dump = json_dumps(reply, JSON_COMPACT);
        json_decref(reply);

        if (dump) {
            size_t len = strlen(dump);
            char *buf = malloc(len + 2);
            if (buf) {
                memcpy(buf, dump, len);
                buf[len] = '\n';
                buf[len + 1] = '\0';
                cctl_tcp_send(plugin->tcp, buf, len + 1);
                free(buf);
            }
            free(dump);
        }
    }

    json_decref(root);
    return 0;
}

int cctl_handle_trans_data(neu_plugin_t *plugin, neu_reqresp_trans_data_t *trans_data)
{
    if (!plugin || !trans_data || !plugin->tcp) {
        return -1;
    }

    const char *device_id = plugin->config.device_id;
    if (!device_id || strlen(device_id) == 0) {
        device_id = trans_data->group;
    }

    json_t *props = json_object();
    if (!props) {
        return -1;
    }

    double current_pressure = 0.0;
    double current_vibration = 0.0;
    bool has_pressure = false;
    bool has_vibration = false;

    if (trans_data->tags) {
        utarray_foreach(trans_data->tags, neu_resp_tag_value_meta_t *, tag_value) {
            neu_json_read_resp_tag_t json_tag = { 0 };
            neu_tag_value_to_json(tag_value, &json_tag);

            if (json_tag.error == 0 && json_tag.name && strlen(json_tag.name) > 0) {
                switch (json_tag.t) {
                case NEU_JSON_INT:
                    json_object_set_new(props, json_tag.name, json_integer((json_int_t) json_tag.value.val_int));
                    if (strcmp(json_tag.name, "pressure") == 0) {
                        current_pressure = (double) json_tag.value.val_int;
                        has_pressure = true;
                    } else if (strcmp(json_tag.name, "vibration") == 0) {
                        current_vibration = (double) json_tag.value.val_int;
                        has_vibration = true;
                    }
                    break;
                case NEU_JSON_DOUBLE:
                case NEU_JSON_FLOAT:
                    json_object_set_new(props, json_tag.name, json_real(json_tag.value.val_double));
                    if (strcmp(json_tag.name, "pressure") == 0) {
                        current_pressure = json_tag.value.val_double;
                        has_pressure = true;
                    } else if (strcmp(json_tag.name, "vibration") == 0) {
                        current_vibration = json_tag.value.val_double;
                        has_vibration = true;
                    }
                    break;
                case NEU_JSON_STR:
                    json_object_set_new(props, json_tag.name, json_string(json_tag.value.val_str ? json_tag.value.val_str : ""));
                    break;
                case NEU_JSON_BOOL:
                    json_object_set_new(props, json_tag.name, json_boolean(json_tag.value.val_bool));
                    break;
                default:
                    break;
                }
            }

            if (json_tag.metas) {
                free(json_tag.metas);
            }
        }
    }

    // Always include runState in properties if not present
    if (!json_object_get(props, "runState")) {
        json_object_set_new(props, "runState", json_string(g_run_state));
    }

    uint64_t now_ms = neu_time_ms();
    json_t *root = json_object();
    json_object_set_new(root, "deviceId", json_string(device_id));
    json_object_set_new(root, "properties", props);
    json_object_set_new(root, "timestamp", json_integer((json_int_t) now_ms));

    char msg_id[64];
    snprintf(msg_id, sizeof(msg_id), "neu_%llu", (unsigned long long) now_ms);
    json_object_set_new(root, "messageId", json_string(msg_id));

    char *dump = json_dumps(root, JSON_COMPACT);
    json_decref(root);

    if (!dump) {
        return -1;
    }

    size_t len = strlen(dump);
    char *buf = malloc(len + 2);
    if (!buf) {
        free(dump);
        return -1;
    }

    memcpy(buf, dump, len);
    buf[len] = '\n';
    buf[len + 1] = '\0';
    free(dump);

    int ret = cctl_tcp_send(plugin->tcp, buf, len + 1);
    free(buf);

    if (ret == 0) {
        NEU_PLUGIN_UPDATE_METRIC(plugin, NEU_METRIC_TRANS_DATA_5S, 1, NULL);
        NEU_PLUGIN_UPDATE_METRIC(plugin, NEU_METRIC_TRANS_DATA_30S, 1, NULL);
        NEU_PLUGIN_UPDATE_METRIC(plugin, NEU_METRIC_TRANS_DATA_60S, 1, NULL);
    }

    // Automated Event Detection & Reporting:
    // If pressure > 2.5 (or >= 150 from integer raw registers), fire overpressure_alarm event!
    static uint64_t last_alarm_ms = 0;
    if (now_ms - last_alarm_ms > 10000) { // Throttle events to once every 10s
        if ((has_pressure && current_pressure > 2.0) || current_pressure >= 100.0) {
            json_t *ev_data = json_object();
            json_object_set_new(ev_data, "pressure", json_real(current_pressure > 10.0 ? current_pressure / 100.0 : current_pressure));
            json_object_set_new(ev_data, "threshold", json_real(2.0));
            json_object_set_new(ev_data, "level", json_string("URGENT"));
            cctl_send_event(plugin, "overpressure_alarm", ev_data);
            last_alarm_ms = now_ms;
        }
    }

    return ret;
}

int cctl_handle_subscribe(neu_plugin_t *plugin, neu_req_subscribe_t *sub_info)
{
    if (!sub_info) {
        return 0;
    }
    if (sub_info->params) {
        free(sub_info->params);
    }
    if (sub_info->static_tags) {
        free(sub_info->static_tags);
    }
    plog_notice(plugin, "CCTL subscribed group: %s from driver: %s", sub_info->group, sub_info->driver);
    return 0;
}

int cctl_handle_unsubscribe(neu_plugin_t *plugin, const char *driver, const char *group)
{
    plog_notice(plugin, "CCTL unsubscribed group: %s from driver: %s", group ? group : "*", driver ? driver : "*");
    return 0;
}

int cctl_handle_timer(void *usr_data)
{
    neu_plugin_t *plugin = usr_data;
    if (!plugin || !plugin->started) {
        return 0;
    }

    pthread_mutex_lock(&plugin->mtx);
    if (!plugin->tcp) {
        pthread_mutex_unlock(&plugin->mtx);
        return 0;
    }

    if (!cctl_tcp_is_connected(plugin->tcp)) {
        plog_debug(plugin, "CCTL attempting reconnect to %s:%d", plugin->config.host, plugin->config.port);
        cctl_tcp_connect(plugin->tcp);
    } else {
        // Send heartbeat ping to CCTL
        uint64_t now_ms = neu_time_ms();
        json_t *root = json_object();
        json_object_set_new(root, "deviceId", json_string(plugin->config.device_id));
        json_t *props = json_object();
        json_object_set_new(props, "runState", json_string(g_run_state));
        json_object_set_new(root, "properties", props);
        json_object_set_new(root, "timestamp", json_integer((json_int_t) now_ms));

        char msg_id[64];
        snprintf(msg_id, sizeof(msg_id), "hb_%llu", (unsigned long long) now_ms);
        json_object_set_new(root, "messageId", json_string(msg_id));

        char *dump = json_dumps(root, JSON_COMPACT);
        json_decref(root);

        if (dump) {
            size_t len = strlen(dump);
            char *buf = malloc(len + 2);
            if (buf) {
                memcpy(buf, dump, len);
                buf[len] = '\n';
                buf[len + 1] = '\0';
                cctl_tcp_send(plugin->tcp, buf, len + 1);
                free(buf);
            }
            free(dump);
        }
    }
    pthread_mutex_unlock(&plugin->mtx);
    return 0;
}
