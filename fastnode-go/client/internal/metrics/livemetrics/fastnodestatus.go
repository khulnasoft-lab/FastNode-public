package livemetrics

var (
	// fastnodestatusAllowed keeps track of all the variables registered with
	// fastnodestatus pkg that are allowed into fastnode_status
	fastnodestatusAllowed = map[string]bool{
		"spyder_suboptimal_settings": true,
	}
)
