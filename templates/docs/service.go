package docs

const SERVICE_TMPL = `
<!DOCTYPE html>
<html lang="en" class="h-full bg-slate-900 text-slate-100">
{{template "HEAD_TMPL"}}
<body class="flex h-full overflow-hidden font-sans antialiased">

<style>
  /* --- Pure-CSS radio switching (no JavaScript) --------------------- */
  /* Hidden but focusable radios drive which section is visible.        */
  .srv-radio {
      appearance: none; -webkit-appearance: none;
      display: block; width: 0; height: 0;
      margin: 0; padding: 0; border: 0;
      opacity: 0; outline: none;
  }

  /* Hide every service_/method_ section by default. */
  section[id^="service_"], section[id^="method_"] { display: none; }

  /* Service: show its section + highlight the sidebar row. */
  body:has(#r-service_{{.ServiceInfo.TsName}}:checked) #service_{{.ServiceInfo.TsName}} { display: block; }
  body:has(#r-service_{{.ServiceInfo.TsName}}:checked) aside label[for="r-service_{{.ServiceInfo.TsName}}"] {
      color: #ffffff;
      background-color: rgba(16, 185, 129, 0.10);
  }

  {{range .ServiceInfo.Methods}}
  body:has(#r-method_{{.Name}}:checked) #method_{{.Name}} { display: block; }
  body:has(#r-method_{{.Name}}:checked) aside label[for="r-method_{{.Name}}"] {
      color: #ffffff;
      background-color: rgba(16, 185, 129, 0.10);
  }
  {{end}}
</style>

<!-- 1. Sidebar Navigation -->
<aside class="w-64 bg-slate-950 border-r border-slate-800 flex flex-col hidden md:flex">
    <div class="p-4 border-b border-slate-800">
        <h1 class="text-lg font-bold text-white tracking-wide">API Version</h1>
        <span class="text-xs text-emerald-400 font-mono">{{.ApiVersion}}</span>
    </div>
    <nav class="flex-1 overflow-y-auto p-4 space-y-6">
        <div>
            <h3 class="text-xs font-semibold text-slate-400 uppercase tracking-wider mb-2">Service</h3>
			<ul class="text-sm font-mono">
                <li>
					<label for="r-service_{{.ServiceInfo.TsName}}" class="block py-1 px-1.5 rounded text-gray-300 font-medium hover:text-white cursor-pointer">{{.ServiceInfo.TsName}}</label>
				</li>
            </ul>
			<br>
            <h3 class="text-xs font-semibold text-slate-400 uppercase tracking-wider mb-2">Service Methods</h3>
			<ul class="text-sm font-mono">
			{{range .ServiceInfo.Methods}}
                <li><label for="r-method_{{.Name}}" class="block py-1 px-1.5 rounded text-gray-300 font-medium hover:text-white cursor-pointer">{{.TsName}}</label></li>
				{{end}}
            </ul>
        </div>
    </nav>
</aside>

<!-- Main Content Wrapper -->
<div class="flex-1 flex flex-col min-w-0 overflow-hidden">

{{template "HEADER_TMPL" .ApiName}}

    <main class="flex-1 overflow-y-auto grid grid-cols-1 xl:grid-cols-1 divide-y xl:divide-y-0 xl:divide-x divide-slate-800 scroll-smooth">
        <!-- <div class="p-6 md:p-6 space-y-16 max-w-3xl"> -->
		<div class="p-4">

            <!-- Service sections: anchor targets for the sidebar links (#TsName) -->
            <section id="service_{{.ServiceInfo.TsName}}" class="scroll-mt-16">
                <input type="radio" name="srv_nav" id="r-service_{{.ServiceInfo.TsName}}" class="srv-radio" tabindex="-1" checked>
                <h2 class="text-2xl font-bold text-white mb-2">{{.ServiceInfo.TsName}}</h2>
				<div class="mt-2">
                    <span class="mr-4 text-sm font-semibold text-slate-300 uppercase">REST Path</span>
                    <span class="ml-4 text-base font-mono text-emerald-400">{{.ServiceInfo.Path}}</span>
				</div>
				<div class="mt-2">
                    <span class="mr-4 text-sm font-semibold text-slate-300 uppercase">Client Lib</span>
                    <span class="ml-4 text-base font-mono text-emerald-400">{{.ServiceInfo.Name}}</span>
				</div>
				<div class="my-8">
				{{range .ServiceInfo.Docs}}<p class="text-slate-100 leading-relaxed mb-4">{{.}}</p>{{end}}
				</div>

                <h2 class="text-xl font-bold text-white">Service Methods</h2>
				{{range .ServiceInfo.Methods}}
					<label for="r-method_{{.Name}}" class="block py-1 px-1.5 rounded text-emerald-300 font-medium hover:text-white cursor-pointer transition-colors">{{.TsName}}</label>
					{{range .Docs}}<p class="text-slate-100 leading-relaxed text-sm mb-2">{{.}}</p>{{end}}
				{{end}}
			</section>

			{{range .ServiceInfo.Methods}}
			<section id="method_{{.Name}}" class="scroll-mt-16">
				<input type="radio" name="srv_nav" id="r-method_{{.Name}}" class="srv-radio" tabindex="-1">
		
				<div class="grid grid-cols-2 xl:grid-cols-2">
					{{template "METHOD_TMPL" .}}
					{{template "EXAMPLE_TMPL" .}}
				</div>
			</section>
			{{end}}
        </div>
    </main>
</div>

</body>
</html>

`
const METHOD_TMPL = `
<div class="border border-slate-800 rounded-lg p-4 mb-4">

	<h2 class="text-xl font-bold text-white mb-2">{{.Name}}</h2>
	<div class="mt-2">
		<span class="mr-4 text-sm font-semibold text-slate-300 uppercase">REST Path</span>
		{{if eq .Method "GET"}}   <span class="ml-4 text-sm bg-blue-300/50 text-blue-200 rounded-md font-bold font-mono">{{.Method}}</span>{{end}}
		{{if eq .Method "POST"}}  <span class="ml-4 text-sm bg-indigo-500/50 text-indigo-200 rounded-md font-bold font-mono">{{.Method}}</span>{{end}}
		{{if eq .Method "PUT"}}   <span class="ml-4 text-sm bg-purple-500/50 text-purple-200 rounded-md font-bold font-mono">{{.Method}}</span>{{end}}
		{{if eq .Method "DELETE"}}<span class="ml-4 text-sm bg-red-500/50 text-red-200 rounded-md font-bold font-mono">{{.Method}}</span>{{end}}
		{{if eq .Method "PATCH"}} <span class="ml-4 text-sm bg-amber-500/50 text-amber-200 rounded-md font-bold font-mono">{{.Method}}</span>{{end}}
		<span class="ml-2 text-base font-mono text-emerald-400">{{.Path}}</span>
	</div>
	<div class="mt-2">
		<span class="mr-4 text-sm font-semibold text-slate-300 uppercase">Client Lib</span>
		<span class="ml-4 text-base font-mono text-emerald-400">{{.TsName}}</span>
	</div>
	<div class="my-8">
	{{range .Docs}}<p class="text-slate-200 leading-relaxed text-base mb-2">{{.}}</p>{{end}}
	</div>

	{{if .PathParams }}
		<h4 class="text-sm font-semibold text-slate-200 uppercase tracking-wider mb-3 mt-6">Path Parameters</h4>
		{{template "PARAMS_TMPL" .PathParams}}
	{{end}}

	{{if .QueryParams }}
		<h4 class="text-sm font-semibold text-slate-200 uppercase tracking-wider mb-3 mt-6">Query Parameters</h4>
		{{template "PARAMS_TMPL" .QueryParams}}
	{{end}}

	{{if .BodyParam }}
		<h4 class="text-sm font-semibold text-slate-200 uppercase tracking-wider mb-3 mt-6">Body Parameter</h4>
		{{template "PARAM_TMPL" .BodyParam}}
	{{end}}

	{{if .FileParam }}
		<h4 class="text-sm font-semibold text-slate-200 uppercase tracking-wider mb-3 mt-6">File Parameter</h4>
		{{template "PARAM_TMPL" .FileParam}}
	{{end}}

	{{if .Return }}
	<h4 class="text-sm font-semibold text-slate-200 uppercase tracking-wider mb-3 mt-6">Return</h4>
	<div class="border-t border-slate-800 divide-y divide-slate-800 text-sm">
		<div class="flex gap-x-2 py-1">
			<span class="w-1/6 font-mono text-emerald-400">Return</span>
			<span class="w-1/6 font-mono text-blue-200">{{.ReturnClass}}</span>
			<span class="text-slate-200">
				{{range .Return.Docs}}<p class="text-slate-100 leading-relaxed text-base mb-2">{{.}}</p>{{end}}
			</span>
		</div>
	</div>
	{{end}}
</div>
`

const EXAMPLE_TMPL = `
<div class="bg-slate-950 p-6 space-y-16 sticky top-0 xl:h-screen xl:overflow-y-auto">
	<div class="space-y-4 pt-12">
		<div class="flex items-center justify-between border-b border-slate-800 pb-2">
			<span class="text-xs font-semibold text-slate-400 uppercase">Request Example</span>
			<span class="text-xs font-mono text-slate-500">HTTP</span>
		</div>
		<pre class="bg-slate-900 border border-slate-800 rounded-lg p-4 font-mono text-xs text-slate-300 overflow-x-auto">
{{.Method}} https://path/to/api/v1/{{.Path}} \
  -H "X-API-KEY: your_api_key"
  -H "X-ACCESS-TOKEN: your_access_token"

		</pre>
	</div>

	<div class="space-y-4 pt-12">
		<div class="flex items-center justify-between border-b border-slate-800 pb-2">
			<span class="text-xs font-semibold text-slate-400 uppercase">Response</span>
			<span class="text-xs font-mono text-emerald-400">200 OK</span>
		</div>
		<pre class="bg-slate-900 border border-slate-800 rounded-lg p-4 font-mono text-xs text-emerald-400 overflow-x-auto">
{
  "object": "list",
  "data": [
    {
      "id": "usr_9J2x",
      "name": "Jane Doe",
      "email": "jane@example.com"
    }
  ],
  "has_more": false
}
		</pre>
	</div>
</div>
`
