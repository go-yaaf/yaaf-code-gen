import { NgModule, ModuleWithProviders } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Services } from './lib/services/services.export';
import { RestUtil } from './utils/rest-util';
import { HttpClientModule } from '@angular/common/http';
import { CoreConfig } from './config';

@NgModule({
  imports: [CommonModule, HttpClientModule]
})
export class CoreLibModule {
  static forRoot(config: CoreConfig): ModuleWithProviders<CoreLibModule> {
    // console.log(config);
    return {
      ngModule: CoreLibModule,
      providers: [
        { provide: 'config', useValue: config },
        RestUtil,
        ...Services
      ]
    };
  }
}
