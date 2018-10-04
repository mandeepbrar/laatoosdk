use platform;
use service::{Service, ServiceRequest};
use utils::{StringMap};
use std::collections::HashMap;
use registry::{Registry, RegistryStore, RegisteredItem};
use redux::{Action, Store, Reducer, Dispatcher, StoreListener};
use storemanager::{StoreManager};
use event::{EventProducer, EventListener};
use std::any::Any;
use std::rc::Rc;
use std::cell::RefCell;
use std::borrow::BorrowMut;

pub struct _Application {
    app_platform: Box<platform::Platform>,
    registries: HashMap<Registry, RegistryStore>,
    dispatchers: HashMap<&'static str, Rc<RefCell<Dispatcher>>>,
    stores: HashMap<&'static str, Rc<RefCell<StoreManager>>>
}

impl _Application {
    pub fn new(pfm: Box<platform::Platform>) -> Box<_Application> {
        let mut app = Box::new(_Application{app_platform: pfm, registries: HashMap::new(), dispatchers: HashMap::new(), stores: HashMap::new()});
        app
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

    pub fn register_store(&mut self, store: Box<Store>, action_type: &'static str) {
        let id = store.get_id();//.clone();
        let mgr = Rc::new(RefCell::new(StoreManager::new(store)));
        self.dispatchers.insert(action_type, mgr.clone());
        self.stores.insert(id, mgr);
    }

    pub fn register_listener(&mut self, store_id: &str, lsnr: StoreListener) {
        match self.stores.get(store_id) {
            Some(stor) => {
                let prod = stor.clone();
                let mut val1 = (*prod).borrow_mut();
                (*val1).register_listener(lsnr);
            },
            None => {}
        }
    }


    pub fn dispatch(&mut self, action: &Action) -> Result<(), String> {
       match self.dispatchers.get(action.get_type()) {
            Some(dispatcher) => {
                let disp = dispatcher.clone();
                let mut val1 = (*disp).borrow_mut();
                let val = (*val1).dispatch(action);
            },
            None => {}
        }
        Ok(())
    }

}
