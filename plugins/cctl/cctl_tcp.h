#ifndef _CCTL_TCP_H_
#define _CCTL_TCP_H_

#include <stdbool.h>
#include <stddef.h>
#include <neuron.h>
#include <neuron/plugin.h>

typedef struct cctl_tcp_client cctl_tcp_client_t;

cctl_tcp_client_t *cctl_tcp_new(neu_plugin_t *plugin, const char *host, int port);
int                cctl_tcp_connect(cctl_tcp_client_t *client);
int                cctl_tcp_send(cctl_tcp_client_t *client, const char *data, size_t len);
bool               cctl_tcp_is_connected(cctl_tcp_client_t *client);
void               cctl_tcp_close(cctl_tcp_client_t *client);
void               cctl_tcp_free(cctl_tcp_client_t *client);

#endif
