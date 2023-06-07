// SPDX-License-Identifier: MIT

#define __OBJC2__ 1
#import <Foundation/Foundation.h>
#import <WebKit/WebKit.h>

typedef struct {
    WKWebView* wv;
    NSWindow* win;
} CocoaWebView;

typedef void (*DispatchFunc) ();

CocoaWebView* create_cocoa(CGFloat x, CGFloat y, CGFloat w, CGFloat h, const char* title) {
    NSRect rect = NSMakeRect(x, y, w, h);
    NSWindowStyleMask style = NSWindowStyleMaskTitled|NSWindowStyleMaskClosable|NSWindowStyleMaskMiniaturizable;
    NSWindow* win = [[NSWindow alloc] initWithContentRect:rect styleMask:style backing: NSBackingStoreBuffered defer: NO];
    [win center];
    
    NSString* t = [NSString stringWithUTF8String:title];
    [win setTitle:t];
    
    WKWebView* wv = [WKWebView alloc];
    [win setContentView:wv];
    
    CocoaWebView* ret = malloc(sizeof(CocoaWebView));
    ret->wv = wv;
    ret->win = win;
    return ret;
}

void set_title(CocoaWebView* v, const char* title) {
    NSString* t = [NSString stringWithUTF8String:title];
    [v->win setTitle: t];
}

void add_user_script(CocoaWebView* wv, const char* js) {
    NSString* str = [NSString stringWithUTF8String:js];
    WKUserScript* script = [WKUserScript alloc];
    [script initWithSource:str injectionTime:WKUserScriptInjectionTimeAtDocumentStart forMainFrameOnly: YES];
    [wv->wv.configuration.userContentController addUserScript:script];
}

void eval(CocoaWebView* wv, const char* js) {
    NSString* str = [NSString stringWithUTF8String:js];
    [wv->wv evaluateJavaScript:str completionHandler: nil];
}

void load(CocoaWebView* wv, const char* url) {
    NSString* str = [NSString stringWithUTF8String:url];
    NSURL* nsurl = [NSURL URLWithString: str];
    NSURLRequest* req = [NSURLRequest requestWithURL:nsurl];
    [wv->wv loadRequest:req];
}

void set_html(CocoaWebView* wv, const char* html) {
    NSString* str = [NSString stringWithUTF8String:html];
    [wv->wv loadHTMLString:str baseURL: nil];
}

void set_position(CocoaWebView* wv, CGFloat x, CGFloat y) {
    NSPoint p = NSMakePoint(x, y);
    [wv->win setFrameTopLeftPoint:p];
}

void set_frame(CocoaWebView* wv, bool display, CGFloat x, CGFloat y, CGFloat w, CGFloat h) {
    NSRect frame = NSMakeRect(x, y, w, h);
    [wv->win setFrame:frame display:display];
}

void set_min_size(CocoaWebView* wv, CGFloat w, CGFloat h) {
    wv->win.minSize = NSMakeSize(w, h);
}

void set_max_size(CocoaWebView* wv, CGFloat w, CGFloat h) {
    wv->win.maxSize = NSMakeSize(w, h);
}

void dispatch_callback(DispatchFunc f) {
    f();
}

void dispatch(DispatchFunc f) {
    dispatch_async_f(dispatch_get_main_queue(), (void*)f, (dispatch_function_t)dispatch_callback);
}

void terminate() {
    [NSApp terminate:nil];
}

void run() {
    [NSApp run];
    [NSApp release];
}
