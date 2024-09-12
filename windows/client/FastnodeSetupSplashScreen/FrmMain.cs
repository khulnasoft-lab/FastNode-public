using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Data;
using System.Diagnostics;
using System.Drawing;
using System.IO;
using System.Linq;
using System.Reflection;
using System.Runtime.InteropServices;
using System.Text;
using System.Windows.Forms;
using FastnodeSetupSplashScreen.Effects;

namespace FastnodeSetupSplashScreen {
    public partial class FrmMain : Form {

        private readonly System.Windows.Controls.MediaElement m_videoPlayer;
        private UInt64? m_fastnodeSetupGoneTimestamp = null;
        private static readonly string k_timedOutFile = "fastnodesetup_splash_screen_timed_out";

        public FrmMain() {
            InitializeComponent();

            m_videoPlayer = new System.Windows.Controls.MediaElement();
            m_videoPlayer.LoadedBehavior = System.Windows.Controls.MediaState.Manual;
            m_videoPlayer.Stretch = System.Windows.Media.Stretch.UniformToFill;
            m_videoPlayer.Effect = new DeeperColorEffect();
            m_videoPlayer.MediaEnded += VideoPlayer_MediaEnded;
            videoElementHost.Child = m_videoPlayer;

            var videoFilePath = WriteVideoFile();
            if (videoFilePath != null) {
                m_videoPlayer.Source = new Uri(videoFilePath);
                m_videoPlayer.Play();
            } else {
                // for some reason or other we couldn't write the video file to disk.
                // as a fallback we just don't play the video.  the text on the bottom looks a little out of place, but it's functional.
                // (this shouldn't happen very often.)
            }

            CenterToScreen();
            BringToFront();

            timeoutTimer.Interval = (int)TimeSpan.FromSeconds(5).TotalMilliseconds;
            timeoutTimer.Enabled = true;

            donePollTimer.Interval = (int)TimeSpan.FromMilliseconds(300).TotalMilliseconds;
            donePollTimer.Enabled = true;
        }

        private void VideoPlayer_MediaEnded(object sender, System.Windows.RoutedEventArgs e) {
            m_videoPlayer.Position = TimeSpan.Zero;
            m_videoPlayer.Play();
        }

        private void timeoutTimer_Tick(object sender, EventArgs e) {
            // here's how timeouts work:
            // - after FastnodeSetup.exe exits, we wait for up to one minute for Fastnode.exe to be running.
            // - if 90 seconds expires, we show an error messagebox and quit.

            if (Process.GetProcessesByName("FastnodeSetup").Length > 0) {
                // FastnodeSetup.exe is still running -- don't start the timeout clock
                return;
            }

            // FastnodeSetup.exe isn't running
            // if this is the first time we've seen this, set m_fastnodeSetupGoneTimestamp
            if(!m_fastnodeSetupGoneTimestamp.HasValue) {
                m_fastnodeSetupGoneTimestamp = GetTickCount64();
            }

            var millisPassedSinceSetupExited = GetTickCount64() - m_fastnodeSetupGoneTimestamp;
            if(millisPassedSinceSetupExited < 90 * 1000) {
                return;
            }

            // FastnodeSetup has been gone for 90 seconds now -- timeout, and disable the timer so
            // the messabe box doesn't appear every timer interval.
            timeoutTimer.Enabled = false;

            // Touch file indicating setup timed out
            TouchTimedOutFile();

            MessageBox.Show(@"Installing Fastnode seems to be taking much longer than expected. Windows Defender may be preventing the installation from finishing. Please disable Windows Defender and try intalling Fastnode again.

If you are still experiencing issues, please see https://github.com/khulnasoft-lab/issue-tracker/issues/158.

This splash screen will now exit.",
                "Fastnode installation timeout", MessageBoxButtons.OK, MessageBoxIcon.Warning);

            Application.Exit();
        }

        private void donePollTimer_Tick(object sender, EventArgs e) {
            // wait for Fastnode.exe to be running, window visible, then exit the application
            foreach (var fastnodeProcess in Process.GetProcessesByName("Fastnode")) {
                try {
                    fastnodeProcess.Refresh();
                    if (fastnodeProcess.MainWindowHandle != IntPtr.Zero) {
                        // Fastnode.exe with visible window!
                        Application.Exit();
                    }
                } catch {
                    // Exceptions can be thrown if e.g. the process is running but doesn't have a main window yet.
                    // Just skip.
                }
            }
        }

        protected override void OnPaintBackground(PaintEventArgs e) {
            base.OnPaintBackground(e);
            e.Graphics.FillRectangle(Brushes.White, new Rectangle(2, 2, 280, 280));
        }

        private static string WriteVideoFile() {
            try {
                // get Local App Data directory
                IntPtr pLocalAppDataPath;
                if (SHGetKnownFolderPath(new Guid("F1B32785-6FBA-4FCF-9D55-7B8E7F157091"), 0, IntPtr.Zero, out pLocalAppDataPath) != 0) {
                    return null;  // error
                }
                var localAppDataPath = Marshal.PtrToStringUni(pLocalAppDataPath);
                Marshal.FreeCoTaskMem(pLocalAppDataPath);

                // get Fastnode directory under it
                var fastnodeDataPath = Path.Combine(localAppDataPath, "Fastnode");
                Directory.CreateDirectory(fastnodeDataPath);

                // write video file
                var videoFilePath = Path.Combine(fastnodeDataPath, "FastnodeSetupSplashScreenVideo.mp4");
                File.WriteAllBytes(videoFilePath, FastnodeSetupSplashScreen.Properties.Resources.FastnodeSetupSplashScreenVideo);

                return videoFilePath;
            } catch {
                // lots of things could cause this, e.g. the file exists and isn't writeable.
                // we support fallback (don't play the video), so just use that.
                return null;
            }
        }

        private static void TouchTimedOutFile()
        {
            try
            {
                // get Local App Data directory
                IntPtr pLocalAppDataPath;
                if (SHGetKnownFolderPath(new Guid("F1B32785-6FBA-4FCF-9D55-7B8E7F157091"), 0, IntPtr.Zero, out pLocalAppDataPath) != 0)
                {
                    return;  // error
                }
                var localAppDataPath = Marshal.PtrToStringUni(pLocalAppDataPath);
                Marshal.FreeCoTaskMem(pLocalAppDataPath);

                // get Fastnode directory under it
                var fastnodeDataPath = Path.Combine(localAppDataPath, "Fastnode");
                Directory.CreateDirectory(fastnodeDataPath);
                var timedOutFilePath = Path.Combine(fastnodeDataPath, k_timedOutFile);
                if (!File.Exists(timedOutFilePath))
                {
                    File.Create(timedOutFilePath);
                }
                File.SetLastWriteTimeUtc(timedOutFilePath, DateTime.UtcNow);
            }
            catch
            {
                // Just skip
                return;
            }
        }

        [DllImport("shell32.dll")]
        private static extern int SHGetKnownFolderPath([MarshalAs(UnmanagedType.LPStruct)] Guid rfid,
            uint dwFlags, IntPtr hToken, out IntPtr pszPath);

        [DllImport("kernel32")]
        private extern static UInt64 GetTickCount64();
    }
}
