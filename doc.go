// Package casow provides an incremental linear constraint solver.
//
// It is a Go port of Rust kasuari, based on the Cassowary constraint solving
// algorithm. Use Solver to add required and weighted constraints, mark edit
// variables, suggest new edit values, and fetch the variables whose solved
// values changed.
package casow
