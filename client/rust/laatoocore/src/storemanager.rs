use redux::{Action, Store, Dispatcher};
use event::{Event, EventListener, EventProducer};
use std::fmt::Debug;

#[derive(Debug)]
pub struct StoreChangeEvent {
    data: (),
}

impl StoreChangeEvent {
    fn new(storeState: ()) -> Self {
        return StoreChangeEvent{data: storeState};
    }
}

impl Event for StoreChangeEvent {

    fn get_type(&self) ->&'static str {
        return "Store Change Event"
    }
    fn get_source(&self) {

    }
    fn get_data(&self) {

    }    
}

#[allow(dead_code)]
pub struct StoreManager {
    state: Box<Store>,
    dispatching: bool,
    lsnrs: Vec<Box<EventListener>>,
}

impl Dispatcher for StoreManager {
    fn dispatch(&mut self, action: &Action) -> Result<(), String> {
        match self.state.reduce(action) {
            Ok(changed) => {
                let storeState = self.state.get_data();
                for ref lsnr in &self.lsnrs {
                    let evt = Box::new(StoreChangeEvent::new(storeState));
                    lsnr.on_event(evt);
                }
            }
            Err(e) => {
                return Err(format!("{}", e));
            }
        }
        Ok(())
    }
}

impl EventProducer for StoreManager {
    fn register_listener(&mut self, lsnr: Box<EventListener>) {
        self.lsnrs.push(lsnr)
    }
}

impl StoreManager {
    pub fn new(data: Box<Store>) -> Self {
        StoreManager{state: data, dispatching: false, lsnrs: vec![]}
    }
}