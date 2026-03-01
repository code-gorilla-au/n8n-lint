package formats

import (
	"testing"

	"github.com/code-gorilla-au/odize"
)

func TestCaseSnake(t *testing.T) {
	group := odize.NewGroup(t, nil)

	err := group.
		Test("should convert camelCase to snake_case", func(t *testing.T) {
			odize.AssertEqual(t, "hello_world", CaseSnake("helloWorld"))
		}).
		Test("should convert PascalCase to snake_case", func(t *testing.T) {
			odize.AssertEqual(t, "hello_world", CaseSnake("HelloWorld"))
		}).
		Test("should convert kebab-case to snake_case", func(t *testing.T) {
			odize.AssertEqual(t, "hello_world", CaseSnake("hello-world"))
		}).
		Test("should convert space separated string to snake_case", func(t *testing.T) {
			odize.AssertEqual(t, "hello_world", CaseSnake("hello world"))
		}).
		Run()
	odize.AssertNoError(t, err)
}

func TestCasePascal(t *testing.T) {
	group := odize.NewGroup(t, nil)

	err := group.
		Test("should convert snake_case to PascalCase", func(t *testing.T) {
			odize.AssertEqual(t, "HelloWorld", CasePascal("hello_world"))
		}).
		Test("should convert camelCase to PascalCase", func(t *testing.T) {
			odize.AssertEqual(t, "HelloWorld", CasePascal("helloWorld"))
		}).
		Test("should convert kebab-case to PascalCase", func(t *testing.T) {
			odize.AssertEqual(t, "HelloWorld", CasePascal("hello-world"))
		}).
		Test("should convert space separated string to PascalCase", func(t *testing.T) {
			odize.AssertEqual(t, "HelloWorld", CasePascal("hello world"))
		}).
		Run()
	odize.AssertNoError(t, err)
}

func TestCaseKebab(t *testing.T) {
	group := odize.NewGroup(t, nil)

	err := group.
		Test("should convert snake_case to kebab-case", func(t *testing.T) {
			odize.AssertEqual(t, "hello-world", CaseKebab("hello_world"))
		}).
		Test("should convert camelCase to kebab-case", func(t *testing.T) {
			odize.AssertEqual(t, "hello-world", CaseKebab("helloWorld"))
		}).
		Test("should convert PascalCase to kebab-case", func(t *testing.T) {
			odize.AssertEqual(t, "hello-world", CaseKebab("HelloWorld"))
		}).
		Test("should convert space separated string to kebab-case", func(t *testing.T) {
			odize.AssertEqual(t, "hello-world", CaseKebab("hello world"))
		}).
		Run()
	odize.AssertNoError(t, err)
}

func TestCaseCamel(t *testing.T) {
	group := odize.NewGroup(t, nil)

	err := group.
		Test("should convert snake_case to camelCase", func(t *testing.T) {
			odize.AssertEqual(t, "helloWorld", CaseCamel("hello_world"))
		}).
		Test("should convert PascalCase to camelCase", func(t *testing.T) {
			odize.AssertEqual(t, "helloWorld", CaseCamel("HelloWorld"))
		}).
		Test("should convert kebab-case to camelCase", func(t *testing.T) {
			odize.AssertEqual(t, "helloWorld", CaseCamel("hello-world"))
		}).
		Test("should convert space separated string to camelCase", func(t *testing.T) {
			odize.AssertEqual(t, "helloWorld", CaseCamel("hello world"))
		}).
		Test("should handle single word", func(t *testing.T) {
			odize.AssertEqual(t, "hello", CaseCamel("Hello"))
		}).
		Test("should handle empty string", func(t *testing.T) {
			odize.AssertEqual(t, "", CaseCamel(""))
		}).
		Run()
	odize.AssertNoError(t, err)
}
