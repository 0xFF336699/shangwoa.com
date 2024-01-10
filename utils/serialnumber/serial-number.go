package serialnumber

type SerialNumber struct {
	serial int64
	service chan chan int64
}

func (this *SerialNumber) GetNextSerialNumber() (int64) {
	ch := make(chan int64)
	this.service <- ch
	serial := <- ch
	return serial
}

func (this *SerialNumber) serve() {
	for{
		select{
		case ch := <- this.service:
			this.serial ++
			ch <- this.serial
		}
	}
}
func NewSerialNumber() *SerialNumber {
	serial := &SerialNumber{
		serial:  0,
		service: make(chan chan int64),
	}
	go serial.serve()
	return serial
}