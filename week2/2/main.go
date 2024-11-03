package main

type TikTak struct {
	n     int
	tikCh chan bool
	takCh chan bool
}

func NewTikTak(n int) *TikTak {
	return &TikTak{
		n:     n,
		tikCh: make(chan bool),
		takCh: make(chan bool),
	}
}

func (t *TikTak) Tik() {
	for i := 0; i < t.n; i++ {
		SayTik()        // Print "Tik"
		t.takCh <- true // Signal Tak
		<-t.tikCh       // Wait for the signal from Tak
	}
}

func (t *TikTak) Tak() {
	for i := 0; i < t.n; i++ {
		<-t.takCh       // Wait for the signal from Tik
		SayTak()        // Print "Tak"
		t.tikCh <- true // Signal Tik to proceed
	}
}
