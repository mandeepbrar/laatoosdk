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
}

pub type StringMap = HashMap<String, StringMapValue>;
//pub type StringMapItem = Value;



impl Constructible for StringMap {
    fn new() -> StringMap {
        let mut string_map: StringMap;
        string_map = HashMap::new();
        return string_map;
    }
}