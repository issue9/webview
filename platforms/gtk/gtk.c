// SPDX-License-Identifier: MIT

#include "_cgo_export.h"
#include "gtk.h"

WebKitUserContentManager* _userContentManager(GtkWidget* wv) {
    return webkit_web_view_get_user_content_manager(WEBKIT_WEB_VIEW(wv));
}

void _add_script(GtkWidget* wv, const char* js) {
    WebKitUserContentManager* m = _userContentManager(wv);
    WebKitUserScript* script = webkit_user_script_new(js, WEBKIT_USER_CONTENT_INJECT_TOP_FRAME, WEBKIT_USER_SCRIPT_INJECT_AT_DOCUMENT_START, NULL, NULL);
    webkit_user_content_manager_add_script(m, script);
}

void _script_message_received(WebKitUserContentManager* m, WebKitJavascriptResult* rslt, gpointer user_data) {
    JSCValue *value = webkit_javascript_result_get_js_value(rslt);
    char* js = jsc_value_to_string(value);
    messageCallback(js);
}

gboolean _dispatch_cb(gpointer data) {
    dispatchCallback();
    return G_SOURCE_REMOVE;
}

void _move(GtkWidget* win, int x, int y) {
    gtk_window_move(GTK_WINDOW(win), x, y);
}

void _set_title(GtkWidget* win, const char* title) {
    gtk_window_set_title(GTK_WINDOW(win), title);
}

void _set_size(GtkWidget* win, int w, int h) {
    GdkGeometry g;
    g.base_width = w;
    g.base_height = h;
    gtk_window_set_geometry_hints(GTK_WINDOW(win), NULL, &g, GDK_HINT_BASE_SIZE);
}

App* create_gtk(bool debug, int x, int y, int w, int h, const char* title) {
    WebKitSettings* settings = webkit_settings_new();
    webkit_settings_set_enable_developer_extras(settings, debug);
    webkit_settings_set_javascript_can_access_clipboard(settings, true);

    GtkWidget* wv = webkit_web_view_new_with_settings(settings);

    WebKitUserContentManager* m = _userContentManager(wv);
    g_signal_connect(m, "script-message-received::external", G_CALLBACK(_script_message_received), NULL);
    webkit_user_content_manager_register_script_message_handler(m, "external");
    _add_script(wv, "window.external={invoke:function(s){window.webkit.messageHandlers.external.postMessage(s);}}");

    GtkWidget* win = gtk_window_new(GTK_WINDOW_TOPLEVEL);
    _move(win, x, y);
    _set_size(win, w, h);
    _set_title(win, title);
    gtk_container_add(GTK_CONTAINER(win), wv);
    gtk_widget_grab_focus(wv);

    App *app = (App*)malloc(sizeof(App));
    app->win = win;
    app->wv = wv;
    return app;
}

void add_script(App* app, const char* js) {
    _add_script(app->wv, js);
}

void set_size(App* app, int w, int h) {
    _set_size(app->win, w, h);
}

void set_fixed_size(App* app, int w,int h) {
    gtk_window_set_resizable(GTK_WINDOW(app->win), true);
    gtk_widget_set_size_request(app->win, w, h);
}

void set_min_size(App* app, int w, int h) {
    GdkGeometry g;
    g.min_width = w;
    g.min_height = h;
    gtk_window_set_geometry_hints(GTK_WINDOW(app->win), NULL, &g, GDK_HINT_MIN_SIZE);
}

void set_max_size(App* app, int w, int h) {
    GdkGeometry g;
    g.max_width = w;
    g.max_height = h;
    gtk_window_set_geometry_hints(GTK_WINDOW(app->win), NULL, &g, GDK_HINT_MAX_SIZE);
}

void move(App* app, int x, int y) {
    _move(app->win, x, y);
}

void set_title(App* app, const char* title) {
    _set_title(app->win, title);
}

void load(App* app, const char* url) {
    webkit_web_view_load_uri(WEBKIT_WEB_VIEW(app->wv), url);
}

void load_html(App* app, const char* html) {
    webkit_web_view_load_html(WEBKIT_WEB_VIEW(app->wv), html, NULL);
}

void eval(App* app, const char* js) {
    webkit_web_view_run_javascript(WEBKIT_WEB_VIEW(app->wv), js, NULL, NULL, NULL);
}

void quit(App *app) {
    gtk_main_quit();
    free(app);
}

void run(App* app) {
    if (gtk_init_check(NULL, NULL)) {
        gtk_widget_show_all(app->win);
        gtk_main();
    }
}

void dispatch() {
    g_idle_add_full(G_PRIORITY_HIGH_IDLE, _dispatch_cb, NULL, NULL);
}
