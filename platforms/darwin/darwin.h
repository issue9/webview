// SPDX-License-Identifier: MIT

#define __OBJC2__ 1
#import <Foundation/Foundation.h>
#import <Cocoa/Cocoa.h>
#import <WebKit/WebKit.h>

@interface AppDelegate: NSObject<NSApplicationDelegate>
// TODO
@end

@implementation AppDelegate

- (NSApplicationTerminateReply)applicationShouldTerminate:(NSApplication *)sender {
    return NSTerminateNow;
}

- (BOOL)applicationShouldTerminateAfterLastWindowClosed:(NSApplication *)sender {
    return YES;
}

@end

typedef struct {
    NSWindow* win;
    WKWebView* wv;
} CocoaWebView;

typedef void (*DispatchFunc) ();

// debug 是否启用调试模式
// x,y,w,h 表示窗口的左上角和宽高
// title 为窗口标题
CocoaWebView* create_cocoa(bool debug, CGFloat x, CGFloat y, CGFloat w, CGFloat h, const char* title) {
    [NSApplication sharedApplication];
    NSApp.delegate = [AppDelegate alloc];
    
    NSWindowStyleMask style = NSWindowStyleMaskTitled|NSWindowStyleMaskClosable|NSWindowStyleMaskMiniaturizable;
    CGRect rect = CGRectMake(x, y, w, h);
    NSWindow* win = [[NSWindow alloc]initWithContentRect:rect styleMask:style backing:NSBackingStoreBuffered defer:NO];
    [win center];
    [win makeKeyAndOrderFront:nil];
    
    WKWebViewConfiguration* config = [[WKWebViewConfiguration alloc] init];
    [config.preferences setValue:[NSNumber numberWithBool:debug] forKey:@"developerExtrasEnabled"];
    WKWebView* wv = [[WKWebView alloc] initWithFrame:rect configuration: config];
    win.contentView=wv;
    
    CocoaWebView* ret = malloc(sizeof(CocoaWebView));
    ret->wv = wv;
    ret->win = win;
    return ret;
}

void set_title(CocoaWebView* wv, const char* title) {
    NSString* t = [NSString stringWithUTF8String:title];
    [wv->win setTitle: t];
}

void add_user_script(CocoaWebView* wv, const char* js) {
    NSString* str = [NSString stringWithUTF8String:js];
    WKUserScript* script = [WKUserScript alloc];
    [script initWithSource:str injectionTime:WKUserScriptInjectionTimeAtDocumentStart forMainFrameOnly:YES];
    [wv->wv.configuration.userContentController addUserScript:script];
}

void eval(CocoaWebView* wv, const char* js) {
    NSString* str = [NSString stringWithUTF8String:js];
    [wv->wv evaluateJavaScript:str completionHandler:nil];
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
    NSPoint p = CGPointMake(x, y);
    [wv->win setFrameTopLeftPoint:p];
}

void set_frame(CocoaWebView* wv, bool display, CGFloat x, CGFloat y, CGFloat w, CGFloat h) {
    NSRect frame = CGRectMake(x, y, w, h);
    [wv->win setFrame:frame display:display];
}

void set_min_size(CocoaWebView* wv, CGFloat w, CGFloat h) {
    wv->win.minSize = CGSizeMake(w, h);
}

void set_max_size(CocoaWebView* wv, CGFloat w, CGFloat h) {
    wv->win.maxSize = CGSizeMake(w, h);
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
}
