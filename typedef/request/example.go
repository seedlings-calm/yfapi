package request

/**
len：等于参数值，例如len=10；
max：小于等于参数值，例如max=10；
min：大于等于参数值，例如min=10；
eq：等于参数值，注意与len不同。对于字符串，eq约束字符串本身的值，而len约束字符串长度。例如eq=10；
ne：不等于参数值，例如ne=10；
gt：大于参数值，例如gt=10；
gte：大于等于参数值，例如gte=10；
lt：小于参数值，例如lt=10；
lte：小于等于参数值，例如lte=10；
oneof：只能是列举出的值其中一个，这些值必须是数值或字符串，以空格分隔，如果字符串中有空格，将字符串用单引号包围，例如oneof=red green。

contains=：包含参数子串，例如contains=email；
containsany：包含参数中任意的 UNICODE 字符，例如containsany=abcd；
containsrune：包含参数表示的 rune 字符，例如containsrune=☻；
excludes：不包含参数子串，例如excludes=email；
excludesall：不包含参数中任意的 UNICODE 字符，例如excludesall=abcd；
excludesrune：不包含参数表示的 rune 字符，excludesrune=☻；
startswith：以参数子串为前缀，例如startswith=hello；
endswith：以参数子串为后缀，例如endswith=bye。

-：跳过该字段，不检验；
|：使用多个约束，只需要满足其中一个，例如rgb|rgba；
required：字段必须设置，不能为默认值；
omitempty：如果字段未设置，则忽略它。


nil：没有错误；
InvalidValidationError：输入参数错误；
ValidationErrors：字段违反约束。


*/

type exampleReq struct {
	Id   string `json:"id" form:"id" validate:"min=6,max=10"` //必须，最短六位 最长10位
	Name string `json:"name" validate:"len=5"`                //长度等于5
	Sex  string `json:"sex" validate:"eq=男"`                  //字符串必须等于男
}
