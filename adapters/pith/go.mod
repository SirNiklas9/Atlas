module github.com/SirNiklas9/atlas/adapters/pith

go 1.25.6

require (
	github.com/SirNiklas9/atlas/engine v0.0.0
	pith v0.0.0
)

require github.com/odvcencio/gotreesitter v0.20.2 // indirect

replace github.com/SirNiklas9/atlas/engine => ../../engine

replace pith => ../../../pith
