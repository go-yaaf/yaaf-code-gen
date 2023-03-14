import { Injectable, Inject } from '@angular/core';
import { RestUtil } from '../../utils/rest-util';
import { CoreConfig } from '../../config';

{{.Methods | addServiceImports}}

/**{{range .Docs}}
 * {{.}}{{end}} 
 */
@Injectable()
export class {{.Name}} {

  // URL to web api
  private baseUrl = '{{.Path}}';

  /**
   * Class constructor
   */
  constructor(@Inject('config') private config: CoreConfig, private rest: RestUtil) {
    this.baseUrl = this.config.api + this.baseUrl;
  }

{{range .Methods}}
  /**{{range .Docs}}
   * {{.}}{{end}}
   */
  {{.Name | toCamelCase}}({{. | handleMethodParams}}) {
    {{. | methodContent}}
  }
{{end}}
}
