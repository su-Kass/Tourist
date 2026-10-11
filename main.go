package main

import (
	"fmt"
	"path/filepath"
)

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

	var rungammaCMD string = commandArgs + " wine " + "\""+filepath.Join(gammaPrefixDir, "ModOrganizer.exe")+"\""
	if termFlag == "-i" || termFlag == "--install" {
		TalkCat("You have run Tourist with the --install flag!")
		TalkCat("This installer has programs that will open")
		TalkCat("When a program opens, make sure to check your terminal for instructions\n")

		TalkCat("This is a disclaimer for you to not get confused!")
		TalkCat("Good luck!")
		MakeOutput("3", "Write \"OK\" to continue to installation:\n")
		var input string
		fmt.Scanln(&input)
		if input != "OK" {
			TalkCat("incorrect input!")
			return
		}
		TalkCat("thank you! :3")
		TalkCat("finding anomaly zip")
		var anomalyZipName, err1 = SearchForFiletype(anomalyDir, "7z") // file for anomaly
		var anomalyZip string = filepath.Join(anomalyDir, anomalyZipName)
		ErrHandle(err1, false)
		if err1 != nil {
			TalkCat("You forgot to put the files for Anomaly into the folders.")
			return
		}
		TalkCat("finding gamma zip")
		var gammaZipName, err2 = SearchForFiletype(gammaDir, "7z") // file for gamma if exists
		var gammaZip string = filepath.Join(gammaDir, gammaZipName)
		var addGAMMA bool = true
		ErrHandle(err2, false)
		if err2 != nil {
			addGAMMA = false
			TalkCat("gamma not located, skipping install")
		}
		TalkCat("searching for modded exes")
		var exeDir = filepath.Join(data.CurDir, "ModdedExecutables")
		var exeZipName, err3 = SearchForFiletype(exeDir, "zip")
		var addExecutables bool = true
		ErrHandle(err3, false)
		if err3 != nil {
			TalkCat("modded exes not located, skipping install")
			addExecutables = false
		}
		TalkCat("setting up the wine prefix!")
		RunShell(commandArgs + "WINEDLLOVERRIDES=\"mscoree,mshtml=\" wineboot") // create prefix
		// winetricks d3dcompiler_43 d3dcompiler_47 d3dx11_43 d3dx9 d3dx10 dxvk vcrun2022 << alternative command
		TalkCat("installing the files needed for running S.T.A.L.K.E.R under wine")
		RunShell(commandArgs + " winetricks -q d3dcompiler_47 d3dx9 d3dx11_43 vcrun2022 dxvk") // setup

		TalkCat("making Anomaly folder")
		MakeDir(anomalyPrefixDir)

		TalkCat("moving Anomaly to the prefix")
		MoveFile(anomalyZip, filepath.Join(anomalyPrefixDir, anomalyZipName))

		anomalyZip = filepath.Join(anomalyPrefixDir, anomalyZipName)

		TalkCat("unzipping anomaly, pls wait")
		TalkCat("when it is done installing then Anomaly will start")
		TalkCat("close it when it starts, it will be the DX8 version...")
		RunShell("7z x '" + anomalyZip + "' -o'" + anomalyPrefixDir + "' -y")

		TalkCat("warming up anomaly")
		RunShell(commandArgs + " wine " + filepath.Join(anomalyPrefixDir, "bin", "AnomalyDX8.exe"))

		if addExecutables {
			TalkCat("adding the custom executables")
			var exeZip = filepath.Join(exeDir, exeZipName)
			MoveFile(exeZip, filepath.Join(anomalyPrefixDir, exeZipName))
			exeZip = filepath.Join(anomalyPrefixDir, exeZipName)
			TalkCat("unzipping")
			RunShell("unzip -o '" + exeZip + "' -d \"" + anomalyPrefixDir + "\"")
		}

		if addGAMMA {
			TalkCat("adding GAMMA")

			MakeDir(gammaPrefixDir)
			TalkCat("moving GAMMA to the prefix")
			MoveFile(gammaZip, filepath.Join(gammaPrefixDir, gammaZipName))
			gammaZip = filepath.Join(gammaPrefixDir, gammaZipName)
			TalkCat("extracting GAMMA directly")
			MakeOutput("1","\n")
			TalkCat("once gamma installs it will open ModOrganizer so that you can point to the Anomaly directory")
			TalkCat("close Mod Organizer once you point it to Anomaly")
			RunShell("7z x '" + gammaZip + "' -o'" + gammaPrefixDir + "' -y")

			TalkCat("MO started")
			RunShell(rungammaCMD)
			TalkCat("MO closed")

			var gammaLauncher string = filepath.Join(data.CurDir, "gamma-launcher")
			TalkCat("downloading gamma-launcher from github")
			RunShell("curl -fL -o "+gammaLauncher+" \"https://github.com/Mord3rca/gamma-launcher/releases/download/v3.1/gamma-launcher\"")
			TalkCat("adding gamma to Anomaly using gamma-launcher")
			RunShell(gammaLauncher+" full-install --anomaly "+"'"+anomalyPrefixDir+"'"+" --gamma "+"'"+gammaPrefixDir+"'")
			TalkCat("verifying md5")
			RunShell(gammaLauncher+" check-md5 --gamma "+"'"+gammaPrefixDir+"'")
		}
		TalkCat("Welcome to the zone stalker!")
		TalkCat("run with '-ra' or '--runanomaly' to run anomaly")
		MakeOutput("1", "\n")
		if addGAMMA {
			TalkCat("GAMMA has been installed as well!")
			TalkCat("run with '-rg' or '--rungamma' to run gamma")
		}
	}
	if termFlag == "-ra" || termFlag == "--runanomaly" {
		RunShell(commandArgs + " wine " + filepath.Join(anomalyPrefixDir, "AnomalyLauncher.exe"))
	}
	if termFlag == "-rg" || termFlag == "--rungamma" {
		MakeOutput("2", rungammaCMD)
		RunShell(rungammaCMD)
	}
}
