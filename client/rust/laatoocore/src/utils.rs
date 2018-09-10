extern crate serde;

use std::collections::HashMap;

pub trait Constructible {
    fn new() -> Self;
}

#[derive(Serialize, Deserialize, Debug, PartialEq)]
pub enum StringMapValue {
    Null,
    Bool(bool),
    Int(i64),
    Float(f64),
    String(String),
    Array(Vec<StringMapValue>),
    Object(HashMap<String, StringMapValue>),
    StringsMap(StringMap),
}

pub type StringMap = HashMap<String, StringMapValue>;
//pub type StringMapItem = Value;
pub type StringsMap = HashMap<String, String>;



impl Constructible for StringMap {
    fn new() -> StringMap {
        let string_map: StringMap;
        string_map = HashMap::new();
        return string_map;
    }
}

impl Constructible for StringsMap {
    fn new() -> StringsMap {
        let strings_map: StringsMap;
        strings_map = HashMap::new();
        return strings_map;
    }
}

/*
#[cfg(target_arch = "wasm32")]
use wasm_bindgen::prelude::*;
#[cfg_attr(target_arch = "wasm32", wasm_bindgen)]
*/