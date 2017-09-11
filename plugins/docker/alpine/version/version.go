package version

import "fmt"

// VERSION indicates which version of the binary is running.
var SNK_PLUGIN_APK_FILE_VERSION string

// GITCOMMIT indicates which git hash the binary was built off of
var SNK_PLUGIN_APK_FILE_VCS_GIT_COMMIT string

// SnkPluginApkFiledVersion is the version of the build
var SnkPluginApkFiledVersion = "undefined"

// SnkPluginApkFiledHeaderValue is the value of the custom SnkPluginApkFileD header
var SnkPluginApkFiledHeaderValue = fmt.Sprintf("Version %s", SnkPluginApkFiledVersion)

// SnkPluginApkFiledUserAgent is the value of the user agent header sent to the backends
var SnkPluginApkFiledUserAgent = fmt.Sprintf("SnkPluginApkFileD Version %s", SnkPluginApkFiledVersion)