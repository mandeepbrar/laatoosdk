/*use std::marker::{Sync, Send};
//use std::error;
use platform;
use service::{Service, ServiceRequest};
use utils::{StringMap};
use std::collections::HashMap;
use registry::{Registry, RegistryStore, RegisteredItem};


pub struct App {
    app_platform: Option<Box<platform::Platform + Sync + Send>>,
    registries: HashMap<Registry, RegistryStore>
}

impl App {
    pub fn new() -> App {
        App{app_platform: Option::None, registries: HashMap::new()}
    }

     #[allow(dead_code)]
    pub fn initialize(&mut self, pfm: Box<platform::Platform + Sync + Send>) {
        self.app_platform = Option::Some(pfm);
      //  self.registries = HashMap::new();

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
    }

    #[allow(dead_code)]
    pub fn execute_service(&self, _service_name: String, _service_request: ServiceRequest, _config: Option<StringMap>) {

    }
}

/*
#[cfg_attr(target_arch = "wasm32", wasm_bindgen)]

#[cfg(target_arch = "wasm32")]
extern crate wasm_bindgen;

#[cfg(target_arch = "wasm32")]
use wasm_bindgen::prelude::*;
*/