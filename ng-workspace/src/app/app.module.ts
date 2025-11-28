import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';

import { AppComponent } from './app.component';
import { CoreLibModule } from '@agentvi/ng-core-lib';
import { HttpClientModule } from '@angular/common/http';

@NgModule({
  declarations: [
    AppComponent
  ],
  imports: [
    BrowserModule,
    HttpClientModule,
    CoreLibModule.forRoot({
      api: 'http://localhost:8080/v1'
    })
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
