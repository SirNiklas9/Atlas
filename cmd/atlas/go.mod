module github.com/SirNiklas9/atlas/cmd/atlas

go 1.25.6

require github.com/SirNiklas9/atlas/adapters/pith v0.0.0

require (
	github.com/SirNiklas9/atlas/engine v0.0.0 // indirect
	github.com/odvcencio/gotreesitter v0.20.2 // indirect
	pith v0.0.0 // indirect
)

replace github.com/SirNiklas9/atlas/adapters/pith => ../../adapters/pith

replace github.com/SirNiklas9/atlas/engine => ../../engine

replace pith => ../../../pith
