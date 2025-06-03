package register

import (
	"log/slog"

	"golang.org/x/sys/windows/registry"
)

func isRegistered(urlScheme, locauthPath string) (bool, error) {
	k, err := registry.OpenKey(registry.CLASSES_ROOT, urlScheme+`\shell\open\command`, registry.QUERY_VALUE)
	if err != nil {
		return false, err
	}
	defer k.Close()

	val, _, err := k.GetStringValue("")

	if err != nil {
		return false, err
	}

	if val == "\""+locauthPath+"\" callback \"%1\"" {
		return true, nil
	}
	return false, nil

}

// RegMe registers app as url Handler in windows
func RegMe(urlScheme, locauthPath string) error {

	registered, err := isRegistered(urlScheme, locauthPath)
	if err != nil {
		return err
	}
	if registered {
		slog.Info("app  is registered as url handler for", "app name", locauthPath, "url scheme", urlScheme)
		return nil
	}
	slog.Info("Creating key " + urlScheme)
	k, _, err := registry.CreateKey(registry.CLASSES_ROOT, urlScheme, registry.WRITE)
	if err != nil {
		return err
	}
	defer k.Close()

	slog.Info("Setting default value")
	err = k.SetStringValue("", "URL:Handles aad tokens")
	if err != nil {
		return err
	}
	err = k.SetStringValue("URL Protocol", "")
	if err != nil {
		return err
	}

	slog.Info("Creating key", "keyname", urlScheme+`\shell`)
	sk, _, err := registry.CreateKey(registry.CLASSES_ROOT, urlScheme+`\shell`, registry.WRITE)
	if err != nil {
		return err
	}
	defer sk.Close()

	slog.Info("Creating key", "keyname", urlScheme+`\shell\open`)
	ok, _, err := registry.CreateKey(registry.CLASSES_ROOT, urlScheme+`\shell\open`, registry.WRITE)
	if err != nil {
		return err
	}
	ok.Close()

	slog.Info("Creating key", "keyname", urlScheme+`\shell\open\command`)
	ck, _, err := registry.CreateKey(registry.CLASSES_ROOT, urlScheme+`\shell\open\command`, registry.WRITE)
	if err != nil {
		return err
	}
	ck.Close()

	err = k.SetStringValue("", "\""+locauthPath+"\" callback \"%1\"")
	if err != nil {
		return err
	}
	return nil
}
