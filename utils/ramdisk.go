package utils

import (
	"os"
	"fmt"
	"flag"
	"strings"
	"os/exec"
)

var ramdiskSizeInGigabytes int

func init() {
	flag.IntVar(&ramdiskSizeInGigabytes, "size", 4, "ramdisk size in gigabytes")
}

// https://bogner.sh/2012/12/os-x-create-a-ram-disk-the-easy-way/
func convertGigabytesToSectors(gigabytes int) int {
	return gigabytes * 1024 * 1024 * 1024 / 512
}

func getRamdiskBackupDir() string {
	ramdiskBackupDir := os.Getenv("RAMDISK_BACKUP_DIR")
	if ramdiskBackupDir == "" {
		fmt.Println("Must set RAMDISK_BACKUP_DIR")
		os.Exit(1)
	}
	fmt.Println("RAMDISK_BACKUP_DIR:", ramdiskBackupDir)
	return ramdiskBackupDir
}

func startRamdisk() {
	sectors := convertGigabytesToSectors(ramdiskSizeInGigabytes)
	device := verboseRamdisk("hdiutil", "attach", "-nomount", fmt.Sprintf("ram://%d", sectors))

	verboseRamdisk("diskutil", "erasevolume", "HFS+", "ramdisk", device)

	ramdiskBackupDir := getRamdiskBackupDir()
	verboseRamdisk("rsync", "-av", ramdiskBackupDir+"/", "/Volumes/ramdisk")
}

func stopRamdisk() {
	syncRamdisk()
	verboseRamdisk("diskutil", "eject", "ramdisk")
}

func getRamdiskStatus() {
	result := verboseRamdisk("diskutil", "list", "ramdisk")
	if strings.Contains(result, "Could not find") {
		fmt.Println("ramdisk stopped")
	} else {
		fmt.Println("ramdisk started")
	}
}

func syncRamdisk() {
	ramdiskBackupDir := getRamdiskBackupDir()
	verboseRamdisk("rsync", "-av", "/Volumes/ramdisk/", ramdiskBackupDir)
}

func verboseRamdisk(args ...string) string {
	fmt.Printf("==> Executing: %s\n", strings.Join(args, " "))
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("==> Error: %s\n", err.Error()))
	}

	if len(output) > 0 {
		fmt.Printf("==> Output: %s\n", string(output))
	}

	return strings.TrimSpace(string(output))
}

// http://www.observium.org/wiki/Persistent_RAM_disk_RRD_storage
func ramdiskMain() {

	flag.Usage = func() {
		fmt.Printf("Usage: ramdisk [options] {start|sync|status|stop}>\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	flag.Parse()
	if len(flag.Args()) == 0 {
		flag.Usage()
	}

	command := flag.Args()[0]

	switch command {
	case "start":
		startRamdisk()
	case "stop":
		stopRamdisk()
	case "status":
		getRamdiskStatus()
	case "sync":
		syncRamdisk()
	default:
		flag.Usage()
	}

}