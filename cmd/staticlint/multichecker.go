package main

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/fieldalignment"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"honnef.co/go/tools/quickfix"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
)

// Проверка вызова os.Exit из функции main пакета main
var ErrOsExitCheckAnalyzer = &analysis.Analyzer{
	Name: "errosexitcheck",
	Doc:  "check for os.Exit call from func main package main",
	Run:  run,
}

// Реализация анализатора
func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {

		if pass.Pkg.Name() != "main" {
			continue
		}
		ast.Inspect(file, func(node ast.Node) bool {
			fn, ok := node.(*ast.FuncDecl) // ast.FuncDecl представляет декларацию функции
			if !ok {
				return true // Не функция, обходим дальше
			}

			if fn.Name.Name != "main" {
				return false // Нашли функцию, но это не main. Дальше проверять нет смысла
			}

			ast.Inspect(fn, func(n ast.Node) bool {
				expr, ok := n.(*ast.CallExpr) // ast.CallExpr представляет вызов функции или метода
				if !ok {
					return true // Это не вызов функции, обходим дальше
				}

				selector, ok := expr.Fun.(*ast.SelectorExpr)
				if !ok {
					return true // Нет вызова метода через ".", обходим дальше
				}

				ident, ok := selector.X.(*ast.Ident) //
				if !ok {
					return true // Не идентификатор, обходим дальше
				}
				if ident.Name == "os" && selector.Sel.Name == "Exit" {
					pass.Reportf(ident.Pos(), "call os.Exit from func main package main")
				}
				return false
			})

			return false

		})
	}
	return nil, nil
}

func main() {
	allPassesChecks := []*analysis.Analyzer{
		assign.Analyzer,         // Пакет assign определяет анализатор, который обнаруживает бесполезные назначения.
		atomic.Analyzer,         // Пакет atomic определяет анализатор, который проверяет наличие распространенных ошибок с помощью пакета sync/atomic.
		bools.Analyzer,          // Пакет bools определяет анализатор, который обнаруживает распространенные ошибки, связанные с логическими операторами.
		copylock.Analyzer,       // Пакет copylock определяет анализатор, который проверяет наличие блокировок, ошибочно переданных по значению.
		errorsas.Analyzer,       // Пакет errorsas определяет анализатор, который проверяет, является ли второй аргумент errors .As - это указатель на тип, реализующий error .
		fieldalignment.Analyzer, // Пакет fieldalignment определяет анализатор, который обнаруживает структуры, которые потребляли бы меньше памяти, если бы их поля были отсортированы.
		httpresponse.Analyzer,   // Пакет httpresponse определяет анализатор, который проверяет наличие ошибок с помощью HTTP-ответов.
		loopclosure.Analyzer,    // Пакет loopclosure определяет анализатор, который проверяет наличие ссылок на переменные окружающего цикла из вложенных функций.
		lostcancel.Analyzer,     // Пакет lostcancel определяет анализатор, который проверяет, не удалось ли вызвать функцию отмены контекста.
		nilfunc.Analyzer,        // Пакет nilfunc определяет анализатор, который проверяет наличие бесполезных сравнений с nil.
		printf.Analyzer,         // Пакет printf определяет анализатор, который проверяет согласованность строк формата Printf и аргументов.
		shadow.Analyzer,         // Пакет shadow определяет анализатор, который проверяет наличие затененных переменных.
		shift.Analyzer,          // Пакет shift определяет анализатор, который проверяет наличие сдвигов, превышающих ширину целого числа.
		sigchanyzer.Analyzer,    // Пакет sigchanyzer определяет анализатор, который обнаруживает неправильное использование небуферизованного сигнала в качестве аргумента signal.Notify.
		sortslice.Analyzer,      // Пакет sortslice определяет анализатор, который проверяет вызовы sort.Slice, которые не используют тип slice в качестве первого аргумента.
		stdmethods.Analyzer,     // Пакет stdmethods определяет анализатор, который проверяет наличие орфографических ошибок в сигнатурах методов, аналогичных хорошо известным интерфейсам.
		stringintconv.Analyzer,  // Пакет stringintconv определяет анализатор, который помечает преобразования типов из целых чисел в строки.
		structtag.Analyzer,      // Пакет structtag определяет анализатор, который проверяет правильность формирования тегов полей struct.
		tests.Analyzer,          // Пакет tests определяют анализатор, который проверяет распространенное ошибочное использование тестов и примеров.
		unmarshal.Analyzer,      // Пакет unmarshal определяет анализатор, который проверяет передачу типов без указателей или без интерфейса функциям unmarshal и decode.
		unreachable.Analyzer,    // Пакет unreachable определяет анализатор, который проверяет наличие недоступного кода.
		unsafeptr.Analyzer,      // Пакет unsafeptr определяет анализатор, который проверяет недопустимые преобразования uintptr в unsafe .Указатель.
		unusedresult.Analyzer,   // Пакет unusedresult определяет анализатор, который проверяет неиспользуемые результаты вызовов определенных чистых функций.
		unusedwrite.Analyzer,    // Пакет unusedwrite проверяет наличие неиспользуемых операций записи в элементы объекта struct или array.
	}
	otherChecks := []*analysis.Analyzer{}

	for _, v := range staticcheck.Analyzers { // Различные злоупотребления стандартной библиотекой
		otherChecks = append(otherChecks, v.Analyzer)
	}

	for _, v := range simple.Analyzers { // Пакет simple содержит проверки, упрощающие код. Все предложения, сделанные в результате этих анализов, призваны привести к объективному упрощению кода.
		otherChecks = append(otherChecks, v.Analyzer)

	}

	for _, v := range quickfix.Analyzers { // Пакет quickfix содержит проверки, реализующие рефакторинг кода.
		otherChecks = append(otherChecks, v.Analyzer)

	}

	for _, v := range stylecheck.Analyzers { // Пакет stylecheck содержит проверки, обеспечивающие соблюдение правил стиля.
		otherChecks = append(otherChecks, v.Analyzer)

	}

	analyzers := []*analysis.Analyzer{}
	analyzers = append(analyzers, allPassesChecks...)
	analyzers = append(analyzers, otherChecks...)
	analyzers = append(analyzers, ErrOsExitCheckAnalyzer) // Проверка вызова os.Exit из функции main пакета main

	multichecker.Main(
		analyzers...,
	)
}
