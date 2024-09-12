package clientapp

import "github.com/khulnasoft-lab/fastnode/fastnode-go/client/internal/windowsui"

// launchOnboarding runs FastnodeOnboarding.exe
func launchOnboarding() error {
	return windowsui.RunOnboarding()
}
