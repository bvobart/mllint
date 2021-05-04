package cqlinters

import "fmt"

// PylintMessage represents a Pylint error / warning message (in JSON)
type PylintMessage struct {
	Path      string      `json:"path" yaml:"path"`
	Line      int         `json:"line" yaml:"line"`
	Column    int         `json:"column" yaml:"column"`
	Type      MessageType `json:"type" yaml:"type"`
	Symbol    string      `json:"symbol" yaml:"symbol"`
	Message   string      `json:"message" yaml:"message"`
	MessageID string      `json:"message-id" yaml:"symbolId"`
}

func (msg PylintMessage) String() string {
	return fmt.Sprintf("%s:%d:%d - %s (%s)", msg.Path, msg.Line, msg.Column, msg.Message, msg.MessageID)
}

// MessageType is the type of Pylint message that is emitted
// See: https://code.visualstudio.com/docs/python/linting#_pylint
type MessageType string

const (
	// TypeConvention indicates a programming standard violation, i.e. stylistic issue.
	TypeConvention MessageType = "convention"

	// TypeRefactor indicates a bad code smell
	TypeRefactor MessageType = "refactor"

	// TypeWarning can include various Python-specific warnings
	TypeWarning MessageType = "warning"

	// TypeError is for "likely code bugs" that will probably definitely give bugs
	TypeError MessageType = "error"

	// TypeFatal is for errors preventing further Pylint processing.
	TypeFatal MessageType = "fatal"
)

// MessageTypes is the list of all message types emitted by Pylint, but then as strings for easy appending to a string array.
var MessageTypes = []string{
	string(TypeConvention),
	string(TypeRefactor),
	string(TypeWarning),
	string(TypeError),
	string(TypeFatal),
}

// Symbols is a list of all 282 linting warnings / errors / messages that Pylint emits
// in its default configuration.
// To regenerate this list, use `pylint --list-msgs-enabled`, then copy the enabled messages
// and use `cut -d " " -f1` to cut off the symbol ID.
var Symbols = []string{
	"abstract-class-instantiated",
	"abstract-method",
	"access-member-before-definition",
	"anomalous-backslash-in-string",
	"anomalous-unicode-escape-in-string",
	"arguments-differ",
	"arguments-out-of-order",
	"assert-on-string-literal",
	"assert-on-tuple",
	"assign-to-new-keyword",
	"assigning-non-slot",
	"assignment-from-no-return",
	"assignment-from-none",
	"astroid-error",
	"attribute-defined-outside-init",
	"bad-classmethod-argument",
	"bad-except-order",
	"bad-exception-context",
	"bad-format-character",
	"bad-format-string",
	"bad-format-string-key",
	"bad-indentation",
	"bad-mcs-classmethod-argument",
	"bad-mcs-method-argument",
	"bad-open-mode",
	"bad-option-value",
	"bad-reversed-sequence",
	"bad-staticmethod-argument",
	"bad-str-strip-call",
	"bad-string-format-type",
	"bad-super-call",
	"bad-thread-instantiation",
	"bare-except",
	"binary-op-exception",
	"blacklisted-name",
	"boolean-datetime",
	"broad-except",
	"c-extension-no-member",
	"catching-non-exception",
	"cell-var-from-loop",
	"chained-comparison",
	"class-variable-slots-conflict",
	"comparison-with-callable",
	"comparison-with-itself",
	"confusing-with-statement",
	"consider-iterating-dictionary",
	"consider-merging-isinstance",
	"consider-swap-variables",
	"consider-using-dict-comprehension",
	"consider-using-enumerate",
	"consider-using-get",
	"consider-using-in",
	"consider-using-join",
	"consider-using-set-comprehension",
	"consider-using-sys-exit",
	"consider-using-ternary",
	"continue-in-finally",
	"cyclic-import",
	"dangerous-default-value",
	"deprecated-method",
	"deprecated-module",
	"dict-iter-missing-items",
	"duplicate-argument-name",
	"duplicate-bases",
	"duplicate-code",
	"duplicate-except",
	"duplicate-key",
	"duplicate-string-formatting-argument",
	"empty-docstring",
	"eval-used",
	"exec-used",
	"expression-not-assigned",
	"f-string-without-interpolation",
	"fatal",
	"fixme",
	"format-combined-specification",
	"format-needs-mapping",
	"function-redefined",
	"global-at-module-level",
	"global-statement",
	"global-variable-not-assigned",
	"global-variable-undefined",
	"implicit-str-concat",
	"import-error",
	"import-outside-toplevel",
	"import-self",
	"inconsistent-mro",
	"inconsistent-quotes",
	"inconsistent-return-statements",
	"inherit-non-class",
	"init-is-generator",
	"invalid-all-object",
	"invalid-bool-returned",
	"invalid-bytes-returned",
	"invalid-characters-in-docstring",
	"invalid-envvar-default",
	"invalid-envvar-value",
	"invalid-format-index",
	"invalid-format-returned",
	"invalid-getnewargs-ex-returned",
	"invalid-getnewargs-returned",
	"invalid-hash-returned",
	"invalid-index-returned",
	"invalid-length-hint-returned",
	"invalid-length-returned",
	"invalid-metaclass",
	"invalid-name",
	"invalid-overridden-method",
	"invalid-repr-returned",
	"invalid-sequence-index",
	"invalid-slice-index",
	"invalid-slots",
	"invalid-slots-object",
	"invalid-star-assignment-target",
	"invalid-str-returned",
	"invalid-unary-operand-type",
	"isinstance-second-argument-not-valid-type",
	"keyword-arg-before-vararg",
	"len-as-condition",
	"line-too-long",
	"literal-comparison",
	"logging-format-interpolation",
	"logging-format-truncated",
	"logging-fstring-interpolation",
	"logging-not-lazy",
	"logging-too-few-args",
	"logging-too-many-args",
	"logging-unsupported-format",
	"lost-exception",
	"method-check-failed",
	"method-hidden",
	"misplaced-bare-raise",
	"misplaced-comparison-constant",
	"misplaced-format-function",
	"misplaced-future",
	"missing-class-docstring",
	"missing-final-newline",
	"missing-format-argument-key",
	"missing-format-attribute",
	"missing-format-string-key",
	"missing-function-docstring",
	"missing-kwoa",
	"missing-module-docstring",
	"missing-parentheses-for-call-in-test",
	"mixed-format-string",
	"mixed-line-endings",
	"multiple-imports",
	"multiple-statements",
	"no-classmethod-decorator",
	"no-else-break",
	"no-else-continue",
	"no-else-raise",
	"no-else-return",
	"no-init",
	"no-member",
	"no-method-argument",
	"no-name-in-module",
	"no-self-argument",
	"no-self-use",
	"no-staticmethod-decorator",
	"no-value-for-parameter",
	"non-ascii-name",
	"non-iterator-returned",
	"non-parent-init-called",
	"non-str-assignment-to-dunder-name",
	"nonexistent-operator",
	"nonlocal-and-global",
	"nonlocal-without-binding",
	"not-a-mapping",
	"not-an-iterable",
	"not-async-context-manager",
	"not-callable",
	"not-context-manager",
	"not-in-loop",
	"notimplemented-raised",
	"parse-error",
	"pointless-statement",
	"pointless-string-statement",
	"possibly-unused-variable",
	"preferred-module",
	"property-with-parameters",
	"protected-access",
	"raise-missing-from",
	"raising-bad-type",
	"raising-format-tuple",
	"raising-non-exception",
	"redeclared-assigned-name",
	"redefine-in-handler",
	"redefined-argument-from-local",
	"redefined-builtin",
	"redefined-outer-name",
	"redundant-keyword-arg",
	"redundant-unittest-assert",
	"reimported",
	"relative-beyond-top-level",
	"repeated-keyword",
	"return-arg-in-generator",
	"return-in-init",
	"return-outside-function",
	"self-assigning-variable",
	"self-cls-assignment",
	"shallow-copy-environ",
	"signature-differs",
	"simplifiable-if-expression",
	"simplifiable-if-statement",
	"simplify-boolean-expression",
	"single-string-used-for-slots",
	"singleton-comparison",
	"star-needs-assignment-target",
	"stop-iteration-return",
	"subprocess-popen-preexec-fn",
	"subprocess-run-check",
	"super-init-not-called",
	"super-with-arguments",
	"superfluous-parens",
	"syntax-error",
	"too-few-format-args",
	"too-few-public-methods",
	"too-many-ancestors",
	"too-many-arguments",
	"too-many-boolean-expressions",
	"too-many-branches",
	"too-many-format-args",
	"too-many-function-args",
	"too-many-instance-attributes",
	"too-many-lines",
	"too-many-locals",
	"too-many-nested-blocks",
	"too-many-public-methods",
	"too-many-return-statements",
	"too-many-star-expressions",
	"too-many-statements",
	"trailing-comma-tuple",
	"trailing-newlines",
	"trailing-whitespace",
	"truncated-format-string",
	"try-except-raise",
	"unbalanced-tuple-unpacking",
	"undefined-all-variable",
	"undefined-loop-variable",
	"undefined-variable",
	"unexpected-keyword-arg",
	"unexpected-line-ending-format",
	"unexpected-special-method-signature",
	"ungrouped-imports",
	"unhashable-dict-key",
	"unidiomatic-typecheck",
	"unnecessary-comprehension",
	"unnecessary-lambda",
	"unnecessary-pass",
	"unnecessary-semicolon",
	"unneeded-not",
	"unpacking-non-sequence",
	"unreachable",
	"unrecognized-inline-option",
	"unsubscriptable-object",
	"unsupported-assignment-operation",
	"unsupported-binary-operation",
	"unsupported-delete-operation",
	"unsupported-membership-test",
	"unused-argument",
	"unused-format-string-argument",
	"unused-format-string-key",
	"unused-import",
	"unused-variable",
	"unused-wildcard-import",
	"used-before-assignment",
	"used-prior-global-declaration",
	"useless-else-on-loop",
	"useless-import-alias",
	"useless-object-inheritance",
	"useless-return",
	"useless-super-delegation",
	"using-constant-test",
	"wildcard-import",
	"wrong-exception-operation",
	"wrong-import-order",
	"wrong-import-position",
	"wrong-spelling-in-comment",
	"wrong-spelling-in-docstring",
	"yield-inside-async-function",
	"yield-outside-function",
}
