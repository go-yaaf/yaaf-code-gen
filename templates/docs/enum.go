package docs

const ENUM_TMPL = `
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
            <h3 class="text-xs font-semibold text-slate-400 uppercase tracking-wider mb-2">Enums</h3>
			<ul class="space-y-1 text-sm">
			{{range .EnumList}}
				<a href="type_{{.Name}}.html" class="block py-1 px-1.5 rounded text-gray-300 font-medium hover:text-white">{{.Name}}</a>
			{{end}}
            </ul>
        </div>
    </nav>
</aside>

<!-- Main Content Wrapper -->
<div class="flex-1 flex flex-col min-w-0 overflow-hidden">
{{template "HEADER_TMPL" .ApiName}}
    <main class="flex-1 overflow-y-auto scroll-smooth grid grid-cols-1 divide-y xl:divide-y-0 xl:divide-x divide-slate-800">
		{{template "ENUM_HEADER_TMPL" .}}
    </main>
</div>

</body>
</html>
`

const ENUM_HEADER_TMPL = `
        <div class="bg-slate-800 p-4 sticky top-0 xl:h-screen xl:overflow-y-auto">
            <section class="space-y-4">
                <h2 class="text-xl font-bold text-white mb-4">{{.EnumInfo.Name}}</h2>
				<div class="my-8">
				{{range .EnumInfo.Docs}}<p class="text-slate-100 leading-relaxed mb-4">{{.}}</p>{{end}}
				</div>
				{{template "ENUM_VALUES_TMPL" .EnumInfo.Values}}
            </section>
        </div>
`

const ENUM_VALUES_TMPL = `
<div class="text-sm mt-3 divide-y divide-slate-800/60">
	<table class="table-auto w-full">
		<tbody>
			{{range . }}
			<tr>
				<td class="w-1/6 align-text-top font-mono text-emerald-400 py-2">{{.Name}}</td>
				<td class="w-1/6 px-6 align-text-top font-mono text-blue-200">{{.Value}}</td>
				<td class="align-text-top">
					{{range .Docs}}
						<span class="text-gray-100 leading-relaxed text-sm mb-2">{{.}}</span><br>
					{{end}}
				</td>
			<tr>
			{{end}}
		</tbody>
	</table>
</div>
`
