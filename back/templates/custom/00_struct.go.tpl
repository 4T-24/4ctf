{{- $alias := .Aliases.Table .Table.Name -}}
{{- $orig_tbl_name := .Table.Name -}}

// {{$alias.DownSingular}}View is the view displayed to clients.
type {{$alias.DownSingular}}View struct {
	{{- range $index, $column := .Table.Columns -}}
	{{- $colAlias := $alias.Column $column.Name -}}
	{{- $orig_col_name := $column.Name -}}

	{{if ignore $orig_tbl_name $orig_col_name $.TagIgnore -}}
		{{if eq (index $column.Type 0) '*'}}
			{{$colAlias}} {{$column.Type}} `{{generateIgnoreTags $.Tags}}boil:"{{$column.Name}}" json:"-" toml:"-" yaml:"-"`
		{{- else -}}
			{{$colAlias}} *{{$column.Type}} `{{generateIgnoreTags $.Tags}}boil:"{{$column.Name}}" json:"-" toml:"-" yaml:"-"`
		{{- end -}}
	{{ else -}}

	{{- /* render column alias and column type */ -}}
	{{if eq (index $column.Type 0) '*'}}
		{{ $colAlias }} {{ $column.Type -}}
	{{- else -}}
		{{ $colAlias }} *{{ $column.Type -}}
	{{- end -}}

	{{- /*
	  handle struct tags
	  StructTagCasing will be replaced with $.StructTagCases
	  however we need to keep this backward compatible
	  $.StructTagCasing will only be used when it's set to "alias"
    */ -}}
	`
	{{- if eq $.StructTagCasing "alias" -}}
	    {{- generateTags $.Tags $colAlias -}}
	    {{- generateTagWithCase "boil" $column.Name $colAlias "alias" false -}}
	    {{- generateTagWithCase "json" $column.Name $colAlias "alias" $column.Nullable -}}
	    {{- generateTagWithCase "toml" $column.Name $colAlias "alias" false -}}
	    {{- trim (generateTagWithCase "yaml" $column.Name $colAlias "alias" $column.Nullable) -}}
		{{- if ne $column.Comment "" -}}
			{{ print " " $column.Comment }}
		{{- end -}}
	{{- else -}}
	    {{- generateTags $.Tags $column.Name }}
	    {{- generateTagWithCase "boil" $column.Name $colAlias $.StructTagCases.Boil false -}}
	    {{- generateTagWithCase "json" $column.Name $colAlias $.StructTagCases.Json $column.Nullable -}}
	    {{- generateTagWithCase "toml" $column.Name $colAlias $.StructTagCases.Toml false -}}
	    {{- trim (generateTagWithCase "yaml" $column.Name $colAlias $.StructTagCases.Yaml $column.Nullable) -}}
		{{- if ne $column.Comment "" -}}
			{{ print " " $column.Comment }}
		{{- end -}}
	{{- end -}}
	`
	{{ end -}}
	{{ end -}}
}

func (o *{{$alias.UpSingular}}) View() *{{$alias.DownSingular}}View {
	return &{{$alias.DownSingular}}View{
		{{- range $index, $column := .Table.Columns -}}
		{{- $colAlias := $alias.Column $column.Name -}}
			{{ print "" }}
			{{if eq (index $column.Type 0) '*'}}
				{{$colAlias}}: o.{{$colAlias}},
			{{- else -}}
				{{$colAlias}}: &o.{{$colAlias}},
			{{- end -}}
		{{ end -}}
		{{ print "" }}
	}
}
