import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders, HttpRequest } from '@angular/common/http';
import { Observable, throwError } from 'rxjs';

const httpOptions = {
  headers: new HttpHeaders({
    'Content-Type':  'application/json',
    'Cache-Control': 'no-cache, no-store',
    'Access-Control-Allow-Origin': '*',
    'Access-Control-Allow-Headers': 'Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, X-API-KEY, X-ACCESS-TOKEN, X-TIMEZONE, accept, origin, Cache-Control, X-Requested-With, Authorization, Content-Disposition, Content-Filename',
    'Access-Control-Exposed-Headers': 'X-API-KEY, X-ACCESS-TOKEN, X-TIMEZONE, Content-Disposition, Content-Filename',
  })
};


// Utility class for all REST services with common functions
@Injectable({
  providedIn: 'root'
})
export class RestUtils {

  // Constructor with injected authentication service
  constructor(private http: HttpClient) { }

  // Upload is HTTP POST action but the body is File object
  upload<T>(file: File, url: string, ...params: string[]) {

    const resourceUrl = this.buildUrl(url, ...params);

    const formData: FormData = new FormData();
    formData.append('fileKey', file, file.name);

    const req = new HttpRequest('POST', resourceUrl, formData, {
      reportProgress: false,
      responseType: 'json',
    });
    return this.http.request<T>(req);
  }

  // Download is HTTP GET action but the content is a blob
  download(fileName: string, url: string, ...params: string[]) {
    const resourceUrl = this.buildUrl(url, ...params);

    let downloadLink = fileName

    // extract file name
    params.forEach(p => {
      let arr = p.split('=');
      if (arr.length > 1) {
        if (arr[0].toLowerCase() === 'filename') {
          downloadLink = arr[1];
        }
      }
    });

    // Set content type for: json / csv / xml / pdf /xslx
    let contentType = this.getMimeType(downloadLink);

    return this.http.get(resourceUrl, {
      responseType: 'blob',
      reportProgress: true,
      observe: 'events',
      headers: new HttpHeaders({ 'Content-Type': contentType })
    });
  }

  // Download2 is an alternative option to download
  download2(fileName: string, url: string, ...params: string[]) {

    let downloadLink = fileName

    // extract file name
    params.forEach(p => {
      let arr = p.split('=');
      if (arr.length > 1) {
        if (arr[0].toLowerCase() === 'filename') {
          downloadLink = arr[1];
        }
      }
    });

    let contentType = this.getMimeType(fileName);

    const link = document.createElement('a');
    link.href = this.buildUrl(url, ...params);
    link.download = downloadLink;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);

  }

  // HTTP GET action
  get<T>(url: string, ...params: string[]): Observable<T> {
    const resourceUrl = this.buildUrl(url, ...params);
    return this.http.get<T>(resourceUrl, httpOptions)
  }

  // HTTP POST action
  post<T>(url: string, body?: string, ...params: string[]): Observable<T> {
    const resourceUrl = this.buildUrl(url, ...params);
    return this.http.post<T>(resourceUrl, body, httpOptions)
  }

  // HTTP PUT action
  put<T>(url: string, body?: string, ...params: string[]): Observable<T> {
    const resourceUrl = this.buildUrl(url, ...params);
    return this.http.put<T>(resourceUrl, body, httpOptions)
  }

  // HTTP PATCH action
  patch<T>(url: string, body?: string, ...params: string[]): Observable<T> {
    const resourceUrl = this.buildUrl(url, ...params);
    return this.http.patch<T>(resourceUrl, body, httpOptions)
  }

  // HTTP DELETE action
  delete<T>(url: string, ...params: string[]): Observable<T> {
    const resourceUrl = this.buildUrl(url, ...params);
    return this.http.delete<T>(resourceUrl, httpOptions)
  }

  // Construct URL with parameters
  private buildUrl(url: string, ...params: string[]) {
    return (params === null) ? url : (params.length === 0) ? url : `${url}${params && params.length > 0 ? '?'+params.join('&') : ''}`;
  }

  // Return MIME type based on file extension
  private getMimeType(fileName: string): string {
    // Set content type for: json / csv / xml / pdf /xslx
    let contentType = 'application/json';
    if (fileName.toLowerCase().endsWith('jpg') || fileName.toLowerCase().endsWith('jpeg')) {
      contentType = 'image/jpeg';
    } else if (fileName.toLowerCase().endsWith('png')) {
      contentType = 'image/png';
    } else if (fileName.toLowerCase().endsWith('tif') || fileName.toLowerCase().endsWith('tiff')) {
      contentType = 'image/tiff';
    } else if (fileName.toLowerCase().endsWith('csv')) {
      contentType = 'text/csv';
    } else if (fileName.toLowerCase().endsWith('xml')) {
      contentType = 'text/xml';
    } else if (fileName.toLowerCase().endsWith('pdf')) {
      contentType = 'application/pdf';
    } else if (fileName.toLowerCase().endsWith('xlsx')) {
      contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
    }
    return contentType
  }
}
