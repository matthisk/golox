package parser

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/matthisk/lox/lexer"
)

func TestInterpreter_WithStmts_TableDriven(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		logs    []string
		wantErr bool
	}{
		// Variable declaration tests
		{"basic var declaration and use", "var a = 1; var b = 2; print a + b;", []string{"3"}, false},
		{"var with string", "var name = \"John\"; print name;", []string{"John"}, false},
		{"var with boolean", "var flag = true; print flag;", []string{"true"}, false},
		{"var with nil", "var empty = nil; print empty;", []string{"<nil>"}, false},
		{"var with expression", "var result = 2 + 3 * 4; print result;", []string{"14"}, false},
		{"var with comparison", "var isGreater = 5 > 3; print isGreater;", []string{"true"}, false},
		{"var with string concat", "var greeting = \"Hello\" + \" World\"; print greeting;", []string{"Hello World"}, false},
		{"var with ternary", "var value = true ? 42 : 0; print value;", []string{"42"}, false},
		{"var with unary", "var negative = -10; print negative;", []string{"-10"}, false},
		{"var with logical not", "var opposite = !false; print opposite;", []string{"true"}, false},

		// Multiple variable operations
		{"multiple vars same line", "var x = 5; var y = 10; var z = x * y; print z;", []string{"50"}, false},
		{"var reuse in expression", "var a = 3; var b = a + a; print b;", []string{"6"}, false},
		{"vars in complex expression", "var x = 2; var y = 3; print (x + y) * (x - y);", []string{"-5"}, false},
		{"vars with grouping", "var a = 2; var b = 3; var c = 4; print a + (b * c);", []string{"14"}, false},
		{"vars in ternary", "var a = 5; var b = 3; print a > b ? a : b;", []string{"5"}, false},
		{"vars in comma expression", "var a = 1; var b = 2; var c = 3; print a, b, c;", []string{"3"}, false},

		// Variable reassignment through redeclaration
		{"var redeclaration", "var x = 1; print x; var x = 2; print x;", []string{"1", "2"}, false},
		{"var redeclaration with different type", "var v = 42; print v; var v = \"hello\"; print v;", []string{"42", "hello"}, false},
		{"assignment error", "x = 2;", []string{}, true},
		{"print assignemnt statement", "var x = 1; print x = 2;", []string{"2"}, false},
		{"var redeclaration with assignment", "var x = 1; print x; x = 2; print x;", []string{"1", "2"}, false},
		{"var redeclaration with assignment and different type", "var v = 42; print v; v = \"hello\"; print v;", []string{"42", "hello"}, false},

		// Variables with all operators
		{"vars with arithmetic", "var a = 10; var b = 3; print a + b; print a - b; print a * b; print a / b;", []string{"13", "7", "30", "3.3333333333333335"}, false},
		{"vars with comparison", "var x = 5; var y = 3; print x > y; print x >= y; print x < y; print x <= y;", []string{"true", "true", "false", "false"}, false},
		{"vars with equality", "var a = 5; var b = 5; var c = 3; print a == b; print a != c;", []string{"true", "true"}, false},

		// Edge cases and complex scenarios
		{"var with decimal", "var pi = 3.14159; print pi * 2;", []string{"6.28318"}, false},
		{"var with large number", "var big = 999999; print big + 1;", []string{"1e+06"}, false},
		{"var with zero", "var zero = 0; print zero == 0;", []string{"true"}, false},
		{"var with empty string", "var empty = \"\"; print empty == \"\";", []string{"true"}, false},

		// Variables in nested expressions
		{"deeply nested vars", "var a = 1; var b = 2; var c = 3; print ((a + b) * c) > (a * (b + c));", []string{"true"}, false},
		{"vars with mixed types", "var num = 42; var str = \"answer\"; var bool = true; print bool ? num : str;", []string{"42"}, false},
		// Print statement tests
		{"print string literal", "print \"hello world\";", []string{"hello world"}, false},
		{"print arithmetic", "print 5 + 5;", []string{"10"}, false},
		{"print complex arithmetic", "print 5 + 5 * 10;", []string{"55"}, false},
		{"multiple print statements", "print 5 + 5 * 10; print \"hello\";", []string{"55", "hello"}, false},
		{"print boolean true", "print true;", []string{"true"}, false},
		{"print boolean false", "print false;", []string{"false"}, false},
		{"print nil", "print nil;", []string{"<nil>"}, false},
		{"print string concatenation", "print \"hello\" + \" \" + \"world\";", []string{"hello world"}, false},
		{"print comparison result", "print 5 > 3;", []string{"true"}, false},
		{"print equality result", "print 5 == 5;", []string{"true"}, false},
		{"print unary negation", "print -42;", []string{"-42"}, false},
		{"print logical not", "print !true;", []string{"false"}, false},
		{"print ternary result", "print true ? \"yes\" : \"no\";", []string{"yes"}, false},
		{"print comma expression", "print 1, 2, 3;", []string{"3"}, false},
		{"print grouped expression", "print (2 + 3) * 4;", []string{"20"}, false},

		// Block statement
		{"simple block scoping", "var x = 0; var y = 1; var z = 2; { var x = 1; var y = 2; { var x = 2; print x; print y; print z; } } print x;", []string{"2", "2", "2", "0"}, false},

		// If statement tests
		{"simple if true condition", "if (true) print \"if executed\";", []string{"if executed"}, false},
		{"if and condition", "if (true && true) print \"if executed\";", []string{"if executed"}, false},
		{"if or condition", "if (false || true) print \"if executed\";", []string{"if executed"}, false},
		{"simple if false condition", "if (false) print \"should not print\";", []string{}, false},
		{"if-else with true condition", "if (true) print \"if block\"; else print \"else block\";", []string{"if block"}, false},
		{"if-else with false condition", "if (false) print \"if block\"; else print \"else block\";", []string{"else block"}, false},
		{"if with expression condition", "if (5 > 3) print \"five is greater\";", []string{"five is greater"}, false},
		{"if with false expression condition", "if (3 > 5) print \"should not print\";", []string{}, false},
		{"if with variable condition", "var condition = true; if (condition) print \"variable is true\";", []string{"variable is true"}, false},
		{"if with arithmetic in condition", "if (2 + 3 == 5) print \"math works\";", []string{"math works"}, false},
		{"if with string comparison", "var name = \"John\"; if (name == \"John\") print \"hello John\";", []string{"hello John"}, false},
		{"if with nil condition", "if (nil) print \"should not print\";", []string{}, false},
		{"if with truthy number", "if (42) print \"number is truthy\";", []string{"number is truthy"}, false},
		{"if with truthy string", "if (\"hello\") print \"string is truthy\";", []string{"string is truthy"}, false},
		{"if with zero condition", "if (0) print \"zero is truthy\";", []string{"zero is truthy"}, false},

		// If statement with blocks
		{"if with block statement", "if (true) { print \"line 1\"; print \"line 2\"; }", []string{"line 1", "line 2"}, false},
		{"if-else with blocks", "if (false) { print \"if block\"; } else { print \"else block\"; }", []string{"else block"}, false},
		{"if with variable in block", "var x = 5; if (true) { var y = 10; print x + y; }", []string{"15"}, false},
		{"if with block scoping", "var x = 1; if (true) { var x = 2; print x; } print x;", []string{"2", "1"}, false},
		{"if with break statement", "var x = 1; if (true) { break; print x; }", []string{}, false},

		// Nested if statements
		{"nested if statements", "if (true) if (true) print \"nested\";", []string{"nested"}, false},
		{"nested if-else", "if (true) if (false) print \"inner if\"; else print \"inner else\";", []string{"inner else"}, false},
		{"nested if with blocks", "if (true) { if (false) print \"inner\"; print \"outer\"; }", []string{"outer"}, false},

		// If-else-if chains
		{"if-else-if chain", "var x = 2; if (x == 1) print \"one\"; else if (x == 2) print \"two\"; else print \"other\";", []string{"two"}, false},
		{"if-else-if chain all false", "var x = 5; if (x == 1) print \"one\"; else if (x == 2) print \"two\"; else print \"other\";", []string{"other"}, false},
		{"long if-else-if chain", "var grade = 85; if (grade >= 90) print \"A\"; else if (grade >= 80) print \"B\"; else if (grade >= 70) print \"C\"; else print \"F\";", []string{"B"}, false},

		// If statements with variables and expressions
		{"if with variable assignment", "var result = \"\"; if (true) result = \"success\"; print result;", []string{"success"}, false},
		{"if modifying outer variable", "var counter = 0; if (true) { counter = counter + 1; } print counter;", []string{"1"}, false},
		{"if with ternary in condition", "if (true ? true : false) print \"ternary worked\";", []string{"ternary worked"}, false},
		{"if with comma expression", "var a = 1; if (a = 2, a == 2) print \"comma and assignment\";", []string{"comma and assignment"}, false},

		// Complex if statement scenarios
		{"multiple if statements", "if (true) print \"first\"; if (false) print \"second\"; if (true) print \"third\";", []string{"first", "third"}, false},
		{"if statements with different types", "if (1) print \"number\"; if (\"\") print \"string\"; if (nil) print \"nil\";", []string{"number", "string"}, false},
		{"if with print in both branches", "var choice = true; if (choice) print \"option A\"; else print \"option B\";", []string{"option A"}, false},

		// While loop tests
		{"while with false condition", "while (false) print \"should not print\";", []string{}, false},
		{"while with counter", "var i = 0; while (i < 3) { print i; i = i + 1; }", []string{"0", "1", "2"}, false},
		{"while with string counter", "var count = 0; while (count < 2) { print \"iteration\"; count = count + 1; }", []string{"iteration", "iteration"}, false},
		{"while with variable condition", "var running = true; var counter = 0; while (running) { print counter; counter = counter + 1; if (counter >= 2) running = false; }", []string{"0", "1"}, false},
		{"while modifying condition inside loop", "var x = 5; while (x > 0) { print x; x = x - 1; }", []string{"5", "4", "3", "2", "1"}, false},
		{"while with complex condition", "var a = 1; var b = 10; while (a < 5 && b > 8) { print a; a = a + 1; b = b - 1; }", []string{"1", "2"}, false},
		{"while with block and variables", "var sum = 0; var i = 1; while (i <= 3) { sum = sum + i; print sum; i = i + 1; }", []string{"1", "3", "6"}, false},
		{"nested while loops", "var outer = 0; while (outer < 2) { var inner = 0; while (inner < 2) { print outer * 10 + inner; inner = inner + 1; } outer = outer + 1; }", []string{"0", "1", "10", "11"}, false},
		{"while with early termination", "var count = 0; while (count < 10) { print count; count = count + 1; if (count == 3) count = 10; }", []string{"0", "1", "2"}, false},
		{"while with nil check", "var value = 5; while (value) { print value; value = value - 1; if (value == 0) value = nil; }", []string{"5", "4", "3", "2", "1"}, false},
		{"while with boolean toggle", "var toggle = true; var count = 0; while (toggle) { print \"on\"; count = count + 1; if (count >= 2) toggle = false; }", []string{"on", "on"}, false},
		{"while with break", "var i = 0; while (i < 3) { print i; break; i = i + 1; }", []string{"0"}, false},
		{"while with continue", "var i = 0; while (i < 3) { i = i + 1; continue; print i; }", []string{}, false},

		// For loop tests (desugared to while loops)
		{"basic for loop with all parts", "for (var i = 0; i < 3; i = i + 1) print i;", []string{"0", "1", "2"}, false},
		{"for loop with expression initializer", "var count = 10; for (count = 5; count > 0; count = count - 1) print count;", []string{"5", "4", "3", "2", "1"}, false},
		{"for loop with no initializer", "var i = 2; for (; i < 5; i = i + 1) print i;", []string{"2", "3", "4"}, false},
		//{"for loop with no condition (infinite)", "var count = 0; for (var i = 0; ; i = i + 1) { print i; count = count + 1; if (count >= 3) i = 10; }", []string{"0", "1", "2"}, false},
		{"for loop with no increment", "for (var x = 1; x <= 3; ) { print x; x = x + 1; }", []string{"1", "2", "3"}, false},
		//{"empty for loop elements", "var i = 0; for (; ; ) { print i; i = i + 1; if (i >= 2) i = 10; }", []string{"0", "1"}, false},
		{"for loop with block statement", "for (var i = 0; i < 2; i = i + 1) { print \"iteration\"; print i; }", []string{"iteration", "0", "iteration", "1"}, false},
		{"for loop with complex expressions", "for (var start = 2 * 3; start < 10 + 5; start = start * 2) print start;", []string{"6", "12"}, false},
		{"for loop with string operations", "for (var msg = \"a\"; msg != \"aaaa\"; msg = msg + \"a\") print msg;", []string{"a", "aa", "aaa"}, false},
		{"nested for loops", "for (var i = 0; i < 2; i = i + 1) for (var j = 0; j < 2; j = j + 1) print i * 10 + j;", []string{"0", "1", "10", "11"}, false},
		{"for loop with variable in all parts", "var step = 2; for (var i = step; i < 10; i = i + step) print i;", []string{"2", "4", "6", "8"}, false},
		{"for loop with arithmetic in condition", "for (var i = 1; i * 2 <= 8; i = i + 1) print i * 2;", []string{"2", "4", "6", "8"}, false},
		{"for loop with boolean condition", "var keepGoing = true; for (var i = 0; keepGoing; i = i + 1) { print i; if (i >= 2) keepGoing = false; }", []string{"0", "1", "2"}, false},
		{"for loop modifying external variable", "var sum = 0; for (var i = 1; i <= 3; i = i + 1) { sum = sum + i; print sum; }", []string{"1", "3", "6"}, false},
		{"for loop with early termination via condition", "for (var i = 0; i < 10; i = i + 1) { print i; if (i == 2) i = 15; }", []string{"0", "1", "2"}, false},
		{"for loop with break", "for (var i = 0; i < 10; i = i + 1) { print i; break; }", []string{"0"}, false},
		// This test hangs because of the desugaring from for to while loop. We do not execute the increment after we execute continue.
		// {"for loop with continue", "for (var i = 0; i <= 1; i = i + 1) { print i; continue; }", []string{"0", "1"}, false},

		// Expression statement tests (no output expected)
		{"expression statement arithmetic", "5 + 5;", []string{}, false},
		{"expression statement string", "\"hello\";", []string{}, false},
		{"expression statement boolean", "true;", []string{}, false},
		{"expression statement comparison", "5 > 3;", []string{}, false},
		{"expression statement function call", "!false;", []string{}, false},

		// Mixed statement tests
		{"print and expression mixed", "print \"start\"; 5 + 5; print \"end\";", []string{"start", "end"}, false},
		{"multiple expression statements", "1 + 1; 2 + 2; 3 + 3;", []string{}, false},
		{"complex mixed statements", "print \"result:\"; print 2 * 3; 4 + 4; print \"done\";", []string{"result:", "6", "done"}, false},

		// Edge cases
		{"empty program", "", []string{}, false},
		{"print empty string", "print \"\";", []string{""}, false},
		{"print zero", "print 0;", []string{"0"}, false},
		{"print negative zero", "print -0;", []string{"-0"}, false},
		{"print large number", "print 999999;", []string{"999999"}, false},
		{"print decimal", "print 3.14159;", []string{"3.14159"}, false},
		{"print simple string", "print \"simple string\";", []string{"simple string"}, false},
		{"returning from top-level code", "return \"simple string\";", []string{}, true},

		// Complex expressions in print statements
		{"print nested ternary", "print true ? (false ? 1 : 2) : 3;", []string{"2"}, false},
		{"print chained comparisons", "print 1 < 2 == 2 > 1;", []string{"true"}, false},
		{"print mixed operators", "print 2 + 3 * 4 == 14;", []string{"true"}, false},
		{"print complex grouping", "print ((1 + 2) * (3 + 4)) / 7;", []string{"3"}, false},

		// Functions - Basic Declaration and Calls
		{"simple function", "fun sayHi(first, last) { print \"Hi, \" + first + \" \" + last + \"!\"; } sayHi(\"Joe\", \"Doe\");", []string{"Hi, Joe Doe!"}, false},
		{"function no parameters", "fun greet() { print \"Hello World!\"; } greet();", []string{"Hello World!"}, false},
		{"function single parameter", "fun square(x) { print x * x; } square(5);", []string{"25"}, false},
		{"function multiple parameters", "fun add(a, b, c) { print a + b + c; } add(1, 2, 3);", []string{"6"}, false},
		{"function with numbers", "fun multiply(x, y) { print x * y; } multiply(4, 7);", []string{"28"}, false},
		{"function with strings", "fun concat(a, b) { print a + b; } concat(\"hello\", \" world\");", []string{"hello world"}, false},
		{"function with booleans", "fun logicalAnd(a, b) { print a && b; } logicalAnd(true, false);", []string{"false"}, false},
		{"function with mixed types", "fun mixed(num, str, bool) { print num; print str; print bool; } mixed(42, \"test\", true);", []string{"42", "test", "true"}, false},

		// Functions - Multiple Calls
		{"function called multiple times", "fun count() { print \"counting\"; } count(); count(); count();", []string{"counting", "counting", "counting"}, false},
		{"function with different arguments", "fun echo(msg) { print msg; } echo(\"first\"); echo(\"second\"); echo(\"third\");", []string{"first", "second", "third"}, false},
		{"function with arithmetic", "fun calc(x) { print x + 10; } calc(1); calc(2); calc(3);", []string{"11", "12", "13"}, false},

		// Functions - Variable Scope
		{"function accessing global variable", "var global = \"accessible\"; fun useGlobal() { print global; } useGlobal();", []string{"accessible"}, false},
		{"function with local variable", "fun withLocal() { var local = \"inside\"; print local; } withLocal();", []string{"inside"}, false},
		{"function parameter shadows global", "var x = \"global\"; fun shadow(x) { print x; } shadow(\"local\");", []string{"local"}, false},
		{"global variable after function", "var x = \"global\"; fun shadow(x) { print x; } shadow(\"local\"); print x;", []string{"local", "global"}, false},
		{"function modifying global", "var counter = 0; fun increment() { counter = counter + 1; print counter; } increment(); increment();", []string{"1", "2"}, false},

		// Functions - Nested Scopes
		{"function with block scope", "fun withBlock() { var outer = \"outer\"; { var inner = \"inner\"; print outer; print inner; } print outer; } withBlock();", []string{"outer", "inner", "outer"}, false},
		{"function parameter in nested block", "fun nested(param) { { print param; var local = param + \" modified\"; print local; } } nested(\"test\");", []string{"test", "test modified"}, false},

		// Functions - Control Flow
		{"function with if statement", "fun conditional(x) { if (x > 0) print \"positive\"; else print \"non-positive\"; } conditional(5); conditional(-3);", []string{"positive", "non-positive"}, false},
		{"function with while loop", "fun countdown(n) { while (n > 0) { print n; n = n - 1; } } countdown(3);", []string{"3", "2", "1"}, false},
		{"function with for loop", "fun forLoop(max) { for (var i = 0; i < max; i = i + 1) print i; } forLoop(3);", []string{"0", "1", "2"}, false},

		// Functions - Complex Bodies
		{"function with multiple statements", "fun complex() { var a = 1; var b = 2; print a + b; print a * b; print a > b; } complex();", []string{"3", "2", "false"}, false},
		{"function with nested function call", "fun inner() { print \"inner\"; } fun outer() { print \"outer start\"; inner(); print \"outer end\"; } outer();", []string{"outer start", "inner", "outer end"}, false},
		{"function calculating factorial", "fun factorial(n) { var result = 1; while (n > 1) { result = result * n; n = n - 1; } print result; } factorial(5);", []string{"120"}, false},

		// Functions - Built-in Functions
		{"clock function exists", "print clock();", []string{}, false}, // We can't predict exact time, so we'll just verify it runs without error

		// Functions - Edge Cases
		{"empty function body", "fun empty() { } empty();", []string{}, false},
		{"function with only variable declaration", "fun varOnly() { var x = 42; } varOnly();", []string{}, false},
		{"function with expression statement", "fun exprStmt() { 1 + 1; } exprStmt();", []string{}, false},

		// Functions - Error Cases (Argument Validation)
		{"function wrong argument count - too few", "fun needsTwo(a, b) { print a + b; } needsTwo(1);", []string{}, true},
		{"function wrong argument count - too many", "fun needsOne(x) { print x; } needsOne(1, 2);", []string{}, true},
		{"function zero params called with args", "fun noParams() { print \"none\"; } noParams(1);", []string{}, true},
		{"function many params called with few", "fun manyParams(a, b, c, d, e) { print a; } manyParams(1, 2);", []string{}, true},
		{"function many params called with many", "fun manyParams(a, b, c, d, e) { print a; } manyParams(1, 2, 3, 4, 5, 6);", []string{}, true},

		// Functions - Error Cases (Call Non-Functions)
		{"call number as function", "var num = 42; num();", []string{}, true},
		{"call string as function", "var str = \"hello\"; str();", []string{}, true},
		{"call boolean as function", "var bool = true; bool();", []string{}, true},
		{"call nil as function", "var nothing = nil; nothing();", []string{}, true},

		// Functions - Error Cases (Undefined Functions)
		{"call undefined function", "undefinedFunction();", []string{}, true},
		{"call function before definition", "callEarly(); fun callEarly() { print \"defined later\"; }", []string{}, true},

		// Functions - Recursive Functions
		{"simple recursive function", "fun fibonacci(n) { if (n <= 1) print n; else { fibonacci(n - 1); fibonacci(n - 2); } } fibonacci(3);", []string{"1", "0", "1"}, false},
		{"simple recursive function with return", "fun fibonacci(n) { if (n <= 1) { return n ; } else { return fibonacci(n - 1) + fibonacci(n - 2); } } for (var i = 0; i < 5; i = i + 1) { print fibonacci(i); }", []string{"0", "1", "1", "2", "3"}, false},
		{"recursive countdown", "fun countdown(n) { if (n > 0) { print n; countdown(n - 1); } } countdown(3);", []string{"3", "2", "1"}, false},

		// Functions - Advanced Scope Tests
		{"function defines local then accesses global", "var x = \"global\"; fun test() { var x = \"local\"; print x; } test(); print x;", []string{"local", "global"}, false},
		{"nested function calls with parameters", "fun outer(a) { fun inner(b) { print a + b; } inner(2); } outer(1);", []string{"3"}, false},
		{"function parameter masking outer parameter", "fun outer(x) { fun inner(x) { print x; } inner(\"inner\"); print x; } outer(\"outer\");", []string{"inner", "outer"}, false},

		// Functions - Multiple Function Definitions
		{"multiple function definitions", "fun first() { print \"1st\"; } fun second() { print \"2nd\"; } first(); second();", []string{"1st", "2nd"}, false},
		{"function redefinition", "fun test() { print \"first\"; } fun test() { print \"second\"; } test();", []string{"second"}, false},
		{"functions calling each other", "fun ping() { print \"ping\"; pong(); } fun pong() { print \"pong\"; } ping();", []string{"ping", "pong"}, false},

		// Functions - Complex Argument Expressions
		{"function with expression arguments", "fun add(a, b) { print a + b; } add(2 * 3, 4 + 5);", []string{"15"}, false},
		{"function with function call arguments", "fun double(x) { print x * 2; } fun getValue() { return 5; } double(getValue());", []string{"10"}, false},
		{"function with ternary arguments", "fun show(x) { print x; } show(true ? \"yes\" : \"no\");", []string{"yes"}, false},
		{"function with comma expression arguments", "fun first(x) { print x; } first(1, 2, 3);", []string{}, true}, // Should error due to wrong arg count

		// Functions - Built-in Functions Extended
		{"clock returns number", "var time = clock(); print time >= 0;", []string{"true"}, false},
		{"clock called multiple times", "var t1 = clock(); var t2 = clock(); print t2 >= t1;", []string{"true"}, false},

		// Functions - Return statements
		{"return statement breaks execution", "fun exitEarly() { print 0; return; print 1; } exitEarly();", []string{"0"}, false},

		// Functions - Local functions and scoping
		{"Local functions and closures", "fun makeCounter() { var i = 0; fun count() { i = i + 1; print i; } return count; } var counter = makeCounter(); counter(); counter();", []string{"1", "2"}, false},
		{"Disallow dynamic scoping", `var a = "global"; { fun showA() { print a; } showA(); var a = "block"; showA(); }`, []string{"global", "global"}, false},

		// Redeclaring errors
		{"Re-declaring in a scope not allowed", "fun scope() { var a = 1; var a = 2; }", []string{}, true},

		// Classes
		{"Print a class' name", "class RoundBall { roll() { return \"rolled\"; } } print RoundBall;", []string{"&{RoundBall}"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logs, err := runLox(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if !SliceEqual(tt.logs, logs) {
				t.Errorf("Interpreter printed unexpected results %v expected %v", logs, tt.logs)
			}
		})
	}
}

func TestInterpreter_TableDriven(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		expected interface{}
		wantErr  bool
	}{
		// Basic arithmetic expressions
		{"simple addition", "1 + 2", float64(3), false},
		{"simple multiplication", "3 * 4", float64(12), false},
		{"string concatenation", `"hello" + "world"`, "helloworld", false},

		// Arithmetic operations
		{"subtraction", "5 - 3", float64(2), false},
		{"division", "8 / 2", float64(4), false},

		// Unary operators
		{"unary minus", "-5", float64(-5), false},
		{"unary not true", "!true", false, false},
		{"unary not false", "!false", true, false},
		{"unary not truthy", "!123", false, false},
		{"unary not nil", "!nil", true, false},

		// Comparison operators (these need number operands)
		{"greater than true", "5 > 3", true, false},
		{"greater than false", "3 > 5", false, false},
		{"greater equal true", "5 >= 5", true, false},
		{"greater equal false", "3 >= 5", false, false},
		{"less than true", "3 < 5", true, false},
		{"less than false", "5 < 3", false, false},
		{"less equal true", "3 <= 3", true, false},
		{"less equal false", "5 <= 3", false, false},

		// Equality operators
		{"equal numbers", "5 == 5", true, false},
		{"not equal numbers", "5 != 3", true, false},
		{"equal strings", `"hello" == "hello"`, true, false},
		{"not equal strings", `"hello" != "world"`, true, false},
		{"equal booleans", "true == true", true, false},
		{"not equal booleans", "true != false", true, false},
		{"nil equality", "nil == nil", true, false},
		{"nil inequality", "nil != nil", false, false},
		{"different types not equal", "5 != true", true, false},

		// Grouping expressions
		{"grouped addition", "(1 + 2)", float64(3), false},
		{"grouped multiplication", "(3 * 4)", float64(12), false},
		{"nested grouping", "((5))", float64(5), false},

		// Mixed operator precedence
		{"addition and multiplication", "2 + 3 * 4", float64(14), false},
		{"multiplication and addition with grouping", "(2 + 3) * 4", float64(20), false},
		{"unary and binary", "-2 + 3", float64(1), false},
		{"comparison and equality", "3 > 2 == true", true, false},

		// Ternary operator
		{"ternary true condition", "true ? 1 : 2", float64(1), false},
		{"ternary false condition", "false ? 1 : 2", float64(2), false},
		{"ternary with nil", "nil ? 1 : 2", float64(2), false},
		{"ternary with truthy", "5 ? 1 : 2", float64(1), false},
		{"nested ternary", "true ? (false ? 1 : 2) : 3", float64(2), false},

		// Comma operator
		{"simple comma", "1, 2", float64(2), false},
		{"comma with expressions", "1 + 2, 3 * 4", float64(12), false},
		{"multiple comma", "1, 2, 3", float64(3), false},

		// Complex nested expressions
		{"complex arithmetic", "2 + 3 * 4 - 1", float64(13), false},
		{"complex with grouping", "(2 + 3) * (4 - 1)", float64(15), false},
		{"complex with ternary", "2 > 1 ? 3 + 4 : 5 * 6", float64(7), false},
		{"complex with comma", "1 + 2, 3 * 4, 5 > 3", true, false},
		{"deeply nested", "((2 + 3) * 4) > 15 ? true : false", true, false},

		// Edge cases with different types
		{"string and number not equal", `"5" != 5`, true, false},
		{"boolean and number not equal", "true != 1", true, false},
		{"nil and false not equal", "nil != false", true, false},

		// More complex truthiness tests
		{"zero is truthy", "!0", false, false},
		{"empty string is truthy", `!""`, false, false},
		{"non-zero number is truthy", "!42", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := runLoxExpression(tt.expr)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if result != tt.expected {
				t.Errorf("Expression %q: expected %v (%T), got %v (%T)",
					tt.expr, tt.expected, tt.expected, result, result)
			}
		})
	}
}

func runLoxExpression(source string) (interface{}, error) {
	lx := lexer.New(bytes.NewBufferString(source))
	lexResult := lexer.Consume(lx)
	if lexResult.Err != nil {
		return nil, lexResult.Err
	}

	parser := New(lexResult.Tokens)
	expr, err := parser.expression()
	if err != nil {
		return nil, err
	}

	interpreter := Interpreter{}
	return interpreter.evaluate(expr)
}

func runLox(source string) ([]string, error) {
	lx := lexer.New(bytes.NewBufferString(source))
	lexResult := lexer.Consume(lx)
	if lexResult.Err != nil {
		return nil, lexResult.Err
	}

	for i := range lexResult.Tokens {
		if lexResult.Tokens[i].Type == lexer.ILLEGAL {
			return nil, fmt.Errorf("Lexer found illegal token found at index %d", i)
		}
	}

	parser := New(lexResult.Tokens)
	stmts, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	printer := &SpyPrinter{}
	interpreter := NewInterpreterWithPrinter(printer)
	resolver := NewResolver(interpreter)
	_, err = resolver.resolveStmts(stmts)
	if err != nil {
		return nil, err
	}

	err = interpreter.Run(stmts)
	if err != nil {
		return nil, err
	}

	return printer.log, nil
}

type SpyPrinter struct {
	log []string
}

func (s *SpyPrinter) Print(value interface{}) {
	s.log = append(s.log, fmt.Sprintf("%v", value))
}

// SliceEqual compares two slices for equality.
// Returns true if both slices have the same length and all elements are equal.
func SliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
