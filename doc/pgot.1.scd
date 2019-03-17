pgot(1)
; This is a scdoc source file for generating a man page
; See https://git.sr.ht/~sircmpwn/scdoc

# NAME

pgot- a command line template processing utility

# SYNOPSIS

_pgot_ [_option_] [_file_] _..._

# OPTIONS

*-d* _string_
	A string of json frontmatter to include

*-i* _paths_
	A colon separated list of paths to search with gotInclude

*-o* _file_
	Output file path (default "-")

# DESCRIPTION

*pgot* (process golang template) is a command line template processing utility.
It directly leverages the powerful golang template package.  *pgot* processes
files containing frontmatter followed by a template.  The frontmatter consists
of JSON formatted data, and the template that follows is free form text
containing template actions.  By convention, files to be processed by pgot are
referred to as "got" (golang template) files and have a ".got" filename
extension although this is not a requirement.

The pgot utility implements a "pgotInclude" keyword that can be defined in
frontmatter to include (aka import) other got files.  Each of these included
files can also include other got files if desired.  The number of nested
includes is limited only by memory and the number of open files allowed by the
operating system (a pretty large number).

Both the frontmatter and the template from the included got file can be
leveraged by the originating got file.  This enables pgot to be used both in
simple text processing applications as well as complex large scale text
processing.  As one example, the author of pgot used it to create a static blog
website with uniform html page headers, footers, and other features while
avoiding repetitive data entry (aka following the DRY principle).

A got file must include a frontmatter JSON section bounded by three semicolons
each on their own line before and after the JSON data.  It is acceptable to
have an empty frontmatter section or template section or both.

# EXAMPLE

A minimal (although admittedly useless) got file consists of the following:

```
;;;
;;;
```

A slightly more relevant got file:

```
;;;
{ "name" : "Jason" }
;;;
Hello world, my name is {{.name}}
```

And a more advanced example leveraged from the golang template website:

```
;;;
{
  "guests" : [
    ["Aunt Mildred", "bone china tea set", true],
    ["Uncle John", "moleskin pants", false],
    ["Cousin Rodney", "", false]
  ]
}
;;;
{{range .guests -}}
Dear {{index . 0}},
{{if index . 2}}
It was a pleasure to see you at the wedding.
{{- else}}
It is a shame you couldn't make it to the wedding.
{{- end}}
{{with index . 1 -}}
Thank you for the lovely {{.}}.
{{end}}
Best wishes,
Josie
---
{{end -}}
```

# BUGS

Bug reporting instructions:++
<https://blog.lenzplace.org/about/contact.html>

# SEE ALSO

*pgot*(5)

Website for pgot:++
<https://git.lenzplace.org/lenzj/pgot>

Other useful references:++
<https://golang.org/pkg/text/template>
<https://en.wikipedia.org/wiki/Template_processor>
<https://en.wikipedia.org/wiki/JSON>

; vim:set noet tw=80:
