//
//  AppDelegate.m
//  FastnodeAutostart
//
//  Created by Tarak Upadhyaya on 10/5/15.
//  Copyright © 2015 Tarak Upadhyaya. All rights reserved.
//

#import "AppDelegate.h"
@import Rollbar;

@interface AppDelegate ()

@end

@implementation AppDelegate

- (void)applicationDidFinishLaunching:(NSNotification *)aNotification {
    // If we are starting up from the .Trash folder, terminate
    NSBundle *mbundle = [NSBundle mainBundle];
    NSLog(@"starting from %@", mbundle.bundlePath);
    if ([mbundle.bundlePath containsString:@"/.Trash/"]) {
        NSLog(@"in trash, terminating");
        [NSApp terminate:self];
    }

    [[NSUserDefaults standardUserDefaults] registerDefaults:@{ @"NSApplicationCrashOnExceptions": @YES }];
    // Rollbar, fastnoded post_client_item
    [Rollbar initWithAccessToken:@"XXXXXXX"];

    NSString *fastnodeBundle = [self bundleNameForApp:@"Fastnode"];
    NSString *menubarPath = [[NSWorkspace sharedWorkspace] absolutePathForAppBundleWithIdentifier:fastnodeBundle];

    NSArray *runningMenuItems = [NSRunningApplication runningApplicationsWithBundleIdentifier:fastnodeBundle];

    // Launch Fastnode if it isn't running
    if ([runningMenuItems count] == 0) {
        NSLog(@"No instances of %@ running, starting using bundle at %@", fastnodeBundle, menubarPath);

        // Pass in --boot flag to Fastnode.app when launching. This lets Fastnode.app know that it has started up
        // via system boot, so it can startup accordingly (e.g launch / not launch the sidebar as it sees fit)
        NSError *error = nil;
        NSURL *url = [NSURL fileURLWithPath:menubarPath];
        NSArray *arguments = [NSArray arrayWithObjects:@"--system-boot", nil];
        [[NSWorkspace sharedWorkspace]
                        launchApplicationAtURL:url
                        options:0
                        configuration:[NSDictionary dictionaryWithObject:arguments
                                                    forKey:NSWorkspaceLaunchConfigurationArguments]
                        error:&error];

        NSLog(@"error: %@", error);

        NSLog(@"started Fastnode.app, terminating");
    }

    [NSApp terminate:self];
}

- (void)applicationWillTerminate:(NSNotification *)aNotification {
    NSLog(@"terminating FastnodeAutostart");
}

- (NSString*) bundleNameForApp: (NSString*)app {
    NSString *prefix = [self bundlePrefix];
    return [NSString stringWithFormat:@"%@.%@", prefix, app];
}


- (NSString*) bundlePrefix {
    NSString *bundle = [[NSBundle mainBundle] bundleIdentifier];
    NSArray *parts = [bundle componentsSeparatedByString:@"."];
    
    NSRange range;
    range.location = 0;
    range.length = 2;
    
    NSArray *prefixParts = [parts subarrayWithRange:range];
    return [prefixParts componentsJoinedByString:@"."];
}

@end
