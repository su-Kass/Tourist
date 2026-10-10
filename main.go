package main

import "path/filepath"

//	type UserData struct {
//		HmeDir	     string
//		CurDir	     string
//		TmpDir	     string
//		AppDir	     string
//		SysType	     string
//		CpuNum	     string
//		Username     string
//	}
func main() {
	var termFlag, _ = ReturnArguments()
	var data UserData = ReturnData()

	var prefixDir string = filepath.Join(data.CurDir, "prefix", "stalker-anomaly")
	var commandArgs = "WINEPREFIX=" + "'" + prefixDir + "'" + " WINEFSYNC=1 WINEDEBUG=-all"
	var anomalyPrefixDir string = filepath.Join(prefixDir, "drive_c", "Anomaly")
	var gammaPrefixDir string = filepath.Join(prefixDir, "drive_c", "GAMMA")

	var anomalyDir = filepath.Join(data.CurDir, "Anomaly")
	var gammaDir = filepath.Join(data.CurDir, "GAMMA")

	if termFlag == "-i" {
		var exeDir = filepath.Join(data.CurDir, "ModdedExecutables")

		var anomalyZipName, err1 = SearchForFiletype(anomalyDir, "7z") // file for anomaly
		var anomalyZip string = filepath.Join(anomalyDir, anomalyZipName)
		ErrHandle(err1, true)

		var gammaZipName, err2 = SearchForFiletype(gammaDir, "7z") // file for gamma if exists
		var gammaZip string = filepath.Join(gammaDir, gammaZipName)
		var addGAMMA bool = true
		ErrHandle(err2, false)
		if err2 != nil {
			addGAMMA = false
		}

		var exeZipName, err3 = SearchForFiletype(exeDir, "zip")
		ErrHandle(err3, true)
		var exeZip = filepath.Join(exeDir, exeZipName)

		RunShell(commandArgs + "WINEDLLOVERRIDES=\"mscoree,mshtml=\" wineboot") // create prefix
		// winetricks d3dcompiler_43 d3dcompiler_47 d3dx11_43 d3dx9 d3dx10 dxvk vcrun2022 << alternative command
		RunShell(commandArgs + " winetricks -q d3dcompiler_47 d3dx9 d3dx11_43 vcrun2022 dxvk") // setup

		MakeDir(anomalyPrefixDir)

		MoveFile(anomalyZip, filepath.Join(anomalyPrefixDir, anomalyZipName))
		MoveFile(exeZip, filepath.Join(anomalyPrefixDir, exeZipName))

		exeZip = filepath.Join(anomalyPrefixDir, exeZipName)
		anomalyZip = filepath.Join(anomalyPrefixDir, anomalyZipName)

		RunShell("7z x '" + anomalyZip + "' -o'" + anomalyPrefixDir + "' -y")
		RunShell("unzip -o '" + exeZip + "' -d \"" + anomalyPrefixDir + "\"")

		RunShell(commandArgs + " wine " + filepath.Join(anomalyPrefixDir, "bin", "AnomalyDX8.exe"))

		if addGAMMA {
			MakeOutput("2", "Adding GAMMA!")

			MakeDir(gammaPrefixDir)
			MoveFile(gammaZip, filepath.Join(gammaPrefixDir, gammaZipName))
			gammaZip = filepath.Join(gammaPrefixDir, gammaZipName)
			RunShell("7z x '" + gammaZip + "' -o'" + gammaPrefixDir + "' -y")

			var gammaLauncher string = filepath.Join(data.CurDir, "gamma-launcher")
			RunShell("curl -fL -o "+gammaLauncher+" \"https://github.com/Mord3rca/gamma-launcher/releases/download/v3.1/gamma-launcher\"")
			RunShell(gammaLauncher+" full-install --anomaly "+"'"+anomalyPrefixDir+"'"+" --gamma "+"'"+gammaPrefixDir+"'")
			RunShell(gammaLauncher+" check-md5 --gamma "+"'"+gammaPrefixDir+"'")
		}
	}
	if termFlag == "-ra" {
		RunShell(commandArgs + " wine " + filepath.Join(anomalyPrefixDir, "AnomalyLauncher.exe"))
	}
	if termFlag == "-rg" {
		var runCMD string = commandArgs + " wine " + "\""+filepath.Join(gammaPrefixDir, "ModOrganizer.exe")+"\""
		MakeOutput("2", runCMD)
		RunShell(runCMD)
	}
}
