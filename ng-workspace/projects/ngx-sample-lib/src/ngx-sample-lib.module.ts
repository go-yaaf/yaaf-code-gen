import { NgModule, ModuleWithProviders } from '@angular/core';
import { CommonModule } from '@angular/common';
import { HttpClientModule } from '@angular/common/http';
import { RestUtils } from './rest-utils';
import { AppConfig } from './config';

@NgModule({
  imports: [CommonModule, HttpClientModule]
})
export class NgxSampleLibModule {
  static forRoot(config: AppConfig): ModuleWithProviders<NgxSampleLibModule> {
    return {
      ngModule: NgxSampleLibModule,
      providers: [
        { provide: 'config', useValue: config }
      ]
    };
  }
}

// Since all services are now marked as providedIn: 'root', no need to include them in providers list
/*
@NgModule({
  imports: [CommonModule, HttpClientModule]
})
export class NgxSampleLibModule {
  static forRoot(config: AppConfig): ModuleWithProviders<NgxSampleLibModule> {
    return {
      ngModule: NgxSampleLibModule,
      providers: [
        { provide: 'config', useValue: config },
        RestUtils,
        ...Services
      ]
    };
  }
}
*/
