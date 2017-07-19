package core

type Codec interface {
	Unmarshal([]byte, interface{}) error
	Marshal(interface{}) ([]byte, error)
}
