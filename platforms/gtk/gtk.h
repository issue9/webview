// SPDX-License-Identifier: MIT

// gtk 文档地址：https://docs.gtk.org/gtk3/index.html
// webview 文档地址：https://webkitgtk.org/reference/webkit2gtk/stable/index.html

#include <JavaScriptCore/JavaScript.h>
#include <gtk/gtk.h>
#include <webkit2/webkit2.h>

typedef struct {
    GtkWidget* win;
    GtkWidget* wv;
} App;

WebKitUserContentManager* userContentManager(App* app) {
    return webkit_web_view_get_user_content_manager(WEBKIT_WEB_VIEW(app->wv));
}

App* create_gtk(bool debug) {
    WebKitSettings* settings = webkit_settings_new();
    webkit_settings_set_enable_developer_extras(settings, debug);
    webkit_settings_set_javascript_can_access_clipboard(settings, true);

    GtkWidget* wv = webkit_web_view_new_with_settings(settings);
    //webkit_user_content_manager_register_script_message_handler(xx, "external");

    GtkWidget* win = gtk_window_new(GTK_WINDOW_TOPLEVEL);
    gtk_container_add(GTK_CONTAINER(win), wv);
    gtk_widget_grab_focus(wv);

    App *app = (App*)malloc(sizeof(App));
    app->win = win;
    app->wv = wv;
    return app;
}

void add_script(App* app, const char* js) {
    WebKitUserContentManager* m = userContentManager(app);
    WebKitUserScript* script = webkit_user_script_new(js, WEBKIT_USER_CONTENT_INJECT_TOP_FRAME, WEBKIT_USER_SCRIPT_INJECT_AT_DOCUMENT_START, NULL, NULL);
    webkit_user_content_manager_add_script(m, script);
}

void set_size(App* app, int w, int h) {
    GdkGeometry g;
    g.base_width = w;
    g.base_height = h;
    gtk_window_set_geometry_hints(GTK_WINDOW(app->win), NULL, &g, GDK_HINT_BASE_SIZE);
}

void set_fixed_size(App*app, int w,int h){
    gtk_window_set_resizable(GTK_WINDOW(app->win), true);
    gtk_widget_set_size_request(app->win, w, h);
}

void set_min_size(App* app, int w, int h) {
    GdkGeometry g;
    g.min_width = w;
    g.min_height = h;
    gtk_window_set_geometry_hints(GTK_WINDOW(app->win), NULL, &g, GDK_HINT_MIN_SIZE);
}

void set_max_size(App*app, int w, int h){
    GdkGeometry g;
    g.max_width = w;
    g.max_height = h;
    gtk_window_set_geometry_hints(GTK_WINDOW(app->win), NULL, &g, GDK_HINT_MAX_SIZE);
}

void set_title(App* app, const char* title) {
    gtk_window_set_title(GTK_WINDOW(app->win), title);
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
