// +build !standalone

package statusicon

import (
	"log"
	"os"

	"github.com/khulnasoft-lab/fastnode/fastnode-go/client/internal/reg"
	"github.com/khulnasoft-lab/fastnode/fastnode-golib/systray"
	"github.com/skratchdot/open-golang/open"
)

func (ui *UI) onBeforeRun() {
	// save the handle to the registry so that we can destroy the icon next time
	h, err := reg.TrayIconHandle()
	if err != nil {
		log.Println("ignoring pre-existing tray icon handle:", err)
		return
	}
	log.Println("destroying pre-existing tray icon with handle", h)
	systray.CleanupHandle(systray.Handle(h))
}

func (ui *UI) onHandleReceived(h systray.Handle) {
	// save the handle to the registry so that we can destroy the icon next time
	log.Println("saving tray icon handle", h)
	err := reg.SetTrayIconHandle(uintptr(h))
	if err != nil {
		log.Println(err)
	}
}

func (ui *UI) onSettingsClicked() {
	open.Run("fastnode://settings")
}

func (ui *UI) onSignedInAsClicked() {
	if _, err := ui.auth.GetUser(); err == nil {
		open.Run("fastnode://settings")
	} else {
		open.Run("fastnode://login")
	}
}

func (ui *UI) onFeedbackClicked() {
	open.Run("fastnode://feedback")
}

func terminate() {
	os.Exit(0)
}
