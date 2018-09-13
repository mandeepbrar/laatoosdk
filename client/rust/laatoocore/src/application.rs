use std::sync::{Mutex};
use app::App;
use std::collections::HashMap;
use std::marker::{Sync, Send};
use platform;
use registry::{Registry, RegisteredItem};

lazy_static! {
    static ref app_obj: Mutex<Box<App>> = Mutex::new(Box::new(App::new()));
}

pub fn initialize(pfm: Box<platform::Platform + Sync + Send>) {
    app_obj.lock().unwrap().initialize(pfm);

}

#[allow(dead_code)]
pub fn register(reg: Registry, item_name: String, item: RegisteredItem) {
    app_obj.lock().unwrap().register(reg, item_name, item);
}

#[allow(dead_code)]
pub fn get_registered_item<'a>(reg: Registry, item_name: String) -> Option<&'a RegisteredItem> {
    let appobj = app_obj.lock();
    let res = appobj.unwrap().get_registered_item(reg, item_name);
    res
}



/*
#[cfg_attr(target_arch = "wasm32", wasm_bindgen)]

#[cfg(target_arch = "wasm32")]
extern crate wasm_bindgen;

#[cfg(target_arch = "wasm32")]
use wasm_bindgen::prelude::*;
*/