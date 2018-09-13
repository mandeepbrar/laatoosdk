use service::{Service};
use std::collections::HashMap;


#[derive(Hash, Eq, PartialEq, Debug)]
pub enum Registry {

}


#[derive(Debug)]
pub enum RegisteredItem {
    String(String),
    Service(Service),
}

pub struct RegistryStore {
    registry: HashMap<String, RegisteredItem>
}

impl RegistryStore {
    pub fn register(&mut self, item_name: String, item: RegisteredItem) {
        self.registry.insert(item_name, item);
    }
    pub fn get_registered_item(&self, item_name: String) -> Option<&RegisteredItem> {
        self.registry.get(&item_name)
    }
}
