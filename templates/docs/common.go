package docs

const PARAM_TMPL = `
<div class="text-sm mt-3 divide-y divide-slate-800/60">
	<table class="table-auto border-spacing-y">
		<tbody>
			<tr>
				<td class="align-text-top font-mono text-emerald-400">{{.Name}}</td>
				<td class="px-6 align-text-top font-mono text-blue-200">{{.Type}}{{if .IsArray}}[]{{end}}</td>
				<td class="align-text-top text-slate-200">
					{{range .Docs}}
						<span class="text-slate-100 leading-relaxed text-base mb-2">{{.}}</span><br>
					{{end}}
				</td>
			<tr>
		</tbody>
	</table>
</div>
`

const PARAMS_TMPL = `
<div class="text-sm mt-3 divide-y divide-slate-800/60">
	<table class="table-auto border-spacing-y-2">
		<tbody>
			{{range . }}
			<tr>
				<td class="align-text-top font-mono text-emerald-400">{{.Name}}</td>
				<td class="px-6 align-text-top font-mono text-blue-200">{{.Type}}{{if .IsArray}}[]{{end}}</td>
				<td class="align-text-top text-slate-200">
					{{range .Docs}}
						<span class="text-slate-100 leading-relaxed text-base mb-2">{{.}}</span><br>
					{{end}}
				</td>
			<tr>
			{{end}}
		</tbody>
	</table>
</div>
`
