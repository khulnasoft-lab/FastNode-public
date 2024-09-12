// +build !standalone

#import <Foundation/Foundation.h>
#import <Cocoa/Cocoa.h>

#include <sys/sysctl.h>
#include <utmpx.h>

#import "_cgo_export.h"
#import "sidebar_darwin.h"


@implementation ActivationObserver

- (id)init {
    // Default value for "wasVisible" should be no
    [[NSUserDefaults standardUserDefaults] registerDefaults:@{@"wasVisible": @NO}];

    NSWorkspace *workspace = [NSWorkspace sharedWorkspace];

    // Add a listener for all future app activations, i.e., when an app is opened, or receives focus. This is needed
    // to launch the sidebar if a user launches Fastnode again (via Applications, Spotlight, etc). Since the app was already
    // running, the app is "activated".
    [[workspace notificationCenter] addObserver:self
                                       selector:@selector(appActivated:)
                                           name:NSWorkspaceDidActivateApplicationNotification
                                         object:workspace];
    return self;
}

- (void)appActivated:(NSNotification *)notification {
    // Listen for com.fastnode.Fastnode becoming activated, and call the go method cgoOnAppActivated

    NSRunningApplication *app = [[notification userInfo] objectForKey:NSWorkspaceApplicationKey];
    if ([[app bundleIdentifier] isEqualToString:bundleNameForApp(@"Fastnode")]) {
        // This calls back into go, launching the sidebar if needed
        cgoOnAppActivated();
    }
}

@end

static ActivationObserver* observer = nil;

void startObserver(char **err) {
    @try {
        @autoreleasepool {
            if (observer != nil) {
                NSLog(@"FastnodeActivationObserver: Observer already running");
                return;
            }

            observer = [[ActivationObserver alloc] init];
        }
    } @catch (NSException* ex) {
        *err = strdup([ex.reason UTF8String]);  // caller must free memory
    }
}


NSRunningApplication* getSidebar() {
    NSString* sidebar = bundleNameForApp(@"FastnodeApp");
    NSArray* apps = [NSRunningApplication runningApplicationsWithBundleIdentifier:sidebar];
    if ([apps count] == 0) {
        return nil;
    } else {
        if ([apps count] > 1) {
            NSLog(@"there are %lu sidebars running", (unsigned long)[apps count]);
        }
        return [apps objectAtIndex:0];
    }
}

bool isRunning(char **err) {
    return getSidebar() != nil;
}

char* appPath() {
    NSArray *components = [NSArray arrayWithObjects:[[NSBundle mainBundle] resourcePath], @"Fastnode.app", nil];
    NSString *fastnodeSidebarAppPath = [NSString pathWithComponents:components];
    return strdup([fastnodeSidebarAppPath UTF8String]); // caller frees memory
}

void launch(char **err) {
    NSArray *components = [NSArray arrayWithObjects:[[NSBundle mainBundle] resourcePath], @"Fastnode.app", nil];
    NSString *fastnodeSidebarAppPath = [NSString pathWithComponents:components];
    NSLog(@"Launching sidebar from %@", fastnodeSidebarAppPath);
    [[NSWorkspace sharedWorkspace] launchApplication:fastnodeSidebarAppPath];
}

void focus(char **err) {
    NSRunningApplication* sidebar = getSidebar();
    if (sidebar == nil) {
        return;
    }
    [sidebar activateWithOptions:NSApplicationActivateIgnoringOtherApps];
}

void quitSidebar(char **err) {
    NSRunningApplication* sidebar = getSidebar();
    if (sidebar == nil) {
        return;
    }
    [sidebar terminate];
}

// --

void setWasVisible(bool val) {
    [[NSUserDefaults standardUserDefaults] setBool:val forKey:@"wasVisible"];
}

bool wasVisible() {
    return [[NSUserDefaults standardUserDefaults] boolForKey:@"wasVisible"];
}

// --

NSString* bundlePrefix() {
    NSString *bundle = [[NSBundle mainBundle] bundleIdentifier];
    NSArray *parts = [bundle componentsSeparatedByString:@"."];

    NSRange range;
    range.location = 0;
    range.length = 2;

    NSArray *prefixParts = [parts subarrayWithRange:range];
    return [prefixParts componentsJoinedByString:@"."];
}

NSString* bundleNameForApp(NSString* app) {
    NSString *prefix = bundlePrefix();
    return [NSString stringWithFormat:@"%@.%@", prefix, app];
}
