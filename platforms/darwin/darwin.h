// SPDX-License-Identifier: MIT

#define __OBJC2__ 1
#import <Foundation/Foundation.h>
#import <Cocoa/Cocoa.h>
#import <WebKit/WebKit.h>

typedef struct {
    NSWindow* win;
    WKWebView* wv;
} App;

@interface AppDelegate: NSObject<NSApplicationDelegate>
@end

@interface AppScriptMessageHandler : NSObject <WKScriptMessageHandler>
@end

App* create_cocoa(bool debug, CGFloat x, CGFloat y, CGFloat w, CGFloat h, const char* title);

void set_title(App* wv, const char* title);

void add_user_script(App* wv, const char* js);

void eval(App* wv, const char* js);

void load(App* wv, const char* url);

void set_html(App* wv, const char* html);

void set_position(App* wv, CGFloat x, CGFloat y);

void set_frame(App* wv, bool display, CGFloat x, CGFloat y, CGFloat w, CGFloat h);

void set_min_size(App* wv, CGFloat w, CGFloat h);

void set_max_size(App* wv, CGFloat w, CGFloat h);

void dispatch();

void terminate();

void run();
