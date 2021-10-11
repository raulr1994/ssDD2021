module example.ssh/lanzar

go 1.17

replace example.ssh/ssh => ../ssh

require example.ssh/ssh v0.0.0-00010101000000-000000000000

require (
	github.com/shiena/ansicolor v0.0.0-20200904210342-c7312218db18 // indirect
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
	golang.org/x/sys v0.0.0-20210615035016-665e8c7367d1 // indirect
)
