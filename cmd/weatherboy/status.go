package main

type HubStatusEvent struct{}
type DeviceStatusEvent struct{}

func HandleDeviceStatusEvent(outb []byte) (*DeviceStatusEvent, error) {
	return new(DeviceStatusEvent), nil
}

func (_ DeviceStatusEvent) String() string {
	return "DeviceStatus"
}

func HandleHubStatusEvent(outb []byte) (*HubStatusEvent, error) {
	return &HubStatusEvent{}, nil
}

func (_ HubStatusEvent) String() string {
	return "HubStatus"
}
