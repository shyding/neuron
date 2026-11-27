#ifndef _CCTL_PLUGIN_H_
#define _CCTL_PLUGIN_H_

#include <pthread.h>
#include <neuron.h>
#include <neuron/plugin.h>
#include <neuron/event/event.h>

#include "cctl_config.h"
#include "cctl_tcp.h"

struct neu_plugin {
    neu_plugin_common_t common;
    neu_events_t *      events;
    pthread_mutex_t     mtx;
    cctl_config_t       config;
    cctl_tcp_client_t * tcp;
    neu_event_timer_t * timer;
    bool                started;
};

#endif
