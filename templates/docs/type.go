package docs

const TYPE_TMPL = `
<!DOCTYPE html>
<html lang="en" class="h-full bg-slate-900 text-slate-100">
{{template "HEAD_TMPL"}}
<body class="flex h-full overflow-hidden font-sans antialiased">

<!-- 1. Sidebar Navigation -->
<aside class="w-64 bg-slate-950 border-r border-slate-800 flex flex-col hidden md:flex">
    <div class="p-4 border-b border-slate-800">
        <h1 class="text-lg font-bold text-white tracking-wide">API Version</h1>
        <span class="text-xs text-emerald-400 font-mono">{{.ApiVersion}}</span>
    </div>
    <nav class="flex-1 overflow-y-auto p-4 space-y-6">
        <div>
			{{range .ClassGroups}}
			<br>
            <h3 class="text-xs font-semibold text-slate-400 uppercase tracking-wider my-4">{{.Name}}</h3>
            <ul class="space-y-1 text-sm">
				{{range .List}}
                    <a href="type_{{.Name}}.html" class="block py-1 px-1.5 rounded text-gray-300 font-medium hover:text-white">{{.Name}}</a>
				{{end}}
            </ul>
			{{end}}
        </div>
    </nav>
</aside>

<!-- Main Content Wrapper -->
<div class="flex-1 flex flex-col min-w-0 overflow-hidden">
{{template "HEADER_TMPL" .ApiName}}
    <main class="flex-1 overflow-y-auto scroll-smooth grid grid-cols-2 divide-y xl:divide-y-0 xl:divide-x divide-slate-800">
		{{template "TYPE_FIELDS_TMPL" .}}
		{{template "TYPE_EXAMPLE_TMPL" .}}
    </main>
</div>

</body>
</html>
`

const TYPE_FIELDS_TMPL = `
        <div class="bg-slate-800 p-4 sticky top-0 xl:h-screen xl:overflow-y-auto">
            <section class="space-y-4">
                <h2 class="text-xl font-bold text-white mb-4">{{.ClassInfo.Name}}</h2>
				<div class="my-8">
				{{range .ClassInfo.Docs}}<p class="text-slate-100 leading-relaxed mb-4">{{.}}</p>{{end}}
				</div>
				{{template "FIELDS_TMPL" .ClassInfo.Fields}}
<!--
				<div class="grid grid-cols-2 xl:grid-cols-2">
					{{template "FIELDS_TMPL" .ClassInfo.Fields}}
				</div>
-->
            </section>
        </div>
`

const FIELDS_TMPL = `
<div class="text-sm mt-3 divide-y divide-slate-800/60">
	<table class="table-auto w-full">
		<tbody>
			{{range . }}
			<tr>
				<td class="align-text-top font-mono text-emerald-400 py-2">{{.Name}}</td>
				<td class="px-6 align-text-top font-mono text-blue-200">
				{{if .IsComplex}}
                    <a href="type_{{.Type}}.html" class="text-gray-200 hover:text-white">{{.Type}}</a>
				{{ else }}
                    <span class="text-gray-200">{{.Type}}</a>
				{{end}}
				<td class="align-text-top">
					{{range .Docs}}
						<span class="text-gray-100 leading-relaxed text-base mb-2">{{.}}</span><br>
					{{end}}
				</td>
				<td>{{if .IsMap}} [Map] {{end}}{{if .IsGeneric}} [Generic] {{end}}</td>
			<tr>
			{{end}}
		</tbody>
	</table>
</div>
`

const TYPE_EXAMPLE_TMPL = `
<div class="bg-slate-950 p-4 sticky top-0 xl:h-screen xl:overflow-y-auto">
	<div class="space-y-4 pt-0">
		<div class="flex items-center justify-between border-b border-slate-800 pb-2">
			<span class="text-xs font-semibold text-slate-400 uppercase">Struct Example</span>
			<span class="text-xs font-mono text-emerald-400">Json</span>
		</div>
		<pre class="bg-slate-900 border border-slate-800 rounded-lg p-4 font-mono text-xs text-emerald-400 overflow-x-auto">
{
  "field_1": "list",
  "id": "usr_9J2x",
  "name": "Jane Doe",
  "email": "jane@example.com"
  "has_more": false
}
		</pre>
	</div>
</div>
`
