package dto

type SensorReadingRequest struct {
	SensorType string `json:"sensor_type"`
	RoomId     string `json:"room_id"`
	SensorId   string `json:"sensor_id"`
	Value      string `json:"value"`
}
