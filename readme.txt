log_fmt はログのフォーマットを正規表現で記載。
その際に item_names でさだめた項目名を使って(?P<項目名>正規表現)の形で指定する。
正規表現は https://github.com/google/re2/wiki/Syntax を参照。
date_fmts と out_date の日付フォーマットは https://pkg.go.dev/time#Layout を参照
