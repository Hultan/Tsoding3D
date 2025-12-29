package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// -------------------------------------------------------------------
// Types that will appear in the generated Go code
// -------------------------------------------------------------------
type Vector3 struct {
	X, Y, Z float64
}

// -------------------------------------------------------------------
// Helper: convert a line like "v 0.25 0.25 0.25" into a Vector3
// -------------------------------------------------------------------
func parseVertex(fields []string) (Vector3, error) {
	if len(fields) < 4 {
		return Vector3{}, fmt.Errorf("not enough components for vertex")
	}
	x, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return Vector3{}, err
	}
	y, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return Vector3{}, err
	}
	z, err := strconv.ParseFloat(fields[3], 64)
	if err != nil {
		return Vector3{}, err
	}
	return Vector3{X: x, Y: y, Z: z}, nil
}

// -------------------------------------------------------------------
// Helper: convert a face definition like "f 1 2 3 4"
//
//	OBJ indices start at 1, we convert them to 0‑based int32.
//	The face may be a triangle, quad, or any polygon.
//
// -------------------------------------------------------------------
func parseFace(fields []string) ([]int32, error) {
	if len(fields) < 4 { // at least "f i j k"
		return nil, fmt.Errorf("face line has too few vertices")
	}
	indices := make([]int32, 0, len(fields)-1)
	for _, f := range fields[1:] {
		// Faces can be written as "v", "v/vt", "v//vn", or "v/vt/vn".
		// We only care about the vertex index before the first '/'.
		parts := strings.Split(f, "/")
		idx, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, err
		}
		// OBJ indices are 1‑based → subtract 1 for Go slices.
		indices = append(indices, int32(idx-1))
	}
	return indices, nil
}

// -------------------------------------------------------------------
// Main conversion routine
// -------------------------------------------------------------------
func convertOBJ(path string) (vertices []Vector3, faces [][]int32, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue // skip empty lines and comments
		}
		fields := strings.Fields(line)
		switch fields[0] {
		case "v":
			v, err := parseVertex(fields)
			if err != nil {
				return nil, nil, fmt.Errorf("line %d: %w", lineNum, err)
			}
			vertices = append(vertices, v)
		case "f":
			f, err := parseFace(fields)
			if err != nil {
				return nil, nil, fmt.Errorf("line %d: %w", lineNum, err)
			}
			faces = append(faces, f)
		// All other prefixes (vt, vn, g, o, s, usemtl, ...) are ignored.
		default:
			// no‑op
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}
	return vertices, faces, nil
}

// -------------------------------------------------------------------
// Pretty‑print the slices as Go source code
// -------------------------------------------------------------------
func emitGoCode(name string, vertices []Vector3, faces [][]int32) string {
	var sb strings.Builder

	// Header
	fmt.Fprintf(&sb, "package %s\n\n", name)

	// Vector type (exported so you can reuse it)
	sb.WriteString(`type Vector3 struct {
    X, Y, Z float64
}

`)

	// Vertices slice
	sb.WriteString("var CubeVertices = []Vector3{\n")
	for _, v := range vertices {
		fmt.Fprintf(&sb, "    {%.6g, %.6g, %.6g},\n", v.X, v.Y, v.Z)
	}
	sb.WriteString("}\n\n")

	// Faces slice
	sb.WriteString("var CubeFaces = [][]int32{\n")
	for _, f := range faces {
		sb.WriteString("    {")
		for i, idx := range f {
			if i > 0 {
				sb.WriteString(", ")
			}
			fmt.Fprintf(&sb, "%d", idx)
		}
		sb.WriteString("},\n")
	}
	sb.WriteString("}\n")

	return sb.String()
}

// -------------------------------------------------------------------
// Entry point
// -------------------------------------------------------------------
func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: obj2go <path-to-obj>")
		os.Exit(1)
	}
	objPath := os.Args[1]

	verts, faces, err := convertOBJ(objPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error converting OBJ: %v\n", err)
		os.Exit(1)
	}

	// Emit Go source to stdout – you can redirect it to a file.
	fmt.Print(emitGoCode("model", verts, faces))
}
