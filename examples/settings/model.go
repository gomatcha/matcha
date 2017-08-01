package settings

import (
	"gomatcha.io/matcha/comm"
	"gomatcha.io/matcha/view/stackview"
)

type App struct {
	Stack     *stackview.Stack
	Wifi      *Wifi
	Bluetooth *Bluetooth

	comm.Relay
	airplaneMode bool
}

func NewApp() *App {
	return &App{
		Stack:        &stackview.Stack{},
		Wifi:         NewWifi(),
		Bluetooth:    NewBluetooth(),
		airplaneMode: false,
	}
}

func (st *App) AirplaneMode() bool {
	return st.airplaneMode
}

func (st *App) SetAirplaneMode(v bool) {
	st.airplaneMode = v

	st.Wifi.SetEnabled(!st.airplaneMode)
	st.Bluetooth.SetEnabled(!st.airplaneMode)

	st.Signal()
}

type Wifi struct {
	comm.Relay
	askToJoin   bool
	enabled     bool
	currentSSID string
	networks    []*WifiNetwork
}

func NewWifi() *Wifi {
	n1 := NewWifiNetwork("XfinityWifi")
	n2 := NewWifiNetwork("Bluestone")
	n3 := NewWifiNetwork("Starbucks")
	n4 := NewWifiNetwork("FastMesh Wifi")

	s := &Wifi{}
	s.SetEnabled(true)
	s.SetCurrentSSID(n4.SSID())
	s.SetNetworks([]*WifiNetwork{n1, n2, n3, n4})
	return s
}

func (s *Wifi) CurrentSSID() string {
	return s.currentSSID
}

func (s *Wifi) SetCurrentSSID(v string) {
	s.currentSSID = v
	s.Signal()
}

func (s *Wifi) Enabled() bool {
	return s.enabled
}

func (s *Wifi) SetEnabled(v bool) {
	s.enabled = v
	s.Signal()
}

func (s *Wifi) AskToJoin() bool {
	return s.askToJoin
}

func (s *Wifi) SetAskToJoin(v bool) {
	s.askToJoin = v
	s.Signal()
}

func (s *Wifi) Networks() []*WifiNetwork {
	return s.networks
}

func (s *Wifi) SetNetworks(n []*WifiNetwork) {
	s.networks = n
	s.Signal()
}

type WifiNetwork struct {
	comm.Relay
	props WifiNetworkProperties
}

func NewWifiNetwork(ssid string) *WifiNetwork {
	return &WifiNetwork{
		props: WifiNetworkProperties{
			SSID:          ssid,
			IPAddress:     "10.0.1.25",
			SubnetMask:    "255.255.255.0",
			Router:        "10.0.1.1",
			DNS:           "10.0.1.1",
			SearchDomains: "hsd1.or.comcast.net.",
			ClientID:      "",
		},
	}
}

func (n *WifiNetwork) SSID() string {
	return n.props.SSID
}

func (n *WifiNetwork) Properties() WifiNetworkProperties {
	return n.props
}

func (n *WifiNetwork) SetProperties(v WifiNetworkProperties) {
	n.props = v
	n.Signal()
}

type WifiNetworkProperties struct {
	SSID   string
	Locked bool
	Signal int

	Kind          int
	IPAddress     string
	SubnetMask    string
	Router        string
	DNS           string
	SearchDomains string
	ClientID      string
	Proxy         int
}

type Bluetooth struct {
	comm.Relay
	enabled bool
	devices []*BluetoothDevice
}

func NewBluetooth() *Bluetooth {
	n1 := NewBluetoothDevice("JBL Charge 3")
	n2 := NewBluetoothDevice("Kevin's AirPods")
	n3 := NewBluetoothDevice("Kevin's Apple Watch")
	n4 := NewBluetoothDevice("Honda Fit")

	return &Bluetooth{
		enabled: true,
		devices: []*BluetoothDevice{n1, n2, n3, n4},
	}
}

func (s *Bluetooth) Enabled() bool {
	return s.enabled
}

func (s *Bluetooth) SetEnabled(v bool) {
	s.enabled = v
	s.Signal()
}

func (s *Bluetooth) Devices() []*BluetoothDevice {
	return s.devices
}

func (s *Bluetooth) SetDevices(v []*BluetoothDevice) {
	s.devices = v
	s.Signal()
}

type BluetoothDevice struct {
	comm.Relay
	ssid      string
	connected bool
}

func NewBluetoothDevice(ssid string) *BluetoothDevice {
	return &BluetoothDevice{
		ssid:      ssid,
		connected: true,
	}
}

func (s *BluetoothDevice) SSID() string {
	return s.ssid
}

func (s *BluetoothDevice) Connected() bool {
	return s.connected
}

func (s *BluetoothDevice) SetConnected(v bool) {
	s.connected = v
	s.Signal()
}
