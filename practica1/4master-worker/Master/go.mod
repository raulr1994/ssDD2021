module example.com/server

go 1.17

replace example.com/com => ../com

require (
	example.com/com v0.0.0-00010101000000-000000000000
	example.ssh/ssh v0.0.0-00010101000000-000000000000
)

require (
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
	golang.org/x/sys v0.0.0-20210615035016-665e8c7367d1 // indirect
)

replace example.ssh/ssh => ../ssh
