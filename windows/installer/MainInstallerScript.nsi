Name "Fastnode"
VIProductVersion "${VERSION}"
VIAddVersionKey "ProductVersion" "${VERSION}"  ; for some reason we need this, too
VIAddVersionKey "FileVersion" "${VERSION}"
VIAddVersionKey "LegalCopyright" "Copyright © Fastnode"
!ifdef WRITE_UNINSTALLER_ONLY
	VIAddVersionKey "FileDescription" "Fastnode Uninstaller"
	VIAddVersionKey "ProductName" "Fastnode Uninstaller"
	VIAddVersionKey "OriginalFilename" "FastnodeUninstallerGenerator.exe"
	VIAddVersionKey "InternalName" "FastnodeUninstallerGenerator"

	OutFile "current_build_bin\out\FastnodeUninstallerGenerator.exe"
!else
	VIAddVersionKey "FileDescription" "Fastnode Setup"
	VIAddVersionKey "ProductName" "Fastnode Setup"
	VIAddVersionKey "OriginalFilename" "FastnodeSetup.exe"
	VIAddVersionKey "InternalName" "FastnodeSetup"

	OutFile "current_build_bin\out\FastnodeSetup.exe"
!endif
Icon "..\tools\artwork\icon\app.ico"
SetCompressor /SOLID lzma
RequestExecutionLevel admin
InstallDir "$PROGRAMFILES64\Fastnode"
BrandingText " "
ShowInstDetails nevershow
ShowUninstDetails nevershow

Var executable_type ; e.g. "installer" "uninstaller" "updater" etc
Var skip_onboarding
Var cmdflags_start
Var cmdflags_substring

!include "MUI.nsh"
!include "LogicLib.nsh"
!include "WordFunc.nsh"
!include "StrFunc.nsh"
${StrLoc} ; must initialize this before it can be used in a Function (a nuance of StrFunc.nsh)
${UnStrLoc}
${StrRep}
${UnStrRep}
!include "FileFunc.nsh"
!include "WinVer.nsh"
!include "GetProcessInfo.nsh"
!include "servicelib.nsh"
!include "x64.nsh"
!include "CPUFeatures.nsh"
!include "NsisIncludes\Debug.nsh"
!include "NsisIncludes\GenerateMachineIDIfAppropriate.nsh"
!include "NsisIncludes\CheckInstallPrereqs.nsh"
!include "NsisIncludes\CheckAlreadyRunningInstallOrUninstall.nsh"
!include "NsisIncludes\KillAllAvailableRunningInstances.nsh"

!define MUI_ICON "..\tools\artwork\icon\app.ico"
!define MUI_UNICON "..\tools\artwork\icon\app.ico"
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_UNPAGE_INSTFILES
!define MUI_CUSTOMFUNCTION_ABORT UserAbortedInstallCallback

Function UserAbortedInstallCallback
	${Debug} "User aborted install"
FunctionEnd

Function un.onInit
	StrCpy $executable_type "uninstaller"
	Call un.CheckAlreadyRunningInstallOrUninstall
FunctionEnd

Function .onInit
	!ifdef WRITE_UNINSTALLER_ONLY
		; this is a build of the setup that is only meant to emit the uninstaller (so we can subsequently sign it)
		System::Call "kernel32::GetCurrentDirectoryW(i ${NSIS_MAX_STRLEN}, t .r0)"
		WriteUninstaller "$0\Uninstaller.exe"
		SetErrorLevel 0
		Quit
	!endif

	; we are a 32 bit installer, uninstaller, and updater, but the main Fastnode binaries are 64-bit, so we
	;   try to standardize on the 64-bit view where possible.
	SetRegView 64

	${StrLoc} $0 $CMDLINE "testprereqsonly" ">"
	${If} $0 != ""
		${Debug} "testprereqsonly command arg set; testing prereqs only.."

		Call CheckInstallPrereqs
		Pop $0
		Pop $1
		${If} $0 == "ok"
			${Debug} "prereqs checked out ok"
			SetErrorLevel 0 ; just pick some rare values
		${Else}
			${Debug} "prereqs checked failed"
			SetErrorLevel 14
		${EndIf}

		Quit
	${EndIf}

	StrCpy $executable_type "installer"

	Call CheckAlreadyRunningInstallOrUninstall

	Call ReadMachineIDOrGenerateIfAppropriate

	; Check for installation prereq's
	${StrLoc} $0 $CMDLINE "skipprereqs" ">"
	${If} $0 == "" ; no match
		Call CheckInstallPrereqs
		Pop $0
		Pop $1
		${If} $0 != "ok"
			${Debug} "Prereq fail reason: $1"

			MessageBox MB_OK|MB_ICONINFORMATION $0
			SetErrorLevel 21
			Quit
		${EndIf}
	${Else}
		${Debug} "Skipping prereqs check due to command line argument..."
	${EndIf}

	; this will make the installer not show any UI (other than MessageBox's)
	; note we don't set this on the uninstaller
	SetSilent silent
FunctionEnd

Section ""
!ifndef WRITE_UNINSTALLER_ONLY ; otherwise don't include an installer section
	; Do this especially before launching Fastnode or any of the executables.
	Call WriteTentativeOrActualMachineIDToRegistry

	; Let the fun begin!
	${Debug} "Copying files..."
	SetOutPath "$INSTDIR"

	; Launch the splash screen!
	File "current_build_bin\in\FastnodeSetupSplashScreen.exe"
	File "current_build_bin\in\FastnodeSetupSplashScreen.exe.config"
	${StrLoc} $0 $CMDLINE "--plugin-launch" ">"
	${If} $0 == "" ; command line flag NOT present -> show the splash screen
		Exec '"$INSTDIR\FastnodeSetupSplashScreen.exe"'
	${EndIf}

	ClearErrors
 	ReadRegDword $0 HKLM "SOFTWARE\Microsoft\VisualStudio\14.0\VC\Runtimes\x64" "Installed"
	IfErrors 0 redist_found
		File "current_build_bin\in\vc_redist.x64.exe"
		ExecWait '"$INSTDIR\vc_redist.x64.exe" /install /passive /quiet /norestart'
		Delete /REBOOTOK "$INSTDIR\vc_redist.x64.exe"

	redist_found:
	File /r "current_build_bin\in\win-unpacked"
	File "current_build_bin\in\FastnodeService.exe"
	File "current_build_bin\in\FastnodeService.exe.config"
	File "current_build_bin\in\tensorflow.dll"
	File "current_build_bin\in\fastnoded.exe"
	File "current_build_bin\in\fastnode-lsp.exe"

	WriteRegStr HKLM "Software\Fastnode\AppData" "InstallPath" "$INSTDIR"

	; Set 'Run' key in registry
	;
	; Note: This is updated by client/internal/autostart/autostart_windows.go.
	; The user has the option to disable autostart through the copilot settings so this
	; ensures the 'Run' key is only set when autostart is enabled.
	WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Run" "Fastnode" '"$INSTDIR\fastnoded.exe" --system-boot'

	; Set protocol handler for electron application
	WriteRegStr HKLM "Software\Classes\fastnode" "" "URL:fastnode"
	WriteRegStr HKLM "Software\Classes\fastnode" "URL Protocol" ""
	WriteRegStr HKLM "Software\Classes\fastnode\shell\open\command" "" '"$INSTDIR\win-unpacked\Fastnode.exe" "%1"'

	; Add 'Program Files' shortcut.  This is particularly (well, somewhat) important for users who disable
	; auto-start.
	; The last argument on Fastnode Local Settings points to fastnoded.exe for the "icon" parameter.
	SetShellVarContext all ; install shortcut for all users
	CreateDirectory "$SMPROGRAMS\Fastnode"
	CreateShortCut "$SMPROGRAMS\Fastnode\Fastnode.lnk" "$INSTDIR\fastnoded.exe"

	; Install service
	; don't forget the trailing ';' in the param list
	!insertmacro SERVICE "create" "FastnodeService" "path=$INSTDIR\FastnodeService.exe;autostart=1;interact=0;display=FastnodeService;description=Fastnode Service maintains your installation of Fastnode to ensure it is always up to date.;"

	; Setup uninstaller
	File "current_build_bin\out\Uninstaller.exe"
	WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Fastnode" "DisplayName" "Fastnode"
	WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Fastnode" "Publisher" "Manhattan Engineering Inc"
	WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Fastnode" "DisplayIcon" "$\"$INSTDIR\FastnodeService.exe$\""
	WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Fastnode" "UninstallString" "$\"$INSTDIR\Uninstaller.exe$\""

	; Launch service and fastnoded
	!insertmacro SERVICE "start" "FastnodeService" ""

	${StrLoc} $skip_onboarding $CMDLINE "--skip-onboarding" ">"
	${If} $skip_onboarding != ""
		${Debug} "skip-onboarding flag specified; creating env var FASTNODE_SKIP_ONBOARDING for fastnoded child process"
		System::Call 'Kernel32::SetEnvironmentVariable(t, t)i ("FASTNODE_SKIP_ONBOARDING", "1").r0'
	${EndIf}

	; This logic takes everything starting at "--" and appends it to fastnoded.exe. The goal is forward along any
	; commandline flags passed into the installer. We first look for `FastnodeSetup.exe` to exclude any occurances
	; of "--" in the path, then we find find the first `--` and forward everything along to `fastnoded.exe`
	${StrLoc} $cmdflags_start $CMDLINE "FastnodeSetup.exe" ">"
	${If} $cmdflags_start != ""
		StrCpy $cmdflags_substring $CMDLINE "" $cmdflags_start
		${StrLoc} $cmdflags_start $cmdflags_substring "--" ">"
		${If} $cmdflags_start != ""
			StrCpy $cmdflags_substring $cmdflags_substring "" $cmdflags_start
		${EndIf}
	${EndIf}

	Exec '"$INSTDIR\fastnoded.exe" $cmdflags_substring'

	${Debug} "Install completed."
!endif
SectionEnd

Section "Uninstall"
	; kill all possible running instances
	Call un.KillAllAvailableRunningInstances
	Sleep 2000  ; This used to not be here, but despite best efforts seems like the RMDir still sometimes needs reboot

	; remove old tray icon if appropriate / possible
	SetRegView 64
	ReadRegDWORD $0 HKCU "Software\Fastnode\AppData" "LastTrayHwnd"
	${If} $0 > 0
		System::Call '*(&l4, i, i, i, i, i, &t64) i(, $0, 1702127979, 0, 0, 0, "") .r0'
		System::Call 'Shell32::Shell_NotifyIcon(i 2, i r0) i.r1'
		System::Free $0
	${EndIf}

	; Note that fastnoded.exe could still be running in other user sessions
	; Thus we'll specify /REBOOTOK when deleting files, and let the user know if they need to reboot

	; stop and uninstall service
	!insertmacro SERVICE "stop" "FastnodeService" ""
	Sleep 2000
	FindProcDLL::WaitProcEnd "FastnodeService.exe" 20000
	Sleep 2000  ; This used to not be here, but despite best efforts seems like the RMDir still sometimes needs reboot
	!insertmacro SERVICE "delete" "FastnodeService" ""

	RMDir /r /REBOOTOK "$INSTDIR"  ; This will delete all of the files, including the uninstaller.

	; the line below is commented out, so that $LOCALAPPDATA\Fastnode will be left behind.
	; this is important so that the editors (at least Atom) know Fastnode has been installed previously
	;   -> they can differentiate an uninstall vs the plugin was installed and it needs to show
	;   the installation wizard.
	; this also mirrors the uninstallation behavior on macOS of leaving behind ~/.fastnode
	; RMDir /r /REBOOTOK "$LOCALAPPDATA\Fastnode"  ; Log files, etc.  We might leave behind ones for other users.

	; delete the 'Program Files' shortcut
	SetShellVarContext all ; uninstall shortcut for all users
	Delete /REBOOTOK "$SMPROGRAMS\Fastnode\Fastnode.lnk"
	Delete /REBOOTOK "$SMPROGRAMS\Fastnode\Fastnode Local Settings.lnk"
	RMDir "$SMPROGRAMS\Fastnode"

	; there might be other AppData's for other users, but we'll unfortunately have
	;   to leave those behind.
	DeleteRegKey HKCU "Software\Fastnode\AppData" ; Don't delete the MachineID
	DeleteRegKey HKLM "Software\Fastnode\AppData" ; Don't delete the MachineID
	SetRegView 32
	DeleteRegKey HKCU "Software\Fastnode\AppData" ; Don't delete the MachineID
	DeleteRegKey HKLM "Software\Fastnode\AppData" ; Don't delete the MachineID
	SetRegView 64

	DeleteRegKey HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Fastnode"  ; note: HKLM!

	; Delete the 'Run' key if it exists
	DeleteRegValue HKCU "Software\Microsoft\Windows\CurrentVersion\Run" "Fastnode"  ; note: HKCU!

	; Delete protocol handler for electron application
	DeleteRegKey HKLM "Software\Classes\fastnode"

	; It gets added here too for some reason; have to to delete it in HKCR as well
	DeleteRegKey HKCR "fastnode"

	IfRebootFlag 0 noreboot
		MessageBox MB_YESNO|MB_ICONINFORMATION "There are some files that will not be deleted until you reboot your computer, probably because another user is running Fastnode.  Would you like to reboot now?" IDNO noreboot
		Reboot
	noreboot:

	${un.Debug} "Uninstall completed."
SectionEnd
