package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Vec struct{ X, Y, Z float64 }

var vecRe = regexp.MustCompile(`\{([^\}]+)\}`) // matches "{x, y, z}"

func parseVec(s string) (Vec, error) {
	parts := strings.Split(strings.TrimSpace(s), ",")
	if len(parts) != 3 {
		return Vec{}, fmt.Errorf("unexpected vector format")
	}
	x, _ := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	y, _ := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	z, _ := strconv.ParseFloat(strings.TrimSpace(parts[2]), 64)
	return Vec{x, y, z}, nil
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("usage: %s <generated_go_file>", os.Args[0])
	}
	path := os.Args[1]
	data, _ := ioutil.ReadFile(path)

	// Extract all vectors inside the CubeVertices literal
	matches := vecRe.FindAll(data, -1)
	var verts []Vec
	for _, m := range matches {
		v, err := parseVec(string(m[1 : len(m)-1]))
		if err != nil {
			continue // ignore non‑vector lines (e.g., faces)
		}
		verts = append(verts, v)
	}
	if len(verts) == 0 {
		log.Fatal("no vertices found")
	}

	// Compute bbox centre
	min, max := verts[0], verts[0]
	for _, v := range verts[1:] {
		if v.X < min.X {
			min.X = v.X
		}
		if v.Y < min.Y {
			min.Y = v.Y
		}
		if v.Z < min.Z {
			min.Z = v.Z
		}
		if v.X > max.X {
			max.X = v.X
		}
		if v.Y > max.Y {
			max.Y = v.Y
		}
		if v.Z > max.Z {
			max.Z = v.Z
		}
	}
	center := Vec{
		X: (min.X + max.X) / 2,
		Y: (min.Y + max.Y) / 2,
		Z: (min.Z + max.Z) / 2,
	}

	// Translate all vertices
	var buf bytes.Buffer
	lastIdx := 0
	for _, loc := range vecRe.FindAllIndex(data, -1) {
		orig := string(data[loc[0]:loc[1]]) // e.g. "{0.123, 0.456, 0.789}"
		vec, _ := parseVec(orig[1 : len(orig)-1])
		newVec := Vec{
			X: vec.X - center.X,
			Y: vec.Y - center.Y,
			Z: vec.Z - center.Z,
		}
		buf.Write(data[lastIdx:loc[0]])
		buf.WriteString(fmt.Sprintf("{%.6g, %.6g, %.6g}", newVec.X, newVec.Y, newVec.Z))
		lastIdx = loc[1]
	}
	buf.Write(data[lastIdx:])

	// Re‑format the Go source (optional but nice)
	formatted, _ := format.Source(buf.Bytes())
	ioutil.WriteFile(path, formatted, 0644)
	fmt.Println("Vertices re‑centred and file overwritten.")
}
