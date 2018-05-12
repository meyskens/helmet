package template

/*
TO DO

    Release.Name: The release name
    Release.Time: The time of the release
    Release.Namespace: The namespace to be released into (if the manifest doesn’t override)
    Release.Service: The name of the releasing service (always Tiller).
    Release.Revision: The revision number of this release. It begins at 1 and is incremented for each helm upgrade.
    Release.IsUpgrade: This is set to true if the current operation is an upgrade or rollback.
    Release.IsInstall: This is set to true if the current operation is an install.

*/

// Release describes a helmet release
type Release struct {
	Name      string
	Service   string
	Namespace string
}

// NewRelease gives a new release object for a given name and namespace
func NewRelease(name, namespace string) Release {
	return Release{
		Name:      name,
		Namespace: namespace,
		Service:   "helmet",
	}
}
