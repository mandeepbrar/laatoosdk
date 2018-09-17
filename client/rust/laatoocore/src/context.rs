use application::{Application};
use platform::{Platform};

#[cfg(target_arch = "wasm32")]
use wasm_bindgen::prelude::*;

#[cfg_attr(target_arch = "wasm32", wasm_bindgen)]
pub struct Context {
    appl: Application
}

impl Context {
    #[allow(dead_code)]
    pub fn new(pfm: Box<Platform>) -> Context {
        Context{appl: Application::new(pfm)}
    }
}