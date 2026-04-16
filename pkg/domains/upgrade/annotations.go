// Package upgrade provides upgrade/migration analysis for code property graphs.
package upgrade

const (
	AnnotVersionConversion = "upgrade:version_conversion"
	AnnotFeatureGate       = "upgrade:feature_gate"
	AnnotPreReleaseAPI     = "upgrade:pre_release_api"
	AnnotMigration         = "upgrade:migration"
	AnnotVersionCheck      = "upgrade:version_check"
)
