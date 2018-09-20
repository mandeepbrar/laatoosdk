//use std::marker::{Sync, Send};
//use std::error;
use platform;
use service::{Service, ServiceRequest};
use utils::{StringMap};
use std::collections::HashMap;
use registry::{Registry, RegistryStore, RegisteredItem};
use action::{Action};
use store::{Store, StoreData};
use reducer::{Reducer};
use std::any::Any;

#[cfg(target_arch = "wasm32")]
use wasm_bindgen::prelude::*;

#[cfg_attr(target_arch = "wasm32", wasm_bindgen)]
pub struct Application {
    app_platform: Box<platform::Platform>,
    registries: HashMap<Registry, RegistryStore>,
    //store: LaatooStore,
}

#[cfg(target_arch = "wasm32")]
#[wasm_bindgen]
extern {
    #[wasm_bindgen(js_namespace = console)]
    fn log(msg: &str);
}

impl Application {
    pub fn new(pfm: Box<platform::Platform>) -> Application {
        Application{app_platform: pfm, registries: HashMap::new()}
    }
    
    #[allow(dead_code)]
    pub fn register(&mut self, registry: Registry, item_name: String, item: RegisteredItem) {
        let registry_store = self.registries.get_mut(&registry).unwrap();
        registry_store.register(item_name, item);
    }
  
    #[allow(dead_code)]
    pub fn get_registered_item(&self, registry: Registry, item_name: String) -> Option<&RegisteredItem> {
        let registry_store = self.registries.get(&registry).unwrap();
        registry_store.get_registered_item(item_name)
    }

    #[allow(dead_code)]
    pub fn execute_service_object(&self, _svc: Service, _service_request: ServiceRequest, _config: Option<StringMap>) {
        /*var method = get_method(service);
        var req = service_request.get_method_object("http");
        var url = this.getURL(service, req);
        return this.HttpCall(url, method, req.params, req.data, req.headers, config);*/
    }

    #[allow(dead_code)]
    pub fn execute_service(&self, _service_name: String, _service_request: ServiceRequest, _config: Option<StringMap>) {

    }

    pub fn register_store(&self, store: Box<Store>, action: Box<Action>) {

    }

    /*pub fn dispatch(&self, action: Action) -> Result<T::Action, String> {
    }*/

}

#[cfg(target_arch = "wasm32")]
#[wasm_bindgen]
impl Application {
    #[allow(dead_code)]
    pub fn js_get_registered_item(&self, registry: String, item_name: String) -> String {
        log(&registry);
        String::from("s: Cow<'a, str>s: Cow<'a, str>")

    }
}
/*
#[cfg_attr(target_arch = "wasm32", wasm_bindgen)]

#[cfg(target_arch = "wasm32")]
extern crate wasm_bindgen;

#[cfg(target_arch = "wasm32")]
use wasm_bindgen::prelude::*;
*/