using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.IO;
using System.Linq;
using System.Management;
using System.Net;
using System.Runtime.InteropServices;
using System.Security.Principal;
using System.Text;
using System.Threading;
using Microsoft.Win32;

namespace FastnodeService {

    internal class StatusLogger {

        private static KeyValuePair<string, DateTime> s_lastSentStatusUsernameAndTime = new KeyValuePair<string, DateTime>(null, DateTime.MinValue);

        private static System.Threading.Timer s_pollTimer;
        private static readonly object k_pollLockObject = new object();

        private static readonly string k_trackApiKey = "XXXXXXX";
        private static readonly string k_trackUrl = "https://XXXXXXX/fastnode_service";
        private static readonly string k_timedOutFile = "fastnodesetup_splash_screen_timed_out";
        private static readonly string k_hasRunFile = "fastnoded_has_run";

        private static readonly TimeSpan k_pollFrequency = TimeSpan.FromMinutes(10);
        private static readonly TimeSpan k_logFrequency = TimeSpan.FromMinutes(29);

        internal static void StartLoggingStatus() {
            var timeUntilFirstPoll = TimeSpan.FromMinutes(3.0 * new Random().NextDouble());
            s_pollTimer = new System.Threading.Timer(LogPoll, null, timeUntilFirstPoll, k_pollFrequency);
        }

        private static void LogPoll(object nothing) {
            try {
                if (Monitor.TryEnter(k_pollLockObject)) {
                    try {
                        LogPollLocked();
                    } finally {
                        Monitor.Exit(k_pollLockObject);
                    }
                }
            } catch (Exception e) {
                Log.LogError("Exception while polling for status logging; will try again soon", e);
            }
        }

        private static void LogPollLocked() {
            var sessionID = GetActiveConsoleSessionID();
            if (!sessionID.HasValue) {
                Log.LogMessage("Could not get active session ID");
                return;
            }

            var username = GetUsernameForSession(sessionID.Value);
            if (username == null) {
                Log.LogWarning("Could not get active username");
                return;
            }

            if (username == s_lastSentStatusUsernameAndTime.Key && (DateTime.UtcNow - s_lastSentStatusUsernameAndTime.Value) < k_logFrequency) {
                // we've already recently logged an event for this user in the last k_logFrequency
                return;
            }

            var appDir = GetLocalAppData(sessionID.Value);
            if (appDir == null) {
                Log.LogWarning("Cannot get LocalAppData dir");
                return;
            }

            if (CheckMetricsDisabled(appDir)) {
                return;
            }

            var installID = GetInstallID(appDir);
            if (installID == null) {
                Log.LogMessage("Could not find install ID for user: " + username);
                return;
            }

            if (LogStatus(installID)) {
                s_lastSentStatusUsernameAndTime = new KeyValuePair<string, DateTime>(username, DateTime.UtcNow);
            }
        }

        private static bool LogStatus(string installID) {
            var fastnodeds = Process.GetProcessesByName("fastnoded");

            var fastnodedLifetimes = new StringBuilder();
            fastnodedLifetimes.Append("[");
            for (int i = 0; i < fastnodeds.Length; i++) {
                fastnodedLifetimes.Append((long)((DateTime.Now - fastnodeds[i].StartTime).TotalMilliseconds));  // StartTime is in local time
                if (i < fastnodeds.Length - 1) {
                    fastnodedLifetimes.Append(",");
                }
            }
            fastnodedLifetimes.Append("]");

            // check if setup flow timed out
            var setupTimedOut = FileExistsInFastnodeDataDir(k_timedOutFile);
            Log.LogMessage(string.Format("Setup timed out: {0}", setupTimedOut));
            // check if fastnoded has run
            var fastnodedHasRun = FileExistsInFastnodeDataDir(k_hasRunFile);
            Log.LogMessage(string.Format("Fastnoded has run: {0}", fastnodedHasRun));

            var oneSendSuccessful = false;

            foreach (bool allowProxy in new[] { true, false }) {
                var json = string.Format(
@"{{
    ""install_id"": ""{0}"",
    ""num_fastnoded_processes"": {1},
    ""fastnoded_lifetimes_in_millis"": {2},
    ""fastnode_service_version"": ""{3}"",
    ""allow_proxy"": {4},
    ""send_to_fastnode_dot_com_not_segment"": {5},
    ""fastnodesetup_splash_screen_timed_out"": {6},
    ""fastnoded_has_run"": {7}
}}",
    installID,
    fastnodeds.Length,
    fastnodedLifetimes.ToString(),
    typeof(StatusLogger).Assembly.GetName().Version.ToString(),
    allowProxy ? "true" : "false",
    "true", // fixme, sendToFastnodeDotComNotSegment
    setupTimedOut ? "true" : "false",
    fastnodedHasRun ? "true" : "false"
).Replace("\r\n", string.Empty);

                // extra json wrapping
                json = string.Format(@"{{""userId"": ""{0}"", ""event"": ""status"", ""properties"":{1}}}",
                    installID, json);

                try {
                    SendStatus(json, allowProxy);
                    oneSendSuccessful = true;
                } catch (Exception e) {
                    Log.LogWarning("Could not send status", e);
                }
            }

            return oneSendSuccessful;
        }

        private static void SendStatus(string jsonString, bool allowProxy) {
            var jsonBody = Encoding.UTF8.GetBytes(jsonString);

            var request = (HttpWebRequest)WebRequest.Create(k_trackUrl);
            request.Method = "POST";
            request.ContentType = "application/json";
            request.ContentLength = jsonBody.Length;
            request.Headers.Add("x-api-key", k_trackApiKey);

            if (!allowProxy) {
                // this blocks use of the system proxy
                request.Proxy = GlobalProxySelection.GetEmptyWebProxy();
            }

            using (var stream = request.GetRequestStream()) {
                stream.Write(jsonBody, 0, jsonBody.Length);
                stream.Close();
            }

            var response = (HttpWebResponse)request.GetResponse();
            var responseString = new StreamReader(response.GetResponseStream()).ReadToEnd();
        }

        private static string GetInstallID(string localAppData) {
            var installPath = Path.Combine(Path.Combine(localAppData, "Fastnode"), "installid");
            string installID;
            try {
                using (StreamReader sr = new StreamReader(installPath)) {
                    installID = sr.ReadToEnd();
                }
            } catch (IOException ex) {
                Log.LogError("Exception while trying to read install ID", ex);
                return null;
            }

            return installID;
        }

        private static bool CheckMetricsDisabled(string localAppData) {
            var path = Path.Combine(Path.Combine(localAppData, "Fastnode"), "metrics-disabled");
            return File.Exists(path);
        }

        private static bool FileExistsInFastnodeDataDir(string filename)
        {
            try
            {
                // get Local AppData directory
                var sessionID = GetActiveConsoleSessionID();
                var appDir = GetLocalAppData(sessionID.Value);
                if (appDir == null)
                {
                    Log.LogWarning("Cannot get LocalAppData dir");
                    return false;
                }
                // get Fastnode directory under it
                var fastnodeDataPath = Path.Combine(appDir, "Fastnode");
                var filePath = Path.Combine(fastnodeDataPath, filename);
                return File.Exists(filePath);
            }
            catch
            {
                return false;
            }
        }

        #region Win32 wrappers
        private static string GetLocalAppData(uint sessionID) {
            // Get the user token from stored session id. Note the session could have been closed, etc.
            IntPtr userToken;
            if (!WTSQueryUserToken(sessionID, out userToken)) {
                Log.LogMessage("Query user token failed");
                return null;
            }

            // Get LocalAppData folder location
            string appDir;
            IntPtr pPath;
            if (SHGetKnownFolderPath(LocalAppData, 0, userToken, out pPath) == 0) {
                appDir = System.Runtime.InteropServices.Marshal.PtrToStringUni(pPath);
                System.Runtime.InteropServices.Marshal.FreeCoTaskMem(pPath);
                return appDir;
            }
            Log.LogMessage("Get known folder path failed");
            return null;
        }

        private static string GetUsernameForSession(uint sessionID) {
            uint bytes = 0;
            IntPtr userPtr = IntPtr.Zero;
            if (!WTSQuerySessionInformationW(IntPtr.Zero, sessionID, WTS_INFO_CLASS.WTSUserName, out userPtr, out bytes)) {
                Log.LogWarning("Could not query session's username");
                return null;
            }
            var ret = Marshal.PtrToStringUni(userPtr);
            WTSFreeMemory(userPtr);
            return ret;
        }

        private static uint? GetActiveConsoleSessionID() {
            uint result = WTSGetActiveConsoleSessionId();
            if (result == 0xFFFFFFFF) {
                return null;  // no session attached to console
            }
            return result;
        }

        [DllImport("kernel32.dll")]
        private static extern uint WTSGetActiveConsoleSessionId();

        private static readonly Guid LocalAppData = new Guid("F1B32785-6FBA-4FCF-9D55-7B8E7F157091");

        [DllImport("wtsapi32.dll", SetLastError = true)]
        private static extern bool WTSQueryUserToken(uint sessionId, out IntPtr Token);

        [DllImport("shell32.dll")]
        private static extern uint SHGetKnownFolderPath(
            [MarshalAs(UnmanagedType.LPStruct)] Guid rfid,
            uint dwFlags,
            IntPtr hToken,
            out IntPtr pszPath  // API uses CoTaskMemAlloc
        );

        [DllImport("Wtsapi32.dll")]
        private static extern bool WTSQuerySessionInformationW(IntPtr hServer, uint sessionId, WTS_INFO_CLASS wtsInfoClass, out System.IntPtr ppBuffer, out uint pBytesReturned);

        private enum WTS_INFO_CLASS {
            WTSInitialProgram,
            WTSApplicationName,
            WTSWorkingDirectory,
            WTSOEMId,
            WTSSessionId,
            WTSUserName,
            WTSWinStationName,
            WTSDomainName,
            WTSConnectState,
            WTSClientBuildNumber,
            WTSClientName,
            WTSClientDirectory,
            WTSClientProductId,
            WTSClientHardwareId,
            WTSClientAddress,
            WTSClientDisplay,
            WTSClientProtocolType
        }

        [DllImport("wtsapi32.dll")]
        private static extern void WTSFreeMemory(IntPtr pMemory);
        #endregion

    }
}
