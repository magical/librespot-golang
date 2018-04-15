// +build android

package librespot

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/librespot-org/librespot-golang/src/Spotify"
)

type Updater interface {
	OnUpdate(device string)
}

func LoginConnection(username string, password string, deviceName string, con io.ReadWriter) (*SpircController, error) {
	s := &session{
		keys:               generateKeys(),
		tcpCon:             con,
		mercuryConstructor: setupMercury,
		shannonConstructor: setupStream,
	}
	s.deviceId = generateDeviceId(deviceName)
	s.deviceName = deviceName

	s.startConnection()
	loginPacket := loginPacketPassword(username, password, s.deviceId)
	return s.doLogin(loginPacket, username)
}

func LoginConnectionSaved(username string, authData []byte, deviceName string, con io.ReadWriter) (*SpircController, error) {
	s := &session{
		keys:               generateKeys(),
		tcpCon:             con,
		mercuryConstructor: setupMercury,
		shannonConstructor: setupStream,
	}
	s.deviceId = generateDeviceId(deviceName)
	s.deviceName = deviceName

	s.startConnection()
	packet := loginPacket(username, authData,
		Spotify.AuthenticationType_AUTHENTICATION_STORED_SPOTIFY_CREDENTIALS.Enum(), s.deviceId)
	return s.doLogin(packet, username)
}

func (c *SpircController) HandleUpdatesCb(cb func(device string)) {
	c.updateChan = make(chan Spotify.Frame, 5)

	go func() {
		for {
			update := <-c.updateChan
			jsonData, err := json.Marshal(update)
			if err != nil {
				fmt.Println("Error marhsaling device json")
			} else {
				cb(string(jsonData))
			}
		}
	}()
}

func (c *SpircController) HandleUpdates(u Updater) {
	c.updateChan = make(chan Spotify.Frame, 5)

	go func() {
		for {
			update := <-c.updateChan
			jsonData, err := json.Marshal(update)
			if err != nil {
				fmt.Println("Error marhsaling device json")
			} else {
				u.OnUpdate(string(jsonData))
			}
		}
	}()
}

func (c *SpircController) ListDevicesJson() (string, error) {
	devices := c.ListDevices()
	jsonData, err := json.Marshal(devices)
	if err != nil {
		return "", nil
	}
	return string(jsonData), nil
}

func (c *SpircController) ListMdnsDevicesJson() (string, error) {
	devices, err := c.ListMdnsDevices()
	if err != nil {
		return "", nil
	}
	jsonData, err := json.Marshal(devices)
	if err != nil {
		return "", nil
	}
	return string(jsonData), nil
}

func (c *SpircController) SuggestJson(term string) (string, error) {
	result, err := c.Suggest(term)
	if err != nil {
		return "", nil
	}
	jsonData, err := json.Marshal(result)
	if err != nil {
		return "", nil
	}
	return string(jsonData), nil
}

func (c *SpircController) SearchJson(term string) (string, error) {
	result, err := c.Search(term)
	if err != nil {
		return "", nil
	}
	jsonData, err := json.Marshal(result)
	if err != nil {
		return "", nil
	}
	return string(jsonData), nil
}
