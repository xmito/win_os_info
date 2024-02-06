# Telegraf inputs plugin for Windows OS
The plugin can be used to collect Windows OS version/update information:
- Windows version and caption
- Last successful update detection date
- Last successful update download date
- Last successful update install date
- Windows update settings

Plugin settings in telegraf's toml configuration:
```toml
[[inputs.win_os_info]]
	## Plugin adds version info such as:
	# ProductName
	# CurrentBuildNumber
	# CurrentVersion
	OsVersion = true

	## Plugin adds update history information such as
	# LastSuccessDetectDate
	# LastSuccessDownloadDate
	# LastSuccessInstallDate
	# TimeZoneKeyName
	UpdateStatus = true

	## Plugin adds update settings such as:
	# AUOptions
	# IncludeRecommendedUpdates
	# ElevateNonAdmins
	# NextDetectionTime
	UpdateSettings = true

	# Configure interval to collect once an hour
	interval = "3600s"
```

After compiling telegraf, you can try to run with the plugin:
> .\telegraf.exe --config .\example.conf --test