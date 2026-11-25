#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <dlfcn.h>
#include <glob.h>
#include <string.h>
#include <neuron/plugin.h>

void *neuron = NULL;
int64_t global_timestamp = 0;

void dump_one(const char *path) {
    void *h = dlopen(path, RTLD_LAZY | RTLD_GLOBAL);
    if (!h) {
        printf("  \"%s\": {\"error\": \"%s\"},\n", path, dlerror());
        return;
    }
    neu_plugin_module_t *m = (neu_plugin_module_t *) dlsym(h, "neu_plugin_module");
    if (!m) {
        printf("  \"%s\": {\"error\": \"no neu_plugin_module\"},\n", path);
        dlclose(h);
        return;
    }
    printf("  \"%s\": {\n", path);
    printf("    \"schema\": \"%s\",\n", m->schema ? m->schema : "");
    printf("    \"module_name\": \"%s\",\n", m->module_name ? m->module_name : "");
    printf("    \"module_descr\": \"%s\",\n", m->module_descr ? m->module_descr : "");
    printf("    \"module_descr_zh\": \"%s\",\n", m->module_descr_zh ? m->module_descr_zh : "");
    printf("    \"type\": %d,\n", m->type);
    printf("    \"kind\": %d,\n", m->kind);
    printf("    \"display\": %s,\n", m->display ? "true" : "false");
    printf("    \"single\": %s,\n", m->single ? "true" : "false");
    if (m->intf_funs) {
        printf("    \"funs\": {\n");
        printf("      \"open\": %s,\n", m->intf_funs->open ? "true" : "false");
        printf("      \"close\": %s,\n", m->intf_funs->close ? "true" : "false");
        printf("      \"init\": %s,\n", m->intf_funs->init ? "true" : "false");
        printf("      \"uninit\": %s,\n", m->intf_funs->uninit ? "true" : "false");
        printf("      \"start\": %s,\n", m->intf_funs->start ? "true" : "false");
        printf("      \"stop\": %s,\n", m->intf_funs->stop ? "true" : "false");
        printf("      \"setting\": %s,\n", m->intf_funs->setting ? "true" : "false");
        printf("      \"request\": %s,\n", m->intf_funs->request ? "true" : "false");
        printf("      \"try_connect\": %s,\n", m->intf_funs->try_connect ? "true" : "false");
        printf("      \"validate_tag\": %s,\n", m->intf_funs->driver.validate_tag ? "true" : "false");
        printf("      \"group_timer\": %s,\n", m->intf_funs->driver.group_timer ? "true" : "false");
        printf("      \"write_tag\": %s,\n", m->intf_funs->driver.write_tag ? "true" : "false");
        printf("      \"write_tags\": %s,\n", m->intf_funs->driver.write_tags ? "true" : "false");
        printf("      \"add_tags\": %s,\n", m->intf_funs->driver.add_tags ? "true" : "false");
        printf("      \"load_tags\": %s,\n", m->intf_funs->driver.load_tags ? "true" : "false");
        printf("      \"del_tags\": %s\n", m->intf_funs->driver.del_tags ? "true" : "false");
        printf("    }\n");
    }
    printf("  },\n");
}

int main(int argc, char **argv) {
    if (argc < 2) return 1;
    printf("{\n");
    for (int i = 1; i < argc; i++) {
        dump_one(argv[i]);
    }
    printf("  \"_end\": {}\n}\n");
    return 0;
}
