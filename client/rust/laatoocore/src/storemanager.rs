use redux::{Action, Store, Dispatcher, StoreListener};
use event::{Event, EventListener, EventProducer};
use std::fmt::Debug;
use std::any::Any;


#[allow(dead_code)]
pub struct StoreManager {
    state: Box<Store>,
    dispatching: bool,
    lsnrs: Vec<StoreListener>,
}

impl Dispatcher for StoreManager {
    fn dispatch(&mut self, action: &Action) -> Result<(), String> {
        match self.state.reduce(action) {
            Ok(changed) => {
                //let storeState = self.state.get_data();
                self.inform_listeners();
            }
            Err(e) => {
                return Err(format!("{}", e));
            }
        }
        Ok(())
    }
}

/*impl EventProducer for StoreManager {
}*/

impl StoreManager {
    pub fn new(data: Box<Store>) -> Self {
        StoreManager{state: data, dispatching: false, lsnrs: vec![]}
    }
    pub fn register_listener(&mut self, lsnr: StoreListener) {
        self.lsnrs.push(lsnr);
    }
    fn inform_listeners(&self) {
        //let k: &'b  = self.state;
        //let evt = StoreChangeEvent::new(k);
       // let evtref: &'b StoreChangeEvent<'b> = &evt;
       // let evtref: Box<Event> = Box::new(evt);
        for ref lsnr in &self.lsnrs {
            lsnr(&self.state);
        }
    }
}

