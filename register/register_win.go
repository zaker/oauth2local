package register

import (
	"log"

	"golang.org/x/sys/windows/registry"
)

// HKEY_CLASSES_ROOT
//    alert
//       URL Protocol = ""
// Under this new key, the URL Protocol string value indicates that this key declares a custom pluggable protocol handler. Without this key, the handler application will not launch. The value should be an empty string.

// Keys should also be added for DefaultIcon and shell. The Default string value of the DefaultIcon key must be the file name to use as an icon for this new URI scheme. The string takes the form "path, iconindex" with a maximum length of MAX_PATH. The name of the first key under the shell key should be an action verb, such as open. Under this key, a command key or a DDEEXEC key indicate how the handler should be invoked. The values under the command and DDEEXEC keys describe how to launch the application handling the new protocol.

// Finally, the Default string value should contain the display name of the new URI scheme. The following example shows how to register an application, alert.exe in this case, to handle the alert scheme.

// HKEY_CLASSES_ROOT
//    alert
//       (Default) = "URL:Alert Protocol"
//       URL Protocol = ""
//       DefaultIcon
//          (Default) = "alert.exe,1"
//       shell
//          open
//             command
//                (Default) = "C:\Program Files\Alert\alert.exe" "%1"

// Windows Registry Editor Version 5.00

// [HKEY_CLASSES_ROOT\rtsp]
// "URL Protocol"=""
// @="URL:Rtsp Streaming Protocol"
// "DefaultIcon"="\"C:\\Program Files (x86)\\VideoLAN\\VLC\\vlc.exe\",1"

// [HKEY_CLASSES_ROOT\rtsp\shell]

// [HKEY_CLASSES_ROOT\rtsp\shell\open]

// [HKEY_CLASSES_ROOT\rtsp\shell\open\command]
// @="\"C:\\Program Files (x86)\\VideoLAN\\VLC\\vlc.exe\" \"%1\""
func RegMe(urlScheme ,locauthPath string) {
	log.Println("Creating key " + urlScheme)
	k, exist, err := registry.CreateKey(registry.CLASSES_ROOT, urlScheme, registry.WRITE)
	if err != nil {
		log.Fatal(err)
	}
	defer k.Close()
	if exist {
		return
	}
	log.Println("Setting default value")
	err = k.SetStringValue("", "URL:Handles aad tokens")
	if err != nil {
		log.Fatal(err)
	}
	err = k.SetStringValue("URL Protocol", "")
	if err != nil {
		log.Fatal(err)
	}

	err = k.SetStringValue("URL Protocol", "")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Creating key", urlScheme+`\shell`)
	k, exist, err = registry.CreateKey(registry.CLASSES_ROOT, urlScheme+`\shell`, registry.WRITE)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Creating key", urlScheme+`\shell\open`)
	k, exist, err = registry.CreateKey(registry.CLASSES_ROOT, urlScheme+`\shell\open`, registry.WRITE)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Creating key", urlScheme+`\shell\open\command`)
	k, exist, err = registry.CreateKey(registry.CLASSES_ROOT, urlScheme+`\shell\open\command`, registry.WRITE)
	if err != nil {
		log.Fatal(err)
	}
	err = k.SetStringValue("", "\""+locauthPath+"\" -r \"%1\"")
	if err != nil {
		log.Fatal(err)
	}

}
