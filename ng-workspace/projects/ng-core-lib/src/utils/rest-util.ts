import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders, HttpRequest, HttpResponse } from '@angular/common/http';  // replaces previous Http service
import { map, catchError } from 'rxjs/operators';
import * as LocalStorageUtil from './localStorage-util';

/**
 * Utility class for all REST services with common functions
 */
@Injectable()
export class RestUtil {

  // Set headers
  private headers = new HttpHeaders().set('Content-Type', 'application/json');

  /**
   * Constructor with injected authentication service
   */
  constructor(private http: HttpClient) { }

  /**
   * Upload is HTTP POST action but the body is File object
   */
  upload(file: File, url: string, ...params: string[]) {

    const resourceUrl = this.buildUrl(url, ...params);

    const formData: FormData = new FormData();
    formData.append('fileKey', file, file.name);

    const req = new HttpRequest('POST', resourceUrl, formData, {
      reportProgress: false,
      responseType: 'json',
    });
    return this.http.request(req);
  }

  /**
   * Download is HTTP GET action but the content is blob
   */
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

    // return this.http.get(resourceUrl, {responseType: 'blob'}).subscribe((data) => {
    //   const downloadURL = window.URL.createObjectURL(data);
    //   const link = document.createElement('a');
    //   link.href = downloadURL;
    //   link.download = downloadLink;
    //
    //   link.click();
    // });

    // Set content type for: json / csv / xml / pdf
    let contentType = 'application/json';
    if (downloadLink.toLowerCase().endsWith('csv')) {
      contentType = 'text/csv';
    } else if (downloadLink.toLowerCase().endsWith('xml')) {
      contentType = 'text/xml';
    } else if (downloadLink.toLowerCase().endsWith('pdf')) {
      contentType = 'application/pdf';
    }

    return this.http.get(resourceUrl, {
      responseType: 'blob',
      reportProgress: true,
      observe: 'events',
      headers: new HttpHeaders({ 'Content-Type': contentType })
    });
  }
  
  /**
   * HTTP GET action
   */
  get(url: string, ...params: string[]) {
    const resourceUrl = this.buildUrl(url, ...params);
    return this.http
      .get(resourceUrl, { headers: this.headers, observe: 'response' })
      .pipe(
        map((res: HttpResponse<any>) => this.processResponse(res)),
        catchError(this.handleError),
      );
  }

  /**
   * HTTP POST action
   */
  post(url: string, body: string, ...params: string[]) {
    const resourceUrl = this.buildUrl(url, ...params);
    return this.http
      .post(resourceUrl, body, { headers: this.headers, observe: 'response' })
      .pipe(
        map((res: HttpResponse<any>) => this.processResponse(res)),
        catchError(this.handleError)
      );
  }

  /**
   * HTTP PUT action
   */
  put(url: string, body: string, ...params: string[]) {
    const resourceUrl = this.buildUrl(url, ...params);
    return this.http
      .put(resourceUrl, body, { headers: this.headers, observe: 'response' })
      .pipe(
        map((res: HttpResponse<any>) => this.processResponse(res)),
        catchError(this.handleError)
      );
  }

  /**
   * HTTP DELETE action
   */
  delete(url: string, ...params: string[]) {
    const resourceUrl = this.buildUrl(url, ...params);
    return this.http
      .delete(resourceUrl, { headers: this.headers, observe: 'response' })
      .pipe(
        map((res: HttpResponse<any>) => this.processResponse(res)),
        catchError(this.handleError)
      );
  }

  /**
   * Construct URL with parameters
   */
  private buildUrl(url: string, ...params: string[]) {
    return (params === null) ? url : (params.length === 0) ? url : `${url}${params && params.length > 0 ? '?'+params.join('&') : ''}`;
  }

  /**
   * Process the response, extract and refresh access token and return the body
   */
  private processResponse(response: HttpResponse<any>) {

    if (response.status === 401) {
      LocalStorageUtil.removeToken();
      throw new Error('Access denied, reset token: ' + response.status);
    } else if (response.status > 400) {
      throw new Error('HTTP status error: ' + response.status);
    }

    // Get access token from header and update authentication service

    const accessToken = response.headers.get('X-ACCESS-TOKEN');

    if ((accessToken !== null) && (accessToken.length > 0)) {
      LocalStorageUtil.setToken(accessToken);
    } 
    
    if (response.body && response.body.code && response.body.code !== 0) {
      throw { code: response.body.code, message: response.body.error };
    }
    
    return response.body;
  }

  /**
   * Error handling
   */
  private handleError(error: any): Promise<any> {
    if (error.code) {
      return Promise.reject(error);
    }
    return Promise.reject(error.message || error);
  }
}
