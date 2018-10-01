#![feature(associated_type_defaults)]
#[macro_use]
extern crate serde_derive;
#[macro_use]
extern crate lazy_static;
#[cfg(target_arch = "wasm32")]
extern crate wasm_bindgen;

pub mod utils;
pub mod platform;
pub mod service;
pub mod application;
pub mod http;
pub mod redux;
//pub mod context;
//mod app;
mod registry;
pub mod error;
pub mod event;
mod test;
mod storemanager;


