// SPDX-License-Identifier: MIT

#import "_cgo_export.h"
#import "darwin.h"

void _add_user_script(WKWebViewConfiguration* config, NSString* js) {
    WKUserScript* script = [WKUserScript alloc];
    [script initWithSource:js injectionTime:WKUserScriptInjectionTimeAtDocumentStart forMainFrameOnly:YES];
    [config.userContentController addUserScript:script];
}

@implementation AppDelegate

- (NSApplicationTerminateReply)applicationShouldTerminate:(NSApplication *)sender {
    return NSTerminateNow;
}

- (BOOL)applicationShouldTerminateAfterLastWindowClosed:(NSApplication *)sender {
    return YES;
}

@end

@implementation AppScriptMessageHandler

- (void)userContentController:(WKUserContentController *)userContentController
      didReceiveScriptMessage:(WKScriptMessage *)message {
    messageCallback((char *) [message.body description].UTF8String);
}

@end

// debug 是否启用调试模式
// x,y,w,h 表示窗口的左上角和宽高
// title 为窗口标题
App* create_cocoa(bool debug, CGFloat x, CGFloat y, CGFloat w, CGFloat h, const char* title) {
    [NSApplication sharedApplication];
    NSApp.delegate = [AppDelegate alloc];
    
    NSWindowStyleMask style = NSWindowStyleMaskTitled|NSWindowStyleMaskClosable|NSWindowStyleMaskMiniaturizable;
    CGRect rect = CGRectMake(x, y, w, h);
    NSWindow* win = [[NSWindow alloc]initWithContentRect:rect styleMask:style backing:NSBackingStoreBuffered defer:NO];
    [win center];
    [win makeKeyAndOrderFront:nil];
    
    WKWebViewConfiguration* config = [[WKWebViewConfiguration alloc] init];
    [config.preferences setValue:[NSNumber numberWithBool:debug] forKey:@"developerExtrasEnabled"];
    _add_user_script(config, @"window.external = {\
        invoke: function(s) {\
            window.webkit.messageHandlers.external.postMessage(s);\
        },\
    };");
    id<WKScriptMessageHandler> messageHandler = [[AppScriptMessageHandler alloc] init];
    [config.userContentController addScriptMessageHandler:messageHandler name:@"external"];
    
    WKWebView* wv = [[WKWebView alloc] initWithFrame:rect configuration: config];
    win.contentView=wv;
    
    App* ret = malloc(sizeof(App));
    ret->wv = wv;
    ret->win = win;
    return ret;
}

void set_title(App* wv, const char* title) {
    NSString* t = [NSString stringWithUTF8String:title];
    [wv->win setTitle: t];
}

void add_user_script(App* wv, const char* js) {
    NSString* str = [NSString stringWithUTF8String:js];
    _add_user_script(wv->wv.configuration, str);
}

void eval(App* wv, const char* js) {
    NSString* str = [NSString stringWithUTF8String:js];
    [wv->wv evaluateJavaScript:str completionHandler:nil];
}

void load(App* wv, const char* url) {
    NSString* str = [NSString stringWithUTF8String:url];
    NSURL* nsurl = [NSURL URLWithString: str];
    NSURLRequest* req = [NSURLRequest requestWithURL:nsurl];
    [wv->wv loadRequest:req];
}

void set_html(App* wv, const char* html) {
    NSString* str = [NSString stringWithUTF8String:html];
    [wv->wv loadHTMLString:str baseURL: nil];
}

void set_position(App* wv, CGFloat x, CGFloat y) {
    NSPoint p = CGPointMake(x, y);
    [wv->win setFrameTopLeftPoint:p];
}

void set_frame(App* wv, bool display, CGFloat x, CGFloat y, CGFloat w, CGFloat h) {
    NSRect frame = CGRectMake(x, y, w, h);
    [wv->win setFrame:frame display:display];
}

void set_min_size(App* wv, CGFloat w, CGFloat h) {
    wv->win.minSize = CGSizeMake(w, h);
}

void set_max_size(App* wv, CGFloat w, CGFloat h) {
    wv->win.maxSize = CGSizeMake(w, h);
}

void dispatch_cb(void* f) {
    dispatchCallback();
}

void dispatch() {
    dispatch_async_f(dispatch_get_main_queue(), nil, (dispatch_function_t)dispatch_cb);
}

void terminate() {
    [NSApp terminate:nil];
}

void run() {
    [NSApp run];
}
