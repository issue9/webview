// SPDX-License-Identifier: MIT

// gtk 文档地址：https://docs.gtk.org/gtk3/index.html
// webview 文档地址：https://webkitgtk.org/reference/webkit2gtk/stable/index.html

#ifndef WEBVIEW_GTK_H
#define WEBVIEW_GTK_H

#include <JavaScriptCore/JavaScript.h>
#include <gtk/gtk.h>
#include <webkit2/webkit2.h>

typedef struct {
    GtkWidget* win;
    GtkWidget* wv;
} App;

void dispatch();

App* create_gtk(bool debug, int x, int y, int w, int h, const char* title);

void add_script(App* app, const char* js);

void set_size(App* app, int w, int h);

void set_fixed_size(App* app, int w,int h);

void set_min_size(App* app, int w, int h);;

void set_max_size(App* app, int w, int h);

void move(App* app, int x, int y);

void set_title(App* app, const char* title);

void load(App* app, const char* url);

void load_html(App* app, const char* html);;

void eval(App* app, const char* js);

void quit(App *app);

void run(App* app);

#endif
