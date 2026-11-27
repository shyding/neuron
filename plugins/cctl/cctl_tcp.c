#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <fcntl.h>
#include <errno.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <netinet/tcp.h>
#include <arpa/inet.h>
#include <netdb.h>
#include <pthread.h>

#include "cctl_tcp.h"
#include "cctl_plugin.h"
#include "cctl_handle.h"
#include <neuron/utils/log.h>

struct cctl_tcp_client {
    neu_plugin_t *   plugin;
    char             host[256];
    int              port;
    int              fd;
    bool             connected;
    pthread_mutex_t  mtx;
    pthread_t        recv_tid;
    bool             recv_running;
    bool             thread_created;
};

static void *cctl_tcp_recv_loop(void *arg)
{
    cctl_tcp_client_t *client = (cctl_tcp_client_t *) arg;
    char buf[2048];
    char line[4096];
    size_t line_len = 0;

    while (client->recv_running) {
        int fd = -1;
        pthread_mutex_lock(&client->mtx);
        fd = client->fd;
        pthread_mutex_unlock(&client->mtx);

        if (fd < 0) {
            usleep(100000);
            continue;
        }

        ssize_t n = recv(fd, buf, sizeof(buf) - 1, 0);
        if (n <= 0) {
            if (n < 0 && (errno == EAGAIN || errno == EWOULDBLOCK || errno == EINTR)) {
                continue;
            }
            // Socket closed or error
            pthread_mutex_lock(&client->mtx);
            if (client->fd >= 0) {
                close(client->fd);
                client->fd = -1;
            }
            client->connected = false;
            if (client->plugin) {
                client->plugin->common.link_state = NEU_NODE_LINK_STATE_DISCONNECTED;
            }
            pthread_mutex_unlock(&client->mtx);
            break;
        }

        for (ssize_t i = 0; i < n; i++) {
            if (buf[i] == '\n' || buf[i] == '\r') {
                line[line_len] = '\0';
                if (line_len > 0) {
                    cctl_handle_downlink_message(client->plugin, line);
                }
                line_len = 0;
            } else if (line_len < sizeof(line) - 1) {
                line[line_len++] = buf[i];
            }
        }
    }
    return NULL;
}

cctl_tcp_client_t *cctl_tcp_new(neu_plugin_t *plugin, const char *host, int port)
{
    cctl_tcp_client_t *client = calloc(1, sizeof(cctl_tcp_client_t));
    if (!client) {
        return NULL;
    }
    client->plugin = plugin;
    strncpy(client->host, host ? host : "127.0.0.1", sizeof(client->host) - 1);
    client->port = port > 0 ? port : 18884;
    client->fd = -1;
    client->connected = false;
    client->recv_running = false;
    client->thread_created = false;
    pthread_mutex_init(&client->mtx, NULL);
    return client;
}

int cctl_tcp_connect(cctl_tcp_client_t *client)
{
    if (!client) {
        return -1;
    }

    pthread_mutex_lock(&client->mtx);
    if (client->connected && client->fd >= 0) {
        pthread_mutex_unlock(&client->mtx);
        return 0;
    }

    if (client->fd >= 0) {
        close(client->fd);
        client->fd = -1;
    }
    client->connected = false;

    int fd = socket(AF_INET, SOCK_STREAM, 0);
    if (fd < 0) {
        pthread_mutex_unlock(&client->mtx);
        return -1;
    }

    // Set TCP_NODELAY
    int flag = 1;
    setsockopt(fd, IPPROTO_TCP, TCP_NODELAY, (char *) &flag, sizeof(int));

    // Resolve address
    struct sockaddr_in serv_addr;
    memset(&serv_addr, 0, sizeof(serv_addr));
    serv_addr.sin_family = AF_INET;
    serv_addr.sin_port = htons(client->port);

    if (inet_pton(AF_INET, client->host, &serv_addr.sin_addr) <= 0) {
        struct hostent *he = gethostbyname(client->host);
        if (!he || !he->h_addr_list[0]) {
            close(fd);
            pthread_mutex_unlock(&client->mtx);
            return -1;
        }
        memcpy(&serv_addr.sin_addr, he->h_addr_list[0], he->h_length);
    }

    // Set socket timeout for connect / send
    struct timeval tv;
    tv.tv_sec = 3;
    tv.tv_usec = 0;
    setsockopt(fd, SOL_SOCKET, SO_SNDTIMEO, (char *) &tv, sizeof(tv));

    if (connect(fd, (struct sockaddr *) &serv_addr, sizeof(serv_addr)) < 0) {
        close(fd);
        pthread_mutex_unlock(&client->mtx);
        return -1;
    }

    client->fd = fd;
    client->connected = true;
    client->recv_running = true;
    if (client->plugin) {
        client->plugin->common.link_state = NEU_NODE_LINK_STATE_CONNECTED;
    }

    if (!client->thread_created) {
        client->thread_created = true;
        pthread_create(&client->recv_tid, NULL, cctl_tcp_recv_loop, client);
    }

    pthread_mutex_unlock(&client->mtx);
    return 0;
}

int cctl_tcp_send(cctl_tcp_client_t *client, const char *data, size_t len)
{
    if (!client || !data || len == 0) {
        return -1;
    }

    pthread_mutex_lock(&client->mtx);
    if (!client->connected || client->fd < 0) {
        pthread_mutex_unlock(&client->mtx);
        // Try reconnect once
        if (cctl_tcp_connect(client) != 0) {
            return -1;
        }
        pthread_mutex_lock(&client->mtx);
    }

    size_t total_sent = 0;
    while (total_sent < len) {
        ssize_t n = send(client->fd, data + total_sent, len - total_sent, MSG_NOSIGNAL);
        if (n <= 0) {
            close(client->fd);
            client->fd = -1;
            client->connected = false;
            if (client->plugin) {
                client->plugin->common.link_state = NEU_NODE_LINK_STATE_DISCONNECTED;
            }
            pthread_mutex_unlock(&client->mtx);
            return -1;
        }
        total_sent += n;
    }

    pthread_mutex_unlock(&client->mtx);
    return 0;
}

bool cctl_tcp_is_connected(cctl_tcp_client_t *client)
{
    if (!client) {
        return false;
    }
    pthread_mutex_lock(&client->mtx);
    bool c = client->connected;
    pthread_mutex_unlock(&client->mtx);
    return c;
}

void cctl_tcp_close(cctl_tcp_client_t *client)
{
    if (!client) {
        return;
    }
    pthread_mutex_lock(&client->mtx);
    client->recv_running = false;
    if (client->fd >= 0) {
        shutdown(client->fd, SHUT_RDWR);
        close(client->fd);
        client->fd = -1;
    }
    client->connected = false;
    if (client->plugin) {
        client->plugin->common.link_state = NEU_NODE_LINK_STATE_DISCONNECTED;
    }
    pthread_mutex_unlock(&client->mtx);

    if (client->thread_created) {
        pthread_join(client->recv_tid, NULL);
        client->thread_created = false;
    }
}

void cctl_tcp_free(cctl_tcp_client_t *client)
{
    if (!client) {
        return;
    }
    cctl_tcp_close(client);
    pthread_mutex_destroy(&client->mtx);
    free(client);
}
