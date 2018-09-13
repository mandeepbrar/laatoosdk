use http;
use utils;

#[derive(Debug)]
pub enum Service {
    Http(String, http::HttpMethod, http::HttpRequest),
}


pub enum ServiceRequest<'a> {
    Http(String, &'a http::HttpRequest),
}

pub enum Response {
    Success(utils::StringMap, utils::StringMap, u32),
    Error(String, utils::StringMap, u32),
}

/*
#[cfg(target_arch = "wasm32")]
extern crate wasm_bindgen;

#[cfg(target_arch = "wasm32")]
use wasm_bindgen::prelude::*;*/