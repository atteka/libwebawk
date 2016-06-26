package libwebawk

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
)

type Match struct {
	tag   string
	class string
}

func consumeMatch(program string) (*Match, int, bool) {
	last := false
	tag := ""
	class := ""
	i := 0

	for {
		if program[i] == '/' || program[i] == '.' || program[i] == '[' {
			break
		}
		i++
	}
	tag = program[:i]
	if program[i] == '[' {
		j := i + 1
		for {
			if program[i] == ']' {
				break
			}
			i++
		}
		class = program[j:i]
		i++
	}
	if program[i] == '/' {
		last = true
	}
	m := Match{tag, class}

	return &m, i + 1, last
}

func consumeAddress(program string) (*Address, int, bool) {
	last := false
	i := 0
	prev := 0
	a := NewAddress()

	for {
		if program[i] == ' ' || program[i] == '}' {
			break
		}
		if program[i] == '.' {
			a.Insert(program[prev:i])
			prev = i + 1
		}
		i++
	}
	a.Insert(program[prev:i])

	if program[i] == '}' {
		last = true
	}
	return a, i + 1, last
}

func ParseWebawkProgram(program string) ([]Match, []*Address, error) {
	match := make([]Match, 0, 10)        // By default match anything
	addresses := make([]*Address, 0, 10) // By default print whole nested sub-block under 'match'
	i := 0

	if program[i] == '/' {
		i++
		for {
			m, l, last := consumeMatch(program[i:])
			i += l
			match = append(match, *m)
			if last {
				break
			}
		}
	}

	for ; i < len(program); i++ {
		if program[i] == '{' {
			break
		}
	}
	if program[i] == '{' {
		i++
		for {
			a, l, last := consumeAddress(program[i:])
			i += l
			addresses = append(addresses, a)
			if last {
				break
			}
		}
	}

	return match, addresses, nil
}

func isMatch(stack []Match, match []Match) bool {
	lenStack := len(stack)
	lenMatch := len(match)
	//fmt.Println(stack)
	if lenMatch > lenStack {
		return false
	}

	for i := 1; i <= lenMatch; i++ {
		if stack[lenStack-i].tag != match[lenMatch-i].tag {
			return false
		}
		if match[lenMatch-i].class != "" && stack[lenStack-i].class != match[lenMatch-i].class {
			return false
		}

	}
	return true
}

func createContext(z *html.Tokenizer, repeatOn string) *Context {
	ctxt := NewContext(repeatOn, "")
	contexts := make([]*Context, 0, 10)
	contexts = append(contexts, ctxt)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			fmt.Println("Error token")
			return nil
		case tt == html.StartTagToken:
			t := z.Token()

			ctxt = ctxt.CreateChild(t.Data, "")
			contexts = append(contexts, ctxt)
		case tt == html.SelfClosingTagToken:

		case tt == html.EndTagToken:
			//t := z.Token()
			if len(contexts) == 1 {
				return ctxt
			}

			contexts = contexts[:len(contexts)-1]
			ctxt = contexts[len(contexts)-1]

		case tt == html.TextToken:
			t := z.Token()
			ctxt.AppendText(t.Data)

		}
	}

}

func executeAction(z *html.Tokenizer, addresses []*Address, repeatOn string) {
	ctxt := createContext(z, repeatOn)
	for i, _ := range addresses {
		fmt.Print(ctxt.GetValue(*addresses[i]) + " ")
	}
	fmt.Println()
}

func createMatchFromToken(t html.Token) Match {
	class := ""

	for _, a := range t.Attr {
		if a.Key == "class" {
			class = a.Val
		}
	}

	return Match{t.Data, class}
}

func Run(body io.Reader, match []Match, addresses []*Address) {
	var stack = make([]Match, 0, 10) // The stack of tags

	z := html.NewTokenizer(body)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			return
		case tt == html.StartTagToken:
			t := z.Token()
			stack = append(stack, createMatchFromToken(t))

			if isMatch(stack, match) {
				executeAction(z, addresses, t.Data)

				stack = stack[:len(stack)-1]
				//return
			}
		case tt == html.EndTagToken:
			stack = stack[:len(stack)-1]
		}
	}
}
