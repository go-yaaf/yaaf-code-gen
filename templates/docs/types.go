package docs

const TYPES_TMPL = `
<!DOCTYPE html>
<html lang="en" class="h-full bg-slate-900 text-slate-100">
{{template "HEAD_TMPL"}}
<style>
  /* --- Pure-CSS radio switching (no JavaScript) ------------------------ */
  /* Hidden radios live inside <main>, so they precede #type_list and      */
  /* #type_fields and can drive them via the ~ sibling combinator.         */
  .type-panel { display: none; }
  .type-item  { border-left: 2px solid transparent; }

  {{range .Types}}
  /* {{.Name}}: show its value panel and highlight its list item */
  #r-{{.TsName}}:checked ~ #type_fields #v-{{.TsName}} { display: block; }
  #r-{{.TsName}}:checked ~ #type_list  label[for="r-{{.TsName}}"] {
      background-color: rgba(16, 185, 129, 0.10);
      border-left-color: #10b981;
  }
  /* Sidebar highlight: the radios live in <main>, which the sidebar        */
  /* precedes, so ~ can't reach it — hop from <body> with :has() instead.   */
  body:has(#r-{{.TsName}}:checked) aside label[for="r-{{.TsName}}"] {
      color: #ffffff;
      background-color: rgba(16, 185, 129, 0.10);
  }
  {{end}}
</style>
<body class="flex h-full overflow-hidden font-sans antialiased">

<!-- 1. Sidebar Navigation -->
<aside class="w-64 bg-slate-950 border-r border-slate-800 flex flex-col hidden md:flex">
    <div class="p-4 border-b border-slate-800">
        <h1 class="text-lg font-bold text-white tracking-wide">API Version</h1>
        <span class="text-xs text-emerald-400 font-mono">{{.Version}}</span>
    </div>
    <nav class="flex-1 overflow-y-auto p-4 space-y-6">
        <div>
            <h3 class="text-xs font-semibold text-slate-400 uppercase tracking-wider mb-2">Types</h3>
            <ul class="space-y-1 text-sm">
				{{range .Types}}
                    <label for="r-{{.TsName}}" class="flex items-center py-1 px-1.5 rounded text-slate-300 hover:text-white cursor-pointer transition-colors">
                        {{.Name}}
                    </label>
				{{end}}
            </ul>
        </div>
    </nav>
</aside>

<!-- Main Content Wrapper -->
<div class="flex-1 flex flex-col min-w-0 overflow-hidden">
{{template "HEADER_TMPL" .Name}}
    <main class="flex-1 overflow-y-auto grid grid-cols-1 xl:grid-cols-2 divide-y xl:divide-y-0 xl:divide-x divide-slate-800">

		<!-- Hidden radio group: one per type, first is the default selection -->
		{{range $i, $e := .Types}}
        <input type="radio" name="type" id="r-{{$e.TsName}}" class="hidden" {{if eq $i 0}}checked{{end}}>
		{{end}}

		<!-- Left Column: Types List (labels drive the radios) -->
        <div id="type_list" class="p-6 md:p-12 space-y-4 max-w-3xl">
			{{range .Types}}
            <label for="r-{{.TsName}}" id="{{.TsName}}"
                   class="type-item block cursor-pointer scroll-mt-16 rounded-lg p-4 -mx-2 hover:bg-slate-800/60 transition-colors">
                <h2 class="text-xl font-bold text-white mb-2">{{.Name}}</h2>
				{{range .Docs}}<p class="text-slate-300 leading-relaxed mb-1">{{.}}</p>{{end}}
            </label>
			{{end}}
        </div>

        <!-- Right Column: Selected Type Fields -->
        <div id="type_fields" class="bg-slate-950 p-6 md:p-12 sticky top-0 xl:h-screen xl:overflow-y-auto">
			{{range .Types}}
            <section id="v-{{.TsName}}" class="type-panel space-y-4">
                <h2 class="text-xl font-bold text-white mb-4">{{.Name}}</h2>
                <div class="divide-y divide-slate-800/60 text-sm">
					{{range .Fields}}
                    <div class="py-3">
                        <!-- Field name and type on one line -->
                        <div class="flex gap-x-2 w-full">
                            <span class="w-1/2 font-mono font-semibold text-emerald-400">{{.Name}}</span>
                            <span class="w-1/2 font-mono text-slate-200">{{.Type}}</span>
                        </div>
                        <!-- Description lines below -->
						{{range .Docs}}
                        <div class="text-slate-400 leading-relaxed mt-1">{{.}}</div>
						{{end}}
                    </div>
					{{end}}
                </div>
            </section>
			{{end}}
        </div>
    </main>
</div>

</body>
</html>
`
