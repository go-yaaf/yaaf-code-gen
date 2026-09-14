package processor_ts

// region TypeScript index file template -------------------------------------------------------------------------------

var indexTsTemplate = `
{{range .}}export * from './{{.}}';
{{end}}
`

// endregion

// region TypeScript client entry point file template ------------------------------------------------------------------

var clientTsTemplate = `
import axios, { AxiosInstance, AxiosRequestConfig } from 'axios';
import { RestUtils } from '../rest-utils';

{{range .}}import { {{ . }} } from './services';
{{end}}


// PulseClient is the main entry point of the client library.
// Create a single global instance of PulseClient in your application and use
// Factory methods to get service endpoints
export class PulseClient {

    private restSvc: RestUtils;
    private axiosInst: AxiosInstance;

    // Constructor
    constructor(apiUrl: string, apiKey?: string, accessToken?: string) {

        // Set axios config
        let axiosRequestConfiguration: AxiosRequestConfig = {
            baseURL: apiUrl,
            responseType: 'json',
            headers: {
                'Content-Type': 'application/json',
                'X-API-KEY': apiKey,
                'X-ACCESS-TOKEN': accessToken,
            },
        };
        // Set axios instance
        this.axiosInst = axios.create(axiosRequestConfiguration);
        this.restSvc = new RestUtils(this.axiosInst);

        // initialize services
		{{range .}}this.{{.}} = new {{.}}(apiUrl, this.restSvc);
		{{end}}
    }

	{{range .}}public readonly {{.}}: {{.}};
	{{end}}
}
`

// endregion
