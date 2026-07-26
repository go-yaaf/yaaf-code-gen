package docs

const HEAD_TMPL = `
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>API Documentation</title>
    {{template "STYLES_TMPL"}}
</head>
`

const HEADER_TMPL = `
    <header class="bg-slate-950 border-b border-slate-800 p-4 flex items-center justify-between">
        <h2 class="text-md font-bold text-white">{{.}}</h2>
        <div class="flex" style="gap: 2rem">
            <a href="index.html" class="text-slate-300 hover:text-white">Services</a>
            <a href="types.html" class="text-slate-300 hover:text-white">Data Types</a>
            <a href="enums.html" class="text-slate-300 hover:text-white">Enums</a>
        </div>
    </header>
`

const INDEX_TMPL = `
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
            <h3 class="text-xs font-semibold text-slate-400 uppercase tracking-wider mb-2">Getting Started</h3>
            <ul class="space-y-1 text-sm">
                <li><a href="#authentication" class="block py-1 text-emerald-300 font-medium">Authentication</a></li>
                <li><a href="#errors" class="block py-1 text-slate-300 hover:text-white">Errors</a></li>
            </ul>
        </div>
        <div>
			{{range .ServiceGroups}}
            <h3 class="text-xs font-semibold text-slate-400 uppercase tracking-wider mt-8">{{.Name}}</h3>
			<ul class="text-sm font-mono">
				{{range .List}}
                <li>
                    <a href="#post-users" class="flex items-center py-1 text-slate-300 hover:text-white">
                        <span class="px-1.5">{{.Name}}</span>
                    </a>
                </li>
				{{end}}
			</ul>
            <div class="tracking-wider mb-4"></div>
			{{end}}
        </div>
    </nav>
</aside>

<!-- Main Content Wrapper -->
<div class="flex-1 flex flex-col min-w-0 overflow-hidden">

{{template "HEADER_TMPL" .ApiName}}

    <main class="flex-1 overflow-y-auto grid grid-cols-1 xl:grid-cols-2 divide-y xl:divide-y-0 xl:divide-x divide-slate-800">
        <div class="p-6 md:p-12 space-y-16 max-w-3xl">
            <section id="authentication" class="scroll-mt-16">
                <h2 class="text-2xl font-bold text-white mb-4">Authentication</h2>
                <p class="text-slate-300 leading-relaxed mb-4">Authenticate your requests using API key and access token (can be also provided as Bearer token) in the Authorization header.</p>
                <div class="bg-slate-950 border border-slate-800 rounded-lg p-3 font-mono text-xs text-slate-300 w-xl">
                    <span>X-API-KEY: </span><span class="mx-8 text-emerald-400">YOUR_API_KEY</span>
                </div>
                <div class="bg-slate-950 border border-slate-800 rounded-lg p-3 font-mono text-xs text-slate-300 w-2xl">
                    <span class="mx-4">X-ACCESS-TOKEN: </span><span class="mx-8 text-emerald-400">YOUR_ACCESS_TOKEN</span><br><br>
                    <span class="my-4 mx-8 text-sm font-bold text-white">or</span><br><br>
                    <span class="mx-4">Authorization: Bearer </span><span class="mx-8 text-emerald-400">YOUR_ACCESS_TOKEN</span>
                </div>
            </section>
            <section id="headers" class="scroll-mt-16">
                <h2 class="text-xl font-bold text-white mb-4">Authentication Headers</h2>
                <p class="text-slate-300 leading-relaxed mb-4">
                        Any consumer MUST include API key provided in the X-API-KEY header.
                        API Keys are used to identify the services consumer (application/service) type and grant access to relevant resources only.
                </p>
                <p class="text-slate-300 leading-relaxed mb-4">
                        After a successful login, any subsequent call MUST include the access token provided in X-ACCESS-TOKEN or Authorization header.
                        The access token is used to identify the user/service who is using API resources.
                </p>
            </section>
            <section id="errors" class="scroll-mt-16">
                <h2 class="text-xl font-bold text-white mb-4">Errors</h2>
                <p class="text-slate-300 leading-relaxed mb-4">
                    Al errors repesented by the error response object in the return body in addition to the HTTP code: 500
                </p>
                <pre class="bg-slate-900 border border-slate-800 rounded-lg p-4 font-mono text-xs text-emerald-400 overflow-x-auto">
{
  "code": -5,
  "error": "detailed error description"
}</pre>
            </section>

        </div>
    </main>
</div>


</body>
</html>
`
