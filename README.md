# Actions Parser

This is a learning project for me to learn a bit about lexing, parsing, and evaluating code. It's focused on the GitHub Actions [context and expression syntax](https://docs.github.com/en/actions/reference/context-and-expression-syntax-for-github-actions).

I was inspired to work specifically on this while I was digging in and learning about the [`nektos/act`](https://github.com/nektos/act) project, and specifically seeing [this issue](https://github.com/nektos/act/issues/104) in there about how the syntax is not quite valid javascript. `nektos/act` uses the Otto JS virtual machine to evaluate these expressions, so was failing (until an [awesome workaround](https://github.com/nektos/act/pull/287) was submitted at least).

Under the hood I didn't want to just fall back on a JS VM, and also didn't want to go right to direct byte-by-byte lexing and parsing, so I built on top of the [Participle](https://github.com/alecthomas/participle) package as it is both well used (lots of stars, etc), pure Go, and also had a lot of useful examples of different languages.

### Shoulders of Giants

Nothing is done in a vacuum, especially stuff that I think of a "deep computer science-y". I read a bunch of blog posts and perused a many existing projects in order to put this together. Ideas came from all over:

- Rob Pike's presentation about [`Lexical Scanning in Go`](https://www.youtube.com/watch?v=HxaD_trXwRE) on Youtube is a great dive into the space
- The awesome work done by @cschleiden on his project [GitHub Actions Parser](https://github.com/cschleiden/github-actions-parser) that uses the JS library [Chevrotain](https://sap.github.io/chevrotain/docs/) (which has a very well written tutorial by the way!)
- Another co-worker @fatih's [blog post on go's lexer/scanner](https://medium.com/@farslan/a-look-at-go-scanner-packages-11710c2655fc) packages is a great read, plus looking at the resulting [HCL](https://github.com/hashicorp/hcl) parsing library that he build and then handed over to Hashicorp
- Go's own [`text/scanner`](https://golang.org/pkg/text/scanner) [`go/scanner`](https://golang.org/pkg/go/scanner/), and [`go/token`](https://golang.org/pkg/go/token/) packages
- The [PEG](https://github.com/pointlander/peg) project, as well as reading up on the [ANTLR project](https://www.antlr.org/) got me thinking about language grammars
- I loved the sentence `Lex is not code that you live in. It is code you write once and then use for a long time. Ok if the code is not clean.` from [this article about using GoYACC](https://about.sourcegraph.com/go/gophercon-2018-how-to-write-a-parser-in-go/)
- [Various](https://blog.gopheracademy.com/advent-2014/parsers-lexers/) [other](https://adampresley.github.io/2015/04/12/writing-a-lexer-and-parser-in-go-part-1.html) [blog](https://blog.gopheracademy.com/advent-2017/parsing-with-antlr4-and-go/) posts as well
- Oh, and [Regex101](https://regex101.com/) for always being super helpful!
