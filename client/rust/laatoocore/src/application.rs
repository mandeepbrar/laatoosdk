#[cfg(target_arch = "wasm32")]
extern crate wasm_bindgen;

#[cfg(target_arch = "wasm32")]
use wasm_bindgen::prelude::*;

use std::sync::{Mutex};
use std::marker::{Sync, Send};
//use std::error;
use platform;
use service::{Service, ServiceRequest};
use utils::{StringMap};


lazy_static! {
    pub static ref Application: Mutex<Box<App>> = Mutex::new(Box::new(App{app_platform: Option::None}));
}


#[cfg_attr(target_arch = "wasm32", wasm_bindgen)]
pub struct App {
    app_platform: Option<Box<platform::Platform + Sync + Send>>
}

impl App {
     #[allow(dead_code)]
    pub fn initialize(&mut self, pfm: Box<platform::Platform + Sync + Send>) {
        self.app_platform = Option::Some(pfm)
    }
    
     #[allow(dead_code)]
    pub fn execute_service_object(_svc: Service, _service_request: ServiceRequest, _config: Option<StringMap>) {
        /*var method = get_method(service);
        var req = service_request.get_method_object("http");
        var url = this.getURL(service, req);
        return this.HttpCall(url, method, req.params, req.data, req.headers, config);*/
    }

    #[allow(dead_code)]
    pub fn execute_service(_service_name: String, _service_request: ServiceRequest, _config: Option<StringMap>) {

    }
}
