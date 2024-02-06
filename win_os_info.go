package win_os_info

import "log"
import(
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
	"golang.org/x/sys/windows/registry"
)



var registryKeyString = map[registry.Key]string {
	registry.CLASSES_ROOT: "HKEY_CLASSES_ROOT",
	registry.CURRENT_USER: "HKEY_CURRENT_USER",
	registry.LOCAL_MACHINE: "HKEY_LOCAL_MACHINE",
	registry.USERS: "HKEY_USERS",
	registry.CURRENT_CONFIG: "HKEY_CURRENT_CONFIG",
	registry.PERFORMANCE_DATA: "HKEY_PERFORMANCE_DATA",
}

//************** WinOsVersion ***************
type WinOsVersion struct {
	ProductName string
	CurrentBuildNumber string
	CurrentVersion string
}

func GetWinOsVersion() (WinOsVersion, error) {
	versionKey := `SOFTWARE\Microsoft\Windows NT\CurrentVersion`
	var k registry.Key
	var win WinOsVersion
	var err error

	if k, err = registry.OpenKey(registry.LOCAL_MACHINE,
		versionKey, registry.QUERY_VALUE); err != nil {
		log.Printf("Error while accessing registry: %s\n%v",
			registryKeyString[registry.LOCAL_MACHINE] + `\` + versionKey, err)
		return win, err
	}
	if win.ProductName, _, err = k.GetStringValue("ProductName"); err != nil {
		log.Printf("Error while obtaining key value: %s\n%v", "ProductName", err)
		return win, err
	}
	if win.CurrentBuildNumber, _, err = k.GetStringValue("CurrentBuildNumber"); err != nil {
		log.Printf("Error while obtaining key value: %s\n%v", "CurrentBuildNumber", err)
		return win, err
	}
	if win.CurrentVersion, _, err = k.GetStringValue("CurrentVersion"); err != nil {
		log.Printf("Error while obtaining key value: %s\n%v", "CurrentVersion", err)
		return win, err
	}

	return win, nil
}

func (w *WinOsVersion) AddFields(acc telegraf.Accumulator,
				 measurement string,
				 tags map[string]string) {
	fields := make(map[string]interface{})
	fields["ProductName"] = w.ProductName
	fields["CurrentBuildNumber"] = w.CurrentBuildNumber
	fields["CurrentVersion"] = w.CurrentVersion

	acc.AddFields(measurement, fields, tags)
}


//************** WinOsResults *****************
type WinOsUpdateStatus struct {
	LastSuccessDetectDate string
	LastSuccessDownloadDate string
	LastSuccessInstallDate string
	TimeZoneKeyName string
}

func GetRegKeyString(rootkey registry.Key, regkey string, access uint32, value string) (string, error) {
	var k registry.Key
	var err error
	var ret string

	if k, err = registry.OpenKey(rootkey, regkey, access); err != nil {
		log.Printf("Error while accessing regkey: %s\n%v", registryKeyString[rootkey] + `\` + regkey, err)
		return "", err
	}
	if ret, _, err = k.GetStringValue(value); err != nil {
		log.Printf("Error while obtaining key value: %s\n%v", value, err)
		return "", err
	}

	return ret, nil
}

func GetWinOsUpdateStatus() (WinOsUpdateStatus, error) {
	installKey := `SOFTWARE\Microsoft\Windows\CurrentVersion\WindowsUpdate\Auto Update\Results\Install`
	downloadKey := `SOFTWARE\Microsoft\Windows\CurrentVersion\WindowsUpdate\Auto Update\Results\Download`
	detectKey := `SOFTWARE\Microsoft\Windows\CurrentVersion\WindowsUpdate\Auto Update\Results\Detect`
	timezonekey := `SYSTEM\CurrentControlSet\Control\TimeZoneInformation`
	var win WinOsUpdateStatus
	var err error

	if win.LastSuccessInstallDate, err = GetRegKeyString(registry.LOCAL_MACHINE,
		installKey, registry.QUERY_VALUE, "LastSuccessTime"); err != nil {
		return win, err
	}

	if win.LastSuccessDownloadDate, err = GetRegKeyString(registry.LOCAL_MACHINE,
		downloadKey, registry.QUERY_VALUE, "LastSuccessTime"); err != nil {
		return win, err
	}

	if win.LastSuccessDetectDate, err = GetRegKeyString(registry.LOCAL_MACHINE,
		detectKey, registry.QUERY_VALUE, "LastSuccessTime"); err != nil {
		return win, err
	}

	if win.TimeZoneKeyName, err = GetRegKeyString(registry.LOCAL_MACHINE,
		timezonekey, registry.QUERY_VALUE, "TimeZoneKeyName"); err != nil {
		return win, err
	}

	return win, nil
}

func (w *WinOsUpdateStatus) AddFields(acc telegraf.Accumulator,
				 measurement string,
				 tags map[string]string) {
	fields := make(map[string]interface{})
	fields["LastSuccessDetectDate"] = w.LastSuccessDetectDate
	fields["LastSuccessDownloadDate"] = w.LastSuccessDownloadDate
	fields["LastSuccessInstallDate"] = w.LastSuccessInstallDate
	fields["TimeZoneKeyName"] = w.TimeZoneKeyName

	acc.AddFields(measurement, fields, tags)
}

//*************** WinOsSettings ***************
type WinOsUpdateSettings struct {
	AUOptions uint64
	IncludeRecommendedUpdates uint64
	ElevateNonAdmins uint64
	NextDetectionTime string
}

func GetWinOsUpdateSettings() (WinOsUpdateSettings, error) {
	autoUpdateKey := `SOFTWARE\Microsoft\Windows\CurrentVersion\WindowsUpdate\Auto Update`
	var win WinOsUpdateSettings
	var k registry.Key
	var err error

	if k, err = registry.OpenKey(registry.LOCAL_MACHINE,
		autoUpdateKey, registry.QUERY_VALUE); err != nil {
		log.Printf("Error while accessing regkey: %s\n%v", registryKeyString[registry.LOCAL_MACHINE] + `\` + autoUpdateKey, err)
	}
	if win.AUOptions, _, err = k.GetIntegerValue("AUOptions");
		err != nil {
		log.Printf("Error while obtaining key value: %s\n%v", "AUOptions", err)
		return win, err
	}
	if win.IncludeRecommendedUpdates, _, err = k.GetIntegerValue("IncludeRecommendedUpdates");
		err != nil {
		log.Printf("Error while obtaining key value: %s\n%v", "IncludeRecommendedUpdates", err)
		return win, err
	}
	if win.ElevateNonAdmins, _, err = k.GetIntegerValue("ElevateNonAdmins");
		err != nil {
		log.Printf("Error while obtaining key value: %s\n%v", "ElevateNonAdmins", err)
		return win, err
	}
	if win.NextDetectionTime, _, err = k.GetStringValue("NextDetectionTime");
		err != nil {
		log.Printf("Error while obtaining key value: %s\n%v", "NextDetectionTime", err)
		return win, err
	}

	return win, nil
}

func (w *WinOsUpdateSettings) AddFields(acc telegraf.Accumulator,
				  measurement string,
				  tags map[string]string) {
	fields := make(map[string]interface{})
	fields["AUOptions"] = w.AUOptions
	fields["IncludeRecommendedUpdates"] = w.IncludeRecommendedUpdates
	fields["ElevateNonAdmins"] = w.ElevateNonAdmins
	fields["NextDetectionTime"] = w.NextDetectionTime

	acc.AddFields(measurement, fields, tags)
}

//**************** WinOsInfo ****************
type WinOsInfo struct {
	OsVersion bool
	UpdateStatus bool
	UpdateSettings bool
}

var WinOsInfoConfig = `
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
`

func (w *WinOsInfo) SampleConfig() string {
	return WinOsInfoConfig
}

func (w *WinOsInfo) Description() string {
	return `##Sends information about Windows system:
		# Windows version and caption
		# Last successful update detection date
		# Last successful update download date
		# Last successful update install date
		# Windows update settings`
}

func (w *WinOsInfo) Gather(acc telegraf.Accumulator) error {
	tags := make(map[string]string)
	if w.OsVersion {
		if version, err := GetWinOsVersion(); err != nil {
			return err
		} else {
			version.AddFields(acc, "win_os_version", tags)
		}
	}

	if w.UpdateStatus {
		if status, err := GetWinOsUpdateStatus(); err != nil {
			return err
		} else {
			status.AddFields(acc, "win_os_update_settings", tags)
		}
	}

	if w.UpdateSettings {
		if settings, err := GetWinOsUpdateSettings(); err != nil {
			return err
		} else {
			settings.AddFields(acc, "win_os_update_settings", tags)
		}
	}

	return nil
}

func init() {
	inputs.Add("win_os_info", func() telegraf.Input{ return &WinOsInfo{} })
}
