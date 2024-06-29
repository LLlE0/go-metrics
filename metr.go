package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"math"
	"os"
	"time"
)

var (
	M  = make(map[string]int)
	Mt = make(map[token.Token]int)
)

type Counter struct {
	Comparison      int
	Assignment      int
	Addition        int
	Multiplication  int
	Division        int
	Subtraction     int
	Brackets        int
	Subfunction     int
	Variables       int
	Constants       int
	Parameters      int
	CurCounter      int
	MaxLevel        int
	MaxNestingLevel int
}

func (c *Counter) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.BinaryExpr:
		Mt[n.Op]++
		switch n.Op {
		case token.EQL, token.NEQ, token.LSS, token.LEQ, token.GTR, token.GEQ:
			c.Comparison++
		case token.ADD, token.ADD_ASSIGN, token.INC:
			c.Addition++
		case token.MUL:
			c.Multiplication++
		case token.QUO:
			c.Division++
		case token.SUB:
			c.Subtraction++
		}

	case *ast.AssignStmt:
		Mt[n.Tok]++
		c.Assignment += len(n.Lhs)
	case *ast.ParenExpr:
		c.Brackets++

	case *ast.FuncDecl:
		c.Subfunction++
	case *ast.Ident:
		if n.Obj != nil && (n.Obj.Kind == ast.Var) {
			M[n.Obj.Name]++
			c.Variables++
		}
	case *ast.BasicLit:

		if n.Kind == token.INT || n.Kind == token.FLOAT {
			c.Constants++
		}
	case *ast.CompositeLit:
		{
			if c.CurCounter > 100 {
				break
			}
			c.CurCounter++
			if c.CurCounter > c.MaxLevel {
				c.MaxLevel = c.CurCounter
			}
			ast.Walk(c, n)
			c.CurCounter--
		}
	case *ast.IncDecStmt:
		{
			if n.Tok == token.INC {
				c.Addition++
			}
			if n.Tok == token.DEC {
				c.Subtraction++
			}
			Mt[n.Tok]++
		}
	}
	return c
}

func main() {
	startTime := time.Now()

	//HERE IS YOUR CODE FROM THE MAIN FILE (functional metrics)
	//foo()

	duration := time.Since(startTime) / time.Millisecond
	fmt.Printf("\n----------------------------\nFunction took %d milliseconds to run.\n", duration)

	//HERE IS YOUR FILE (to do a code analysis)
	z, _ := os.ReadFile("main.go")

	src := string(z)
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	var counter Counter
	ast.Walk(&counter, node)

	for i := 0; i < len(z); i++ {
		switch z[i] {
		case '{':
			counter.CurCounter++
			if counter.CurCounter > counter.MaxNestingLevel {
				counter.MaxNestingLevel = counter.CurCounter
			}
		case '}':
			counter.CurCounter--
		}
	}
	counter.Addition += counter.MaxLevel
	fmt.Println("Comparison:", counter.Comparison)
	fmt.Println("Assignment:", counter.Assignment)
	fmt.Println("Addition:", counter.Addition)
	fmt.Println("Multiplication:", counter.Multiplication)
	fmt.Println("Division:", counter.Division)
	fmt.Println("Subtraction:", counter.Subtraction)
	fmt.Println("Brackets:", counter.Brackets)
	fmt.Println("Subfunction:", counter.Subfunction)
	fmt.Println("Variables:", counter.Variables)
	fmt.Println("Constants:", counter.Constants)
	fmt.Println("Parameters:", counter.Parameters)
	fmt.Println("Max level:", counter.MaxLevel)
	fmt.Println("Max nesting level:", counter.MaxNestingLevel)
	fmt.Println("-------------------------------")

	size := 0
	size2 := 0
	keys := make([]string, 0, len(M))
	for k, v := range M {
		keys = append(keys, k)
		size += v
	}
	keysops := make([]string, 0, len(Mt))
	for k, v := range Mt {
		keysops = append(keysops, string(rune(k)))
		size2 += v
	}
	fmt.Println(M)
	fmt.Println(Mt)
	fmt.Println("n = ", len(keys), "+", len(keysops), "=", len(keys)+len(keysops))
	fmt.Println("N = ", size, "+", size2, "=", size+size2)
	fmt.Println("V = ", float64(size+size2), "* log2(", float64(len(keys)+len(keysops)), ") =", float64(size+size2)*math.Log2(float64(len(keys)+len(keysops))))
}
