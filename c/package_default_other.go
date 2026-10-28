//go:build !darwin && !windows && !(linux && (amd64 || arm64) && !android)

package c

const LLGoPackage = "decl"
