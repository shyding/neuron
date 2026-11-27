#ifndef _CCTL_HANDLE_H_
#define _CCTL_HANDLE_H_

#include <neuron.h>
#include <neuron/plugin.h>
#include <jansson.h>

int  cctl_handle_trans_data(neu_plugin_t *plugin, neu_reqresp_trans_data_t *trans_data);
int  cctl_handle_subscribe(neu_plugin_t *plugin, neu_req_subscribe_t *sub_info);
int  cctl_handle_unsubscribe(neu_plugin_t *plugin, const char *driver, const char *group);
int  cctl_handle_timer(void *usr_data);

int  cctl_handle_downlink_message(neu_plugin_t *plugin, const char *msg_str);
int  cctl_send_event(neu_plugin_t *plugin, const char *event_id, json_t *event_data);

#endif
