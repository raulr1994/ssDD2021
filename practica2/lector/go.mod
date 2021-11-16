module example.com/lector

go 1.17

replace example.com/RA => ../ra

replace example.ms/ms => ../ms

replace example.gestorfichero/gestorfichero => ../gestorfichero

require (
	example.com/RA v0.0.0-00010101000000-000000000000
	example.gestorfichero/gestorfichero v0.0.0-00010101000000-000000000000
	github.com/DistributedClocks/GoVector v0.0.0-20210402100930-db949c81a0af
)

require (
	example.ms/ms v0.0.0-00010101000000-000000000000 // indirect
	github.com/daviddengcn/go-colortext v1.0.0 // indirect
	github.com/vmihailenco/msgpack/v5 v5.1.4 // indirect
	github.com/vmihailenco/tagparser v0.1.2 // indirect
)
