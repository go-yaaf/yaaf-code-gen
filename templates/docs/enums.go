package docs

const ENUMS_TMPL = `
<!DOCTYPE html>
<html lang="en" class="h-full bg-slate-900 text-slate-100">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>API Documentation</title>
	<link rel="stylesheet" href="styles.css">


<style>
  /* --- Pure-CSS radio switching + scroll-to (no JavaScript) ------------ */
  .enum-panel { display: none; }
  .enum-item  { border-left: 2px solid transparent; }

  /* Each radio lives INSIDE its list <label>. Clicking any label (sidebar  */
  /* or list) focuses the radio, and the browser scrolls the focused        */
  /* control into view — this is what scrolls the list, with no JS.         */
  /* Must NOT use display:none / visibility:hidden (both kill focus);       */
  /* appearance:none + height:0 drops it from layout, opacity:0 hides it.   */
  .enum-radio {
      appearance: none; -webkit-appearance: none;
      display: block; width: 0; height: 0;
      margin: 0; padding: 0; border: 0;
      opacity: 0; outline: none;
      scroll-margin-top: 4rem;   /* land below the sticky header on scroll */
  }

  {{range .EnumList}}
  /* {{.Name}}: show value panel + highlight list row. Radios are no longer */
  /* preceding siblings, so :has() replaces the ~ sibling combinator.       */
  main:has(#r-{{.TsName}}:checked) #v-{{.TsName}} { display: block; }
  .enum-item:has(#r-{{.TsName}}:checked) {
      background-color: rgba(16, 185, 129, 0.10);
      border-left-color: #10b981;
  }
  body:has(#r-{{.TsName}}:checked) aside label[for="r-{{.TsName}}"] {
      color: #ffffff;
      background-color: rgba(16, 185, 129, 0.10);
  }
  {{end}}
</style>
</head>
<body class="flex h-full overflow-hidden font-sans antialiased">

<!-- 1. Sidebar Navigation -->
<aside class="w-64 bg-slate-950 border-r border-slate-800 flex flex-col hidden md:flex">
    <div class="p-4 border-b border-slate-800">
        <h1 class="text-lg font-bold text-white tracking-wide">API Version</h1>
        <span class="text-xs text-emerald-400 font-mono">{{.ApiVersion}}</span>
    </div>
    <nav class="flex-1 overflow-y-auto p-4 space-y-6">
        <div>
            <h3 class="text-xs font-semibold text-slate-400 uppercase tracking-wider mb-2">Enums</h3>
			<ul class="space-y-1 text-sm">
			{{range .EnumList}}
				<label for="r-{{.TsName}}" class="block py-1 px-1.5 rounded text-gray-300 hover:text-white cursor-pointer">
					{{.Name}}
				</label>
				{{end}}
            </ul>
        </div>
    </nav>
</aside>

<!-- Main Content Wrapper -->
<div class="flex-1 flex flex-col min-w-0 overflow-hidden">
{{template "HEADER_TMPL" .ApiName}}
    <main class="flex-1 overflow-y-auto scroll-smooth grid grid-cols-1 xl:grid-cols-2 divide-y xl:divide-y-0 xl:divide-x divide-slate-800">

		<!-- Left Column: Enums List. Each label wraps its own radio, which is  -->
		<!-- both the selection state and the scroll anchor (see .enum-radio).  -->
        <div id="enum_list" class="p-6 md:p-12 space-y-4 max-w-3xl">
			{{range $i, $e := .EnumList}}
            <label for="r-{{$e.TsName}}" id="{{$e.TsName}}"
                   class="enum-item block cursor-pointer scroll-mt-16 rounded-lg p-4 -mx-2 hover:bg-slate-800/60 transition-colors">
                <input type="radio" name="enum" id="r-{{$e.TsName}}" class="enum-radio" tabindex="-1" {{if eq $i 0}}checked{{end}}>
                <h2 class="text-xl font-bold text-white mb-2">{{$e.Name}}</h2>
				{{range $e.Docs}}<p class="text-slate-300 leading-relaxed mb-1">{{.}}</p>{{end}}
            </label>
			{{end}}
        </div>

        <!-- Right Column: Selected Enum Values -->
        <div id="enum_values" class="bg-slate-950 p-6 md:p-12 sticky top-0 xl:h-screen xl:overflow-y-auto">
			{{range .EnumList}}
            <section id="v-{{.TsName}}" class="enum-panel space-y-4">
                <h2 class="text-xl font-bold text-white mb-4">{{.Name}}</h2>
                <div class="divide-y divide-slate-800/60 text-sm">
					{{range .Values}}
                    <div class="py-3">
                        <!-- Field name and type on one line -->
						<div class="flex gap-x-2 w-full">
                            <span class="w-1/2 font-mono font-semibold text-emerald-400">{{.Name}}</span>
                            <span class="w-1/2 font-mono text-slate-200">{{.Value}}</span>
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
